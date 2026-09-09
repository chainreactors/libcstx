"""The generated core messages are the complete protobuf boundary.

Plugin structure is carried by the registered schema document and by the
schema-named ``EntityValue`` fields. No generated plugin message or
``google.protobuf.Any`` participates in these round trips.
"""

import cstxpy
from cstxpy.proto import cstx_pb2 as cstx


def _ip_graph(value: str) -> bytes:
    node = cstx.Node(
        id=f"ip:{value}",
        sources=["test"],
        value=cstx.EntityValue(
            node_type="ip",
            fields=[cstx.EntityField(name="ip", text=value)],
        ),
    )
    return cstx.Graph(nodes=[node]).SerializeToString()


def test_generated_messages_are_the_wire_contract() -> None:
    payload = cstx.ParserPayload(
        plugin="easm",
        artifact="gogo",
        data=b'{"ip":"1.1.1.1"}',
        content_type="application/json",
    )

    decoded = cstx.ParserPayload.FromString(payload.SerializeToString())

    assert decoded.plugin == "easm"
    assert decoded.artifact == "gogo"
    assert decoded.data == b'{"ip":"1.1.1.1"}'
    assert decoded.content_type == "application/json"


def test_python_runtime_reads_the_schema_named_payload() -> None:
    runtime = cstxpy.CSTX()
    runtime.extensions.enable("easm")
    runtime.graph.add_nodes(_ip_graph("typed"))

    node = cstx.Node.FromString(runtime.graph.node("ip:typed"))
    fields = {field.name: field for field in node.value.fields}

    assert node.value.node_type == "ip"
    assert fields["ip"].text == "typed"


def test_python_runtime_exposes_typed_graph_stats() -> None:
    runtime = cstxpy.CSTX()
    runtime.extensions.enable("easm")
    runtime.graph.add_nodes(_ip_graph("stats"))

    stats = cstx.GraphStats.FromString(runtime.graph.stats())

    assert stats.nodes_by_type["ip"] == 1


def test_python_runtime_exposes_anchor_catalog_proto() -> None:
    runtime = cstxpy.CSTX()
    runtime.extensions.enable("easm")

    catalog = cstx.GraphAnchorCatalog.FromString(runtime.graph.find_anchors("threat"))

    assert catalog.anchors == []
