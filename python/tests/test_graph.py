"""Python binding checks for the protobuf-only graph boundary."""

import ast
import inspect
from pathlib import Path

import cstxpy
import pytest
from cstxpy import CSTX, CSTXError, CSTXGraph, Extensions, GraphCursor, NodeFlags, Repository
from cstxpy.proto import cstx_pb2 as cstx, sco_pb2 as easm
from google.protobuf.any_pb2 import Any
from google.protobuf.json_format import ParseDict
from google.protobuf.struct_pb2 import Struct


def _ip(value: str, *, flags: int = 0, annotations: dict | None = None) -> cstx.Node:
    entity = easm.Ip(ip=value)
    node = cstx.Node(
        id=f"ip:{value}", sources=["test"],
        entity=Any(type_url="type.googleapis.com/easm.Ip", value=(entity).SerializeToString()),
    )
    for number in sorted(n for n in cstx.NodeFlag.values() if n > 0):
        if flags & (1 << (number - 1)):
            node.flags.append(number)
    if annotations:
        ParseDict(annotations, node.annotations)
    return node


def _contain(source: str, target: str) -> cstx.Relationship:
    relation = cstx.Relationship(
        id=f"relationship:{source}:contain:{target}",
        source_id=source,
        target_id=target,
        sources=["test"],
    )
    relation.relation.CopyFrom(Any(
        type_url="type.googleapis.com/easm.Contain", value=b""
    ))
    return relation


def _graph(nodes=(), relationships=()) -> bytes:
    return (cstx.Graph(
        nodes=list(nodes), relationships=list(relationships)
    )).SerializeToString()


def _window(limit: int = 1024, page: int = 1, order: int = 0) -> bytes:
    return (cstx.QueryWindow(limit=limit, page=page, order=order)).SerializeToString()


def _node_rows(cursor: GraphCursor) -> list[cstx.Node]:
    return [cstx.Node.FromString(payload) for payload in cursor]


def _register(runtime: CSTX) -> None:
    # The built-in extension ships its own schema document; enabling it is
    # the whole registration step.
    runtime.extensions.enable("easm")


def test_graph_methods_accept_and_return_only_protobuf() -> None:
    runtime = CSTX()
    _register(runtime)
    assert runtime.graph.add_nodes(_graph([_ip("1.1.1.1")])) == 1
    assert runtime.graph.add_nodes(_graph([_ip("1.1.1.1")])) == 0

    node = cstx.Node.FromString(runtime.graph.node("ip:1.1.1.1"))
    entity = easm.Ip.FromString(node.entity.value)
    assert entity.ip == "1.1.1.1"

    assert runtime.graph.add_nodes(
        _graph([_ip("2.2.2.2")])
    ) == 1
    assert runtime.graph.add_relationships(
        _graph(relationships=[_contain("ip:1.1.1.1", "ip:2.2.2.2")])
    ) == 1
    edge = cstx.Relationship.FromString(
        runtime.graph.relationship("relationship:ip:1.1.1.1:contain:ip:2.2.2.2")
    )
    assert edge.source_id == "ip:1.1.1.1"

    stats = cstx.GraphStats.FromString(runtime.graph.stats())
    assert stats.nodes_by_type["ip"] == 2
    assert stats.relationships_by_type["contain"] == 1
    change = cstx.GraphChangeSet.FromString(runtime.last_change())
    assert list(change.added_relationship_ids) == [edge.id]


def test_typed_cursor_page_and_next_share_one_transport() -> None:
    runtime = CSTX()
    _register(runtime)
    runtime.graph.add_nodes(_graph([_ip("2.2.2.2"), _ip("1.1.1.1", flags=NodeFlags.HONEYPOT)]))
    cursor = runtime.graph.nodes(
        (cstx.NodeFilter(flags_any=[cstx.NodeFlag.NODE_FLAG_HONEYPOT])).SerializeToString(),
        _window(order=cstx.SortOrder.SORT_ORDER_ID_ASC),
    )
    row = cstx.Node.FromString(next(cursor))
    assert row.id == "ip:1.1.1.1"
    assert cursor.next() is None

    all_rows = runtime.graph.nodes(b"", _window(order=cstx.SortOrder.SORT_ORDER_ID_ASC))
    page = cstx.GraphResultPage.FromString(all_rows.page(limit=10, page=1))
    assert [item.id for item in page.nodes.values] == ["ip:1.1.1.1", "ip:2.2.2.2"]


def test_cursor_invalidation_and_typed_stats() -> None:
    runtime = CSTX()
    _register(runtime)
    runtime.graph.add_nodes(_graph([_ip("1.1.1.1"), _ip("2.2.2.2")]))
    cursor = runtime.graph.nodes(b"", _window())
    assert cstx.Node.FromString(next(cursor)).id == "ip:1.1.1.1"
    runtime.graph.add_nodes(_graph([_ip("3.3.3.3")]))
    with pytest.raises(CSTXError) as error:
        cursor.next()
    assert error.value.code == "CURSOR_INVALIDATED"

    filtered = cstx.GraphStats.FromString(
        runtime.graph.stats(exclude_mask=NodeFlags.HONEYPOT)
    )
    assert filtered.nodes_by_type["ip"] == 3


def test_algorithms_return_typed_pages() -> None:
    runtime = CSTX()
    _register(runtime)
    runtime.graph.add_nodes(_graph([_ip("a"), _ip("b")]))
    runtime.graph.add_relationships(_graph(relationships=[_contain("ip:a", "ip:b")]))
    assert runtime.graph.analyze(cstxpy.Algorithm.is_dag()) is True
    paths = runtime.graph.analyze(
        cstxpy.Algorithm.shortest_paths("ip:a", "ip:b", direction="both")
    )
    page = cstx.GraphResultPage.FromString(paths.page(limit=10, page=1))
    assert [list(item.node_ids) for item in page.paths.values] == [["ip:a", "ip:b"]]


def test_extension_introspection_is_protobuf() -> None:
    runtime = CSTX()
    _register(runtime)
    catalog = cstx.ExtensionCatalog.FromString(runtime.extensions.list())
    assert any(item.name == "easm" for item in catalog.extensions)
    schema = cstx.NodeType.FromString(runtime.extensions.schema("ip"))
    assert schema.type_url


def test_rag_uses_protobuf_plan_and_results() -> None:
    runtime = CSTX()
    _register(runtime)
    runtime.graph.add_nodes(_graph([_ip("1.1.1.1")]))
    # easm.Ip marks `ip` as non-semantic — an identity value is not useful
    # for semantic recall — so the record needs a field that is indexed.
    runtime.graph.add_nodes(_graph([
        cstx.Node(
            id="ip:1.1.1.1", sources=["test"],
            entity=Any(
                type_url="type.googleapis.com/easm.Ip",
                value=(easm.Ip(ip="1.1.1.1", as_name="Example Networks")).SerializeToString(),
            ),
        )
    ]))
    session = runtime.graph.rag().index(
        (cstx.RagIndexPlan(
            commit="working",
            mode=cstx.RagIndexMode.RAG_INDEX_FULL,
        )).SerializeToString()
    )
    record = cstx.RagRecord.FromString(next(session.pending("v1")))
    assert record.id == "node/ip:1.1.1.1"
    retrieval = runtime.graph.rag().retrieve(
        (cstx.RagQuery(text="1.1.1.1", limit=5)).SerializeToString()
    )
    plan = cstx.RecallPlan.FromString(retrieval.requests())
    assert len(plan.queries) == 2
    result = cstx.RagResult.FromString(
        retrieval.complete((cstx.RecallResults()).SerializeToString())
    )
    assert result is not None


def test_every_exported_binding_has_documentation() -> None:
    public_members = {
        CSTX: ("extensions", "graph", "repo", "closed", "project_id", "close", "last_change", "__enter__", "__exit__"),
        Extensions: ("register", "enable", "list", "info", "contains", "schema", "schemas", "has_native_artifact", "anchor_concepts"),
        CSTXGraph: ("rag", "add_nodes", "replace_nodes", "relationship", "node", "nodes", "relationships", "neighbors", "query", "analyze", "add_relationships", "delete_nodes", "delete_relationships", "union", "merge", "difference", "contains", "node_count", "relationship_count", "stats", "ingest", "find_node", "patch_node_extras", "node_types", "link", "update_node_flags", "degree", "subgraph", "query_subgraph", "induced_subgraph", "filter", "filter_with_reasons", "find_anchors", "elevate"),
        Repository: ("resolve", "head", "checkout", "commit", "diff", "log", "history", "branch", "merge", "stat", "delta"),
        GraphCursor: ("kind", "page", "next", "closed", "close", "__iter__", "__next__", "__enter__", "__exit__"),
        NodeFlags: ("all_mask", "default_exclude_mask"),
    }
    assert inspect.getdoc(cstxpy.CSTXError)
    for api_type, members in public_members.items():
        assert inspect.getdoc(api_type)
        for member in members:
            assert inspect.getdoc(getattr(api_type, member)), f"{api_type.__name__}.{member}"


def test_type_stub_documents_every_exported_class_and_method() -> None:
    stub = Path(cstxpy.__file__).with_name("_cstxpy.pyi")
    tree = ast.parse(stub.read_text(encoding="utf-8"), filename=str(stub))
    for node in ast.walk(tree):
        if isinstance(node, (ast.ClassDef, ast.FunctionDef, ast.AsyncFunctionDef)):
            assert ast.get_docstring(node), f"missing stub documentation: {node.name}"
