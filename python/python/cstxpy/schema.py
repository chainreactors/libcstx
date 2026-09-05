"""Runtime schema — the single structure contract shared by every extension.

`make codegen` projects each ``.proto`` into one ``<extension>.schema.json``
describing node types, their identity, their columns and the relation types
that connect them. That JSON is what this module loads.

protobuf stays on the serialization boundary: nothing here needs ``protoc``,
a descriptor pool, or a per-extension generated lookup table. The built-in
EASM extension and a third-party one are loaded by the same two calls, so
neither is privileged.
"""

from __future__ import annotations

import json
import dataclasses
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Dict, Mapping, Optional, Tuple

__all__ = (
    "FieldSchema",
    "NodeSchema",
    "RelationSchema",
    "FlagSchema",
    "ExtensionSchema",
    "registry",
    "node_type_url",
    "node_type_from_url",
    "relation_type_url",
    "relation_type_from_url",
    "relation_types",
)

TYPE_URL_PREFIX = "type.googleapis.com/"
_SCHEMA_DIR = Path(__file__).resolve().parent / "schemas"
_SUPPORTED_SCHEMA_VERSION = 1


def _message_short_name(type_url: str) -> str:
    """Bare message name from a type URL, a full name, or a short name."""
    return type_url.rsplit("/", 1)[-1].rsplit(".", 1)[-1]


@dataclass(frozen=True)
class FieldSchema:
    """One column, described exactly as the source ``.proto`` declared it."""

    name: str
    number: int
    type: str
    repeated: bool = False
    optional: bool = False
    semantic: bool = True
    semantic_label: str = ""
    #: The field's value domain in order, lowest first; empty when it declares
    #: none. Carried so a caller sees the same contract the runtime compares
    #: by — the words belong to the extension, not to whoever reads them.
    ordered_values: Tuple[str, ...] = ()
    #: The column this field lands in when the proto type does not imply it.
    #: Only ``"json"``: text on the wire, a document in the column — a type's
    #: declared bag for values it has no column for.
    column: str = ""

    @classmethod
    def parse(cls, data: Mapping[str, Any]) -> "FieldSchema":
        name = str(data["name"])
        return cls(
            name=name,
            number=int(data["number"]),
            type=str(data["type"]),
            repeated=bool(data.get("repeated", False)),
            optional=bool(data.get("optional", False)),
            semantic=bool(data.get("semantic", True)),
            semantic_label=str(data.get("semantic_label") or name),
            ordered_values=tuple(str(value) for value in data.get("ordered_values", ())),
            column=str(data.get("column") or ""),
        )


@dataclass(frozen=True)
class NodeSchema:
    """One node type: its message, identity contract and columns."""

    node_type: str
    message: str
    value_field: Optional[str] = None
    #: The column carrying this type's display label, when it declares one.
    label_field: Optional[str] = None
    identity_field: Optional[str] = None
    identity_format: Optional[str] = None
    fields: Tuple[FieldSchema, ...] = ()

    @classmethod
    def parse(cls, node_type: str, data: Mapping[str, Any]) -> "NodeSchema":
        identity = data.get("identity") or {}
        return cls(
            node_type=node_type,
            message=str(data["message"]),
            value_field=data.get("value_field"),
            label_field=data.get("label_field"),
            identity_field=identity.get("field"),
            identity_format=identity.get("format"),
            fields=tuple(FieldSchema.parse(item) for item in data.get("fields", ())),
        )

    @property
    def type_url(self) -> str:
        return f"{TYPE_URL_PREFIX}{self.message}"

    def field(self, name: str) -> Optional[FieldSchema]:
        for item in self.fields:
            if item.name == name:
                return item
        return None

    def semantic_fields(self) -> Tuple[FieldSchema, ...]:
        return tuple(item for item in self.fields if item.semantic)


@dataclass(frozen=True)
class RelationSchema:
    """One relation type and the message that carries its edge payload."""

    relation_type: str
    message: str

    @classmethod
    def parse(cls, relation_type: str, data: Mapping[str, Any]) -> "RelationSchema":
        return cls(relation_type=relation_type, message=str(data["message"]))

    @property
    def type_url(self) -> str:
        return f"{TYPE_URL_PREFIX}{self.message}"


@dataclass(frozen=True)
class FlagSchema:
    """One flag an extension declares, and the bit it owns.

    The bit is identity, like a field number: it is what a stored mask means.
    A runtime holds the mechanism — a 64-bit mask — and never the vocabulary,
    so `honeypot` is EASM's word and lives in EASM's document.
    """

    bit: int
    default_exclude: bool = False

    @classmethod
    def parse(cls, data: Mapping[str, Any]) -> "FlagSchema":
        return cls(
            bit=int(data["bit"]),
            default_exclude=bool(data.get("default_exclude", False)),
        )


@dataclass(frozen=True)
class ExtensionSchema:
    """Every node and relation type contributed by one extension."""

    extension: str
    nodes: Dict[str, NodeSchema] = dataclasses.field(default_factory=dict)
    relations: Dict[str, RelationSchema] = dataclasses.field(default_factory=dict)
    #: Named judgements this extension makes about a node, keyed by name.
    flags: Dict[str, FlagSchema] = dataclasses.field(default_factory=dict)

    @classmethod
    def parse(cls, data: Mapping[str, Any]) -> "ExtensionSchema":
        version = int(data.get("schema_version", 0))
        if version != _SUPPORTED_SCHEMA_VERSION:
            raise ValueError(
                f"unsupported schema_version {version}; "
                f"cstxpy understands {_SUPPORTED_SCHEMA_VERSION}"
            )
        return cls(
            extension=str(data["extension"]),
            nodes={
                name: NodeSchema.parse(name, item)
                for name, item in (data.get("nodes") or {}).items()
            },
            relations={
                name: RelationSchema.parse(name, item)
                for name, item in (data.get("relations") or {}).items()
            },
            flags={
                name: FlagSchema.parse(item)
                for name, item in (data.get("flags") or {}).items()
            },
        )

    @classmethod
    def from_json(cls, text: str) -> "ExtensionSchema":
        return cls.parse(json.loads(text))


class SchemaRegistry:
    """A read view of the schemas a runtime holds.

    This is a projection, not an authority. Which extension may claim a node
    type, how a short message name resolves when two packages spell it the
    same way, which proto types a document may declare — every one of those
    is decided in Rust, and a document that reaches this class has already
    been accepted there. Deciding any of it a second time here is how the two
    ended up disagreeing about the same document.

    Lookups are flat on purpose: a caller asking for ``node_type_url("domain")``
    should not have to know whether ``domain`` came from the built-in extension
    or from one registered at runtime.
    """

    def __init__(self) -> None:
        self._extensions: Dict[str, ExtensionSchema] = {}
        self._nodes: Dict[str, NodeSchema] = {}
        self._relations: Dict[str, RelationSchema] = {}
        self._nodes_by_message: Dict[str, NodeSchema] = {}
        self._relations_by_message: Dict[str, RelationSchema] = {}

    def _install(self, schema: ExtensionSchema) -> ExtensionSchema:
        """Record one extension, replacing any earlier version of itself."""
        self._extensions[schema.extension] = schema
        self._reindex()
        return schema

    def _remove(self, extension: str) -> None:
        """Drop one extension from the view."""
        if self._extensions.pop(extension, None) is not None:
            self._reindex()

    def _reindex(self) -> None:
        self._nodes.clear()
        self._relations.clear()
        self._nodes_by_message.clear()
        self._relations_by_message.clear()
        # Short names are built only where they are unambiguous. Rust resolves
        # a bare message name to nothing when two packages both spell it that
        # way; answering with whichever landed first would be a different
        # answer to the same question.
        node_short: Dict[str, list] = {}
        relation_short: Dict[str, list] = {}
        for schema in self._extensions.values():
            for node in schema.nodes.values():
                self._nodes[node.node_type] = node
                self._nodes_by_message[node.message] = node
                node_short.setdefault(_message_short_name(node.message), []).append(node)
            for relation in schema.relations.values():
                self._relations[relation.relation_type] = relation
                self._relations_by_message[relation.message] = relation
                relation_short.setdefault(
                    _message_short_name(relation.message), []
                ).append(relation)
        for short, nodes in node_short.items():
            if len(nodes) == 1 and short not in self._nodes_by_message:
                self._nodes_by_message[short] = nodes[0]
        for short, relations in relation_short.items():
            if len(relations) == 1 and short not in self._relations_by_message:
                self._relations_by_message[short] = relations[0]

    # ── extensions ──

    def extensions(self) -> Tuple[str, ...]:
        return tuple(self._extensions)

    def extension(self, name: str) -> Optional[ExtensionSchema]:
        return self._extensions.get(name)

    # ── nodes ──

    def node(self, node_type: str) -> Optional[NodeSchema]:
        return self._nodes.get(node_type)

    def node_types(self) -> Tuple[str, ...]:
        return tuple(self._nodes)

    def node_type_url(self, node_type: str) -> Optional[str]:
        node = self._nodes.get(node_type)
        return node.type_url if node else None

    def node_type_from_url(self, type_url: str) -> Optional[str]:
        node = self._by_message(self._nodes_by_message, type_url)
        return node.node_type if node else None

    # ── relations ──

    def relation(self, relation_type: str) -> Optional[RelationSchema]:
        return self._relations.get(relation_type)

    def relation_types(self) -> Tuple[str, ...]:
        return tuple(self._relations)

    def relation_type_url(self, relation_type: str) -> Optional[str]:
        relation = self._relations.get(relation_type)
        return relation.type_url if relation else None

    def relation_type_from_url(self, type_url: str) -> Optional[str]:
        relation = self._by_message(self._relations_by_message, type_url)
        return relation.relation_type if relation else None

    @staticmethod
    def _by_message(index: Dict[str, Any], type_url: str) -> Any:
        """Exact message name first, bare name only as a fallback.

        The index carries both spellings, but the bare name is present only
        when it is unambiguous. Shortening the query before looking made the
        fully-qualified entries unreachable, so `acme.Asset` and `easm.Asset`
        answered as one.
        """
        full = str(type_url or "").rsplit("/", 1)[-1]
        return index.get(full) or index.get(_message_short_name(type_url))


registry = SchemaRegistry()


def load_schema(source: "str | Path | Mapping[str, Any]") -> ExtensionSchema:
    """Read one schema document into the view.

    Internal. This is not a registration entry: it validates nothing about
    whether the document may be registered, and a runtime that rejects a
    document will never call it. The one public way to declare types is
    ``runtime.extensions.register`` — a second entry here is how a type became
    visible to Python and unknown to the core.

    Accepts a path to a ``.schema.json``, its raw text, or an already-parsed
    mapping — whichever an extension has on hand.
    """
    if isinstance(source, Mapping):
        schema = ExtensionSchema.parse(source)
    else:
        text = str(source)
        # A document's own text is not a path, and asking the filesystem about
        # it raises rather than answering: a JSON document is longer than any
        # filename a kernel will accept. Decide by shape, not by stat.
        if not text.lstrip().startswith("{"):
            path = Path(text)
            if path.is_file():
                text = path.read_text(encoding="utf-8")
        schema = ExtensionSchema.from_json(text)
    return registry._install(schema)


def load_bundled(name: str) -> ExtensionSchema:
    """Read a schema shipped inside cstxpy (the built-in extensions).

    Internal, and the one thing that has to happen before any runtime exists:
    the EASM model classes are built from this document at import time, when
    there is nothing to project from. The file is byte-identical to the one
    the core compiles in, pinned by `tests/test_schema_document_copies_agree`.
    """
    return load_schema(_SCHEMA_DIR / f"{name}.schema.json")


def bundled_names() -> Tuple[str, ...]:
    if not _SCHEMA_DIR.is_dir():
        return ()
    return tuple(sorted(p.name[: -len(".schema.json")] for p in _SCHEMA_DIR.glob("*.schema.json")))


# ── module-level facade over the registry ──
#
# These mirror the queries call sites make. They read through the registry, so
# a schema registered at runtime answers them exactly like a bundled one.


def node_type_url(node_type: str) -> Optional[str]:
    return registry.node_type_url(node_type)


def node_type_from_url(type_url: str) -> Optional[str]:
    return registry.node_type_from_url(type_url)


def relation_type_url(relation_type: str) -> Optional[str]:
    return registry.relation_type_url(relation_type)


def relation_type_from_url(type_url: str) -> Optional[str]:
    return registry.relation_type_from_url(type_url)


def relation_types() -> Tuple[str, ...]:
    return registry.relation_types()


def project(contract: "Any") -> None:
    """Replace the view with the schemas a runtime holds.

    Takes the documents out of an ``ExtensionContract`` the runtime handed
    back, so what Python can see is exactly what the core accepted. Rejected
    registrations never reach here, which is why this class needs no conflict
    policy of its own.

    Bundled extensions the contract does not mention are kept: they were read
    at import time, before any runtime existed.
    """
    for name, definition in contract.extensions.items():
        document = definition.schema
        if document:
            load_schema(document)


# The built-in extensions are read once, at import, because the model bases
# are built from them at class-definition time. Everything else arrives by
# projection from a runtime.
for _bundled in bundled_names():
    load_bundled(_bundled)
