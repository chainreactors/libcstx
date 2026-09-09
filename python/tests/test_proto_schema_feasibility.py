"""The transport messages exercise every schema feature at the wire edge."""

from cstxpy.proto import cstx_pb2 as cstx
from google.protobuf.json_format import MessageToDict, ParseDict
from google.protobuf.struct_pb2 import Struct


def test_retired_payload_surface_is_absent() -> None:
    assert not hasattr(cstx, "PayloadFormat")
    assert "payload_format" not in cstx.RuntimeConfig.DESCRIPTOR.fields_by_name
    assert "entity" not in cstx.Node.DESCRIPTOR.fields_by_name
    assert "relation" not in cstx.Relationship.DESCRIPTOR.fields_by_name


def test_graph_page_oneofs_optional_map_repeated_and_enum_round_trip() -> None:
    node = cstx.Node(
        id="app:app-1",
        value=cstx.EntityValue(
            node_type="app",
            fields=[
                cstx.EntityField(name="app_id", text="app-1"),
                cstx.EntityField(name="frameworks", list=cstx.StringList(values=["react", "nginx"])),
                cstx.EntityField(name="status_code", number=200),
                cstx.EntityField(name="url", text="https://example.test"),
            ],
        ),
        sources=["fixture", "scanner"],
        # easm declares `threat_present` at bit 4 and `internal` at bit 6.
        flags_mask=(1 << 4) | (1 << 6),
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
    assert decoded.nodes.values[0].flags_mask == (1 << 4) | (1 << 6)
    assert dict(decoded.query.nodes_by_type) == {"app": 11}
    annotations = MessageToDict(
        decoded.nodes.values[0].annotations, preserving_proto_field_name=True
    )
    assert annotations["nested"]["enabled"] is True
    carried = {field.name: field for field in decoded.nodes.values[0].value.fields}
    assert carried["url"].text == "https://example.test"
    assert list(carried["frameworks"].list.values) == ["react", "nginx"]
    assert carried["status_code"].number == 200
    # Message equality, not byte equality: this page carries map fields
    # (`nodes_by_type`, and the Struct's own), and protobuf does not promise a
    # stable order for those. Round-tripping the message is the claim; making
    # it about bytes would be a claim the format does not support.
    assert cstx.GraphResultPage.FromString(encoded) == decoded
    assert cstx.GraphResultPage.FromString(
        (decoded).SerializeToString()
    ) == decoded


def test_absence_and_zero_are_distinguishable_in_a_payload() -> None:
    """A payload says which fields the producer sent, not which are non-zero.

    proto3 gives a scalar no presence, which is why a generated per-type message
    needed `optional` on every field that could legitimately be zero. A payload
    names its fields, so a field the producer set to zero is in the list and a
    field it never set is not — the distinction is structural.
    """
    unset = cstx.EntityValue(
        node_type="app", fields=[cstx.EntityField(name="app_id", text="app-2")]
    )
    zero = cstx.EntityValue(
        node_type="app",
        fields=[
            cstx.EntityField(name="app_id", text="app-2"),
            cstx.EntityField(name="status_code", number=0),
        ],
    )

    assert [field.name for field in unset.fields] == ["app_id"]
    assert [field.name for field in zero.fields] == ["app_id", "status_code"]
    restored = cstx.EntityValue.FromString((zero).SerializeToString())
    assert restored.fields[1].WhichOneof("value") == "number"
    assert restored.fields[1].number == 0

    # Field 99 (varint) is unknown to EntityValue. The parser must skip it while
    # preserving every field it does know, which is what keeps a payload written
    # by a newer build readable by an older one.
    decoded = cstx.EntityValue.FromString((zero).SerializeToString() + b"\x98\x06\x01")
    assert decoded.node_type == "app"
    assert [field.name for field in decoded.fields] == ["app_id", "status_code"]
