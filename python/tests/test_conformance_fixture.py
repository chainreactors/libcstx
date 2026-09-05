"""Cross-language checks backed by the canonical protobuf graph fixture.

Rust and Go run this same file through their own bindings and assert the same
ids, counts and query result. That claim only means something if all three use
the fixture as written: this test used to remap every record onto `easm.Ip`
and `easm.Contain` before handing it over, which made it a test of easm rather
than of the fixture, and left the three languages agreeing with nobody.

The fixture declares its own type in a schema document, so there is nothing to
remap — and no generated class for `conformance.Asset` in any language.
"""

import json
from pathlib import Path

import cstxpy
from cstxpy.proto import cstx_pb2 as cstx
from google.protobuf.any_pb2 import Any

FIXTURE = json.loads(
    (
        Path(__file__).resolve().parents[3] / "tests/fixtures/conformance.json"
    ).read_text(encoding="utf-8")
)


def _graph(fixture: dict) -> bytes:
    """The fixture's records as the boundary takes them, named by the document."""
    nodes = []
    for item in fixture["nodes"]:
        entity = cstx.EntityValue(node_type=item["type"])
        for name, value in item["model"].items():
            entity.fields.add(name=name, text=value)
        nodes.append(cstx.Node(id=item["id"], sources=item["sources"], value=entity))

    relations = fixture["document"]["relations"]
    relationships = [
        cstx.Relationship(
            id=item["id"],
            source_id=item["source_id"],
            target_id=item["target_id"],
            sources=item["sources"],
            # A relation type is a field-less marker: the document names the
            # message and the payload is empty.
            relation=Any(
                type_url="type.googleapis.com/"
                + relations[item["relation_type"]]["message"],
                value=b"",
            ),
        )
        for item in fixture["relationships"]
    ]
    return (cstx.Graph(nodes=nodes, relationships=relationships)).SerializeToString()


def _register(runtime: cstxpy.CSTX) -> None:
    contract = cstx.ExtensionContract(contract_version=1)
    definition = contract.extensions["conformance"]
    definition.name = "conformance"
    definition.schema = json.dumps(FIXTURE["document"])
    runtime.extensions.register(contract.SerializeToString())


def test_conformance_fixture_uses_typed_protobuf_transport() -> None:
    runtime = cstxpy.CSTX()
    _register(runtime)
    runtime.graph.add_nodes(_graph(FIXTURE))

    expected = FIXTURE["expected"]
    assert runtime.graph.node_count() == expected["node_count"]
    assert runtime.graph.relationship_count() == expected["relationship_count"]
    page = cstx.GraphResultPage.FromString(
        runtime.graph.query(
            (cstx.GraphQuery(expression=FIXTURE["query"])).SerializeToString()
        ).page(limit=100, page=1)
    )
    assert [item.id for item in page.nodes.values] == expected["node_ids"]
