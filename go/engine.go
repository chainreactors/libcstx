package cstx

import (
	"context"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
	"google.golang.org/protobuf/types/known/structpb"
)

// engine is the internal boundary between the typed facade and the transport
// implementation. The native build implements it over the cstx-ffi C ABI;
// neither layer owns independent semantics.
type engine interface {
	close() error

	lastChange(context.Context) (*cstxproto.GraphChangeSet, error)

	extensionRegister(context.Context, *cstxproto.ExtensionContract) error
	extensionExportContract(context.Context) (cstxproto.ExtensionContract, error)
	extensionEnable(context.Context, string) error
	extensionList(context.Context) (*cstxproto.ExtensionCatalog, error)
	extensionInfo(context.Context, string) (*cstxproto.ExtensionInfo, error)
	extensionContains(context.Context, string) (bool, error)
	extensionSchema(context.Context, string) (cstxproto.NodeType, error)
	extensionSchemas(context.Context) (cstxproto.NodeTypeCatalog, error)
	extensionParsesArtifact(context.Context, string) (bool, error)
	extensionAnchorConcepts(context.Context) (cstxproto.AnchorConceptCatalog, error)

	graphParse(context.Context, *cstxproto.ParserPayload) (*cstxproto.Graph, uint64, error)
	graphAddNodes(context.Context, []*cstxproto.Node) (uint64, error)
	graphReplaceNodes(context.Context, []*cstxproto.Node) (uint64, error)
	graphAddRelationships(context.Context, []*cstxproto.Relationship) (uint64, error)
	graphAddRelationship(context.Context, *cstxproto.Relationship) (*cstxproto.Relationship, error)
	graphDeleteNodes(context.Context, []string) (uint64, error)
	graphDeleteRelationships(context.Context, []string) (uint64, error)
	graphNode(context.Context, string) (*cstxproto.Node, error)
	graphFindNode(context.Context, string) (*cstxproto.Node, error)
	graphRelationship(context.Context, string) (*cstxproto.Relationship, error)
	graphContains(context.Context, string) (bool, error)
	graphNodeCount(context.Context) (uint64, error)
	graphRelationshipCount(context.Context) (uint64, error)
	graphNodeTypes(context.Context) ([]string, error)
	graphDegree(context.Context, string, string) (uint64, error)
	graphStats(context.Context) (*cstxproto.GraphStats, error)
	graphNodes(context.Context, *cstxproto.NodeQuery) (graphCursor, error)
	graphRelationships(context.Context, *cstxproto.RelationshipQuery) (graphCursor, error)
	graphNeighbors(context.Context, *cstxproto.NeighborQuery) (graphCursor, error)
	graphQuery(context.Context, *cstxproto.GraphQuery) (graphCursor, error)
	graphAnalyze(context.Context, *cstxproto.Algorithm, *string) (uint8, bool, graphCursor, error)
	graphLink(context.Context, []string, string) (*cstxproto.GraphLinkResult, error)
	graphUpdateNodeFlags(context.Context, *cstxproto.NodeFlagChange) (uint64, error)
	graphPatchNodeAnnotations(context.Context, *cstxproto.NodeAnnotationUpdate) (uint64, error)
	graphFindAnchors(context.Context, string) (*cstxproto.GraphAnchorCatalog, error)

	// Derived-graph operations return an independently owned engine, matching
	// the Python methods that return a new CSTX.
	graphSubgraph(context.Context, []string, uint32) (engine, error)
	graphQuerySubgraph(context.Context, *cstxproto.GraphQuery) (engine, error)
	graphInducedSubgraph(context.Context, []string, []string) (engine, error)
	graphFilter(context.Context, *cstxproto.NodeFilter) (engine, error)
	graphFilterWithReasons(context.Context, *cstxproto.NodeFilter) (engine, *cstxproto.GraphProjectionReport, error)
	graphElevate(context.Context, string) (engine, error)
	graphUnion(context.Context, engine) (engine, error)
	graphDifference(context.Context, engine, string) (engine, error)
	graphMerge(context.Context, engine) (uint64, error)

	ragIndex(context.Context, *cstxproto.RagIndexPlan) (ragIndexSession, error)
	ragRetrieve(context.Context, *cstxproto.RagQuery) (ragRetrieval, error)

	repoResolve(context.Context, string) (string, error)
	repoHead(context.Context, string) (*string, error)
	repoCheckout(context.Context, string, bool) (*cstxproto.Commit, error)
	repoCommit(context.Context, string, string, *string, *structpb.Struct) (*cstxproto.Commit, error)
	repoPrepare(context.Context, string, string, *string, *structpb.Struct, *int64) (*cstxproto.PublicationPlan, error)
	repoAccept(context.Context, string) error
	repoDiscard(context.Context) error
	repoSynchronize(context.Context, *cstxproto.RepositoryState) error
	repoContains(context.Context, string) (bool, error)
	repoMissing(context.Context, *cstxproto.RepositoryObjectPlan) (*cstxproto.ObjectSelection, error)
	repoReleaseTransientObjects(context.Context) error
	repoDiff(context.Context, string, string, *uint64, cstxproto.DiffDetail) (*cstxproto.GraphDiff, error)
	repoLog(context.Context, string, int) (*cstxproto.CommitLog, error)
	repoEntities(context.Context, string, []string) (*cstxproto.Graph, error)
	repoHistory(context.Context, string, string, *int) (*cstxproto.EntityHistory, error)
	repoBranch(context.Context, string, string) (string, error)
	repoMerge(context.Context, string, string, *string, *string) (*cstxproto.Commit, error)
	repoStat(context.Context, string, uint64, uint64) (*cstxproto.GraphStats, error)
	repoDelta(context.Context, string, *int64, *int64) (*cstxproto.GraphChangeSummary, error)
}

type graphCursor interface {
	page(context.Context, int, int) (*cstxproto.GraphResultPage, error)
	close()
}

// ragIndexSession is a retained projection. It is a handle rather than a value
// because the records it holds are streamed in bounded pages instead of being
// materialized at once.
type ragIndexSession interface {
	metadata(context.Context) (*cstxproto.RagIndexResult, error)
	pending(context.Context, int, int) (*cstxproto.RagRecordPage, error)
	deletes(context.Context) ([]string, error)
	records(context.Context) (ragRecordIterator, error)
	close()
}

type ragRecordIterator interface {
	next(context.Context) (*cstxproto.RagRecord, bool, error)
	close()
}

// ragRetrieval is a suspended retrieval: the plan is read with requests, the
// embedder's answers are handed back through complete, and completing consumes
// the retrieval.
type ragRetrieval interface {
	requests(context.Context) (*cstxproto.RecallPlan, error)
	complete(context.Context, *cstxproto.RecallResults) (*cstxproto.RagResult, error)
	close()
}
