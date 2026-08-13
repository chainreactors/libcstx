package cstx

import "context"

// engine is the internal boundary between the typed facade and the transport
// implementation. The native build implements it over the cstx-ffi C ABI;
// neither layer owns independent semantics.
type engine interface {
	close() error

	lastChange(context.Context) (ChangeSet, error)

	schemaImport(context.Context, SchemaContract) error
	schemaExport(context.Context) (SchemaContract, error)
	schemaRegister(context.Context, string, map[string]any, string) error
	schemaRegisterJoinRule(context.Context, JoinRuleSpec) error
	schemaContains(context.Context, string) (bool, error)
	schemaGet(context.Context, string) (map[string]any, error)
	schemaList(context.Context) ([]map[string]any, error)
	schemaLoadPlugin(context.Context, string) error
	schemaLoadAllPlugins(context.Context) error
	schemaAvailablePlugins(context.Context) ([]string, error)
	schemaPluginArtifacts(context.Context, string) ([]string, error)
	schemaHasNativeArtifact(context.Context, string) (bool, error)
	schemaAnchorConcepts(context.Context) ([]AnchorConcept, error)

	graphAddNodes(context.Context, []Node) (uint64, error)
	graphAddEdges(context.Context, []Edge) (uint64, error)
	graphIngest(context.Context, string, []byte) (uint64, error)
	graphNode(context.Context, string) (Node, error)
	graphContains(context.Context, string) (bool, error)
	graphNodeCount(context.Context) (uint64, error)
	graphEdgeCount(context.Context) (uint64, error)
	graphStats(context.Context) (GraphStats, error)
	graphNodes(context.Context, NodeFilter, CollectionOptions) (graphCursor, error)
	graphEdges(context.Context, EdgeFilter, CollectionOptions) (graphCursor, error)
	graphNeighbors(context.Context, string, string, CollectionOptions) (graphCursor, error)
	graphQuery(context.Context, string, QueryOptions) (graphCursor, error)
	graphAnalyze(context.Context, any, *string) (uint8, bool, graphCursor, error)
	graphSubgraph(context.Context, []string, uint32) (engine, error)

	repoResolve(context.Context, string) (string, error)
	repoHead(context.Context, string) (*string, error)
	repoCheckout(context.Context, string, bool) (Commit, error)
	repoCommit(context.Context, string, string, *string, any) (Commit, error)
	repoDiff(context.Context, string, string, *int) (GraphDiff, error)
	repoDiffStat(context.Context, string, string) (Delta, error)
	repoLog(context.Context, string, int) ([]map[string]any, error)
	repoHistory(context.Context, string, string, *int) ([]map[string]any, error)
	repoBranch(context.Context, string, string) (string, error)
	repoMerge(context.Context, string, string, *string, *string) (Commit, error)
	repoStat(context.Context, string, uint64, uint64) (GraphStats, error)
	repoDelta(context.Context, string, *int64, *int64) (Delta, error)
}

type graphCursor interface {
	page(context.Context, int, int) (CursorPage, error)
	close()
}
