package cstx

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
	"github.com/chainreactors/libcstx/go/proto/easmproto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

var testContext = context.Background()

//go:embed testdata/conformance.json
var conformanceFixture []byte

func stringPtr(value string) *string { return &value }

func openRuntime(t *testing.T) *CSTX {
	t.Helper()
	rt, err := Open(testContext, &cstxproto.RuntimeConfig{ProjectId: "sdk-go-test"})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.Close(); err != nil {
			t.Fatalf("close: %v", err)
		}
	})
	// The built-in extension ships its own schema document; enabling it is
	// the whole registration step.
	if err := rt.Extensions.Enable(testContext, "easm"); err != nil {
		t.Fatalf("enable easm: %v", err)
	}
	return rt
}

func domainNode(value string) *cstxproto.Node {
	entity, err := anypb.New(&easmproto.Domain{Host: value})
	if err != nil {
		panic(err)
	}
	id := "domain:" + value
	return &cstxproto.Node{
		Id:      &id,
		Sources: []string{"test"},
		Entity:  entity,
	}
}

func usesRelationship(source, target string) *cstxproto.Relationship {
	id := "relationship:" + source + ":uses:" + target
	return &cstxproto.Relationship{
		Id:       &id,
		SourceId: source,
		TargetId: target,
		Sources:  []string{"test"},
		Relation: &anypb.Any{TypeUrl: "type.googleapis.com/easm.Uses"},
	}
}

func domainValue(t *testing.T, node *cstxproto.Node) string {
	t.Helper()
	var domain easmproto.Domain
	if node == nil || node.Entity == nil || node.Entity.UnmarshalTo(&domain) != nil {
		t.Fatalf("invalid domain node: %+v", node)
	}
	return domain.Host
}

func addDomain(t *testing.T, rt *CSTX, value string) uint64 {
	t.Helper()
	affected, err := rt.Graph.AddNodes(testContext, []*cstxproto.Node{domainNode(value)})
	if err != nil {
		t.Fatalf("add node %s: %v", value, err)
	}
	return affected
}

func TestRepositoryExternalPersistenceRoundTrip(t *testing.T) {
	writer := openRuntime(t)
	addDomain(t, writer, "persisted.example")

	prepared, err := writer.Repo.Prepare(
		testContext,
		"external persistence",
		"main",
		nil,
		&structpb.Struct{Fields: map[string]*structpb.Value{"source": structpb.NewStringValue("go-test")}},
		nil,
	)
	if err != nil {
		t.Fatalf("prepare: %v", err)
	}
	if prepared.Commit.Id == "" || prepared.IndexRoot == "" || len(prepared.Objects) == 0 {
		t.Fatalf("incomplete prepared payload: %+v", prepared)
	}

	objects := make(map[string]*cstxproto.RepositoryState_Object, len(prepared.Objects))
	var commitObject, indexObject *cstxproto.RepositoryState_Object
	for _, object := range prepared.Objects {
		stored := &cstxproto.RepositoryState_Object{Id: object.Id, Payload: append([]byte(nil), object.Payload...)}
		objects[object.Id] = stored
		if object.Kind == cstxproto.RepositoryObjectKind_REPOSITORY_OBJECT_KIND_COMMIT && object.Id == prepared.Commit.Id {
			commitObject = stored
		}
		if object.Id == prepared.IndexRoot {
			indexObject = stored
		}
	}
	if commitObject == nil || indexObject == nil {
		t.Fatal("prepared payload does not contain commit and index-root envelopes")
	}
	if err := writer.Repo.Accept(testContext, prepared.Commit.Id); err != nil {
		t.Fatalf("accept: %v", err)
	}

	reader := openRuntime(t)
	head := prepared.Commit.Id
	if err := reader.Repo.Synchronize(testContext, &cstxproto.RepositoryState{
		Objects: []*cstxproto.RepositoryState_Object{commitObject, indexObject},
	}); err != nil {
		t.Fatalf("synchronize commit objects: %v", err)
	}
	if err := reader.Repo.Synchronize(testContext, &cstxproto.RepositoryState{
		Refs:    []*cstxproto.RepositoryState_Ref{{Name: "main", CommitId: &head}},
		Indexes: []*cstxproto.RepositoryState_Index{{CommitId: head, IndexRoot: prepared.IndexRoot}},
	}); err != nil {
		t.Fatalf("synchronize commit frontier: %v", err)
	}

	for {
		missing, err := reader.Repo.Missing(testContext, &cstxproto.RepositoryObjectPlan{Kind: cstxproto.RepositoryPlanKind_REPOSITORY_PLAN_TREE, CommitId: head})
		if err != nil {
			t.Fatalf("plan missing tree: %v", err)
		}
		if len(missing.ObjectIds) == 0 {
			break
		}
		batch := make([]*cstxproto.RepositoryState_Object, 0, len(missing.ObjectIds))
		for _, id := range missing.ObjectIds {
			object, ok := objects[id]
			if !ok {
				t.Fatalf("planner requested unknown object %s", id)
			}
			batch = append(batch, object)
		}
		if err := reader.Repo.Synchronize(testContext, &cstxproto.RepositoryState{Objects: batch}); err != nil {
			t.Fatalf("hydrate tree: %v", err)
		}
	}
	if _, err := reader.Repo.Checkout(testContext, "main", true); err != nil {
		t.Fatalf("checkout hydrated main: %v", err)
	}
	node, err := reader.Graph.Node(testContext, "domain:persisted.example")
	if err != nil || domainValue(t, node) != "persisted.example" {
		t.Fatalf("restored node: %+v err=%v", node, err)
	}
	if err := reader.Repo.ReleaseTransientObjects(testContext); err != nil {
		t.Fatalf("release transient objects: %v", err)
	}
}

func TestGraphDeleteNodesCascadesAndCommits(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "delete-a.example")
	addDomain(t, rt, "delete-b.example")
	addDomain(t, rt, "keep.example")
	relationships := []*cstxproto.Relationship{
		usesRelationship("domain:delete-a.example", "domain:delete-b.example"),
		usesRelationship("domain:delete-b.example", "domain:keep.example"),
	}
	if _, err := rt.Graph.AddRelationships(testContext, relationships); err != nil {
		t.Fatalf("add relationships: %v", err)
	}
	base, err := rt.Repo.Commit(testContext, "base", "main", nil, nil)
	if err != nil {
		t.Fatalf("commit base: %v", err)
	}
	cursor, err := rt.Graph.Nodes(testContext, &cstxproto.NodeQuery{})
	if err != nil {
		t.Fatalf("open cursor: %v", err)
	}
	defer cursor.Close()

	affected, err := rt.Graph.DeleteNodes(testContext, []string{"domain:delete-b.example"})
	if err != nil || affected != 3 {
		t.Fatalf("delete node: affected=%d err=%v", affected, err)
	}
	if _, err := cursor.Page(testContext, 10, 1); !IsCode(err, CodeCursorInvalidated) {
		t.Fatalf("expected cursor invalidation, got %v", err)
	}
	if count, _ := rt.Graph.NodeCount(testContext); count != 2 {
		t.Fatalf("node count after delete=%d", count)
	}
	if count, _ := rt.Graph.RelationshipCount(testContext); count != 0 {
		t.Fatalf("relationship count after cascade=%d", count)
	}
	change, err := rt.LastChange(testContext)
	if err != nil || !reflect.DeepEqual(change.RemovedNodeIds, []string{"domain:delete-b.example"}) || len(change.RemovedRelationshipIds) != 2 {
		t.Fatalf("delete change=%+v err=%v", change, err)
	}
	head, err := rt.Repo.Commit(testContext, "delete", "main", &base.Id, nil)
	if err != nil {
		t.Fatalf("commit delete: %v", err)
	}
	diff, err := rt.Repo.Diff(testContext, base.Id, head.Id, nil, cstxproto.DiffDetail_DIFF_DETAIL_ENTITIES)
	if err != nil || !reflect.DeepEqual(diff.Removed.NodeIds, []string{"domain:delete-b.example"}) || len(diff.Removed.RelationshipIds) != 2 {
		t.Fatalf("delete diff=%+v err=%v", diff, err)
	}
}

func TestGraphDeleteIsAtomicOnMissingID(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "present.example")
	affected, err := rt.Graph.DeleteNodes(testContext, []string{"domain:present.example", "domain:missing.example"})
	if err == nil || affected != 0 {
		t.Fatalf("expected atomic validation failure: affected=%d err=%v", affected, err)
	}
	if count, _ := rt.Graph.NodeCount(testContext); count != 1 {
		t.Fatalf("failed delete changed graph: count=%d", count)
	}
}

func TestExtensionsSchemaSurface(t *testing.T) {
	rt := openRuntime(t)
	if err := rt.Extensions.Enable(testContext, "easm"); err != nil {
		t.Fatalf("enable easm: %v", err)
	}
	builder := newExtensionBuilder("sdk-test", "1")
	builder.Rule(&cstxproto.JoinRule{LeftTypeUrl: "type.googleapis.com/easm.Domain", RightTypeUrl: "type.googleapis.com/easm.Domain", RelationshipTypeUrl: "type.googleapis.com/easm.Uses", LeftKey: "domain", RightKey: "domain"})
	if err := rt.Extensions.Register(testContext, builder.Build()); err != nil {
		t.Fatalf("register schema contract: %v", err)
	}
	contains, err := rt.Extensions.Contains(testContext, "domain")
	if err != nil || !contains {
		t.Fatalf("contains: %v %v", contains, err)
	}
	schema, err := rt.Extensions.Schema(testContext, "domain")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if schema.TypeUrl != "type.googleapis.com/easm.Domain" {
		t.Fatalf("unexpected schema type: %v", schema.TypeUrl)
	}
	// value_field comes from the generated schema document, not from a
	// hand-written restatement of it.
	if schema.Metadata == nil || schema.Metadata.Fields["value_field"].GetStringValue() != "host" {
		t.Fatalf("unexpected schema metadata: %v", schema.Metadata)
	}
	list, err := rt.Extensions.Schemas(testContext)
	if err != nil || len(list.Schemas) == 0 {
		t.Fatalf("list: %v %v", list, err)
	}
	easmInfo, err := rt.Extensions.Info(testContext, "easm")
	if err != nil || !slices.Contains(easmInfo.Artifacts, "gogo") {
		t.Fatalf("easm artifacts: %+v err=%v", easmInfo, err)
	}
	concepts, err := rt.Extensions.AnchorConcepts(testContext)
	if err != nil || len(concepts.Concepts) == 0 || concepts.Concepts[0].Name == "" {
		t.Fatalf("anchor concepts: %+v err=%v", concepts, err)
	}
	if err := rt.Extensions.Enable(testContext, "easm"); err != nil {
		t.Fatalf("enable easm plugin: %v", err)
	}
	hasGogo, err := rt.Extensions.HasNativeArtifact(testContext, "gogo")
	if err != nil || !hasGogo {
		t.Fatalf("has native gogo artifact: %v err=%v", hasGogo, err)
	}
	gogo := []byte(`{"ip":"192.0.2.1","port":"80","protocol":"tcp","status":"200"}` + "\n")
	if result, err := rt.Graph.Ingest(testContext, "easm", "gogo", gogo); err != nil || result.RecordsParsed == 0 {
		t.Fatalf("ingest gogo: result=%v err=%v", result, err)
	}
}

func TestExtensions(t *testing.T) {
	// This one watches the disabled -> enabled transition, so it must not
	// use the fixture that enables easm up front.
	rt, err := Open(testContext, &cstxproto.RuntimeConfig{ProjectId: "sdk-go-extensions"})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = rt.Close() })
	ctx := testContext

	items, err := rt.Extensions.List(ctx)
	if err != nil || len(items.Extensions) == 0 || items.Extensions[0].Name != "easm" || items.Extensions[0].Enabled {
		t.Fatalf("extension list: %+v err=%v", items, err)
	}
	if err := rt.Extensions.Enable(ctx, "easm"); err != nil {
		t.Fatalf("enable easm: %v", err)
	}
	info, err := rt.Extensions.Info(ctx, "easm")
	if err != nil || !info.Enabled || info.Kind != "native" {
		t.Fatalf("extension info: %+v err=%v", info, err)
	}

	builder := newExtensionBuilder("sdk-external", "")
	inputSchema, _ := structpb.NewStruct(map[string]any{"type": "object"})
	builder.Parser("report", &cstxproto.ParserType{Artifact: "report", InputSchema: inputSchema})
	builder.Rule(&cstxproto.JoinRule{LeftTypeUrl: "type.googleapis.com/easm.Domain", RightTypeUrl: "type.googleapis.com/easm.Domain", RelationshipTypeUrl: "type.googleapis.com/easm.Uses", LeftKey: "domain", RightKey: "domain"})
	if err := rt.Extensions.Register(ctx, builder.Build()); err != nil {
		t.Fatalf("register external extension: %v", err)
	}
	contract, err := rt.Extensions.ExportContract(ctx)
	if err != nil {
		t.Fatalf("export extension contract: %v", err)
	}
	definition, ok := contract.Extensions["sdk-external"]
	if !ok || definition.Parsers["report"] == nil {
		t.Fatalf("exported contract lost sdk-external parser: %+v", contract.Extensions)
	}
	external, err := rt.Extensions.Info(ctx, "sdk-external")
	if err != nil || external.Kind != "external" || !external.Enabled || !slices.Contains(external.Artifacts, "report") {
		t.Fatalf("external extension info: %+v err=%v", external, err)
	}
}

func TestGraphMutationAndCursors(t *testing.T) {
	rt := openRuntime(t)

	if affected := addDomain(t, rt, "example.com"); affected != 1 {
		t.Fatalf("first add affected=%d", affected)
	}
	addDomain(t, rt, "www.example.com")

	if ok, err := rt.Graph.Contains(testContext, "domain:example.com"); err != nil || !ok {
		t.Fatalf("contains: %v %v", ok, err)
	}
	if count, err := rt.Graph.NodeCount(testContext); err != nil || count != 2 {
		t.Fatalf("node count: %v %v", count, err)
	}

	node, err := rt.Graph.Node(testContext, "domain:example.com")
	if err != nil {
		t.Fatalf("node: %v", err)
	}
	if node.GetId() != "domain:example.com" || domainValue(t, node) != "example.com" {
		t.Fatalf("unexpected node: %+v", node)
	}

	if _, err := rt.Graph.Node(testContext, "domain:missing.example"); !IsCode(err, CodeNotFound) {
		t.Fatalf("expected NOT_FOUND, got %v", err)
	}

	// Re-adding identical content is a no-op: zero affected and live cursors
	// stay valid.
	cursor, err := rt.Graph.Nodes(testContext, &cstxproto.NodeQuery{Filter: &cstxproto.NodeFilter{NodeTypes: []string{"domain"}}, Window: &cstxproto.QueryWindow{Order: cstxproto.SortOrder_SORT_ORDER_ID_ASC}})
	if err != nil {
		t.Fatalf("nodes: %v", err)
	}
	if affected := addDomain(t, rt, "example.com"); affected != 0 {
		t.Fatalf("no-op add affected=%d", affected)
	}

	page, err := cursor.Page(testContext, 1024, 1)
	if err != nil {
		t.Fatalf("cursor page: %v", err)
	}
	resultNodes := page.GetNodes().GetValues()
	ids := make([]string, len(resultNodes))
	for index, node := range resultNodes {
		ids[index] = node.GetId()
	}
	if err := cursor.Close(); err != nil {
		t.Fatalf("cursor close: %v", err)
	}
	if len(ids) != 2 || ids[0] != "domain:example.com" {
		t.Fatalf("unexpected ids: %v", ids)
	}

	stats, err := rt.Graph.Stats(testContext)
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.NodesByType["domain"] != 2 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestCursorInvalidation(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "example.com")

	cursor, err := rt.Graph.Nodes(testContext, &cstxproto.NodeQuery{})
	if err != nil {
		t.Fatalf("nodes: %v", err)
	}
	defer cursor.Close()
	if _, err := cursor.Page(testContext, 1, 1); err != nil {
		t.Fatalf("expected first page: %v", err)
	}
	addDomain(t, rt, "other.example.com")
	if _, err := cursor.Page(testContext, 1, 1); !IsCode(err, CodeCursorInvalidated) {
		t.Fatalf("expected CURSOR_INVALIDATED, got %v", err)
	}
}

func TestGraphRelationshipsNeighborsAndQuery(t *testing.T) {
	rt := openRuntime(t)

	addDomain(t, rt, "example.com")
	addDomain(t, rt, "www.example.com")
	affected, err := rt.Graph.AddRelationships(testContext, []*cstxproto.Relationship{usesRelationship("domain:www.example.com", "domain:example.com")})
	if err != nil {
		t.Fatalf("add edge: %v", err)
	}
	if affected != 1 {
		t.Fatalf("add edge affected=%d", affected)
	}
	if count, err := rt.Graph.RelationshipCount(testContext); err != nil || count != 1 {
		t.Fatalf("edge count: %v %v", count, err)
	}

	relationships, err := rt.Graph.Relationships(testContext, &cstxproto.RelationshipQuery{Filter: &cstxproto.RelationshipFilter{SourceId: stringPtr("domain:www.example.com")}})
	if err != nil {
		t.Fatalf("edges: %v", err)
	}
	defer relationships.Close()
	relationshipPage, err := relationships.Page(testContext, 10, 1)
	if err != nil {
		t.Fatalf("edge page: %v", err)
	}
	relationshipItems := relationshipPage.GetRelationships().GetValues()
	if len(relationshipItems) != 1 {
		t.Fatalf("decode relationship page: items=%v", relationshipItems)
	}
	if relationshipItems[0].GetTargetId() != "domain:example.com" {
		t.Fatalf("unexpected relationship: %+v", relationshipItems[0])
	}

	neighbors, err := rt.Graph.Neighbors(testContext, &cstxproto.NeighborQuery{NodeId: "domain:www.example.com", Direction: cstxproto.Direction_DIRECTION_OUT})
	if err != nil {
		t.Fatalf("neighbors: %v", err)
	}
	defer neighbors.Close()
	neighborPage, err := neighbors.Page(testContext, 10, 1)
	if err != nil {
		t.Fatalf("neighbor page: %v", err)
	}
	neighborItems := neighborPage.GetNodes().GetValues()
	if len(neighborItems) != 1 || neighborItems[0].GetId() != "domain:example.com" {
		t.Fatalf("unexpected neighbor: items=%v", neighborItems)
	}

	matches, err := rt.Graph.Query(testContext, &cstxproto.GraphQuery{Expression: "domain"})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer matches.Close()
	matchPage, err := matches.Page(testContext, 10, 1)
	if err != nil {
		t.Fatalf("query page: %v", err)
	}
	if len(matchPage.GetNodes().GetValues()) != 2 {
		t.Fatalf("expected 2 query matches, got %d", len(matchPage.GetNodes().GetValues()))
	}
}

func TestGraphAlgorithmsAndCommunityCursor(t *testing.T) {
	rt := openRuntime(t)
	for _, value := range []string{"a.example", "b.example", "c.example", "d.example"} {
		addDomain(t, rt, value)
	}
	_, err := rt.Graph.AddRelationships(testContext, []*cstxproto.Relationship{
		usesRelationship("domain:a.example", "domain:b.example"),
		usesRelationship("domain:b.example", "domain:c.example"),
	})
	if err != nil {
		t.Fatalf("add algorithm fixture edges: %v", err)
	}

	bfs, _, err := rt.Graph.Analyze(testContext, &cstxproto.Algorithm{Kind: &cstxproto.Algorithm_Bfs{Bfs: &cstxproto.BfsAlgorithm{SeedId: "domain:a.example", Depth: 2, Direction: cstxproto.Direction_DIRECTION_OUT}}}, nil)
	if err != nil {
		t.Fatalf("bfs: %v", err)
	}
	defer bfs.Close()
	if bfs.Kind() != CursorKindNodes {
		t.Fatalf("bfs kind=%q", bfs.Kind())
	}
	bfsPage, err := bfs.Page(testContext, 2, 1)
	if err != nil {
		t.Fatalf("bfs page: %v", err)
	}
	if bfsPage.Total == nil || *bfsPage.Total != 2 || bfsPage.GetHasNext() {
		t.Fatalf("unexpected bfs page metadata: %+v", bfsPage)
	}
	bfsNodes := bfsPage.GetNodes().GetValues()
	if len(bfsNodes) != 2 || bfsNodes[0].GetId() != "domain:b.example" {
		t.Fatalf("unexpected bfs rows: nodes=%v", bfsNodes)
	}

	components, _, err := rt.Graph.Analyze(testContext, &cstxproto.Algorithm{Kind: &cstxproto.Algorithm_Parameterless{Parameterless: cstxproto.ParameterlessAlgorithm_PARAMETERLESS_WEAK_COMPONENTS}}, nil)
	if err != nil {
		t.Fatalf("weak components: %v", err)
	}
	defer components.Close()
	componentPage, err := components.Page(testContext, 10, 1)
	if err != nil {
		t.Fatalf("component page: %v", err)
	}
	if components.Kind() != CursorKindComponents || len(componentPage.GetComponents().GetValues()) != 4 {
		t.Fatalf("unexpected component cursor: kind=%q page=%+v", components.Kind(), componentPage)
	}
	componentSummary := componentPage.GetComponent()
	if componentSummary.GetComponentCount() != 2 || componentSummary.GetProjection() != "undirected" {
		t.Fatalf("component summary=%+v", componentSummary)
	}

	_, isDAG, err := rt.Graph.Analyze(testContext, &cstxproto.Algorithm{Kind: &cstxproto.Algorithm_Parameterless{Parameterless: cstxproto.ParameterlessAlgorithm_PARAMETERLESS_IS_DAG}}, nil)
	if err != nil || isDAG == nil || !*isDAG {
		t.Fatalf("is dag=%v err=%v", isDAG, err)
	}
	order, _, err := rt.Graph.Analyze(testContext, &cstxproto.Algorithm{Kind: &cstxproto.Algorithm_Parameterless{Parameterless: cstxproto.ParameterlessAlgorithm_PARAMETERLESS_TOPOLOGICAL_ORDER}}, nil)
	if err != nil || order == nil {
		t.Fatalf("topological order result=%v err=%v", order, err)
	}
	defer order.Close()
	orderPage, err := order.Page(testContext, 10, 1)
	if err != nil {
		t.Fatalf("topological page: %v", err)
	}
	if len(orderPage.GetNodes().GetValues()) != 4 {
		t.Fatalf("topological row count=%d", len(orderPage.GetNodes().GetValues()))
	}

	core, _, err := rt.Graph.Analyze(testContext, &cstxproto.Algorithm{Kind: &cstxproto.Algorithm_Parameterless{Parameterless: cstxproto.ParameterlessAlgorithm_PARAMETERLESS_CORE_NUMBERS}}, nil)
	if err != nil {
		t.Fatalf("core numbers: %v", err)
	}
	defer core.Close()
	corePage, err := core.Page(testContext, 10, 1)
	if err != nil || len(corePage.GetScores().GetValues()) != 4 || core.Kind() != CursorKindNodeScores {
		t.Fatalf("core page=%+v kind=%q err=%v", corePage, core.Kind(), err)
	}
	coreRow := corePage.GetScores().GetValues()[0]
	if coreRow.GetMetric() != "core_number" || coreRow.GetNodeId() == "" {
		t.Fatalf("core row=%+v", coreRow)
	}

	communities, _, err := rt.Graph.Analyze(testContext, &cstxproto.Algorithm{Kind: &cstxproto.Algorithm_Leiden{Leiden: &cstxproto.LeidenAlgorithm{Resolution: 1, MinCommunitySize: 1}}}, nil)
	if err != nil {
		t.Fatalf("analyze leiden: %v", err)
	}
	defer communities.Close()
	if communities.Kind() != CursorKindCommunities {
		t.Fatalf("unexpected community cursor kind: %q", communities.Kind())
	}
	assignmentPage, err := communities.Page(testContext, 2, 1)
	if err != nil {
		t.Fatalf("community assignment page: %v", err)
	}
	communitySummary := assignmentPage.GetCommunity()
	if communitySummary.GetAlgorithm() != "leiden" || communitySummary.GetProjection() != "undirected" || communitySummary.GetTotalCommunities() == 0 {
		t.Fatalf("unexpected community summary: %+v", communitySummary)
	}
	if assignmentPage.Total == nil || *assignmentPage.Total != 4 || len(assignmentPage.GetCommunities().GetValues()) != 2 || !assignmentPage.GetHasNext() {
		t.Fatalf("unexpected assignment page: %+v", assignmentPage)
	}
	assignment := assignmentPage.GetCommunities().GetValues()[0]
	if assignment.GetNodeId() == "" {
		t.Fatalf("community assignment=%+v", assignment)
	}
}

func TestGraphQueryOptionsCrossFFIBoundary(t *testing.T) {
	rt := openRuntime(t)

	external := domainNode("external.example.com")
	internal := domainNode("internal.example.com")
	internal.Flags = []cstxproto.NodeFlag{cstxproto.NodeFlag_NODE_FLAG_INTERNAL}
	if affected, err := rt.Graph.AddNodes(testContext, []*cstxproto.Node{external, internal}); err != nil || affected != 2 {
		t.Fatalf("add query option nodes: affected=%d err=%v", affected, err)
	}

	collect := func(query *cstxproto.GraphQuery) []string {
		t.Helper()
		cursor, err := rt.Graph.Query(testContext, query)
		if err != nil {
			t.Fatalf("query with options %+v: %v", query, err)
		}
		defer cursor.Close()
		limit := 1024
		pageNumber := 1
		if query.Options != nil && query.Options.Window != nil && query.Options.Window.Limit != nil {
			limit = int(*query.Options.Window.Limit)
			pageNumber = int(query.Options.Window.Page)
		}
		page, err := cursor.Page(testContext, limit, pageNumber)
		if err != nil {
			t.Fatalf("query cursor with options %+v: %v", query, err)
		}
		nodes := page.GetNodes().GetValues()
		ids := make([]string, len(nodes))
		for index, node := range nodes {
			ids[index] = node.GetId()
		}
		return ids
	}

	one := uint64(1)
	if ids := collect(&cstxproto.GraphQuery{Expression: "domain", Options: &cstxproto.QueryOptions{Window: &cstxproto.QueryWindow{Limit: &one, Page: 2}}}); !reflect.DeepEqual(ids, []string{internal.GetId()}) {
		t.Fatalf("second query page returned %v", ids)
	}
	if ids := collect(&cstxproto.GraphQuery{Expression: "domain", Options: &cstxproto.QueryOptions{ResultFilter: &cstxproto.NodeFilter{FlagsAll: []cstxproto.NodeFlag{cstxproto.NodeFlag_NODE_FLAG_INTERNAL}}}}); !reflect.DeepEqual(ids, []string{internal.GetId()}) {
		t.Fatalf("include-mask query returned %v", ids)
	}
	if ids := collect(&cstxproto.GraphQuery{Expression: "domain", Options: &cstxproto.QueryOptions{ResultFilter: &cstxproto.NodeFilter{FlagsNone: []cstxproto.NodeFlag{cstxproto.NodeFlag_NODE_FLAG_INTERNAL}}}}); !reflect.DeepEqual(ids, []string{external.GetId()}) {
		t.Fatalf("exclude-mask query returned %v", ids)
	}
}

func TestRepositoryRoundTrip(t *testing.T) {
	rt := openRuntime(t)

	addDomain(t, rt, "example.com")
	commit, err := rt.Repo.Commit(testContext, "initial", "main", nil, &structpb.Struct{Fields: map[string]*structpb.Value{"origin": structpb.NewStringValue("test")}})
	if err != nil {
		t.Fatalf("commit: %v", err)
	}
	if commit.Id == "" {
		t.Fatalf("unexpected commit: %+v", commit)
	}
	if commit.Metadata.GetFields()["origin"].GetStringValue() != "test" {
		t.Fatalf("commit metadata: %#v", commit.Metadata)
	}

	head, err := rt.Repo.Head(testContext, "main")
	if err != nil || head == nil || *head != commit.Id {
		t.Fatalf("head: %v %v", head, err)
	}
	resolved, err := rt.Repo.Resolve(testContext, "main")
	if err != nil || resolved != commit.Id {
		t.Fatalf("resolve: %s %v", resolved, err)
	}
	if _, err := rt.Repo.Branch(testContext, "initial", "main"); err != nil {
		t.Fatalf("branch: %v", err)
	}

	addDomain(t, rt, "www.example.com")
	second, err := rt.Repo.Commit(testContext, "second", "main", &commit.Id, nil)
	if err != nil {
		t.Fatalf("second commit: %v", err)
	}
	diff, err := rt.Repo.Diff(testContext, commit.Id, second.Id, nil, cstxproto.DiffDetail_DIFF_DETAIL_ENTITIES)
	if err != nil || len(diff.Added.NodeIds) != 1 || diff.Stats.AddedNodes != 1 {
		t.Fatalf("diff: %+v %v", diff, err)
	}
	counted, err := rt.Repo.Diff(testContext, commit.Id, second.Id, nil, cstxproto.DiffDetail_DIFF_DETAIL_COUNTS)
	if err != nil || counted.Stats.AddedNodes != 1 || len(counted.Added.NodeIds) != 0 {
		t.Fatalf("counted diff: %+v %v", counted, err)
	}
	log, err := rt.Repo.Log(testContext, "main", 10)
	if err != nil || len(log.Commits) != 2 {
		t.Fatalf("log: %+v %v", log, err)
	}
	history, err := rt.Repo.History(testContext, "domain:www.example.com", "main", nil)
	if err != nil || len(history.Changes) != 1 {
		t.Fatalf("history: %+v %v", history, err)
	}
	if stat, err := rt.Repo.Stat(testContext, "main", 0, 0); err != nil || stat.NodesByType["domain"] != 2 {
		t.Fatalf("stat: %+v %v", stat, err)
	}
	if _, err := rt.Repo.Delta(testContext, "main", nil, nil); err != nil {
		t.Fatalf("delta: %v", err)
	}
	if _, err := rt.Repo.Checkout(testContext, "initial", true); err != nil {
		t.Fatalf("checkout initial: %v", err)
	}
	if count, _ := rt.Graph.NodeCount(testContext); count != 1 {
		t.Fatalf("initial node count: %d", count)
	}
	if _, err := rt.Repo.Checkout(testContext, "main", true); err != nil {
		t.Fatalf("checkout main: %v", err)
	}
}

func TestLastChange(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "example.com")

	change, err := rt.LastChange(testContext)
	if err != nil {
		t.Fatalf("last change: %v", err)
	}
	if len(change.AddedNodeIds) != 1 || Affected(change) != 1 {
		t.Fatalf("unexpected change set: %+v", change)
	}
}

func TestServicesAndCursorsHonorCanceledContext(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "example.com")
	canceled, cancel := context.WithCancel(testContext)
	cancel()
	if _, err := rt.Graph.NodeCount(canceled); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled service call, got %v", err)
	}

	cursor, err := rt.Graph.Nodes(testContext, &cstxproto.NodeQuery{})
	if err != nil {
		t.Fatalf("nodes: %v", err)
	}
	defer cursor.Close()
	if _, err := cursor.Page(canceled, 1, 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled cursor error, got %v", err)
	}
}

func TestConformanceFixtureMatchesGoContract(t *testing.T) {
	// Rust and Python run this same file and assert the same ids, counts and
	// query result. That only means something if all three use it as written:
	// this used to decode the fixture's schema and then overwrite every field
	// of it with easm's `domain`, so what it actually checked was easm.
	fixture := loadConformanceFixture(t)

	rt, err := Open(testContext, &cstxproto.RuntimeConfig{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer rt.Close()
	contract := newExtensionBuilder("conformance", "1.0").
		Schema(string(fixture.Document)).Build()
	if err := rt.Extensions.Register(testContext, contract); err != nil {
		t.Fatalf("register: %v", err)
	}

	nodes := make([]*cstxproto.Node, 0, len(fixture.Nodes))
	for _, item := range fixture.Nodes {
		entity := &cstxproto.EntityValue{NodeType: item.Type}
		for _, name := range sortedKeys(item.Model) {
			entity.Fields = append(entity.Fields, &cstxproto.EntityField{
				Name:  name,
				Value: &cstxproto.EntityField_Text{Text: item.Model[name]},
			})
		}
		id := item.ID
		nodes = append(nodes, &cstxproto.Node{
			Id: &id, Sources: item.Sources, Value: entity,
		})
	}
	if _, err := rt.Graph.AddNodes(testContext, nodes); err != nil {
		t.Fatalf("add nodes: %v", err)
	}

	var document struct {
		Relations map[string]struct {
			Message string `json:"message"`
		} `json:"relations"`
	}
	if err := json.Unmarshal(fixture.Document, &document); err != nil {
		t.Fatalf("document: %v", err)
	}
	edges := make([]*cstxproto.Relationship, 0, len(fixture.Relationships))
	for _, item := range fixture.Relationships {
		id := item.ID
		edges = append(edges, &cstxproto.Relationship{
			Id: &id, SourceId: item.SourceID, TargetId: item.TargetID,
			Sources: item.Sources,
			// A relation type is a field-less marker: the document names the
			// message and the payload is empty.
			Relation: &anypb.Any{
				TypeUrl: "type.googleapis.com/" + document.Relations[item.RelationType].Message,
			},
		})
	}
	if _, err := rt.Graph.AddRelationships(testContext, edges); err != nil {
		t.Fatalf("add relationships: %v", err)
	}

	if count, err := rt.Graph.NodeCount(testContext); err != nil || count != fixture.Expected.NodeCount {
		t.Fatalf("node count: got %d err=%v", count, err)
	}
	if count, err := rt.Graph.RelationshipCount(testContext); err != nil || count != fixture.Expected.RelationshipCount {
		t.Fatalf("relationship count: got %d err=%v", count, err)
	}
	cursor, err := rt.Graph.Query(testContext, &cstxproto.GraphQuery{Expression: fixture.Query})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer cursor.Close()
	page, err := cursor.Page(testContext, 100, 1)
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	ids := make([]string, 0, len(page.GetNodes().GetValues()))
	for _, node := range page.GetNodes().GetValues() {
		ids = append(ids, node.GetId())
	}
	if !reflect.DeepEqual(ids, fixture.Expected.NodeIDs) {
		t.Fatalf("query ids = %v; want %v", ids, fixture.Expected.NodeIDs)
	}
}

func stringSliceMismatch(left, right []string) bool {
	if len(left) != len(right) {
		return true
	}
	for i := range left {
		if left[i] != right[i] {
			return true
		}
	}
	return false
}
