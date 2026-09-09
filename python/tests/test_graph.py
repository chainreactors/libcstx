"""Python binding checks for the protobuf-only graph boundary."""

import ast
import inspect
from pathlib import Path

import cstxpy
import pytest
from cstxpy import CSTX, CSTXError, CSTXGraph, Extensions, GraphCursor, NodeFlags, Repository
from cstxpy.proto import cstx_pb2 as cstx
from google.protobuf.json_format import ParseDict
from google.protobuf.struct_pb2 import Struct


def _value(node_type: str, **fields: str) -> cstx.EntityValue:
    """A payload named from the schema document — the only spelling there is."""
    return cstx.EntityValue(
        node_type=node_type,
        fields=[cstx.EntityField(name=name, text=text) for name, text in sorted(fields.items())],
    )


def _ip(value: str, *, flags: int = 0, annotations: dict | None = None) -> cstx.Node:
    node = cstx.Node(
        id=f"ip:{value}", sources=["test"],
        value=_value("ip", ip=value),
    )
    node.flags_mask = flags
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
    relation.value.relationship_type = "contain"
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
    carried = {field.name: field.text for field in node.value.fields}
    assert carried["ip"] == "1.1.1.1"

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
        (cstx.NodeFilter(flags_any_mask=1 << 0)).SerializeToString(),  # easm: `honeypot`
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
    # The document marks `ip` as non-semantic — an identity value is not useful
    # for semantic recall — so the record needs a field that is indexed.
    runtime.graph.add_nodes(_graph([
        cstx.Node(
            id="ip:1.1.1.1", sources=["test"],
            value=_value("ip", ip="1.1.1.1", as_name="Example Networks"),
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
    """Every member the stub declares is documented on the live object.

    The member list is read out of `_cstxpy.pyi` rather than restated here.
    It used to be a hand-written tuple per class, and it drifted twice --
    `Extensions.export_contract` and `CSTXGraph.add_relationship` were both
    absent, so neither was ever checked for a docstring. What makes reading
    the stub sound is `python_stub_declares_every_binding_member` in
    cstx-abi's `runtime_contract` test: it fails when the stub omits a member
    the binding exposes, so the stub is the complete list by construction.
    """
    stub = Path(cstxpy.__file__).with_name("_cstxpy.pyi")
    tree = ast.parse(stub.read_text(encoding="utf-8"), filename=str(stub))

    assert inspect.getdoc(cstxpy.CSTXError)
    checked = 0
    for node in tree.body:
        if not isinstance(node, ast.ClassDef):
            continue
        api_type = getattr(cstxpy, node.name, None)
        if api_type is None or node.name == "CSTXError":
            continue
        assert inspect.getdoc(api_type), node.name
        for member in node.body:
            if not isinstance(member, (ast.FunctionDef, ast.AsyncFunctionDef)):
                continue
            # `__init__` documentation belongs to the class for a pyclass, and
            # `__exit__` is plumbing pyo3 generates no docstring for.
            if member.name in {"__init__", "__enter__", "__exit__"}:
                continue
            attribute = getattr(api_type, member.name, None)
            assert attribute is not None, f"{node.name}.{member.name} missing from binding"
            assert inspect.getdoc(attribute), f"{node.name}.{member.name}"
            checked += 1

    # NodeFlags is pure Python and has no stub of its own.
    for member in ("all_mask", "default_exclude_mask", "bit", "mask", "names"):
        assert inspect.getdoc(getattr(NodeFlags, member)), f"NodeFlags.{member}"

    assert checked > 80, f"stub walk only reached {checked} members"


def test_type_stub_documents_every_exported_class_and_method() -> None:
    stub = Path(cstxpy.__file__).with_name("_cstxpy.pyi")
    tree = ast.parse(stub.read_text(encoding="utf-8"), filename=str(stub))
    for node in ast.walk(tree):
        if isinstance(node, (ast.ClassDef, ast.FunctionDef, ast.AsyncFunctionDef)):
            assert ast.get_docstring(node), f"missing stub documentation: {node.name}"


def test_flag_declared_above_bit_six_round_trips_through_protobuf() -> None:
    """An extension's flag at bit 40 survives the boundary, filter included.

    This is the regression the `uint64` mask exists for. `Node.flags` used to
    be a seven-value enum, and the Python encoder derived its bit positions
    from that enum, so a flag any extension declared above bit 6 was dropped
    in silence -- written, accepted, and simply not there on read. The schema
    contract has always said extensions own bits 0-55.
    """
    import json

    document = {
        "schema_version": 1,
        "extension": "acme",
        "nodes": {
            "acme_asset": {
                "message": "acme.Asset",
                "value_field": "asset_id",
                "identity": {"field": "asset_id"},
                "fields": [{"name": "asset_id", "number": 1, "type": "string"}],
            }
        },
        "flags": {"quarantined": {"bit": 40, "default_exclude": True}},
    }
    high = 1 << 40

    db = cstxpy.CSTX("flag-bit-40", 1024)
    try:
        contract = cstx.ExtensionContract(contract_version=1)
        definition = contract.extensions["acme"]
        definition.name = "acme"
        definition.schema = json.dumps(document)
        db.extensions.register(contract.SerializeToString())

        node = cstx.Node(
            id="acme_asset:a1",
            flags_mask=high,
            value=cstx.EntityValue(
                node_type="acme_asset",
                fields=[cstx.EntityField(name="asset_id", text="a1")],
            ),
        )
        db.graph.add_nodes(cstx.Graph(nodes=[node]).SerializeToString())

        stored = cstx.Node.FromString(db.graph.node("acme_asset:a1"))
        assert stored.flags_mask == high

        cursor = db.graph.nodes(
            (cstx.NodeFilter(flags_any_mask=high)).SerializeToString(),
            (cstx.QueryWindow(page=1)).SerializeToString(),
        )
        page = cstx.GraphResultPage.FromString(cursor.page(limit=10, page=1))
        assert [row.id for row in page.nodes.values] == ["acme_asset:a1"]

        # And a bit nobody declared passes through untouched: the registry
        # decides what a bit means, so the boundary has nothing to reject.
        undeclared = cstx.Node(
            id="acme_asset:a2",
            flags_mask=1 << 55,
            value=cstx.EntityValue(
                node_type="acme_asset",
                fields=[cstx.EntityField(name="asset_id", text="a2")],
            ),
        )
        db.graph.add_nodes(cstx.Graph(nodes=[undeclared]).SerializeToString())
        assert cstx.Node.FromString(db.graph.node("acme_asset:a2")).flags_mask == 1 << 55
    finally:
        db.close()
