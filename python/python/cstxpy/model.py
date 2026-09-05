"""Pydantic model bases built from the runtime schema.

The high-level ``cstx`` package mixes its runtime model (``Element``/``SCO``)
into these bases. They are derived from :mod:`cstxpy.schema` at import time
rather than emitted by a per-extension code generator, so an extension that
registers a schema at runtime gets model bases on exactly the same terms as
the built-in one — and, since a node's payload crosses the boundary named by
that same schema, its models serialize on the same terms too.

protobuf is not involved here. ``cstxpy.proto`` carries the generated messages
used to serialize, and nothing in this module imports them.
"""

from __future__ import annotations

from typing import Any, Dict, List, Optional, Tuple, Type

from pydantic import BaseModel, ConfigDict, Field, create_model

from cstxpy.schema import FieldSchema, NodeSchema, registry

__all__ = ("base_model", "base_model_name", "python_type", "annotation_for")

_INT_TYPES = frozenset(
    {"int64", "sint64", "sfixed64", "int32", "uint32", "sint32", "fixed32", "sfixed32"}
)
_FLOAT_TYPES = frozenset({"double", "float"})

# One shared config for every derived base: extensions may send fields the
# schema does not declare, and scanner output routinely types numbers as text.
_MODEL_CONFIG = ConfigDict(extra="allow", coerce_numbers_to_str=True)

# Keyed by node type, but the schema it was built from is kept alongside: an
# extension re-registering a changed schema must not keep handing out a base
# built from the old one. NodeSchema is frozen and replaced wholesale, so an
# identity check is enough.
_cache: Dict[str, "Tuple[NodeSchema, Type[BaseModel]]"] = {}


#: Column kinds a document may declare, and the Python type each arrives as.
#: Read, not derived: the document is where a type says what it stores.
_DECLARED_COLUMNS = {"json": dict}


def python_type(field: FieldSchema) -> type:
    """The Python type one schema field arrives as.

    A field that declares its column is read, not guessed — `column: "json"`
    is text on the wire and a document in the column, and only the document
    knows that. Everything else follows protobuf's own type table: `int32` is
    an integer because protobuf says so.

    The column a field lands in is decided in exactly one place, Rust's
    ``FieldSchema::column_type``. Deciding it a second time here is what let
    ``int32``/``uint32``/``sint32``/``double`` fields be built as strings that
    the core then refused, so this mirrors that decision and
    ``test_derived_annotation_survives_the_column_it_lands_in`` walks every
    declarable type through the real path to keep the two honest.
    """
    declared = _DECLARED_COLUMNS.get(getattr(field, "column", "") or "")
    if declared is not None:
        return declared
    if field.type in _INT_TYPES:
        return int
    if field.type in _FLOAT_TYPES:
        return float
    if field.type == "bool":
        return bool
    return str


def annotation_for(field: FieldSchema) -> Tuple[Any, Any]:
    """``(annotation, default)`` for one schema field.

    proto3 scalar semantics: a field that carries no presence is always
    populated, so it defaults to its zero value rather than to ``None``. Only
    repeated fields distinguish "absent" from "empty".
    """
    if field.repeated:
        # Repeated fields of any element type are stored as string lists;
        # registration refuses a repeated non-string for that reason.
        return (Optional[List[str]], Field(default=None))
    if python_type(field) is dict:
        # A declared bag is absent until something goes in it. Unlike a proto3
        # scalar it has no zero value that means "unset" — `{}` would claim the
        # producer sent an empty bag.
        return (Optional[Dict[str, Any]], Field(default=None))
    annotation = python_type(field)
    if not field.optional:
        return (annotation, ...)
    return (annotation, Field(default=annotation()))


def base_model_name(node: NodeSchema) -> str:
    """``easm.Subdomain`` -> ``SubdomainBase``."""
    return f"{node.message.rsplit('.', 1)[-1]}Base"


def base_model(node_type: str) -> Type[BaseModel]:
    """Return (and memoize) the pydantic base for one node type."""
    node = registry.node(node_type)
    if node is None:
        raise KeyError(f"unknown node type: {node_type}")
    cached = _cache.get(node_type)
    if cached is not None and cached[0] is node:
        return cached[1]
    model = create_model(
        base_model_name(node),
        __config__=_MODEL_CONFIG,
        __doc__=f"Schema base for the {node.node_type!r} node type.",
        **{field.name: annotation_for(field) for field in node.fields},
    )
    _cache[node_type] = (node, model)
    return model


_by_class_name: Dict[str, str] = {}


def _index() -> Dict[str, str]:
    """``SubdomainBase`` -> ``subdomain``, refreshed as extensions register."""
    for node_type in registry.node_types():
        node = registry.node(node_type)
        if node is not None:
            _by_class_name[base_model_name(node)] = node_type
    return _by_class_name


def __getattr__(name: str) -> Type[BaseModel]:
    """Resolve ``DomainBase`` and friends without generating a module of stubs."""
    node_type = _index().get(name)
    if node_type is None:
        raise AttributeError(f"module {__name__!r} has no attribute {name!r}")
    return base_model(node_type)


def __dir__() -> List[str]:
    return [*__all__, *_index()]
