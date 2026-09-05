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
	extensionHasNativeArtifact(context.Context, string) (bool, error)
	extensionAnchorConcepts(context.Context) (cstxproto.AnchorConceptCatalog, error)

	graphAddNodes(context.Context, []*cstxproto.Node) (uint64, error)
	graphReplaceNodes(context.Context, []*cstxproto.Node) (uint64, error)
	graphAddRelationships(context.Context, []*cstxproto.Relationship) (uint64, error)
	graphDeleteNodes(context.Context, []string) (uint64, error)
	graphDeleteRelationships(context.Context, []string) (uint64, error)
	graphIngest(context.Context, string, string, []byte) (cstxproto.GraphIngestResult, error)
	graphNode(context.Context, string) (*cstxproto.Node, error)
	graphRelationship(context.Context, string) (*cstxproto.Relationship, error)
	graphContains(context.Context, string) (bool, error)
	graphNodeCount(context.Context) (uint64, error)
	graphRelationshipCount(context.Context) (uint64, error)
	graphStats(context.Context) (*cstxproto.GraphStats, error)
	graphNodes(context.Context, *cstxproto.NodeQuery) (graphCursor, error)
	graphRelationships(context.Context, *cstxproto.RelationshipQuery) (graphCursor, error)
	graphNeighbors(context.Context, *cstxproto.NeighborQuery) (graphCursor, error)
	graphQuery(context.Context, *cstxproto.GraphQuery) (graphCursor, error)
	graphAnalyze(context.Context, *cstxproto.Algorithm, *string) (uint8, bool, graphCursor, error)
	graphSubgraph(context.Context, []string, uint32) (engine, error)

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
