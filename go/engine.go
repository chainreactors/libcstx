package cstx

import "context"

// engine is the internal boundary between the typed facade and the transport
// implementation. The native build implements it over the cstx-ffi C ABI; the
// stub build returns ErrNativeUnavailable. Neither layer owns semantics.
type engine interface {
	close() error

	lastChange(context.Context) (ChangeSet, error)

	schemaRegister(context.Context, string, map[string]any, string) error
	schemaContains(context.Context, string) (bool, error)
	schemaGet(context.Context, string) (map[string]any, error)
	schemaList(context.Context) ([]map[string]any, error)
	schemaLoadPlugin(context.Context, string) error
	schemaLoadAllPlugins(context.Context) error
	schemaAvailablePlugins(context.Context) ([]string, error)

	graphAddNodes(context.Context, []Node) (uint64, error)
	graphAddEdges(context.Context, []Edge) (uint64, error)
	graphIngest(context.Context, string, []byte) (uint64, error)
	graphNode(context.Context, string) (Node, error)
	graphContains(context.Context, string) (bool, error)
	graphNodeCount(context.Context) (uint64, error)
	graphEdgeCount(context.Context) (uint64, error)
	graphStats(context.Context) (GraphStats, error)
	graphNodes(context.Context, NodeFilter, CollectionOptions) (nodeIter, error)
	graphEdges(context.Context, EdgeFilter, CollectionOptions) (edgeIter, error)
	graphNeighbors(context.Context, string, string, CollectionOptions) (nodeIter, error)
	graphQuery(context.Context, string, QueryOptions) (nodeIter, error)
	graphSubgraph(context.Context, []string, uint32) (engine, error)

	repoCommit(context.Context, string, string, any) (Commit, error)
	repoDiff(context.Context, string, string) (GraphDiff, error)
	repoDump(context.Context, string) ([]byte, error)
	repoLoad(context.Context, []byte, string) (uint64, error)
	repoDumpJSON(context.Context) ([]byte, error)
	repoLoadJSON(context.Context, []byte) (uint64, error)
	repoSnapshotFingerprint(context.Context) (string, error)
	repoHead(context.Context, string) (*string, error)
	repoRefs(context.Context) ([]Ref, error)
}

// nodeIter yields one node at a time from the native runtime. ok is false at
// exhaustion; close is idempotent.
type nodeIter interface {
	next(context.Context) (node Node, ok bool, err error)
	close()
}

type edgeIter interface {
	next(context.Context) (edge Edge, ok bool, err error)
	close()
}
