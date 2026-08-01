"""Typed public surface for the low-level CSTX v0.3 Rust runtime.

The binding deliberately exposes ordinary ``dict``, ``list``, ``bytes``,
``int``, keyword arguments, and iterators.  It does not introduce public node,
filter, or options wrapper classes.  ``*_json`` methods are explicit transport
fast paths for data that is already JSON or must leave CSTX as JSON; they are
not replacements for the native Python APIs.
"""

from __future__ import annotations

from typing import Any, Iterator

__version__: str


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
    """Canonical field associated with the failure when known."""
    expected: str | None
    """Expected type or value description when available."""
    actual: str | None
    """Observed type or value description when available."""


class NodeCursor(Iterator[dict[str, Any]]):
    """Lazy node dictionaries that avoid materializing a full Python result set."""

    @property
    def closed(self) -> bool:
        """Whether the cursor was explicitly closed."""
        ...

    def close(self) -> None:
        """Release cursor state; repeated calls are safe."""
        ...

    def __enter__(self) -> NodeCursor:
        """Use the cursor as a context manager for deterministic cleanup."""
        ...

    def __exit__(self, *args: Any) -> None:
        """Close the cursor when its context exits."""
        ...


class EdgeCursor(Iterator[dict[str, Any]]):
    """Lazy edge dictionaries that postpone Python object construction."""

    @property
    def closed(self) -> bool:
        """Whether the cursor was explicitly closed."""
        ...

    def close(self) -> None:
        """Release cursor state; repeated calls are safe."""
        ...

    def __enter__(self) -> EdgeCursor:
        """Use the cursor as a context manager for deterministic cleanup."""
        ...

    def __exit__(self, *args: Any) -> None:
        """Close the cursor when its context exits."""
        ...


class Schemas:
    """Schema/plugin namespace sharing state with its owning ``CSTX`` runtime."""

    def register(
        self,
        node_type: str,
        schema: dict[str, Any],
        value_field: str | None = None,
    ) -> None:
        """Register canonical validation metadata without a Python schema wrapper."""
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

    def bfs(self, seed_id: str, depth: int = 0, reverse: bool = False) -> list[str]:
        """Return IDs reachable through directed breadth-first traversal."""
        ...

    def shortest_paths(
        self, start_id: str, end_id: str, max_depth: int = 5
    ) -> list[list[str]]:
        """Return shortest undirected paths between two nodes."""
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
        offset: int = 0,
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

    def add_edges(self, edges: list[dict[str, Any]]) -> int:
        """Atomically mutate native relationship dictionaries."""
        ...

    def node(self, node_id: str) -> dict[str, Any]:
        """Return one node dictionary or raise ``CSTXError(NOT_FOUND)``."""
        ...

    def find_node(self, identifier: str) -> dict[str, Any] | None:
        """Resolve a node by ID, value, or extras.name."""
        ...

    def patch_node_extras(
        self, node_ids: list[str], patch: dict[str, Any]
    ) -> int:
        """Merge contextual fields into selected node extras."""
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

    def difference(
        self, other: CSTXGraph, node_type: str | None = None
    ) -> CSTX:
        """Return nodes present only in this graph."""
        ...

    def is_path_expression(self, expression: str) -> bool:
        """Classify an expression using the native query parser."""
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
        self, exclude_mask: int = 0, include_mask: int = 0
    ) -> dict[str, dict[str, int]]:
        """Return aggregate counts, optionally filtered by CSTX flags."""
        ...

    def nodes(
        self,
        types: list[str] | None = None,
        ids: list[str] | None = None,
        sources: list[str] | None = None,
        flags_all: int = 0,
        flags_any: int = 0,
        flags_none: int = 0,
        limit: int | None = None,
        page_size: int = 1024,
        order: str = "unspecified",
    ) -> NodeCursor:
        """Create the normal lazy native-dictionary read path.

        Keyword filters avoid public filter/options wrapper objects. ``page_size``
        bounds incremental materialization and ``order`` accepts ``unspecified``,
        ``id_asc``, or ``id_desc``.
        """
        ...

    def nodes_page(
        self,
        node_type: str | None = None,
        name_pattern: str | None = None,
        exclude_mask: int = 0,
        include_mask: int = 0,
        limit: int = 500,
        offset: int = 0,
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
        page_size: int = 1024,
        order: str = "unspecified",
    ) -> EdgeCursor:
        """Create a lazy native-dictionary relationship cursor using keywords."""
        ...

    def neighbors(
        self,
        node_id: str,
        direction: str = "out",
        limit: int | None = None,
        page_size: int = 1024,
        order: str = "unspecified",
    ) -> NodeCursor:
        """Lazily traverse neighboring nodes without an eager Python result list."""
        ...

    def query(
        self,
        expression: str,
        limit: int | None = None,
        offset: int = 0,
        exclude_mask: int = 0,
        include_mask: int = 0,
        page_size: int = 1024,
        order: str = "unspecified",
    ) -> NodeCursor:
        """Execute the graph DSL and return the normal lazy dictionary path."""
        ...


class Repository:
    """Snapshot/version namespace separated from graph mutation concerns."""

    def index_graph(self) -> dict[str, Any]:
        """Derive canonical content hashes and CAS payloads from the current graph."""
        ...

    def commit(
        self,
        message: str,
        ref_name: str = "main",
        metadata: Any | None = None,
    ) -> dict[str, Any]:
        """Commit current graph content to a named in-memory repository ref."""
        ...

    def diff(self, base_ref: str, head_ref: str) -> dict[str, Any]:
        """Return added, removed, and modified IDs between two refs."""
        ...

    def dump(self, compression: str = "zstd1") -> bytes:
        """Return a compact snapshot byte container for storage or transfer."""
        ...

    def load(self, data: bytes, compression: str = "zstd1") -> int:
        """Replace graph state from a compact snapshot and return bytes consumed."""
        ...

    def dump_json(self) -> bytes:
        """Return canonical JSON snapshot bytes for interoperability/debugging."""
        ...

    def load_json(self, data: bytes) -> int:
        """Load existing canonical JSON bytes without Python object transcoding."""
        ...

    def snapshot_fingerprint(self) -> str:
        """Return a stable hash of canonical snapshot content."""
        ...

    def head(self, ref_name: str = "main") -> str | None:
        """Return one named repository ref without exposing Merkle internals."""
        ...

    def refs(self) -> list[tuple[str, str]]:
        """List repository refs without exposing raw Merkle primitives."""
        ...

    def validate_ref(self, ref_name: str) -> None:
        """Validate a ref and its complete loaded object graph."""
        ...

    def set_ref(self, ref_name: str, commit_hash: str) -> None:
        """Set a ref to a loaded commit after native integrity validation."""
        ...

    def delete_ref(self, ref_name: str) -> None:
        """Delete a named repository ref."""
        ...

    def fork(self, source: str, target: str) -> str:
        """Fork one loaded ref to another name."""
        ...

    def checkout_entries(
        self,
        ref_name: str,
        types: list[str] | None = None,
    ) -> dict[str, dict[str, str]]:
        """Return pre-hashed entries for a ref and optional element types."""
        ...

    def commit_entries(
        self,
        entries: dict[str, dict[str, str]],
        ref_name: str,
        message: str,
        metadata: dict[str, Any],
        created_at: int,
    ) -> dict[str, Any]:
        """Commit pre-hashed entries for an external CAS storage adapter."""
        ...

    def commit_delta(
        self,
        delta_entries: dict[str, dict[str, str]],
        removed_node_ids: list[str],
        removed_edge_ids: list[str],
        ref_name: str,
        message: str,
        metadata: dict[str, Any],
        created_at: int,
    ) -> dict[str, Any]:
        """Commit a pre-hashed delta for an external CAS storage adapter."""
        ...

    def merge(self, source: str, target: str, created_at: int) -> dict[str, Any]:
        """Merge one loaded ref into another."""
        ...

    def merge_base(self, ref_a: str, ref_b: str) -> str | None:
        """Return the nearest common ancestor of two loaded refs."""
        ...

    def log(self, ref_name: str, limit: int = 50) -> list[dict[str, Any]]:
        """Return commit history for a loaded ref, newest first."""
        ...

    def tree_stats(self, ref_name: str) -> dict[str, int]:
        """Return element counts by type for a loaded ref."""
        ...

    def cherry_pick(
        self,
        source_ref: str,
        target_ref: str,
        commit_hashes: list[str],
        created_at: int,
    ) -> dict[str, Any]:
        """Apply selected loaded commits to a target ref."""
        ...

    def import_commit(self, hash: str, data: dict[str, Any]) -> None:
        """Import one canonical commit object into native repository state."""
        ...

    def export_commit(self, hash: str) -> dict[str, Any] | None:
        """Export one canonical commit object from native repository state."""
        ...

    def export_tree_objects(
        self,
        root_hash: str,
    ) -> dict[str, dict[str, dict[str, str]]]:
        """Export all canonical Merkle objects reachable from a tree root."""
        ...

    def import_tree_objects(
        self,
        objects: dict[str, dict[str, dict[str, str]]],
    ) -> None:
        """Import canonical Merkle objects fetched by a storage adapter."""
        ...

    def tree_entries(
        self,
        root_hash: str,
        types: list[str] | None = None,
    ) -> dict[str, dict[str, str]]:
        """Return content-hash entries for a loaded tree root."""
        ...

    def find_tree_entry(self, root_hash: str, element_id: str) -> str | None:
        """Find one element's content hash in a loaded tree root."""
        ...

    def diff_tree_entries(
        self,
        base_root: str,
        head_root: str,
    ) -> dict[str, dict[str, tuple[str | None, str | None]]]:
        """Return the raw content-hash diff between two loaded tree roots."""
        ...

    def tree_root_stats(self, root_hash: str) -> dict[str, int]:
        """Return element counts by type for a loaded tree root."""
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
        model_revision: str,
        page_size: int = 512,
    ) -> RagRecordCursor:
        """Stream projected records through a bounded native cursor."""
        ...

    def pending_json(
        self,
        model_revision: str,
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


class Bm25Recall:
    """In-memory BM25 recall extension and default lexical implementation."""

    @property
    def id(self) -> str:
        """Return the stable extension identifier."""
        ...

    def apply(self, delta: dict[str, Any]) -> dict[str, Any]:
        """Atomically apply one graph projection delta."""
        ...

    def recall(self, request: dict[str, Any]) -> dict[str, Any]:
        """Return ranked record candidates for one recall request."""
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
