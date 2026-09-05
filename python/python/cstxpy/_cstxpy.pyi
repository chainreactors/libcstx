"""Typed public surface for the low-level CSTX Rust runtime."""

from __future__ import annotations

from typing import Any, Iterator

__version__: str


def decode_object(envelope: bytes) -> tuple[bytes, str]:
    """Validate an object envelope once and return its ID and kind."""
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


class Algorithm:
    """Typed graph algorithm request built at the Python boundary."""

    @staticmethod
    def bfs(seed_id: str, depth: int = 0, direction: str = "out",
            max_visited_nodes: int | None = None, timeout_ms: int | None = None) -> Algorithm:
        """Build a breadth-first traversal request."""
        ...
    @staticmethod
    def weak_components() -> Algorithm:
        """Build a weakly connected-components request."""
        ...
    @staticmethod
    def strong_components() -> Algorithm:
        """Build a strongly connected-components request."""
        ...
    @staticmethod
    def cycle_basis() -> Algorithm:
        """Build a cycle-basis request."""
        ...
    @staticmethod
    def bridges() -> Algorithm:
        """Build a bridge-edge request."""
        ...
    @staticmethod
    def articulation_points() -> Algorithm:
        """Build an articulation-point request."""
        ...
    @staticmethod
    def core_numbers() -> Algorithm:
        """Build a core-number request."""
        ...
    @staticmethod
    def is_dag() -> Algorithm:
        """Build a directed-acyclic-graph check request."""
        ...
    @staticmethod
    def topological_order() -> Algorithm:
        """Build a topological-order request."""
        ...
    @staticmethod
    def betweenness(include_endpoints: bool = False, normalized: bool = True,
                    top_k: int | None = None) -> Algorithm:
        """Build a betweenness-centrality request."""
        ...
    @staticmethod
    def closeness(wf_improved: bool = True, top_k: int | None = None) -> Algorithm:
        """Build a closeness-centrality request."""
        ...
    @staticmethod
    def leiden(resolution: float = 1.0, min_community_size: int = 2,
               top_k: int | None = None) -> Algorithm:
        """Build a Leiden community-detection request."""
        ...
    @staticmethod
    def shortest_paths(start_id: str, end_id: str, direction: str = "out",
                       max_depth: int = 0, limit: int = 10,
                       max_visited_nodes: int | None = None,
                       timeout_ms: int | None = None) -> Algorithm:
        """Build a shortest-path enumeration request."""
        ...


class GraphCursor(Iterator[bytes]):
    """Unified graph-result cursor with one-based ``limit + page`` pagination."""

    @property
    def kind(self) -> str:
        """Logical row shape emitted by this cursor."""
        ...

    def next(self) -> bytes | None:
        """Return the next Node or Relationship protobuf row."""
        ...

    def page(self, limit: int = 1024, page: int = 1) -> bytes:
        """Materialize one page as a typed ``GraphResultPage`` protobuf."""
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


class Extensions:
    """Unified extension lifecycle and schema namespace."""

    def register(self, contract: bytes) -> None:
        """Atomically register one serialized ExtensionContract protobuf."""
        ...

    def enable(self, name: str) -> None:
        """Explicitly enable one linked native Rust extension."""
        ...

    def list(self) -> bytes:
        """List metadata as an ``ExtensionCatalog`` protobuf."""
        ...

    def info(self, name: str) -> bytes:
        """Return metadata as an ``ExtensionInfo`` protobuf."""
        ...

    def contains(self, node_type: str) -> bool:
        """Check schema existence."""
        ...

    def schema(self, node_type: str) -> bytes:
        """Return one retained schema as a ``NodeType`` protobuf."""
        ...

    def schemas(self) -> bytes:
        """Return retained schemas as a ``NodeTypeCatalog`` protobuf."""
        ...

    def has_native_artifact(self, artifact: str) -> bool:
        """Return whether an enabled native parser supports an artifact."""
        ...

    def anchor_concepts(self) -> bytes:
        """List native concepts as an ``AnchorConceptCatalog`` protobuf."""
        ...


class CSTXGraph:
    """Rust-owned in-memory graph handle."""

    def rag(self) -> Rag:
        """Return the GraphRAG extension bound to this graph."""
        ...

    def add_nodes(self, data: bytes) -> int:
        """Add a serialized semantic ``Graph`` protobuf."""
        ...

    def replace_nodes(self, data: bytes) -> int:
        """Replace graph contents from a serialized semantic ``Graph`` protobuf."""
        ...

    def add_relationships(self, data: bytes) -> int:
        """Add relationships from a serialized semantic ``Graph`` protobuf."""
        ...

    def node(self, node_id: str) -> bytes:
        """Return one semantic ``Node`` protobuf."""
        ...

    def relationship(self, relationship_id: str) -> bytes:
        """Return one semantic ``Relationship`` protobuf."""
        ...

    def find_node(self, identifier: str) -> bytes:
        """Resolve an identifier and return a semantic ``Node`` protobuf."""
        ...

    def nodes(self, filter: bytes, window: bytes) -> GraphCursor:
        """Create a cursor from serialized ``NodeFilter`` and ``QueryWindow``."""
        ...

    def relationships(self, filter: bytes, window: bytes) -> GraphCursor:
        """Create a cursor from serialized ``RelationshipFilter`` and ``QueryWindow``."""
        ...

    def neighbors(self, data: bytes) -> GraphCursor:
        """Create a cursor from serialized ``NeighborQuery``."""
        ...

    def query(self, data: bytes) -> GraphCursor:
        """Create a cursor from serialized ``GraphQuery``."""
        ...

    def ingest(self, data: bytes) -> bytes:
        """Ingest serialized semantic ``ParserPayload`` bytes."""
        ...

    def link(self, selection: bytes, data_source: str) -> bytes:
        """Run linker rules from a ``GraphSelection`` protobuf payload."""
        ...

    def update_node_flags(self, data: bytes) -> int:
        """Atomically update selected nodes from a serialized ``NodeFlagChange``."""
        ...

    def analyze(
        self, algorithm: Algorithm, selection: str | None = None
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
        self, node_ids: list[str], relationship_ids: list[str] | None = None
    ) -> CSTX:
        """Materialize selected nodes and optional relationships."""
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

    def find_anchors(self, concept_name: str) -> bytes:
        """Find native anchors as a ``GraphAnchorCatalog`` protobuf."""
        ...

    def elevate(self, concept_name: str) -> CSTX:
        """Return an elevated graph handle."""
        ...

    def delete_nodes(self, node_ids: list[str]) -> int:
        """Atomically remove nodes and their incident relationships."""
        ...

    def delete_relationships(self, relationship_ids: list[str]) -> int:
        """Atomically remove relationships by stable CSTX ID."""
        ...

    def patch_node_extras(self, data: bytes) -> int:
        """Merge annotations from a serialized ``NodeAnnotationUpdate``."""
        ...

    def add_relationship(self, data: bytes) -> bytes:
        """Create or merge one relationship from a serialized protobuf."""
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

    def relationship_count(self) -> int:
        """Return the current number of relationships."""
        ...

    def stats(
        self,
        exclude_mask: int = 0,
        include_mask: int = 0,
        selection: str | None = None,
    ) -> bytes:
        """Return aggregate counts as a ``GraphStats`` protobuf payload."""
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
        metadata: bytes | None = None,
        timestamp: int | None = None,
    ) -> bytes:
        """Prepare a serialized ``PublicationPlan`` protobuf."""
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
        metadata: bytes | None = None,
        timestamp: int | None = None,
    ) -> bytes:
        """Prepare a serialized merge ``PublicationPlan`` protobuf."""
        ...

    def _synchronize(self, data: bytes) -> None:
        """Synchronize from a serialized ``RepositoryState`` protobuf."""
        ...

    def _missing(self, plan: bytes) -> bytes:
        """Return a serialized ``ObjectSelection`` protobuf for one plan."""
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
    ) -> bytes:
        """Replace the working tree and return a serialized ``Commit`` protobuf."""
        ...

    def commit(
        self,
        message: str,
        ref_name: str = "main",
        expected_head: str | None = None,
        metadata: bytes | None = None,
        timestamp: int | None = None,
    ) -> bytes:
        """Commit the working tree and return a serialized ``Commit`` protobuf."""
        ...

    def diff(
        self,
        base: str,
        head: str,
        limit: int | None = None,
        detail: str = "entities",
    ) -> bytes:
        """Compare two revisions and return a serialized ``GraphDiff`` protobuf.

        ``limit`` bounds the reported entity IDs; ``detail="counts"`` drops them
        entirely. ``stats`` counts the whole range either way.
        """
        ...

    def log(
        self,
        revision: str = "main",
        limit: int = 50,
    ) -> bytes:
        """Return a serialized ``CommitLog`` protobuf."""
        ...

    def history(
        self,
        entity_id: str,
        revision: str = "main",
        limit: int | None = None,
    ) -> bytes:
        """Return a serialized ``EntityHistory`` protobuf."""
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
        metadata: bytes | None = None,
        timestamp: int | None = None,
    ) -> bytes:
        """Merge one branch and return a serialized ``Commit`` protobuf."""
        ...

    def stat(
        self,
        revision: str = "main",
        exclude_mask: int = 0,
        include_mask: int = 0,
    ) -> bytes:
        """Return persisted aggregates as a serialized ``GraphStats`` protobuf."""
        ...

    def delta(
        self,
        revision: str = "main",
        start_timestamp: int | None = None,
        end_timestamp: int | None = None,
    ) -> bytes:
        """Return a time-bounded delta as a serialized ``GraphChangeSummary`` protobuf."""
        ...


class Rag:
    """Graph-owned projection and retrieval planner."""

    def index(self, data: bytes) -> RagIndexSession:
        """Project a serialized ``RagIndexPlan`` protobuf into a retained session."""
        ...

    def retrieve(self, data: bytes) -> RagRetrieval:
        """Suspend retrieval from a serialized ``RagQuery`` protobuf."""
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

    def deletes(self) -> list[str]:
        """Return deleted record IDs."""
        ...
    def close(self) -> None:
        """Release the retained projection."""
        ...

    @property
    def closed(self) -> bool:
        """Report whether the retained projection was released."""
        ...


class RagRecordCursor(Iterator[bytes]):
    """Bounded iterator over serialized ``RagRecord`` protobuf messages."""

    def __iter__(self) -> RagRecordCursor:
        """Return this cursor as its iterator."""
        ...

    def __next__(self) -> bytes:
        """Return the next projected record protobuf or raise StopIteration."""
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

    def requests(self) -> bytes:
        """Return a serialized ``RecallPlan`` protobuf."""
        ...

    def complete(self, data: bytes) -> bytes:
        """Fuse serialized ``RecallResults`` and return a ``RagResult`` protobuf."""
        ...


class CSTX:
    """Single owner of shared schema, graph and repository state."""

    def __init__(
        self,
        project_id: str = "default",
        cursor_page_size: int = 1024,
        payload_format: int = 0,
    ) -> None:
        """Open an in-memory runtime with bounded cursor materialization.

        ``payload_format`` is the ``cstx.PayloadFormat`` number: 0 returns node
        payloads as the stored ``Any``, 1 returns them as ``EntityValue``,
        which a caller with no generated message type can still read.
        """
        ...

    @property
    def extensions(self) -> Extensions:
        """Return the unified extension namespace for this runtime."""
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

    def last_change(self) -> bytes:
        """Return the most recent mutation as serialized GraphChangeSet protobuf."""
        ...

    def __enter__(self) -> CSTX:
        """Use the runtime as a context manager without creating another owner."""
        ...

    def __exit__(self, *args: Any) -> None:
        """Close the shared runtime when its context exits."""
        ...


