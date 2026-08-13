package cstx

import "context"

// Graph is the graph namespace of a CSTX runtime. It owns graph data, graph
// queries, and ingest; the repository lifecycle lives elsewhere.
type Graph struct{ eng engine }

// AddNodes atomically adds or merges nodes and returns the number of elements
// actually changed. A no-op write reports zero and does not invalidate
// cursors.
func (g *Graph) AddNodes(ctx context.Context, nodes []Node) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphAddNodes(ctx, nodes)
}

// AddEdges atomically adds or merges relationships.
func (g *Graph) AddEdges(ctx context.Context, edges []Edge) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphAddEdges(ctx, edges)
}

// Ingest feeds one linked native-plugin payload into the shared graph.
func (g *Graph) Ingest(ctx context.Context, source string, data []byte) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphIngest(ctx, source, data)
}

// Node returns one node or a *Error with CodeNotFound.
func (g *Graph) Node(ctx context.Context, nodeID string) (Node, error) {
	if err := contextError(ctx); err != nil {
		return Node{}, err
	}
	return g.eng.graphNode(ctx, nodeID)
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

// EdgeCount returns the current number of relationships.
func (g *Graph) EdgeCount(ctx context.Context) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return g.eng.graphEdgeCount(ctx)
}

// Stats returns small aggregate counts.
func (g *Graph) Stats(ctx context.Context) (GraphStats, error) {
	if err := contextError(ctx); err != nil {
		return GraphStats{}, err
	}
	return g.eng.graphStats(ctx)
}

// Nodes creates a lazy cursor over nodes matching the filter. The zero
// filter and options select everything with runtime defaults.
func (g *Graph) Nodes(ctx context.Context, filter NodeFilter, options CollectionOptions) (*GraphCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	cursor, err := g.eng.graphNodes(ctx, filter, options.normalize())
	if err != nil {
		return nil, err
	}
	return &GraphCursor{inner: cursor, kind: CursorKindNodes}, nil
}

// Edges creates a lazy cursor over relationships matching the filter.
func (g *Graph) Edges(ctx context.Context, filter EdgeFilter, options CollectionOptions) (*GraphCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	cursor, err := g.eng.graphEdges(ctx, filter, options.normalize())
	if err != nil {
		return nil, err
	}
	return &GraphCursor{inner: cursor, kind: CursorKindEdges}, nil
}

// Neighbors lazily traverses neighboring nodes. Direction is "out", "in",
// or "both".
func (g *Graph) Neighbors(ctx context.Context, nodeID, direction string, options CollectionOptions) (*GraphCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	cursor, err := g.eng.graphNeighbors(ctx, nodeID, direction, options.normalize())
	if err != nil {
		return nil, err
	}
	return &GraphCursor{inner: cursor, kind: CursorKindNodes}, nil
}

// Query executes the graph DSL and returns a lazy cursor over terminal nodes.
func (g *Graph) Query(ctx context.Context, expression string, options QueryOptions) (*GraphCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	options.Collection = options.Collection.normalize()
	cursor, err := g.eng.graphQuery(ctx, expression, options)
	if err != nil {
		return nil, err
	}
	return &GraphCursor{inner: cursor, kind: CursorKindNodes}, nil
}

// Analyze executes one graph algorithm. The result is nil, bool, or
// *GraphCursor according to the selected algorithm.
func (g *Graph) Analyze(ctx context.Context, algorithm map[string]any, selection ...string) (any, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if len(selection) > 1 {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.analyze", Message: "at most one selection expression is allowed"}
	}
	var selected *string
	if len(selection) == 1 {
		selected = &selection[0]
	}
	kind, boolean, cursor, err := g.eng.graphAnalyze(ctx, algorithm, selected)
	if err != nil {
		return nil, err
	}
	switch kind {
	case 0:
		return nil, nil
	case 1:
		return boolean, nil
	case 2:
		return &GraphCursor{inner: cursor, kind: algorithmCursorKind(algorithm)}, nil
	default:
		return nil, &Error{Code: CodeInternal, Operation: "graph.analyze", Message: "unknown algorithm result kind"}
	}
}

func algorithmCursorKind(algorithm map[string]any) CursorKind {
	switch algorithm["name"] {
	case "weak_components", "strong_components":
		return CursorKindComponents
	case "cycle_basis":
		return CursorKindCycles
	case "bridges":
		return CursorKindNodePairs
	case "core_numbers", "betweenness", "closeness":
		return CursorKindNodeScores
	case "shortest_paths":
		return CursorKindPaths
	case "leiden":
		return CursorKindCommunities
	default:
		return CursorKindNodes
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
