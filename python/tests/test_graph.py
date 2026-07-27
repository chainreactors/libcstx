import ast
import inspect
import json
from pathlib import Path

import cstxpy
import pytest

from cstxpy import (
    CSTX,
    CSTXError,
    EdgeCursor,
    Graph,
    NodeCursor,
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
    assert value.graph.add_node(node("1.1.1.1")) == 1
    assert value.graph.add_node(node("1.1.1.1")) == 0
    assert value.graph.node_count() == 1
    assert value.graph.node("ip:1.1.1.1")["model"]["ip"] == "1.1.1.1"
    assert value.last_change()["updated_node_ids"] == []


def test_cursor_filter_order_close_and_invalidation():
    value = db()
    value.graph.add_nodes([node("2.2.2.2"), node("1.1.1.1", flags=NodeFlags.HONEYPOT)])
    cursor = value.graph.nodes(order="id_asc")
    assert next(cursor)["id"] == "ip:1.1.1.1"
    value.graph.add_node(node("3.3.3.3"))
    with pytest.raises(CSTXError) as error:
        next(cursor)
    assert error.value.code == "CURSOR_INVALIDATED"
    selected = list(value.graph.nodes(flags_any=NodeFlags.HONEYPOT))
    assert [item["id"] for item in selected] == ["ip:1.1.1.1"]
    selected_cursor = value.graph.nodes(limit=1)
    selected_cursor.close()
    assert selected_cursor.closed
    assert list(selected_cursor) == []


def test_prefetched_cursor_page_is_invalidated_before_returning_stale_items():
    value = db()
    value.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2"), node("3.3.3.3")])
    cursor = value.graph.nodes(page_size=3, order="id_asc")
    assert next(cursor)["id"] == "ip:1.1.1.1"
    assert next(cursor)["id"] == "ip:2.2.2.2"
    value.graph.add_node(node("4.4.4.4"))

    with pytest.raises(CSTXError) as error:
        next(cursor)
    assert error.value.code == "CURSOR_INVALIDATED"


def test_unchanged_edge_does_not_invalidate_live_cursor():
    value = db()
    value.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2"), node("3.3.3.3")])
    relation = edge("ip:1.1.1.1", "ip:2.2.2.2")
    value.graph.add_edge(relation)
    cursor = value.graph.nodes(page_size=3)
    next(cursor)
    assert value.graph.add_edge(relation) == 0

    assert next(cursor)["id"] == "ip:2.2.2.2"


def test_edges_neighbors_and_direct_json():
    value = db()
    value.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2")])
    relation = edge("ip:1.1.1.1", "ip:2.2.2.2")
    assert value.graph.add_edge(relation) == 1
    assert value.graph.edge_count() == 1
    assert list(value.graph.edges())[0]["id"] == relation["id"]
    assert list(value.graph.neighbors("ip:1.1.1.1"))[0]["id"] == "ip:2.2.2.2"
    assert json.loads(value.graph.edges_json()) == list(value.graph.edges())
    assert json.loads(value.graph.neighbors_json("ip:1.1.1.1")) == list(
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
    assert value.graph.patch_node_extras(
        ["ip:2.2.2.2"], {"scope": "new", "run": "r1"}
    ) == 1
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
    assert relation["id"] == (
        "relationship:ip:1.1.1.1:related:ip:2.2.2.2:qualified"
    )
    assert value.graph.is_path_expression("ip -> ip")
    assert cstxpy.is_path_expression('ip[country="CN"]')


def test_native_graph_union_and_difference_are_deterministic():
    left = db()
    right = db()
    left.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2")])
    right.graph.add_nodes([node("2.2.2.2"), node("3.3.3.3")])

    union = left.graph.union_snapshot(right.graph)
    assert sorted(item["id"] for item in union["nodes"]) == [
        "ip:1.1.1.1",
        "ip:2.2.2.2",
        "ip:3.3.3.3",
    ]
    difference = left.graph.difference_snapshot(right.graph)
    assert [item["id"] for item in difference["nodes"]] == ["ip:1.1.1.1"]
    assert difference["edges"] == []


def test_trusted_cas_bytes_are_hydrated_natively():
    value = db()
    first = node("1.1.1.1")
    second = node("2.2.2.2")
    first["model"]["__node_type__"] = "ip"
    second["model"]["__node_type__"] = "ip"
    assert value.graph.hydrate_cas_nodes(
        [
            (first["id"], first["type"], json.dumps(first).encode()),
            (second["id"], second["type"], json.dumps(second).encode()),
        ]
    ) == 2
    relation = edge(first["id"], second["id"])
    assert value.graph.hydrate_cas_edges(
        [(relation["id"], "edge:related", json.dumps(relation).encode())]
    ) == 1
    assert value.graph.node_count() == 2
    assert value.graph.edge_count() == 1


def test_direct_dict_cursor_matches_canonical_json_for_all_column_types():
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
    value.graph.add_node(item)

    assert list(value.graph.nodes()) == json.loads(value.graph.nodes_json())


def test_repo_json_round_trip_and_close():
    value = db()
    value.graph.add_node(node("1.1.1.1"))
    snapshot = value.repo.dump_json()
    restored = db()
    assert restored.repo.load_json(snapshot) == len(snapshot)
    assert restored.graph.node("ip:1.1.1.1")["id"] == "ip:1.1.1.1"
    restored.close()
    with pytest.raises(CSTXError) as error:
        restored.graph.node_count()
    assert error.value.code == "NOT_INITIALIZED"


def test_not_found_has_stable_error_context():
    value = db()
    with pytest.raises(CSTXError) as error:
        value.graph.node("ip:missing")
    assert error.value.code == "NOT_FOUND"
    assert error.value.operation == "graph.node"


def test_repository_errors_have_stable_context():
    value = db()

    with pytest.raises(CSTXError) as error:
        value.repo.load_json(b"not-json")
    assert error.value.code == "CORRUPT_DATA"
    assert error.value.operation == "repo.load"


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
            "contains",
            "get",
            "list",
            "load_plugin",
            "load_all_plugins",
            "available_plugins",
            "has_native_artifact",
            "anchor_concepts",
        ),
        Graph: (
            "rag",
            "add_node",
            "add_nodes",
            "add_nodes_json",
            "hydrate_cas_nodes",
            "add_edge",
            "add_edges",
            "add_edges_json",
            "hydrate_cas_edges",
            "ingest",
            "node",
            "find_node",
            "patch_node_extras",
            "create_relationship",
            "is_path_expression",
            "union_snapshot",
            "difference_snapshot",
            "contains",
            "node_count",
            "edge_count",
            "stats",
            "nodes",
            "edges",
            "neighbors",
            "nodes_json",
            "edges_json",
            "neighbors_json",
            "query",
            "query_json",
            "ingest_native",
            "node_ids",
            "node_types",
            "nodes_by_ids",
            "link_nodes",
            "update_node_flags",
            "bfs",
            "shortest_paths",
            "degree",
            "subgraph_ids",
            "query_node_ids",
            "query_trace_ids",
            "induced_snapshot",
            "filter_snapshot",
            "find_anchors",
            "elevate_snapshot",
        ),
        Repository: (
            "index_graph",
            "commit",
            "commit_entries",
            "commit_delta",
            "diff",
            "dump",
            "load",
            "dump_json",
            "load_json",
            "snapshot_fingerprint",
            "head",
            "refs",
            "validate_ref",
            "set_ref",
            "delete_ref",
            "fork",
            "checkout_entries",
            "merge",
            "merge_base",
            "log",
            "tree_stats",
            "cherry_pick",
            "import_commit",
            "export_commit",
            "import_tree_objects",
            "export_tree_objects",
            "tree_entries",
            "find_tree_entry",
            "diff_tree_entries",
            "tree_root_stats",
        ),
        NodeCursor: ("closed", "close", "__iter__", "__next__", "__enter__", "__exit__"),
        EdgeCursor: ("closed", "close", "__iter__", "__next__", "__enter__", "__exit__"),
        NodeFlags: ("all_mask", "default_exclude_mask"),
    }

    assert inspect.getdoc(CSTXError)
    for api_type, members in public_members.items():
        assert inspect.getdoc(api_type), api_type.__name__
        for member in members:
            assert inspect.getdoc(getattr(api_type, member)), f"{api_type.__name__}.{member}"
        discovered = {
            name
            for name, value in api_type.__dict__.items()
            if not name.startswith("_")
            and (callable(getattr(api_type, name, None)) or inspect.isdatadescriptor(value))
        }
        assert discovered <= set(members), (
            f"undocumented API added to {api_type.__name__}: {discovered - set(members)}"
        )


def test_json_api_docs_define_the_fast_path_boundary():
    """JSON suffixes must stay justified as transport paths, not duplicate APIs."""
    for member in (
        Graph.add_nodes_json,
        Graph.add_edges_json,
        Graph.nodes_json,
        Graph.edges_json,
        Graph.neighbors_json,
        Graph.query_json,
        Repository.dump_json,
        Repository.load_json,
    ):
        documentation = inspect.getdoc(member).lower()
        assert "json" in documentation
        assert any(
            reason in documentation
            for reason in ("fast path", "transport", "existing", "interoper")
        )


def test_type_stub_documents_every_exported_class_and_method():
    """IDE-visible signatures must carry the same explanation as runtime help()."""
    stub = Path(cstxpy.__file__).with_name("_cstxpy.pyi")
    tree = ast.parse(stub.read_text(encoding="utf-8"), filename=str(stub))
    for node in ast.walk(tree):
        if isinstance(node, (ast.ClassDef, ast.FunctionDef, ast.AsyncFunctionDef)):
            assert ast.get_docstring(node), f"missing stub documentation: {node.name}"
