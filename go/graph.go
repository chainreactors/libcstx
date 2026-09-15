package cstx

import (
	"context"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// Graph is the graph namespace of a CSTX runtime. It owns graph data, native
// ingestion, and queries; the repository lifecycle lives elsewhere.
type Graph struct{ eng engine }

// Parse runs one artifact through the native extension parser that owns it and
// returns a graph batch plus the number of records the parser read. The graph
// is not mutated: the batch enters the same merge/compute/link path as a parser
// implemented in any other language, so callers pass it to AddNodes and then
// Link.
//
// The extension providing the parser must be enabled first, or the call reports
// CodeNotFound.
func (g *Graph) Parse(ctx context.Context, payload *cstxproto.ParserPayload) (*cstxproto.Graph, uint64, error) {
	if err := contextError(ctx); err != nil {
		return nil, 0, err
	}
	if payload == nil {
		return nil, 0, &Error{Code: CodeInvalidArgument, Operation: "graph.parse", Message: "payload must not be nil"}
	}
	return g.eng.graphParse(ctx, payload)
}

// AddNodes atomically adds or merges nodes and returns the number of elements
// actually changed. A no-op write reports zero and does not invalidate
// cursors.
func (g *Graph) AddNodes(ctx context.Context, nodes []*cstxproto.Node) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphAddNodes(ctx, nodes)
}

// ReplaceNodes atomically writes each node as its current state and returns the
// number of elements actually changed.
//
// AddNodes merges: fields fill in, sources accumulate, and two different values
// under one annotations key are kept as both. That is what aggregating sightings of
// one entity needs. ReplaceNodes is for records that have a current value — an
// oracle that moved from "future" to "intent" has one status — where merging
// would silently keep the old value alongside the new one. Restating an
// unchanged record still reports zero and writes no history.
func (g *Graph) ReplaceNodes(ctx context.Context, nodes []*cstxproto.Node) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphReplaceNodes(ctx, nodes)
}

// AddRelationships atomically adds or merges generated protobuf relationships.
func (g *Graph) AddRelationships(ctx context.Context, relationships []*cstxproto.Relationship) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphAddRelationships(ctx, relationships)
}

// AddRelationship adds or merges one generated protobuf relationship and
// returns its canonical stored representation.
func (g *Graph) AddRelationship(ctx context.Context, relationship *cstxproto.Relationship) (*cstxproto.Relationship, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if relationship == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.add_relationship", Message: "relationship must not be nil"}
	}
	return g.eng.graphAddRelationship(ctx, relationship)
}

// Link derives relationships between the selected nodes from the registered
// extension schemas and reports what it created or updated. dataSource names
// the origin recorded on the new relationships, so an unlinked batch stays
// traceable to the artifact it came from.
func (g *Graph) Link(ctx context.Context, nodeIDs []string, dataSource string) (*cstxproto.GraphLinkResult, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return g.eng.graphLink(ctx, nodeIDs, dataSource)
}

// DeleteNodes atomically removes nodes and all incident relationships.
func (g *Graph) DeleteNodes(ctx context.Context, nodeIDs []string) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphDeleteNodes(ctx, nodeIDs)
}

// DeleteRelationships atomically removes relationships by stable CSTX ID.
func (g *Graph) DeleteRelationships(ctx context.Context, relationshipIDs []string) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphDeleteRelationships(ctx, relationshipIDs)
}

// Node returns one node or a *Error with CodeNotFound.
func (g *Graph) Node(ctx context.Context, nodeID string) (*cstxproto.Node, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return g.eng.graphNode(ctx, nodeID)
}

// FindNode resolves a human-readable identifier — an IP, a domain, a URL —
// to its node. Node takes a stable CSTX ID and FindNode takes the thing a
// user actually typed; both report CodeNotFound when nothing matches.
func (g *Graph) FindNode(ctx context.Context, identifier string) (*cstxproto.Node, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return g.eng.graphFindNode(ctx, identifier)
}

// Relationship returns one generated protobuf relationship or CodeNotFound.
func (g *Graph) Relationship(ctx context.Context, relationshipID string) (*cstxproto.Relationship, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return g.eng.graphRelationship(ctx, relationshipID)
}

// Contains reports node existence without materializing the node.
func (g *Graph) Contains(ctx context.Context, nodeID string) (bool, error) {
	if err := contextError(ctx); err != nil {
		return false, err
	}
	return g.eng.graphContains(ctx, nodeID)
}

// NodeCount returns the current number of nodes.
func (g *Graph) NodeCount(ctx context.Context) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphNodeCount(ctx)
}

// RelationshipCount returns the current number of relationships.
func (g *Graph) RelationshipCount(ctx context.Context) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphRelationshipCount(ctx)
}

// NodeTypes returns the distinct node types currently present in the graph.
// It reports what the data holds, not what the registered extensions declare —
// use Extensions.Schemas for the declared catalog.
func (g *Graph) NodeTypes(ctx context.Context) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return g.eng.graphNodeTypes(ctx)
}

// Degree returns one node's relationship count in the given direction, which
// is "out", "in", or "both". An empty direction means "both".
func (g *Graph) Degree(ctx context.Context, nodeID, direction string) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	if direction == "" {
		direction = "both"
	}
	return g.eng.graphDegree(ctx, nodeID, direction)
}

// UpdateNodeFlags atomically applies one flag change to the selected nodes and
// returns how many changed. Flag bits are declared by extensions, not by this
// SDK; read them with FlagRegistry.
func (g *Graph) UpdateNodeFlags(ctx context.Context, change *cstxproto.NodeFlagChange) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	if change == nil {
		return 0, &Error{Code: CodeInvalidArgument, Operation: "graph.update_node_flags", Message: "change must not be nil"}
	}
	return g.eng.graphUpdateNodeFlags(ctx, change)
}

// PatchNodeAnnotations merges an annotation patch into the selected nodes and
// returns how many changed. An empty selection patches every node.
func (g *Graph) PatchNodeAnnotations(ctx context.Context, update *cstxproto.NodeAnnotationUpdate) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	if update == nil {
		return 0, &Error{Code: CodeInvalidArgument, Operation: "graph.patch_node_annotations", Message: "update must not be nil"}
	}
	return g.eng.graphPatchNodeAnnotations(ctx, update)
}

// FindAnchors returns the nodes anchoring a named concept. The concept names
// come from Extensions.AnchorConcepts.
func (g *Graph) FindAnchors(ctx context.Context, conceptName string) (*cstxproto.GraphAnchorCatalog, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return g.eng.graphFindAnchors(ctx, conceptName)
}

// Stats returns small aggregate counts.
func (g *Graph) Stats(ctx context.Context) (*cstxproto.GraphStats, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return g.eng.graphStats(ctx)
}

// Nodes creates a lazy cursor over nodes matching the filter. The zero
// filter and options select everything with runtime defaults.
func (g *Graph) Nodes(ctx context.Context, query *cstxproto.NodeQuery) (*GraphCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if query == nil {
		query = &cstxproto.NodeQuery{}
	}
	cursor, err := g.eng.graphNodes(ctx, query)
	if err != nil {
		return nil, err
	}
	return &GraphCursor{inner: cursor, kind: CursorKindNodes}, nil
}

// Relationships creates a lazy cursor over generated protobuf relationships.
func (g *Graph) Relationships(ctx context.Context, query *cstxproto.RelationshipQuery) (*GraphCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if query == nil {
		query = &cstxproto.RelationshipQuery{}
	}
	cursor, err := g.eng.graphRelationships(ctx, query)
	if err != nil {
		return nil, err
	}
	return &GraphCursor{inner: cursor, kind: CursorKindRelationships}, nil
}

// Neighbors lazily traverses neighboring nodes. Direction is "out", "in",
// or "both".
func (g *Graph) Neighbors(ctx context.Context, query *cstxproto.NeighborQuery) (*GraphCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if query == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.neighbors", Message: "query must not be nil"}
	}
	cursor, err := g.eng.graphNeighbors(ctx, query)
	if err != nil {
		return nil, err
	}
	return &GraphCursor{inner: cursor, kind: CursorKindNodes}, nil
}

// Query executes the graph DSL and returns a lazy cursor over terminal nodes.
func (g *Graph) Query(ctx context.Context, query *cstxproto.GraphQuery) (*GraphCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if query == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.query", Message: "query must not be nil"}
	}
	cursor, err := g.eng.graphQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	return &GraphCursor{inner: cursor, kind: CursorKindNodes}, nil
}

// Analyze executes one generated protobuf algorithm. A boolean result is
// returned through boolean; collection results use cursor. Both are nil when
// the algorithm has no value result.
func (g *Graph) Analyze(ctx context.Context, algorithm *cstxproto.Algorithm, selection *string) (cursor *GraphCursor, boolean *bool, err error) {
	if err := contextError(ctx); err != nil {
		return nil, nil, err
	}
	if algorithm == nil {
		return nil, nil, &Error{Code: CodeInvalidArgument, Operation: "graph.analyze", Message: "algorithm must not be nil"}
	}
	kind, value, nativeCursor, err := g.eng.graphAnalyze(ctx, algorithm, selection)
	if err != nil {
		return nil, nil, err
	}
	switch kind {
	case 0:
		return nil, nil, nil
	case 1:
		return nil, &value, nil
	case 2:
		return &GraphCursor{inner: nativeCursor, kind: algorithmCursorKind(algorithm)}, nil, nil
	default:
		return nil, nil, &Error{Code: CodeInternal, Operation: "graph.analyze", Message: "unknown algorithm result kind"}
	}
}

// Subgraph returns an independently owned runtime containing nodes reachable
// from the selected seeds within depth hops. The caller must close it.
func (g *Graph) Subgraph(ctx context.Context, seedIDs []string, depth uint32) (*CSTX, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	eng, err := g.eng.graphSubgraph(ctx, seedIDs, depth)
	if err != nil {
		return nil, err
	}
	return wrapRuntime(eng, "derived"), nil
}

// QuerySubgraph runs a graph DSL query and materializes exactly the nodes and
// relationships it traversed as an independently owned runtime. Query returns
// a cursor over terminal nodes; this keeps the whole traced path. The caller
// must close the result.
func (g *Graph) QuerySubgraph(ctx context.Context, query *cstxproto.GraphQuery) (*CSTX, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if query == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.query_subgraph", Message: "query must not be nil"}
	}
	eng, err := g.eng.graphQuerySubgraph(ctx, query)
	if err != nil {
		return nil, err
	}
	return wrapRuntime(eng, "derived"), nil
}

// InducedSubgraph materializes the named nodes as an independently owned
// runtime. A nil relationshipIDs keeps every relationship between those nodes;
// a non-nil one keeps only the relationships named. The caller must close the
// result.
func (g *Graph) InducedSubgraph(ctx context.Context, nodeIDs, relationshipIDs []string) (*CSTX, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	eng, err := g.eng.graphInducedSubgraph(ctx, nodeIDs, relationshipIDs)
	if err != nil {
		return nil, err
	}
	return wrapRuntime(eng, "derived"), nil
}

// Filter projects the graph through a node filter and returns the result as an
// independently owned runtime. The caller must close it.
func (g *Graph) Filter(ctx context.Context, filter *cstxproto.NodeFilter) (*CSTX, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if filter == nil {
		filter = &cstxproto.NodeFilter{}
	}
	eng, err := g.eng.graphFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return wrapRuntime(eng, "derived"), nil
}

// FilterWithReasons is Filter plus a report naming every excluded node and why
// it was excluded, and whether the projection reused an existing one. Use it
// when a filtered result has to be explainable; use Filter when it does not.
// The caller must close the returned runtime.
func (g *Graph) FilterWithReasons(ctx context.Context, filter *cstxproto.NodeFilter) (*CSTX, *cstxproto.GraphProjectionReport, error) {
	if err := contextError(ctx); err != nil {
		return nil, nil, err
	}
	if filter == nil {
		filter = &cstxproto.NodeFilter{}
	}
	eng, report, err := g.eng.graphFilterWithReasons(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	return wrapRuntime(eng, "derived"), report, nil
}

// Elevate promotes a named concept's anchors into a graph organized around
// that concept and returns it as an independently owned runtime. The caller
// must close it.
func (g *Graph) Elevate(ctx context.Context, conceptName string) (*CSTX, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	eng, err := g.eng.graphElevate(ctx, conceptName)
	if err != nil {
		return nil, err
	}
	return wrapRuntime(eng, "derived"), nil
}

// Union returns a new runtime holding everything in this graph and in other.
// Neither input is modified; the caller must close the result.
func (g *Graph) Union(ctx context.Context, other *Graph) (*CSTX, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if other == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.union", Message: "other must not be nil"}
	}
	eng, err := g.eng.graphUnion(ctx, other.eng)
	if err != nil {
		return nil, err
	}
	return wrapRuntime(eng, "derived"), nil
}

// Difference returns a new runtime holding the nodes present in this graph but
// not in other. A non-empty nodeType restricts the comparison to that type.
// Neither input is modified; the caller must close the result.
func (g *Graph) Difference(ctx context.Context, other *Graph, nodeType string) (*CSTX, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if other == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.difference", Message: "other must not be nil"}
	}
	eng, err := g.eng.graphDifference(ctx, other.eng, nodeType)
	if err != nil {
		return nil, err
	}
	return wrapRuntime(eng, "derived"), nil
}

// Merge folds other into this graph in place and returns the number of
// elements changed. Union builds a third graph and leaves both inputs alone;
// Merge writes into this one.
func (g *Graph) Merge(ctx context.Context, other *Graph) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	if other == nil {
		return 0, &Error{Code: CodeInvalidArgument, Operation: "graph.merge", Message: "other must not be nil"}
	}
	return g.eng.graphMerge(ctx, other.eng)
}
