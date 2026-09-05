import cstxpy
from cstxpy.proto import cstx_pb2 as cstx, sco_pb2 as easm
from google.protobuf.any_pb2 import Any


def _ip_graph(value: str) -> bytes:
    entity = easm.Ip(ip=value)
    node = cstx.Node(
        id=f"ip:{value}",
        sources=["test"],
        entity=Any(type_url="type.googleapis.com/easm.Ip", value=(entity).SerializeToString()),
    )
    return (cstx.Graph(nodes=[node])).SerializeToString()


def test_generated_messages_are_the_wire_contract():
    payload = cstx.ParserPayload(
        plugin="easm",
        artifact="gogo",
        data=b'{"ip":"1.1.1.1"}',
        content_type="application/json",
    )
    decoded = cstx.ParserPayload.FromString((payload).SerializeToString())
    assert decoded.plugin == "easm"
    assert decoded.artifact == "gogo"
    assert decoded.data == b'{"ip":"1.1.1.1"}'
    assert decoded.content_type == "application/json"


def test_python_runtime_uses_typed_domain_methods():
    runtime = cstxpy.CSTX()
    runtime.extensions.enable("easm")
    runtime.graph.add_nodes(_ip_graph("typed"))
    node = cstx.Node.FromString(runtime.graph.node("ip:typed"))
    entity = easm.Ip.FromString(node.entity.value)
    assert entity.ip == "typed"


def test_python_runtime_exposes_typed_graph_stats():
    runtime = cstxpy.CSTX()
    runtime.extensions.enable("easm")
    runtime.graph.add_nodes(_ip_graph("stats"))
    stats = cstx.GraphStats.FromString(runtime.graph.stats())
    assert stats.nodes_by_type["ip"] == 1


def test_python_runtime_exposes_anchor_catalog_proto():
    runtime = cstxpy.CSTX()
    runtime.extensions.enable("easm")
    catalog = cstx.GraphAnchorCatalog.FromString(
        runtime.graph.find_anchors("threat")
    )
    assert catalog.anchors == []
