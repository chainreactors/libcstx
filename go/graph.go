package cstx

import (
	"context"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// Graph is the graph namespace of a CSTX runtime. It owns graph data, native
// ingestion, and queries; the repository lifecycle lives elsewhere.
type Graph struct{ eng engine }

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
// under one extras key are kept as both. That is what aggregating sightings of
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

// Ingest parses one artifact through a registered plugin. The raw parser bytes
// are nested in ParserPayload and cross the FFI only as protobuf.
func (g *Graph) Ingest(ctx context.Context, plugin, artifact string, data []byte) (cstxproto.GraphIngestResult, error) {
	if err := contextError(ctx); err != nil {
		return cstxproto.GraphIngestResult{}, err
	}
	return g.eng.graphIngest(ctx, plugin, artifact, data)
}

// Node returns one node or a *Error with CodeNotFound.
func (g *Graph) Node(ctx context.Context, nodeID string) (*cstxproto.Node, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return g.eng.graphNode(ctx, nodeID)
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
