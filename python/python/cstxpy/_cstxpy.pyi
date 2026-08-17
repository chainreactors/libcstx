"""Typed public surface for the low-level CSTX Rust runtime.

The binding deliberately exposes ordinary ``dict``, ``list``, ``bytes``,
``int``, keyword arguments, and iterators.  It does not introduce public node,
filter, or options wrapper classes.  ``*_json`` methods are explicit transport
fast paths for data that is already JSON or must leave CSTX as JSON; they are
not replacements for the native Python APIs.
"""

from __future__ import annotations

from typing import Any, Iterator

__version__: str


def _object_id(envelope: bytes) -> bytes:
    """Validate one internal object envelope and return its 32-byte ObjectId."""
    ...


def _object_kind(envelope: bytes) -> str:
    """Validate one internal object envelope and return its closed object kind."""
    ...


def _verify_object(envelope: bytes) -> tuple[bytes, str]:
    """Validate once and return the internal ObjectId and object kind."""
    ...


def is_path_expression(expression: str) -> bool:
    """Classify CSTX path syntax with the native parser."""
    ...


class CSTXError(Exception):
    """Stable CSTX error whose attributes can be handled without parsing text."""

    code: str
    """Machine-readable error category such as ``NOT_FOUND``."""
    operation: str
    """Public operation that failed, for example ``graph.add_nodes``."""
    item_index: int | None
    """Index of the failing batch item when the error is item-specific."""
    field: str | None
    """CSTX field associated with the failure when known."""
    expected: str | None
    """Expected type or value description when available."""
    actual: str | None
    """Observed type or value description when available."""


class GraphCursor(Iterator[dict[str, Any]]):
    """Unified graph-result cursor with one-based ``limit + page`` pagination."""

    @property
    def kind(self) -> str:
        """Logical row shape emitted by this cursor."""
        ...

    def page(self, limit: int = 1024, page: int = 1) -> dict[str, Any]:
        """Materialize one page without rerunning the originating operation."""
        ...

    @property
    def closed(self) -> bool:
        """Whether the cursor was explicitly closed."""
        ...

    def close(self) -> None:
        """Release cursor state; repeated calls are safe."""
        ...

    def __enter__(self) -> GraphCursor:
        """Use the cursor as a context manager for deterministic cleanup."""
        ...

    def __exit__(self, *args: Any) -> None:
        """Close the cursor when its context exits."""
        ...


class Schemas:
    """Schema/plugin namespace sharing state with its owning ``CSTX`` runtime."""

    def import_schema(self, schema: dict[str, Any]) -> None:
        """Atomically validate and register a portable schema contract."""
        ...

    def export_schema(self) -> dict[str, Any]:
        """Export the complete portable schema contract."""
        ...

    def register(
        self,
        node_type: str,
        schema: dict[str, Any],
        value_field: str | None = None,
    ) -> None:
        """Register CSTX validation metadata without a Python schema wrapper."""
        ...

    def register_join_rule(self, rule: dict[str, Any]) -> None:
        """Register a native linker rule using the KeyExpr DSL."""
        ...

    def contains(self, node_type: str) -> bool:
        """Check schema existence without materializing the schema dictionary."""
        ...

    def get(self, node_type: str) -> dict[str, Any]:
        """Return one retained schema as an ordinary dictionary."""
        ...

    def list(self) -> list[dict[str, Any]]:
        """Return retained schemas in deterministic node-type order."""
        ...

    def load_plugin(self, name: str) -> None:
        """Load one linked native plugin into the shared graph engine."""
        ...

    def load_all_plugins(self) -> None:
        """Load every linked native plugin into the shared graph engine."""
        ...

    def available_plugins(self) -> list[str]:
        """List linked plugins without changing runtime state."""
        ...

    def plugin_artifacts(self, name: str) -> list[str]:
        """List artifacts provided by one linked plugin."""
        ...

    def has_native_artifact(self, artifact: str) -> bool:
        """Return whether a linked native parser supports this artifact."""
        ...

    def anchor_concepts(self) -> list[tuple[str, list[str]]]:
        """List native anchor concepts and member node types."""
        ...


class CSTXGraph:
    """Rust-owned in-memory graph handle."""

    def rag(self) -> Rag:
        """Return the GraphRAG extension bound to this graph."""
        ...

    def ingest_native(
        self, plugin: str, artifact: str, data: bytes
    ) -> dict[str, Any]:
        """Ingest plugin bytes and return detailed native mutation statistics."""
        ...

    def link(self, node_ids: list[str], data_source: str) -> dict[str, Any]:
        """Run native linker rules for selected nodes."""
        ...

    def update_node_flags(
        self,
        node_ids: list[str],
        add: int = 0,
        remove: int = 0,
        set_to: int | None = None,
    ) -> int:
        """Atomically update selected nodes' native flag bitsets."""
        ...

    def analyze(
        self, algorithm: dict[str, Any], selection: str | None = None
    ) -> bool | GraphCursor | None:
        """Execute one typed graph algorithm."""
        ...

    def degree(self, node_id: str, direction: str = "both") -> int:
        """Return one node's directional degree."""
        ...

    def subgraph(
        self, seed_ids: list[str] | None = None, depth: int = 0
    ) -> CSTX:
        """Materialize a seed neighborhood, or copy the full graph."""
        ...

    def query_subgraph(
        self,
        expression: str,
        limit: int | None = None,
        page: int = 1,
        exclude_mask: int = 0,
        include_mask: int = 0,
    ) -> CSTX:
        """Execute a traced query and materialize its exact subgraph in Rust."""
        ...

    def induced_subgraph(
        self, node_ids: list[str], edge_ids: list[str] | None = None
    ) -> CSTX:
        """Materialize selected nodes and optional edges."""
        ...

    def filter(
        self,
        exclude_mask: int = 0,
        include_mask: int = 0,
        excluded_ids: list[str] = [],
    ) -> CSTX:
        """Return a filtered graph handle."""
        ...

    def filter_with_reasons(
        self,
        exclude_mask: int = 0,
        include_mask: int = 0,
        excluded_ids: list[str] = [],
    ) -> tuple[CSTX, list[tuple[str, str]], bool]:
        """Return a filtered graph, exclusions with reasons, and reuse status."""
        ...

    def find_anchors(self, concept_name: str) -> list[dict[str, Any]]:
        """Find native anchor instances by concept name."""
        ...

    def elevate(self, concept_name: str) -> CSTX:
        """Return an elevated graph handle."""
        ...

    def add_nodes(self, nodes: list[dict[str, Any]]) -> int:
        """Atomically mutate native dictionaries without a JSON round trip."""
        ...

    def replace_nodes(self, nodes: list[dict[str, Any]]) -> int:
        """Atomically overwrite native dictionaries instead of merging them."""
        ...

    def add_edges(self, edges: list[dict[str, Any]]) -> int:
        """Atomically mutate native relationship dictionaries."""
        ...

    def delete_nodes(self, node_ids: list[str]) -> int:
        """Atomically remove nodes and their incident relationships."""
        ...

    def delete_edges(self, edge_ids: list[str]) -> int:
        """Atomically remove relationships by stable CSTX ID."""
        ...

    def node(self, node_id: str) -> dict[str, Any]:
        """Return one node dictionary or raise ``CSTXError(NOT_FOUND)``."""
        ...

    def edge(self, edge_id: str) -> dict[str, Any]:
        """Return one relationship dictionary or raise ``CSTXError(NOT_FOUND)``."""
        ...

    def find_node(self, identifier: str) -> dict[str, Any] | None:
        """Resolve a node by ID, value, or extras.name."""
        ...

    def patch_node_extras(
        self, node_ids: list[str] | None, patch: dict[str, Any]
    ) -> int:
        """Merge contextual fields into selected node extras; None selects all."""
        ...

    def create_relationship(
        self,
        source_id: str,
        target_id: str,
        relation: str,
        sources: list[str] = [],
        attrs: dict[str, Any] | None = None,
        identity_key: str | None = None,
    ) -> dict[str, Any]:
        """Create or merge a relationship with Rust-owned identity."""
        ...

    def union(self, other: CSTXGraph) -> CSTX:
        """Return the native union with another graph."""
        ...

    def merge(self, other: CSTXGraph) -> int:
        """Atomically merge another graph into this graph."""
        ...

    def difference(
        self, other: CSTXGraph, node_type: str | None = None
    ) -> CSTX:
        """Return nodes present only in this graph."""
        ...

    def node_types(self) -> list[str]:
        """Return distinct node types present in the graph."""
        ...

    def contains(self, node_id: str) -> bool:
        """Check node existence without materializing a node dictionary."""
        ...

    def node_count(self) -> int:
        """Return the current number of nodes."""
        ...

    def edge_count(self) -> int:
        """Return the current number of relationships."""
        ...

    def stats(
        self,
        exclude_mask: int = 0,
        include_mask: int = 0,
        selection: str | None = None,
    ) -> dict[str, dict[str, int]]:
        """Return aggregate counts for an optional query selection and flag masks."""
        ...

    def nodes(
        self,
        types: list[str] | None = None,
        ids: list[str] | None = None,
        sources: list[str] | None = None,
        name_contains: str | None = None,
        flags_all: int = 0,
        flags_any: int = 0,
        flags_none: int = 0,
        limit: int | None = None,
        page: int = 1,
        order: str = "unspecified",
    ) -> GraphCursor:
        """Create a unified cursor over matching node dictionaries.

        Keyword filters avoid public filter/options wrapper objects. ``order``
        accepts ``unspecified``, ``id_asc``, or ``id_desc``.
        """
        ...

    def nodes_page(
        self,
        node_type: str | None = None,
        name_pattern: str | None = None,
        exclude_mask: int = 0,
        include_mask: int = 0,
        limit: int = 500,
        page: int = 1,
    ) -> dict[str, Any]:
        """Return one bounded node page with exact totals and type counts."""
        ...

    def edges(
        self,
        source_id: str | None = None,
        target_id: str | None = None,
        relations: list[str] | None = None,
        sources: list[str] | None = None,
        limit: int | None = None,
        page: int = 1,
        order: str = "unspecified",
    ) -> GraphCursor:
        """Create a unified cursor over matching relationship dictionaries."""
        ...

    def neighbors(
        self,
        node_id: str,
        direction: str = "out",
        limit: int | None = None,
        page: int = 1,
        order: str = "unspecified",
    ) -> GraphCursor:
        """Return a unified cursor over neighboring nodes."""
        ...

    def query(
        self,
        expression: str,
        limit: int | None = None,
        page: int = 1,
        types: list[str] | None = None,
        ids: list[str] | None = None,
        name_contains: str | None = None,
        exclude_mask: int = 0,
        include_mask: int = 0,
        order: str = "unspecified",
    ) -> GraphCursor:
        """Execute the graph DSL once and return a unified result cursor."""
        ...


class Repository:
    """Git-like repository over one shared working tree."""

    def _release_transient_objects(self) -> None:
        """Drop immutable objects retained only for the completed operation."""
        ...

    def _contains(self, object: bytes) -> bool:
        """Return whether this native session already holds one object."""
        ...

    def _prepare(
        self,
        message: str,
        ref_name: str,
        expected_head: str | None = None,
        metadata: dict[str, Any] | None = None,
        timestamp: int | None = None,
    ) -> tuple[dict[str, Any], bytes, list[tuple[bytes, str, bytes]]]:
        """Prepare one complete commit payload for external publication."""
        ...

    def _accept(self, commit: bytes) -> None:
        """Accept a prepared commit after external persistence succeeds."""
        ...

    def _discard(self) -> None:
        """Discard the currently prepared repository transaction."""
        ...

    def _prepare_merge(
        self,
        source: str,
        target: str = "main",
        expected_head: str | None = None,
        message: str | None = None,
        metadata: dict[str, Any] | None = None,
        timestamp: int | None = None,
    ) -> tuple[dict[str, Any], bytes, list[tuple[bytes, str, bytes]]]:
        """Prepare one complete merge payload for external publication."""
        ...

    def _synchronize(
        self,
        objects: list[tuple[bytes, bytes]],
        refs: list[tuple[str, bytes | None]],
        indexes: list[tuple[bytes, bytes]],
    ) -> None:
        """Synchronize externally persisted objects, refs, and indexes."""
        ...

    def _missing_tree(self, commit: bytes) -> list[bytes]:
        """Return graph-tree objects missing from the native object set."""
        ...

    def _object_closure(self, commit: bytes) -> list[bytes]:
        """Return every object this commit and its ancestry are built from."""
        ...

    def _missing_stat(self, commit: bytes) -> list[bytes]:
        """Return the graph root needed for persisted statistics."""
        ...

    def _missing_merge(
        self,
        source: bytes,
        target: bytes | None = None,
    ) -> list[bytes]:
        """Return the commit frontier or graph objects needed by merge."""
        ...

    def _missing_delta(
        self,
        commit: bytes,
        start_timestamp: int | None = None,
        end_timestamp: int | None = None,
    ) -> list[bytes]:
        """Return index nodes needed for a time-bounded delta."""
        ...

    def _missing_prepare(self, commit: bytes) -> list[bytes]:
        """Return index nodes needed to prepare the working journal."""
        ...

    def _missing_history(self, commit: bytes, entity: str) -> list[bytes]:
        """Return index nodes needed for one entity history."""
        ...

    def _missing_commits(self, commit: bytes, limit: int) -> list[bytes]:
        """Return index nodes needed for a bounded commit log."""
        ...

    def _missing_diff(
        self,
        base: bytes,
        head: bytes,
        detail: str = "entities",
    ) -> list[bytes]:
        """Return index or graph objects needed for a revision diff.

        A limit never narrows the plan, so the request carries only the detail
        level: ``"counts"`` skips the pages a page summary already answers for.
        """
        ...

    def _commits(self, commit: bytes, limit: int) -> list[bytes]:
        """Return bounded first-parent commit objects for synchronization."""
        ...

    def resolve(self, revision: str) -> str:
        """Resolve a branch or object ID to a commit ID."""
        ...

    def head(self, ref_name: str = "main") -> str | None:
        """Return the commit currently referenced by a branch."""
        ...

    def checkout(
        self,
        revision: str = "main",
        force: bool = False,
    ) -> dict[str, Any]:
        """Replace the working tree with one committed graph."""
        ...

    def commit(
        self,
        message: str,
        ref_name: str = "main",
        expected_head: str | None = None,
        metadata: Any | None = None,
        timestamp: int | None = None,
    ) -> dict[str, Any]:
        """Commit the working tree and atomically advance one branch."""
        ...

    def diff(
        self,
        base: str,
        head: str,
        limit: int | None = None,
        detail: str = "entities",
    ) -> dict[str, Any]:
        """Compare two revisions.

        ``limit`` bounds the reported entity IDs; ``detail="counts"`` drops them
        entirely. ``stats`` counts the whole range either way.
        """
        ...

    def log(
        self,
        revision: str = "main",
        limit: int = 50,
    ) -> list[dict[str, Any]]:
        """Return first-parent commits newest first."""
        ...

    def history(
        self,
        entity_id: str,
        revision: str = "main",
        limit: int | None = None,
    ) -> list[dict[str, Any]]:
        """Return indexed changes for one node or relationship."""
        ...

    def branch(self, name: str, start_point: str = "main") -> str:
        """Create a branch at one revision."""
        ...

    def merge(
        self,
        source: str,
        target: str = "main",
        expected_head: str | None = None,
        message: str | None = None,
        metadata: Any | None = None,
        timestamp: int | None = None,
    ) -> dict[str, Any]:
        """Merge one branch or commit into a target branch."""
        ...

    def stat(
        self,
        revision: str = "main",
        exclude_mask: int = 0,
        include_mask: int = 0,
    ) -> dict[str, dict[str, int]]:
        """Return persisted graph aggregates without loading graph state."""
        ...

    def delta(
        self,
        revision: str = "main",
        start_timestamp: int | None = None,
        end_timestamp: int | None = None,
    ) -> dict[str, Any]:
        """Return a time-bounded commit-index delta without loading graph state."""
        ...


class Rag:
    """Graph-owned projection and retrieval planner."""

    def index(self, request: dict[str, Any]) -> RagIndexSession:
        """Project graph changes into a retained deterministic index session."""
        ...

    def retrieve(self, query: dict[str, Any]) -> RagRetrieval:
        """Suspend retrieval until external recall batches are supplied."""
        ...


class RagIndexSession:
    """Retained deterministic projection with bounded record streaming."""

    @property
    def operation_id(self) -> str:
        """Return the content-derived idempotent operation identity."""
        ...

    def commit(self) -> str:
        """Return the commit targeted by this projection."""
        ...

    @property
    def full(self) -> bool:
        """Report whether this projection replaces the complete index."""
        ...

    @property
    def upsert_count(self) -> int:
        """Return the number of projected records to upsert."""
        ...

    @property
    def delete_count(self) -> int:
        """Return the number of projected record IDs to delete."""
        ...

    def pending(
        self,
        _model_revision: str,
        page_size: int = 512,
    ) -> RagRecordCursor:
        """Stream projected records through a bounded native cursor."""
        ...

    def pending_json(
        self,
        _model_revision: str,
        batch_size: int = 512,
    ) -> bytes:
        """Return one projected-record batch as JSON bytes."""
        ...

    def deletes_json(self) -> bytes:
        """Return deleted record IDs as JSON bytes."""
        ...

    def close(self) -> None:
        """Release the retained projection."""
        ...

    @property
    def closed(self) -> bool:
        """Report whether the retained projection was released."""
        ...


class RagRecordCursor(Iterator[dict[str, Any]]):
    """Bounded iterator over projected record dictionaries."""

    def __iter__(self) -> RagRecordCursor:
        """Return this cursor as its iterator."""
        ...

    def __next__(self) -> dict[str, Any]:
        """Return the next projected record or raise StopIteration."""
        ...

    def close(self) -> None:
        """Release the cursor early."""
        ...

    @property
    def closed(self) -> bool:
        """Report whether the cursor was explicitly closed."""
        ...


class RagRetrieval:
    """Suspended retrieval bound to one graph generation and checkpoint."""

    def requests(self) -> list[dict[str, Any]]:
        """Return recall requests required to complete this retrieval."""
        ...

    def requests_json(self) -> bytes:
        """Return recall requests as JSON bytes."""
        ...

    def complete(self, batches: list[dict[str, Any]]) -> dict[str, Any]:
        """Fuse recall batches and build the structured graph context."""
        ...

    def complete_json(self, batches: bytes) -> bytes:
        """Fuse JSON batches and return the result as JSON bytes."""
        ...


class CSTX:
    """Single owner of shared schema, graph and repository state."""

    def __init__(
        self,
        project_id: str = "default",
        cursor_page_size: int = 1024,
    ) -> None:
        """Open an in-memory runtime with bounded cursor materialization."""
        ...

    @property
    def schemas(self) -> Schemas:
        """Return the lightweight schema namespace for this runtime."""
        ...

    @property
    def graph(self) -> CSTXGraph:
        """Return the Rust-owned graph sharing this runtime's state."""
        ...

    @property
    def repo(self) -> Repository:
        """Return the lightweight repository namespace for this runtime."""
        ...

    @property
    def closed(self) -> bool:
        """Whether the shared runtime has invalidated retained handles."""
        ...

    @property
    def project_id(self) -> str:
        """Immutable project namespace supplied at construction."""
        ...

    def close(self) -> None:
        """Close shared state and invalidate retained services/cursors."""
        ...

    def last_change(self) -> dict[str, Any]:
        """Return IDs changed by the most recent committed mutation."""
        ...

    def __enter__(self) -> CSTX:
        """Use the runtime as a context manager without creating another owner."""
        ...

    def __exit__(self, *args: Any) -> None:
        """Close the shared runtime when its context exits."""
        ...


class NodeFlags:
    """Discoverable namespace of engine-compatible integer bit constants.

    Graph APIs still accept ordinary ``int`` values; this class is not a node
    flag wrapper and cannot create instances.
    """

    NONE: int
    HONEYPOT: int
    NOISE: int
    FALSE_POSITIVE: int
    MANUAL_IGNORED: int
    THREAT_PRESENT: int
    HISTORIC_VULNERABLE: int
    INTERNAL: int

    @staticmethod
    def all_mask() -> int:
        """Return a mask containing every currently defined node flag."""
        ...

    @staticmethod
    def default_exclude_mask() -> int:
        """Return the engine's standard default-exclusion mask."""
        ...
