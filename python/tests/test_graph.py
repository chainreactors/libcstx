import ast
import inspect
import json
from pathlib import Path

import cstxpy
import pytest

from cstxpy import (
    CSTX,
    CSTXGraph,
    CSTXError,
    GraphCursor,
    NodeFlags,
    Repository,
    Schemas,
)


SCHEMA = {"properties": {"ip": {"type": "string"}}}


def node(value: str, *, flags: int = 0) -> dict:
    return {
        "id": f"ip:{value}",
        "type": "ip",
        "value": value,
        "model": {"ip": value, "cstx_flags": flags},
        "sources": ["test"],
        "extras": {},
    }


def edge(source: str, target: str) -> dict:
    return {
        "id": f"relationship:{source}:related:{target}",
        "source_id": source,
        "target_id": target,
        "relation_type": "related",
        "sources": ["test"],
        "attrs": {},
    }


def db() -> CSTX:
    value = CSTX()
    value.schemas.register("ip", SCHEMA, "ip")
    return value


def test_root_services_and_native_values():
    value = db()
    assert value.graph.add_nodes([node("1.1.1.1")]) == 1
    assert value.graph.add_nodes([node("1.1.1.1")]) == 0
    assert value.graph.node_count() == 1
    assert value.graph.node("ip:1.1.1.1")["model"]["ip"] == "1.1.1.1"
    assert value.last_change()["updated_node_ids"] == []


def test_easm_metadata_requires_explicit_loading():
    value = CSTX()
    assert "easm" in value.schemas.available_plugins()
    assert "gogo" in value.schemas.plugin_artifacts("easm")
    assert not value.schemas.has_native_artifact("gogo")

    value.schemas.load_plugin("easm")
    assert value.schemas.has_native_artifact("gogo")


def test_cursor_filter_order_close_and_invalidation():
    value = db()
    value.graph.add_nodes([node("2.2.2.2"), node("1.1.1.1", flags=NodeFlags.HONEYPOT)])
    cursor = value.graph.nodes(order="id_asc")
    assert next(cursor)["id"] == "ip:1.1.1.1"
    value.graph.add_nodes([node("3.3.3.3")])
    with pytest.raises(CSTXError) as error:
        next(cursor)
    assert error.value.code == "CURSOR_INVALIDATED"
    selected = list(value.graph.nodes(flags_any=NodeFlags.HONEYPOT))
    assert [item["id"] for item in selected] == ["ip:1.1.1.1"]
    selected_cursor = value.graph.nodes(limit=1)
    selected_cursor.close()
    assert selected_cursor.closed
    assert list(selected_cursor) == []


def test_noop_cstx_flag_filter_does_not_materialize_snapshot():
    value = db()
    value.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2", flags=NodeFlags.INTERNAL)])

    result = value.graph.filter(exclude_mask=NodeFlags.HONEYPOT)

    assert isinstance(result, CSTX)
    assert result.graph.node_count() == 2


def test_stats_and_query_subgraph_accept_cstx_flag_masks_without_parallel_apis():
    value = db()
    value.graph.add_nodes(
        [
            node("1.1.1.1"),
            node("2.2.2.2", flags=NodeFlags.HONEYPOT),
        ]
    )

    assert value.graph.stats()["nodes"] == {"ip": 2}
    assert value.graph.stats(exclude_mask=NodeFlags.HONEYPOT)["nodes"] == {"ip": 1}

    all_matches = value.graph.query_subgraph("ip")
    assert all_matches.graph.node_count() == 2
    assert all_matches.graph.edge_count() == 0

    filtered = value.graph.query_subgraph("ip", exclude_mask=NodeFlags.HONEYPOT)
    assert filtered.graph.node_count() == 1
    assert filtered.graph.find_node("ip:1.1.1.1") is not None

    assert not hasattr(value.graph, "stats_filtered")
    assert not hasattr(value.graph, "query_trace_ids_filtered")
    for internal_name in (
        "node_ids",
        "nodes_by_ids",
        "link_nodes",
        "subgraph_ids",
        "query_node_ids",
        "query_trace_ids",
        "induced_snapshot",
        "_node_ids",
        "_nodes_by_ids",
        "_link_nodes",
        "_subgraph_ids",
        "_query_node_ids",
        "_query_trace_ids",
        "_induced_snapshot",
    ):
        assert not hasattr(value.graph, internal_name)


def test_stats_selection_counts_induced_edges_without_materializing_a_subgraph():
    value = db()
    selected_a = node("1.1.1.1")
    selected_a["extras"] = {"flow_ids": ["flow-a"]}
    selected_b = node("2.2.2.2")
    selected_b["extras"] = {"flow_ids": ["flow-a"]}
    outside = node("3.3.3.3")
    outside["extras"] = {"flow_ids": ["flow-b"]}
    value.graph.add_nodes([selected_a, selected_b, outside])
    value.graph.add_edges(
        [
            edge("ip:1.1.1.1", "ip:2.2.2.2"),
            edge("ip:2.2.2.2", "ip:3.3.3.3"),
        ]
    )

    stats = value.graph.stats(selection='*[flow_ids=="flow-a"]')

    assert stats["nodes"] == {"ip": 2}
    assert stats["edges"] == {"related": 1}
    assert stats["sources"] == {"test": 2}


def test_analyze_leiden_returns_complete_selected_partition():
    value = db()
    node_ids = ["ip:a1", "ip:a2", "ip:a3", "ip:b1", "ip:b2", "ip:b3"]
    value.graph.add_nodes([node(node_id.removeprefix("ip:")) for node_id in node_ids])
    value.graph.add_edges(
        [
            edge("ip:a1", "ip:a2"),
            edge("ip:a2", "ip:a3"),
            edge("ip:a3", "ip:a1"),
            edge("ip:b1", "ip:b2"),
            edge("ip:b2", "ip:b3"),
            edge("ip:b3", "ip:b1"),
        ]
    )

    result = value.graph.analyze({"name": "leiden", "resolution": 1.0})
    assignment_page = result.page(limit=100, page=1)
    summary = assignment_page["summary"]

    assert summary["algorithm"] == "leiden"
    assert summary["projection"] == "undirected"
    assert summary["resolution"] == 1.0
    assert summary["num_communities"] >= 2
    assert summary["total_communities"] == summary["num_communities"]
    assert summary["communities_truncated"] is False
    assert sum(summary["community_sizes"].values()) == len(node_ids)
    assignments = {
        item["node_id"]: item["community"] for item in assignment_page["items"]
    }
    assert set(assignments) == set(node_ids)
    assert assignments["ip:a1"] != assignments["ip:b1"]

    limited = value.graph.analyze(
        {"name": "leiden", "resolution": 1.0, "top_k": 1},
        'ip[ip!="a1"]',
    )
    limited_page = limited.page(limit=100, page=1)
    limited_summary = limited_page["summary"]
    assert limited_summary["num_communities"] == 1
    assert limited_summary["communities_truncated"] is True
    assert len(limited_summary["community_sizes"]) == 1
    assert all(item["node_id"] != "ip:a1" for item in limited_page["items"])

    with pytest.raises(CSTXError):
        value.graph.analyze({"name": "leiden", "resolution": 0.0})
    with pytest.raises(CSTXError):
        value.graph.analyze({"name": "leiden", "top_k": 0})


def test_analyze_is_the_single_typed_algorithm_atom():
    value = db()
    value.graph.add_nodes([node(name) for name in ("a", "b", "c", "d")])
    value.graph.add_edges(
        [
            edge("ip:a", "ip:b"),
            edge("ip:a", "ip:c"),
            edge("ip:b", "ip:d"),
            edge("ip:c", "ip:d"),
        ]
    )

    assert value.graph.analyze({"name": "is_dag"}) is True
    weak = value.graph.analyze({"name": "weak_components"})
    assert weak.page(limit=10, page=1)["summary"]["projection"] == "undirected"
    strong = value.graph.analyze({"name": "strong_components"})
    assert "projection" not in strong.page(limit=10, page=1)["summary"]
    paths = value.graph.analyze(
        {
            "name": "shortest_paths",
            "start_id": "ip:a",
            "end_id": "ip:d",
            "direction": "both",
            "limit": 10,
        }
    )
    assert isinstance(paths, GraphCursor)
    assert paths.kind == "paths"
    assert paths.page(limit=10, page=1)["items"] == [
        {"node_ids": ["ip:a", "ip:b", "ip:d"]},
        {"node_ids": ["ip:a", "ip:c", "ip:d"]},
    ]

    with pytest.raises(ValueError):
        value.graph.analyze({"name": "is_dag", "unexpected": True})


def test_prefetched_cursor_page_is_invalidated_before_returning_stale_items():
    value = db()
    value.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2"), node("3.3.3.3")])
    cursor = value.graph.nodes(order="id_asc")
    assert next(cursor)["id"] == "ip:1.1.1.1"
    assert next(cursor)["id"] == "ip:2.2.2.2"
    value.graph.add_nodes([node("4.4.4.4")])

    with pytest.raises(CSTXError) as error:
        next(cursor)
    assert error.value.code == "CURSOR_INVALIDATED"


def test_unchanged_edge_does_not_invalidate_live_cursor():
    value = db()
    value.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2"), node("3.3.3.3")])
    relation = edge("ip:1.1.1.1", "ip:2.2.2.2")
    value.graph.add_edges([relation])
    cursor = value.graph.nodes()
    next(cursor)
    assert value.graph.add_edges([relation]) == 0

    assert next(cursor)["id"] == "ip:2.2.2.2"


def test_edges_neighbors_and_direct_json():
    value = db()
    value.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2")])
    relation = edge("ip:1.1.1.1", "ip:2.2.2.2")
    assert value.graph.add_edges([relation]) == 1
    assert value.graph.edge_count() == 1
    assert list(value.graph.edges())[0]["id"] == relation["id"]
    assert list(value.graph.neighbors("ip:1.1.1.1"))[0]["id"] == "ip:2.2.2.2"
    assert json.loads(value.graph._edges_json()) == list(value.graph.edges())
    assert json.loads(value.graph._neighbors_json("ip:1.1.1.1")) == list(
        value.graph.neighbors("ip:1.1.1.1")
    )


def test_native_graph_semantics_own_lookup_context_and_relationship_identity():
    value = db()
    first = node("2.2.2.2")
    first["extras"] = {"name": "shared", "scope": "old"}
    second = node("1.1.1.1")
    second["extras"] = {"name": "shared"}
    duplicate_value = node("3.3.3.3")
    duplicate_value["value"] = "1.1.1.1"
    value.graph.add_nodes([first, second, duplicate_value])

    assert value.graph.find_node("ip:2.2.2.2")["id"] == "ip:2.2.2.2"
    assert value.graph.find_node("1.1.1.1")["id"] == "ip:1.1.1.1"
    assert value.graph.find_node("shared")["id"] == "ip:1.1.1.1"
    assert (
        value.graph.patch_node_extras(["ip:2.2.2.2"], {"scope": "new", "run": "r1"})
        == 1
    )
    assert value.graph.node("ip:2.2.2.2")["extras"] == {
        "name": "shared",
        "scope": ["old", "new"],
        "run": "r1",
    }

    relation = value.graph.create_relationship(
        "ip:1.1.1.1",
        "ip:2.2.2.2",
        "related",
        ["test"],
        {"weight": 1},
        "qualified",
    )
    assert relation["id"] == ("relationship:ip:1.1.1.1:related:ip:2.2.2.2:qualified")
    assert cstxpy.is_path_expression("ip -> ip")
    assert cstxpy.is_path_expression('ip[country="CN"]')


def test_patch_node_extras_none_selects_all_nodes():
    value = db()
    value.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2")])

    assert value.graph.patch_node_extras(None, {"task_ids": ["task-1"]}) == 2
    assert value.graph.node("ip:1.1.1.1")["extras"]["task_ids"] == ["task-1"]
    assert value.graph.node("ip:2.2.2.2")["extras"]["task_ids"] == ["task-1"]


def test_native_graph_union_and_difference_are_deterministic():
    left = db()
    right = db()
    left.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2")])
    right.graph.add_nodes([node("2.2.2.2"), node("3.3.3.3")])

    union = left.graph.union(right.graph)
    assert sorted(item["id"] for item in union.graph.nodes()) == [
        "ip:1.1.1.1",
        "ip:2.2.2.2",
        "ip:3.3.3.3",
    ]
    difference = left.graph.difference(right.graph)
    assert [item["id"] for item in difference.graph.nodes()] == ["ip:1.1.1.1"]
    assert list(difference.graph.edges()) == []


def test_native_graph_merge_is_in_place_and_reports_exact_changes():
    target = db()
    source = db()
    first = node("1.1.1.1")
    second = node("2.2.2.2")
    first["sources"] = ["target"]
    target.graph.add_nodes([first])
    merged_first = node("1.1.1.1")
    merged_first["sources"] = ["source-a", "source-b"]
    merged_first["extras"] = {"scope": "source"}
    second["sources"] = ["source-b"]
    source.graph.add_nodes([merged_first, second])
    relationship = edge(merged_first["id"], second["id"])
    relationship["sources"] = ["source-a", "source-b"]
    relationship["attrs"] = {"confidence": 90}
    source.graph.add_edges([relationship])

    assert target.graph.merge(source.graph) == 3
    assert target.graph.node_count() == 2
    assert target.graph.edge_count() == 1
    assert source.graph.node_count() == 2
    assert source.graph.edge_count() == 1
    assert set(target.graph.node(first["id"])["sources"]) == {
        "target",
        "source-a",
        "source-b",
    }
    merged_edge = target.graph.edge(relationship["id"])
    assert merged_edge["source_id"] == first["id"]
    assert merged_edge["target_id"] == second["id"]
    assert set(merged_edge["sources"]) == {"source-a", "source-b"}
    assert merged_edge["attrs"] == {"confidence": 90}
    assert target.last_change() == {
        "added_node_ids": [second["id"]],
        "updated_node_ids": [first["id"]],
        "removed_node_ids": [],
        "added_edge_ids": [relationship["id"]],
        "updated_edge_ids": [],
        "removed_edge_ids": [],
        "reset": False,
    }
    assert target.graph.merge(target.graph) == 0


def test_direct_dict_cursor_matches_cstx_json_for_all_column_types():
    value = CSTX(cursor_page_size=2)
    value.schemas.register(
        "mixed",
        {
            "properties": {
                "name": {"type": "string"},
                "optional": {"type": "string"},
                "count": {"type": "integer"},
                "ratio": {"type": "number"},
                "active": {"type": "boolean"},
                "tags": {"type": "array", "items": {"type": "string"}},
                "ports": {"type": "array", "items": {"type": "integer"}},
                "details": {"type": "object"},
            }
        },
        "name",
    )
    item = {
        "id": "mixed:one",
        "type": "mixed",
        "value": "one",
        "model": {
            "name": "one",
            "optional": None,
            "count": 7,
            "ratio": 1.5,
            "active": True,
            "tags": ["a", "b"],
            "ports": [80, 443],
            "details": {"nested": [1, 2]},
        },
        "sources": ["test"],
        "extras": {"scope": "unit"},
    }
    value.graph.add_nodes([item])

    assert list(value.graph.nodes()) == json.loads(value.graph._nodes_json())


def test_repo_checkout_and_close():
    value = db()
    value.graph.add_nodes([node("1.1.1.1")])
    commit = value.repo.commit("initial")
    value.graph.add_nodes([node("2.2.2.2")])
    value.repo.checkout(commit["id"], force=True)
    assert value.graph.node("ip:1.1.1.1")["id"] == "ip:1.1.1.1"
    assert value.graph.node_count() == 1
    value.close()
    with pytest.raises(CSTXError) as error:
        value.graph.node_count()
    assert error.value.code == "NOT_INITIALIZED"


def test_repo_prepare_uses_bytes_for_object_transport():
    value = db()
    value.graph.add_nodes([node("1.1.1.1")])
    try:
        commit, index_root, objects = value.repo._prepare(
            "transport", "main", None, {}, 1
        )
        assert isinstance(commit, dict)
        assert isinstance(index_root, bytes)
        assert objects
        assert all(
            isinstance(object_id, bytes) and isinstance(envelope, bytes)
            for object_id, _kind, envelope in objects
        )
    finally:
        value.repo._discard()
        value.close()


def test_repo_stats_is_one_time_range_delta():
    value = db()
    value.graph.add_nodes([node("1.1.1.1")])
    value.repo.commit("first", timestamp=100)
    value.graph.add_nodes([node("2.2.2.2")])
    value.repo.commit("second", timestamp=200)

    assert value.repo.delta(start_timestamp=100, end_timestamp=199)["added_nodes"] == 1
    assert value.repo.delta(start_timestamp=101, end_timestamp=200)["added_nodes"] == 1
    assert "bucket" not in inspect.signature(value.repo.delta).parameters
    with pytest.raises(CSTXError, match="start_timestamp"):
        value.repo.delta(start_timestamp=201, end_timestamp=200)


def test_not_found_has_stable_error_context():
    value = db()
    with pytest.raises(CSTXError) as error:
        value.graph.node("ip:missing")
    assert error.value.code == "NOT_FOUND"
    assert error.value.operation == "graph.node"


def test_repository_errors_have_stable_context():
    value = db()

    with pytest.raises(CSTXError) as error:
        value.repo.checkout("missing")
    assert error.value.code == "NOT_FOUND"
    assert error.value.operation == "repo.checkout"


def test_every_supported_api_has_runtime_documentation():
    """Keep newly exported APIs explainable through Python help(), not only source."""
    public_members = {
        CSTX: (
            "schemas",
            "graph",
            "repo",
            "closed",
            "project_id",
            "close",
            "last_change",
            "__enter__",
            "__exit__",
        ),
        Schemas: (
            "register",
            "register_join_rule",
            "import_schema",
            "export_schema",
            "contains",
            "get",
            "list",
            "load_plugin",
            "load_all_plugins",
            "available_plugins",
            "plugin_artifacts",
            "has_native_artifact",
            "anchor_concepts",
        ),
        CSTXGraph: (
            "rag",
            "analyze",
            "add_nodes",
            "add_edges",
            "node",
            "edge",
            "create_relationship",
            "union",
            "merge",
            "difference",
            "contains",
            "node_count",
            "edge_count",
            "stats",
            "nodes",
            "nodes_page",
            "edges",
            "neighbors",
            "query",
            "ingest_native",
            "find_node",
            "patch_node_extras",
            "node_types",
            "link",
            "update_node_flags",
            "analyze",
            "degree",
            "subgraph",
            "query_subgraph",
            "induced_subgraph",
            "filter",
            "filter_with_reasons",
            "find_anchors",
            "elevate",
        ),
        Repository: (
            "resolve",
            "head",
            "checkout",
            "commit",
            "diff",
            "diff_stat",
            "log",
            "history",
            "branch",
            "merge",
            "stat",
            "delta",
        ),
        GraphCursor: (
            "kind",
            "page",
            "closed",
            "close",
            "__iter__",
            "__next__",
            "__enter__",
            "__exit__",
        ),
        NodeFlags: ("all_mask", "default_exclude_mask"),
    }

    assert inspect.getdoc(CSTXError)
    for api_type, members in public_members.items():
        assert inspect.getdoc(api_type), api_type.__name__
        for member in members:
            assert inspect.getdoc(getattr(api_type, member)), (
                f"{api_type.__name__}.{member}"
            )
        discovered = {
            name
            for name, value in api_type.__dict__.items()
            if not name.startswith("_")
            and (
                callable(getattr(api_type, name, None))
                or inspect.isdatadescriptor(value)
            )
        }
        assert discovered <= set(members), (
            f"undocumented API added to {api_type.__name__}: {discovered - set(members)}"
        )


def test_type_stub_documents_every_exported_class_and_method():
    """IDE-visible signatures must carry the same explanation as runtime help()."""
    stub = Path(cstxpy.__file__).with_name("_cstxpy.pyi")
    tree = ast.parse(stub.read_text(encoding="utf-8"), filename=str(stub))
    for node in ast.walk(tree):
        if isinstance(node, (ast.ClassDef, ast.FunctionDef, ast.AsyncFunctionDef)):
            assert ast.get_docstring(node), f"missing stub documentation: {node.name}"
