"""The generated Python bindings exercise every schema feature at the wire edge."""

from cstxpy.proto import cstx_pb2 as cstx, sco_pb2 as easm
from google.protobuf.any_pb2 import Any
from google.protobuf.json_format import MessageToDict, ParseDict
from google.protobuf.struct_pb2 import Struct


def test_graph_page_oneofs_optional_map_repeated_and_enum_round_trip() -> None:
    entity = easm.App(
        app_id="app-1",
        url="https://example.test",
        frameworks=["react", "nginx"],
        status_code=200,
    )
    node = cstx.Node(
        id="app:app-1",
        entity=Any(
            type_url="type.googleapis.com/easm.App",
            value=(entity).SerializeToString(),
        ),
        sources=["fixture", "scanner"],
        flags=[cstx.NodeFlag.NODE_FLAG_THREAT_PRESENT, cstx.NodeFlag.NODE_FLAG_INTERNAL],
        annotations=ParseDict(
            {"nested": {"enabled": True}, "labels": ["a", "b"]}, Struct()
        ),
    )
    page = cstx.GraphResultPage(
        page=2,
        limit=10,
        total=11,
        has_next=True,
        nodes=cstx.NodePage(values=[node]),
        query=cstx.QuerySummary(nodes_by_type={"app": 11}),
    )

    encoded = (page).SerializeToString()
    decoded = cstx.GraphResultPage.FromString(encoded)
    result_kind = decoded.WhichOneof("result")
    summary_kind = decoded.WhichOneof("summary")

    assert result_kind == "nodes"
    assert summary_kind == "query"
    assert list(decoded.nodes.values[0].flags) == [
        cstx.NodeFlag.NODE_FLAG_THREAT_PRESENT,
        cstx.NodeFlag.NODE_FLAG_INTERNAL,
    ]
    assert dict(decoded.query.nodes_by_type) == {"app": 11}
    annotations = MessageToDict(
        decoded.nodes.values[0].annotations, preserving_proto_field_name=True
    )
    assert annotations["nested"]["enabled"] is True
    assert easm.App.FromString(decoded.nodes.values[0].entity.value).url == "https://example.test"
    # Message equality, not byte equality: this page carries map fields
    # (`nodes_by_type`, and the Struct's own), and protobuf does not promise a
    # stable order for those. Round-tripping the message is the claim; making
    # it about bytes would be a claim the format does not support.
    assert cstx.GraphResultPage.FromString(encoded) == decoded
    assert cstx.GraphResultPage.FromString(
        (decoded).SerializeToString()
    ) == decoded


def test_optional_presence_and_unknown_wire_fields_are_safe() -> None:
    unset = easm.App(app_id="app-2")
    set_value = easm.App(app_id="app-2", status_code=0)

    assert not unset.HasField("status_code")
    assert set_value.status_code == 0
    assert easm.App.FromString(set_value.SerializeToString()).HasField("status_code")

    # Field 99 (varint) is intentionally unknown to App.  The generated
    # parser must ignore it while preserving all known fields.
    decoded = easm.App.FromString(set_value.SerializeToString() + b"\x98\x06\x01")
    assert decoded.app_id == "app-2"
    assert decoded.status_code == 0
