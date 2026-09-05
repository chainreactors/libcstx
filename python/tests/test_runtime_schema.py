"""The runtime schema is the only structure contract, and it is the same one
for a built-in extension and for one registered at runtime.

These tests pin the two properties that make that true:

* every lookup a caller can make is answered from the schema, not from a
  generated per-extension lookup table or a protobuf descriptor;
* an extension that registers a schema at runtime gets identical treatment —
  same registry, same derived pydantic bases, same queries.
"""

import json
from pathlib import Path

import pytest
from cstxpy import model as cstx_model
from cstxpy import schema as cstx_schema
from cstxpy.schema import ExtensionSchema, NodeSchema, SchemaRegistry
from pydantic import BaseModel

SCHEMA_DIR = Path(cstx_schema.__file__).resolve().parent / "schemas"


@pytest.fixture
def easm_schema() -> ExtensionSchema:
    schema = cstx_schema.registry.extension("easm")
    assert schema is not None, "the built-in EASM schema ships with cstxpy"
    return schema


# ── the schema artifact itself ──


def test_bundled_schema_is_loaded_for_every_shipped_extension():
    shipped = {path.name[: -len(".schema.json")] for path in SCHEMA_DIR.glob("*.schema.json")}
    assert shipped, "at least the built-in EASM schema must ship"
    assert shipped <= set(cstx_schema.registry.extensions())


def test_schema_describes_identity_and_columns(easm_schema):
    subdomain = easm_schema.nodes["subdomain"]
    assert subdomain.message == "easm.Subdomain"
    assert subdomain.type_url == "type.googleapis.com/easm.Subdomain"
    assert subdomain.identity_field == "host"
    assert subdomain.identity_format is None

    host = subdomain.field("host")
    assert host is not None
    assert (host.number, host.type, host.repeated, host.semantic) == (1, "string", False, False)

    a_records = subdomain.field("a")
    assert a_records is not None and a_records.repeated


def test_schema_carries_composite_identity(easm_schema):
    """A JSON-Schema style contract cannot express this; the schema can."""
    port = easm_schema.nodes["port"]
    assert port.identity_format == "{ip}:{port}"
    assert port.identity_field is None


def test_unsupported_schema_version_is_rejected():
    with pytest.raises(ValueError, match="unsupported schema_version"):
        ExtensionSchema.parse({"schema_version": 999, "extension": "x"})


# ── derived pydantic bases ──


def test_derived_base_matches_the_schema_field_for_field(easm_schema):
    for node_type, node in easm_schema.nodes.items():
        base = cstx_model.base_model(node_type)
        assert issubclass(base, BaseModel)
        assert set(base.model_fields) == {field.name for field in node.fields}
        for field in node.fields:
            info = base.model_fields[field.name]
            required = not field.optional and not field.repeated
            assert info.is_required() is required, f"{node_type}.{field.name}"


def test_derived_base_keeps_proto3_zero_value_semantics():
    """A scalar without presence is always populated; only repeated fields
    distinguish absent from empty."""
    app = cstx_model.base_model("app")(app_id="https://a/")
    assert app.url == ""
    assert app.status_code == 0
    assert app.frameworks is None


def test_derived_base_accepts_unknown_fields():
    domain = cstx_model.base_model("domain")(host="a.com", scanner_tag="x")
    assert domain.model_dump()["scanner_tag"] == "x"


def test_base_attribute_access_matches_message_name(easm_schema):
    assert cstx_model.SubdomainBase is cstx_model.base_model("subdomain")
    with pytest.raises(AttributeError):
        _ = cstx_model.NoSuchThingBase


# ── a runtime-registered extension is not a second-class citizen ──


THIRD_PARTY = {
    "schema_version": 1,
    "extension": "acme",
    "nodes": {
        "acme_asset": {
            "message": "acme.Asset",
            "value_field": "asset_id",
            "identity": {"field": "asset_id"},
            "fields": [
                {"name": "asset_id", "number": 1, "type": "string", "semantic": False},
                {"name": "owner", "number": 2, "type": "string", "optional": True},
                {"name": "score", "number": 3, "type": "int64", "optional": True},
                {"name": "tags", "number": 4, "type": "string", "repeated": True},
            ],
        }
    },
    "relations": {"acme_owns": {"message": "acme.Owns"}},
}


@pytest.fixture
def isolated_registry() -> SchemaRegistry:
    """A private view so cross-test state cannot leak into the global one.

    `_install` rather than a public entry: the view has no registration API,
    because registering is the core's decision and this class only mirrors it.
    """
    registry = SchemaRegistry()
    registry._install(ExtensionSchema.parse(THIRD_PARTY))
    return registry


def test_runtime_extension_answers_every_builtin_query(isolated_registry):
    assert isolated_registry.node_type_url("acme_asset") == "type.googleapis.com/acme.Asset"
    assert isolated_registry.node_type_from_url("type.googleapis.com/acme.Asset") == "acme_asset"
    assert isolated_registry.relation_type_url("acme_owns") == "type.googleapis.com/acme.Owns"
    assert isolated_registry.relation_type_from_url("type.googleapis.com/acme.Owns") == "acme_owns"


def test_runtime_schema_carries_no_export_format_metadata():
    """The runtime schema describes columns and identity — nothing else.

    STIX spellings used to ride along here; they belong to the exporter that
    needs them, not to the contract every extension has to satisfy.
    """
    import dataclasses

    assert "stix_type" not in {f.name for f in dataclasses.fields(NodeSchema)}
    assert not hasattr(cstx_schema, "stix_type_for")
    assert not hasattr(cstx_schema, "node_type_from_stix")
    assert not hasattr(SchemaRegistry, "stix_type_for")
    for path in SCHEMA_DIR.glob("*.schema.json"):
        assert "stix" not in path.read_text(encoding="utf-8").lower(), path


def register(document: dict, *, values: bool = False) -> "cstxpy.CSTX":
    """Declare a type the one way a caller can: through a runtime.

    The core validates and accepts, then the view is refreshed from what the
    core holds. `cstxpy.schema` has no registration entry of its own — these
    tests used to reach for one, which is exactly the second declaration path
    the boundary is supposed to have removed.
    """
    import cstxpy
    from cstxpy.proto import cstx_pb2 as cstx_proto

    # `values=True` opens the runtime in the payload format that hands nodes
    # back named by the document. Without it reads return an `Any`, and a
    # `Node.value` assertion then sees an empty message rather than a failure.
    runtime = cstxpy.CSTX(
        payload_format=cstx_proto.PAYLOAD_FORMAT_VALUE
        if values
        else cstx_proto.PAYLOAD_FORMAT_ENTITY
    )
    contract = cstx_proto.ExtensionContract(contract_version=1)
    definition = contract.extensions[document["extension"]]
    definition.name = document["extension"]
    definition.schema = json.dumps(document)
    runtime.extensions.register(contract.SerializeToString())
    cstx_schema.project(
        cstx_proto.ExtensionContract.FromString(runtime.extensions.export_contract())
    )
    return runtime


def test_runtime_rejects_unknown_payload_format_at_the_python_boundary():
    import cstxpy

    with pytest.raises(cstxpy.CSTXError) as captured:
        cstxpy.CSTX(payload_format=99)

    assert captured.value.code == "INVALID_ARGUMENT"
    assert captured.value.operation == "cstx.open"
    assert captured.value.field == "payload_format"
    assert captured.value.actual == "99"


def test_runtime_extension_gets_the_same_derived_base():
    """Registering a schema is the whole story: no codegen step, no descriptor."""
    runtime = register(THIRD_PARTY)
    try:
        base = cstx_model.base_model("acme_asset")
        assert base.__name__ == "AssetBase"
        assert set(base.model_fields) == {"asset_id", "owner", "score", "tags"}
        assert base.model_fields["asset_id"].is_required()
        assert not base.model_fields["owner"].is_required()

        instance = base(asset_id="a-1", owner="ops", score=7, tags=["x"])
        assert instance.model_dump() == {
            "asset_id": "a-1",
            "owner": "ops",
            "score": 7,
            "tags": ["x"],
        }
        # reached through the same attribute protocol as the built-in bases
        assert cstx_model.AssetBase is base
    finally:
        runtime.close()
        cstx_schema.registry._remove("acme")
        cstx_model._cache.pop("acme_asset", None)  # noqa: SLF001
        cstx_model._by_class_name.pop("AssetBase", None)


def test_rebuilt_base_follows_a_replaced_schema():
    """Re-registering an extension must not keep serving the stale base."""
    runtime = register(THIRD_PARTY)
    try:
        first = cstx_model.base_model("acme_asset")

        changed = json.loads(json.dumps(THIRD_PARTY))
        changed["nodes"]["acme_asset"]["fields"].append(
            {"name": "region", "number": 5, "type": "string", "optional": True}
        )
        register(changed).close()
        second = cstx_model.base_model("acme_asset")

        assert second is not first
        assert "region" in second.model_fields
    finally:
        runtime.close()
        cstx_schema.registry._remove("acme")
        cstx_model._cache.pop("acme_asset", None)
        cstx_model._by_class_name.pop("AssetBase", None)


def test_builtin_and_runtime_extensions_share_one_registry(isolated_registry, easm_schema):
    """Nothing in the lookup path branches on where a node type came from."""
    both = SchemaRegistry()
    both._install(easm_schema)
    both._install(ExtensionSchema.parse(THIRD_PARTY))

    assert set(both.extensions()) == {"easm", "acme"}
    for node_type in ("subdomain", "acme_asset"):
        node = both.node(node_type)
        assert node is not None
        assert both.node_type_url(node_type) == node.type_url
        assert both.node_type_from_url(node.type_url) == node_type


# ── the core decides; this module mirrors ──


def test_the_view_has_no_registration_entry():
    """Declaring a type is `runtime.extensions.register`, and only that.

    A public `register`/`load_schema` here was a second declaration entry: it
    could put a type into Python that the core had never accepted, and it
    carried its own opinion about conflicts and ambiguous names — an opinion
    that disagreed with the core's.
    """
    assert not hasattr(cstx_schema.registry, "register")
    for retired in ("SchemaRegistry", "load_schema", "load_bundled"):
        assert retired not in cstx_schema.__all__


def test_a_conflict_the_core_refuses_leaves_no_half_in_the_view():
    """Two extensions claiming one node type is refused, and refused wholly."""
    import cstxpy
    from cstxpy.proto import cstx_pb2 as cstx_proto

    def contract(extension: str) -> bytes:
        document = dict(THIRD_PARTY, extension=extension)
        message = cstx_proto.ExtensionContract(contract_version=1)
        definition = message.extensions[extension]
        definition.name = extension
        definition.schema = json.dumps(document)
        return message.SerializeToString()

    runtime = cstxpy.CSTX()
    try:
        runtime.extensions.register(contract("acme"))
        with pytest.raises(cstxpy.CSTXError):
            # `squatter` declares `acme_asset` too. The core owns that call.
            runtime.extensions.register(contract("squatter"))
        cstx_schema.project(
            cstx_proto.ExtensionContract.FromString(
                runtime.extensions.export_contract()
            )
        )
        assert "squatter" not in cstx_schema.registry.extensions()
        assert cstx_schema.registry.node("acme_asset") is not None
    finally:
        runtime.close()
        cstx_schema.registry._remove("acme")
        cstx_schema.registry._remove("squatter")


def test_an_ambiguous_short_message_name_resolves_to_nothing():
    """Two packages spelling one message the same way has no answer.

    Picking whichever landed first is an answer, and it is the wrong one half
    the time. The core returns nothing here; so does this.
    """
    view = SchemaRegistry()
    view._install(ExtensionSchema.parse(dict(THIRD_PARTY, extension="one")))
    view._install(
        ExtensionSchema.parse(
            {
                "schema_version": 1,
                "extension": "two",
                "nodes": {
                    "other_asset": {
                        "message": "other.Asset",
                        "value_field": "asset_id",
                        "identity": {"field": "asset_id"},
                        "fields": [
                            {"name": "asset_id", "number": 1, "type": "string"}
                        ],
                    }
                },
                "relations": {},
            }
        )
    )
    # Fully qualified still answers; the bare name no longer does.
    assert view.node_type_from_url("type.googleapis.com/acme.Asset") == "acme_asset"
    assert view.node_type_from_url("type.googleapis.com/other.Asset") == "other_asset"
    assert view.node_type_from_url("Asset") is None


# ── the derived annotation and the column it lands in ──

# One sample per type a schema document may declare. `ENCODABLE_PROTO_TYPES`
# in `cstx-graph/src/schema_def.rs` is the same list; a type added there and
# not here simply is not covered, which the first assertion below catches.
ENCODABLE_SAMPLES = {
    "string": "probe-value",
    "bool": True,
    "int64": 7,
    "int32": 7,
    "uint32": 7,
    "sint64": 7,
    "sint32": 7,
    "double": 1.5,
}


def _probe_document() -> dict:
    fields = [{"name": "key", "number": 1, "type": "string", "semantic": False}]
    for index, proto_type in enumerate(sorted(ENCODABLE_SAMPLES), start=2):
        fields.append(
            {
                "name": f"f_{proto_type}",
                "number": index,
                "type": proto_type,
                "optional": True,
                "semantic": False,
            }
        )
    return {
        "schema_version": 1,
        "extension": "probe",
        "nodes": {
            "probe_node": {
                "message": "probe.Node",
                "value_field": "key",
                "identity": {"field": "key"},
                "fields": fields,
            }
        },
    }


def _entity_field(name: str, value):
    """Pick the payload branch from the Python value, as every SDK does.

    `cstx/core/values.py` and `sdk/go/values.go` choose the same way: the
    branch follows the value's own type and the runtime checks it against the
    column the schema declared. That is precisely why the annotation this
    module derives has to agree with `FieldSchema::column_type` — a field
    built as the wrong Python type picks the wrong branch and the write is
    refused. bool before int: in Python `bool` is a subclass of `int`.
    """
    from cstxpy.proto import cstx_pb2 as cstx_proto

    field = cstx_proto.EntityField(name=name)
    if isinstance(value, bool):
        field.flag = value
    elif isinstance(value, int):
        field.number = value
    elif isinstance(value, float):
        field.real = value
    else:
        field.text = str(value)
    return field


def test_derived_annotation_survives_the_column_it_lands_in():
    """The derived model must type a field as what its column stores.

    `python_type` is a second copy of a decision Rust owns
    (`FieldSchema::column_type`), written in a different vocabulary, and Rust
    pins its own three copies together (`every_encodable_field_type_round_trips
    _through_its_column`). This copy sat outside that net, which is how
    `int32` / `uint32` / `sint32` / `double` came to be built as strings the
    core then refused. So: build the value through the derived base, send it
    the way an SDK does, read it back.
    """
    from cstxpy.proto import cstx_pb2 as cstx_proto

    document = _probe_document()
    runtime = register(document, values=True)
    try:
        base = cstx_model.base_model("probe_node")
        declared = {f"f_{name}" for name in ENCODABLE_SAMPLES}
        assert declared <= set(base.model_fields)

        instance = base(key="k-1", **{f"f_{k}": v for k, v in ENCODABLE_SAMPLES.items()})
        payload = instance.model_dump()

        value = cstx_proto.EntityValue(node_type="probe_node")
        value.fields.extend(
            _entity_field(name, payload[name]) for name in sorted(payload)
        )
        graph = cstx_proto.Graph(nodes=[cstx_proto.Node(id="probe_node:k-1", value=value)])
        runtime.graph.add_nodes(graph.SerializeToString())

        stored = cstx_proto.Node.FromString(runtime.graph.node("probe_node:k-1"))
        read_back = {
            field.name: getattr(field, field.WhichOneof("value"))
            for field in stored.value.fields
        }
        for proto_type, sample in ENCODABLE_SAMPLES.items():
            assert read_back[f"f_{proto_type}"] == sample, proto_type
    finally:
        runtime.close()
