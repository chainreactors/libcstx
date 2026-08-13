package cstx

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"testing"
)

var testContext = context.Background()

//go:embed testdata/conformance.json
var conformanceFixture []byte

var domainSchema = map[string]any{"properties": map[string]any{"domain": map[string]any{"type": "string"}}}

func openRuntime(t *testing.T) *CSTX {
	t.Helper()
	rt, err := Open(testContext, Config{ProjectID: "sdk-go-test"})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.Close(); err != nil {
			t.Fatalf("close: %v", err)
		}
	})
	if err := rt.Schemas.Register(testContext, "domain", domainSchema, "domain"); err != nil {
		t.Fatalf("register schema: %v", err)
	}
	return rt
}

func domainNode(value string) Node {
	return Node{
		ID:      "domain:" + value,
		Type:    "domain",
		Value:   value,
		Model:   map[string]any{"domain": value, "cstx_flags": 0},
		Sources: []string{"test"},
		Extras:  map[string]any{},
	}
}

func relatedEdge(source, target string) Edge {
	return Edge{
		ID:           "relationship:" + source + ":related:" + target,
		SourceID:     source,
		TargetID:     target,
		RelationType: "related",
		Sources:      []string{"test"},
		Attrs:        map[string]any{},
	}
}

func addDomain(t *testing.T, rt *CSTX, value string) uint64 {
	t.Helper()
	affected, err := rt.Graph.AddNodes(testContext, []Node{domainNode(value)})
	if err != nil {
		t.Fatalf("add node %s: %v", value, err)
	}
	return affected
}

func TestSchemas(t *testing.T) {
	rt := openRuntime(t)
	valueField := "domain"
	contract := SchemaContract{
		Format: "cstx.schema",
		Plugins: map[string]PluginSchemaContract{
			"sdk-test": {
				Version: "1",
				SCO: map[string]SCOSchemaContract{
					"domain": {
						Schema:     domainSchema,
						ValueField: &valueField,
						Metadata:   map[string]any{},
					},
				},
				SRO:     map[string]SROSchemaContract{},
				Parsers: map[string]ParserSchemaContract{},
			},
		},
	}
	if err := rt.Schemas.Import(testContext, contract); err != nil {
		t.Fatalf("import schema contract: %v", err)
	}
	exported, err := rt.Schemas.Export(testContext)
	if err != nil || exported.Format != "cstx.schema" || exported.Plugins["sdk-test"].Version != "1" {
		t.Fatalf("export schema contract: %+v err=%v", exported, err)
	}
	if err := rt.Schemas.RegisterJoinRule(testContext, JoinRuleSpec{
		LeftType: "domain", RightType: "domain", Relation: "related",
		LeftKey: "domain", RightKey: "domain",
	}); err != nil {
		t.Fatalf("register join rule: %v", err)
	}

	contains, err := rt.Schemas.Contains(testContext, "domain")
	if err != nil || !contains {
		t.Fatalf("contains: %v %v", contains, err)
	}
	schema, err := rt.Schemas.Get(testContext, "domain")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if schema["node_type"] != "domain" || schema["value_field"] != "domain" {
		t.Fatalf("unexpected schema: %v", schema)
	}
	if _, ok := schema["schema"].(map[string]any)["properties"]; !ok {
		t.Fatalf("unexpected schema body: %v", schema)
	}
	list, err := rt.Schemas.List(testContext)
	if err != nil || len(list) == 0 {
		t.Fatalf("list: %v %v", list, err)
	}
	plugins, err := rt.Schemas.AvailablePlugins(testContext)
	if err != nil {
		t.Fatalf("available plugins: %v", err)
	}
	if !reflect.DeepEqual(plugins, []string{"easm"}) {
		t.Fatalf("available plugins: %v", plugins)
	}
	artifacts, err := rt.Schemas.PluginArtifacts(testContext, "easm")
	if err != nil || !slices.Contains(artifacts, "gogo") {
		t.Fatalf("easm artifacts: %v err=%v", artifacts, err)
	}
	concepts, err := rt.Schemas.AnchorConcepts(testContext)
	if err != nil || len(concepts) == 0 || concepts[0].Name == "" {
		t.Fatalf("anchor concepts: %+v err=%v", concepts, err)
	}
	if err := rt.Schemas.LoadPlugin(testContext, "easm"); err != nil {
		t.Fatalf("load easm plugin: %v", err)
	}
	hasGogo, err := rt.Schemas.HasNativeArtifact(testContext, "gogo")
	if err != nil || !hasGogo {
		t.Fatalf("has native gogo artifact: %v err=%v", hasGogo, err)
	}
	gogo := []byte(`{"ip":"192.0.2.1","port":"80","protocol":"tcp","status":"200"}` + "\n")
	if affected, err := rt.Graph.Ingest(testContext, "gogo", gogo); err != nil || affected == 0 {
		t.Fatalf("ingest gogo: affected=%d err=%v", affected, err)
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
	if node.ID != "domain:example.com" || node.Type != "domain" {
		t.Fatalf("unexpected node: %+v", node)
	}
	if node.Model["domain"] != "example.com" {
		t.Fatalf("unexpected model: %+v", node.Model)
	}

	if _, err := rt.Graph.Node(testContext, "domain:missing.example"); !IsCode(err, CodeNotFound) {
		t.Fatalf("expected NOT_FOUND, got %v", err)
	}

	// Re-adding identical content is a no-op: zero affected and live cursors
	// stay valid.
	cursor, err := rt.Graph.Nodes(testContext, NodeFilter{Types: []string{"domain"}}, CollectionOptions{Order: OrderIDAsc})
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
	nodes, err := page.Nodes()
	if err != nil {
		t.Fatalf("decode nodes: %v", err)
	}
	ids := make([]string, len(nodes))
	for index, node := range nodes {
		ids[index] = node.ID
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
	if stats.Nodes["domain"] != 2 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestCursorInvalidation(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "example.com")

	cursor, err := rt.Graph.Nodes(testContext, NodeFilter{}, CollectionOptions{})
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

func TestGraphEdgesNeighborsAndQuery(t *testing.T) {
	rt := openRuntime(t)

	addDomain(t, rt, "example.com")
	addDomain(t, rt, "www.example.com")
	affected, err := rt.Graph.AddEdges(testContext, []Edge{relatedEdge("domain:www.example.com", "domain:example.com")})
	if err != nil {
		t.Fatalf("add edge: %v", err)
	}
	if affected != 1 {
		t.Fatalf("add edge affected=%d", affected)
	}
	if count, err := rt.Graph.EdgeCount(testContext); err != nil || count != 1 {
		t.Fatalf("edge count: %v %v", count, err)
	}

	edges, err := rt.Graph.Edges(testContext, EdgeFilter{SourceID: "domain:www.example.com"}, CollectionOptions{})
	if err != nil {
		t.Fatalf("edges: %v", err)
	}
	defer edges.Close()
	edgePage, err := edges.Page(testContext, 10, 1)
	if err != nil {
		t.Fatalf("edge page: %v", err)
	}
	edgeItems, err := edgePage.Edges()
	if err != nil || len(edgeItems) != 1 {
		t.Fatalf("decode edge page: items=%v err=%v", edgeItems, err)
	}
	if edgeItems[0].RelationType != "related" || edgeItems[0].TargetID != "domain:example.com" {
		t.Fatalf("unexpected edge: %+v", edgeItems[0])
	}

	neighbors, err := rt.Graph.Neighbors(testContext, "domain:www.example.com", "out", CollectionOptions{})
	if err != nil {
		t.Fatalf("neighbors: %v", err)
	}
	defer neighbors.Close()
	neighborPage, err := neighbors.Page(testContext, 10, 1)
	if err != nil {
		t.Fatalf("neighbor page: %v", err)
	}
	neighborItems, err := neighborPage.Nodes()
	if err != nil || len(neighborItems) != 1 || neighborItems[0].ID != "domain:example.com" {
		t.Fatalf("unexpected neighbor: items=%v err=%v", neighborItems, err)
	}

	matches, err := rt.Graph.Query(testContext, "domain", QueryOptions{})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer matches.Close()
	matchPage, err := matches.Page(testContext, 10, 1)
	if err != nil {
		t.Fatalf("query page: %v", err)
	}
	if len(matchPage.Items) != 2 {
		t.Fatalf("expected 2 query matches, got %d", len(matchPage.Items))
	}
}

func TestGraphAlgorithmsAndCommunityCursor(t *testing.T) {
	rt := openRuntime(t)
	for _, value := range []string{"a.example", "b.example", "c.example", "d.example"} {
		addDomain(t, rt, value)
	}
	_, err := rt.Graph.AddEdges(testContext, []Edge{
		relatedEdge("domain:a.example", "domain:b.example"),
		relatedEdge("domain:b.example", "domain:c.example"),
	})
	if err != nil {
		t.Fatalf("add algorithm fixture edges: %v", err)
	}

	bfsResult, err := rt.Graph.Analyze(testContext, map[string]any{
		"name": "bfs", "seed_id": "domain:a.example", "depth": 2, "direction": "out",
	})
	if err != nil {
		t.Fatalf("bfs: %v", err)
	}
	bfs := bfsResult.(*GraphCursor)
	defer bfs.Close()
	if bfs.Kind() != CursorKindNodes {
		t.Fatalf("bfs kind=%q", bfs.Kind())
	}
	bfsPage, err := bfs.Page(testContext, 2, 1)
	if err != nil {
		t.Fatalf("bfs page: %v", err)
	}
	if bfsPage.Total == nil || *bfsPage.Total != 2 || bfsPage.HasNext {
		t.Fatalf("unexpected bfs page metadata: %+v", bfsPage)
	}
	bfsNodes, err := bfsPage.Nodes()
	if err != nil || len(bfsNodes) != 2 || bfsNodes[0].ID != "domain:b.example" {
		t.Fatalf("unexpected bfs rows: nodes=%v err=%v", bfsNodes, err)
	}

	componentsResult, err := rt.Graph.Analyze(testContext, map[string]any{"name": "weak_components"})
	if err != nil {
		t.Fatalf("weak components: %v", err)
	}
	components := componentsResult.(*GraphCursor)
	defer components.Close()
	componentPage, err := components.Page(testContext, 10, 1)
	if err != nil {
		t.Fatalf("component page: %v", err)
	}
	if components.Kind() != CursorKindComponents || len(componentPage.Items) != 4 {
		t.Fatalf("unexpected component cursor: kind=%q page=%+v", components.Kind(), componentPage)
	}
	var componentSummary struct {
		ComponentCount uint64 `json:"component_count"`
		Projection     string `json:"projection"`
	}
	if err := json.Unmarshal(componentPage.Summary, &componentSummary); err != nil || componentSummary.ComponentCount != 2 || componentSummary.Projection != "undirected" {
		t.Fatalf("component summary=%+v err=%v", componentSummary, err)
	}

	isDAGResult, err := rt.Graph.Analyze(testContext, map[string]any{"name": "is_dag"})
	if err != nil || !isDAGResult.(bool) {
		t.Fatalf("is dag=%v err=%v", isDAGResult, err)
	}
	orderResult, err := rt.Graph.Analyze(testContext, map[string]any{"name": "topological_order"})
	if err != nil || orderResult == nil {
		t.Fatalf("topological order result=%v err=%v", orderResult, err)
	}
	order := orderResult.(*GraphCursor)
	defer order.Close()
	orderPage, err := order.Page(testContext, 10, 1)
	if err != nil {
		t.Fatalf("topological page: %v", err)
	}
	if len(orderPage.Items) != 4 {
		t.Fatalf("topological row count=%d", len(orderPage.Items))
	}

	coreResult, err := rt.Graph.Analyze(testContext, map[string]any{"name": "core_numbers"})
	if err != nil {
		t.Fatalf("core numbers: %v", err)
	}
	core := coreResult.(*GraphCursor)
	defer core.Close()
	corePage, err := core.Page(testContext, 10, 1)
	if err != nil || len(corePage.Items) != 4 || core.Kind() != CursorKindNodeScores {
		t.Fatalf("core page=%+v kind=%q err=%v", corePage, core.Kind(), err)
	}
	var coreRow struct {
		NodeID string  `json:"node_id"`
		Metric string  `json:"metric"`
		Score  float64 `json:"score"`
	}
	if err := json.Unmarshal(corePage.Items[0], &coreRow); err != nil || coreRow.Metric != "core_number" || coreRow.NodeID == "" {
		t.Fatalf("core row=%+v err=%v", coreRow, err)
	}

	communityResult, err := rt.Graph.Analyze(testContext, map[string]any{
		"name": "leiden", "resolution": 1.0, "min_community_size": 1,
	})
	if err != nil {
		t.Fatalf("analyze leiden: %v", err)
	}
	communities := communityResult.(*GraphCursor)
	defer communities.Close()
	if communities.Kind() != CursorKindCommunities {
		t.Fatalf("unexpected community cursor kind: %q", communities.Kind())
	}
	assignmentPage, err := communities.Page(testContext, 2, 1)
	if err != nil {
		t.Fatalf("community assignment page: %v", err)
	}
	var communitySummary struct {
		Algorithm        string `json:"algorithm"`
		Projection       string `json:"projection"`
		TotalCommunities uint64 `json:"total_communities"`
	}
	if err := json.Unmarshal(assignmentPage.Summary, &communitySummary); err != nil || communitySummary.Algorithm != "leiden" || communitySummary.Projection != "undirected" || communitySummary.TotalCommunities == 0 {
		t.Fatalf("unexpected community summary: %+v err=%v", communitySummary, err)
	}
	if assignmentPage.Total == nil || *assignmentPage.Total != 4 || len(assignmentPage.Items) != 2 || !assignmentPage.HasNext {
		t.Fatalf("unexpected assignment page: %+v", assignmentPage)
	}
	var assignment struct {
		NodeID    string `json:"node_id"`
		Community uint32 `json:"community"`
	}
	if err := json.Unmarshal(assignmentPage.Items[0], &assignment); err != nil || assignment.NodeID == "" {
		t.Fatalf("community assignment=%+v err=%v", assignment, err)
	}
}

func TestGraphQueryOptionsCrossFFIBoundary(t *testing.T) {
	rt := openRuntime(t)

	external := domainNode("external.example.com")
	internal := domainNode("internal.example.com")
	internal.Model["cstx_flags"] = FlagInternal
	if affected, err := rt.Graph.AddNodes(testContext, []Node{external, internal}); err != nil || affected != 2 {
		t.Fatalf("add query option nodes: affected=%d err=%v", affected, err)
	}

	collect := func(options QueryOptions) []string {
		t.Helper()
		cursor, err := rt.Graph.Query(testContext, "domain", options)
		if err != nil {
			t.Fatalf("query with options %+v: %v", options, err)
		}
		defer cursor.Close()
		limit := 1024
		pageNumber := 1
		if options.Collection.Limit != nil {
			limit = *options.Collection.Limit
			pageNumber = options.Collection.Page
		}
		page, err := cursor.Page(testContext, limit, pageNumber)
		if err != nil {
			t.Fatalf("query cursor with options %+v: %v", options, err)
		}
		nodes, err := page.Nodes()
		if err != nil {
			t.Fatalf("decode query page: %v", err)
		}
		ids := make([]string, len(nodes))
		for index, node := range nodes {
			ids[index] = node.ID
		}
		return ids
	}

	one := 1
	if ids := collect(QueryOptions{Collection: CollectionOptions{Limit: &one, Page: 2}}); !reflect.DeepEqual(ids, []string{internal.ID}) {
		t.Fatalf("second query page returned %v", ids)
	}
	if ids := collect(QueryOptions{IncludeMask: FlagInternal}); !reflect.DeepEqual(ids, []string{internal.ID}) {
		t.Fatalf("include-mask query returned %v", ids)
	}
	if ids := collect(QueryOptions{ExcludeMask: FlagInternal}); !reflect.DeepEqual(ids, []string{external.ID}) {
		t.Fatalf("exclude-mask query returned %v", ids)
	}
}

func TestRepositoryRoundTrip(t *testing.T) {
	rt := openRuntime(t)

	addDomain(t, rt, "example.com")
	commit, err := rt.Repo.Commit(testContext, "initial", "main", nil, map[string]any{"origin": "test"})
	if err != nil {
		t.Fatalf("commit: %v", err)
	}
	if commit.ID == "" {
		t.Fatalf("unexpected commit: %+v", commit)
	}

	head, err := rt.Repo.Head(testContext, "main")
	if err != nil || head == nil || *head != commit.ID {
		t.Fatalf("head: %v %v", head, err)
	}
	resolved, err := rt.Repo.Resolve(testContext, "main")
	if err != nil || resolved != commit.ID {
		t.Fatalf("resolve: %s %v", resolved, err)
	}
	if _, err := rt.Repo.Branch(testContext, "initial", "main"); err != nil {
		t.Fatalf("branch: %v", err)
	}

	addDomain(t, rt, "www.example.com")
	second, err := rt.Repo.Commit(testContext, "second", "main", &commit.ID, nil)
	if err != nil {
		t.Fatalf("second commit: %v", err)
	}
	diff, err := rt.Repo.Diff(testContext, commit.ID, second.ID, nil)
	if err != nil || len(diff.Added["domain"]) != 1 {
		t.Fatalf("diff: %+v %v", diff, err)
	}
	diffStat, err := rt.Repo.DiffStat(testContext, commit.ID, second.ID)
	if err != nil || diffStat.AddedNodes != 1 {
		t.Fatalf("diff stat: %+v %v", diffStat, err)
	}
	log, err := rt.Repo.Log(testContext, "main", 10)
	if err != nil || len(log) != 2 {
		t.Fatalf("log: %+v %v", log, err)
	}
	history, err := rt.Repo.History(testContext, "domain:www.example.com", "main", nil)
	if err != nil || len(history.Entries) != 1 {
		t.Fatalf("history: %+v %v", history, err)
	}
	if stat, err := rt.Repo.Stat(testContext, "main", 0, 0); err != nil || stat.Nodes["domain"] != 2 {
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
	if len(change.AddedNodeIDs) != 1 || change.Affected() != 1 {
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

	cursor, err := rt.Graph.Nodes(testContext, NodeFilter{}, CollectionOptions{})
	if err != nil {
		t.Fatalf("nodes: %v", err)
	}
	defer cursor.Close()
	if _, err := cursor.Page(canceled, 1, 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled cursor error, got %v", err)
	}
}

func TestConformanceFixtureMatchesGoContract(t *testing.T) {
	var fixture struct {
		Schema struct {
			NodeType   string         `json:"node_type"`
			JSONSchema map[string]any `json:"json_schema"`
			ValueField string         `json:"value_field"`
		} `json:"schema"`
		Nodes    []Node `json:"nodes"`
		Edges    []Edge `json:"edges"`
		Query    string `json:"query"`
		Expected struct {
			NodeIDs   []string `json:"node_ids"`
			NodeCount uint64   `json:"node_count"`
			EdgeCount uint64   `json:"edge_count"`
		} `json:"expected"`
	}
	if err := json.Unmarshal(conformanceFixture, &fixture); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}

	rt, err := Open(testContext, Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer rt.Close()
	if err := rt.Schemas.Register(
		testContext,
		fixture.Schema.NodeType,
		fixture.Schema.JSONSchema,
		fixture.Schema.ValueField,
	); err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, err := rt.Graph.AddNodes(testContext, fixture.Nodes); err != nil {
		t.Fatalf("add nodes: %v", err)
	}
	if _, err := rt.Graph.AddEdges(testContext, fixture.Edges); err != nil {
		t.Fatalf("add edges: %v", err)
	}
	if count, err := rt.Graph.NodeCount(testContext); err != nil || count != fixture.Expected.NodeCount {
		t.Fatalf("node count: got %d err=%v", count, err)
	}
	if count, err := rt.Graph.EdgeCount(testContext); err != nil || count != fixture.Expected.EdgeCount {
		t.Fatalf("edge count: got %d err=%v", count, err)
	}
	cursor, err := rt.Graph.Query(testContext, fixture.Query, QueryOptions{})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer cursor.Close()
	page, err := cursor.Page(testContext, 1024, 1)
	if err != nil {
		t.Fatalf("page query: %v", err)
	}
	nodes, err := page.Nodes()
	if err != nil {
		t.Fatalf("decode query: %v", err)
	}
	ids := make([]string, len(nodes))
	for index, node := range nodes {
		ids[index] = node.ID
	}
	if stringSliceMismatch(ids, fixture.Expected.NodeIDs) {
		t.Fatalf("query IDs: got %v want %v", ids, fixture.Expected.NodeIDs)
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
