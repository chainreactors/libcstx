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
func (g *Graph) Nodes(ctx context.Context, filter NodeFilter, options CollectionOptions) (*NodeCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	iter, err := g.eng.graphNodes(ctx, filter, options.normalize())
	if err != nil {
		return nil, err
	}
	return &NodeCursor{iter: iter}, nil
}

// Edges creates a lazy cursor over relationships matching the filter.
func (g *Graph) Edges(ctx context.Context, filter EdgeFilter, options CollectionOptions) (*EdgeCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	iter, err := g.eng.graphEdges(ctx, filter, options.normalize())
	if err != nil {
		return nil, err
	}
	return &EdgeCursor{iter: iter}, nil
}

// Neighbors lazily traverses neighboring nodes. Direction is "out", "in",
// or "both".
func (g *Graph) Neighbors(ctx context.Context, nodeID, direction string, options CollectionOptions) (*NodeCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	iter, err := g.eng.graphNeighbors(ctx, nodeID, direction, options.normalize())
	if err != nil {
		return nil, err
	}
	return &NodeCursor{iter: iter}, nil
}

// Query executes the graph DSL and returns a lazy cursor over terminal nodes.
func (g *Graph) Query(ctx context.Context, expression string, options QueryOptions) (*NodeCursor, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	options.Collection = options.Collection.normalize()
	iter, err := g.eng.graphQuery(ctx, expression, options)
	if err != nil {
		return nil, err
	}
	return &NodeCursor{iter: iter}, nil
}
