//go:build cstx_native

package cstx

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

var testContext = context.Background()

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
	if _, err := rt.Schemas.AvailablePlugins(testContext); err != nil {
		t.Fatalf("available plugins: %v", err)
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

	var ids []string
	for cursor.Next(testContext) {
		ids = append(ids, cursor.Node().ID)
	}
	if err := cursor.Err(); err != nil {
		t.Fatalf("cursor: %v", err)
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
	if !cursor.Next(testContext) {
		t.Fatalf("expected first node, err=%v", cursor.Err())
	}
	addDomain(t, rt, "other.example.com")
	if cursor.Next(testContext) {
		t.Fatal("expected cursor to stop after mutation")
	}
	if !IsCode(cursor.Err(), CodeCursorInvalidated) {
		t.Fatalf("expected CURSOR_INVALIDATED, got %v", cursor.Err())
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
	if !edges.Next(testContext) {
		t.Fatalf("expected one edge, err=%v", edges.Err())
	}
	if edges.Edge().RelationType != "related" || edges.Edge().TargetID != "domain:example.com" {
		t.Fatalf("unexpected edge: %+v", edges.Edge())
	}
	if edges.Next(testContext) {
		t.Fatal("expected exactly one edge")
	}

	neighbors, err := rt.Graph.Neighbors(testContext, "domain:www.example.com", "out", CollectionOptions{})
	if err != nil {
		t.Fatalf("neighbors: %v", err)
	}
	defer neighbors.Close()
	if !neighbors.Next(testContext) || neighbors.Node().ID != "domain:example.com" {
		t.Fatalf("unexpected neighbor, err=%v", neighbors.Err())
	}

	matches, err := rt.Graph.Query(testContext, "domain", QueryOptions{})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer matches.Close()
	count := 0
	for matches.Next(testContext) {
		count++
	}
	if err := matches.Err(); err != nil {
		t.Fatalf("query cursor: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 query matches, got %d", count)
	}
}

func TestRepositoryRoundTrip(t *testing.T) {
	rt := openRuntime(t)

	addDomain(t, rt, "example.com")
	commit, err := rt.Repo.Commit(testContext, "main", "initial", map[string]any{"origin": "test"})
	if err != nil {
		t.Fatalf("commit: %v", err)
	}
	if commit.ID == "" || commit.Tree == "" {
		t.Fatalf("unexpected commit: %+v", commit)
	}

	head, err := rt.Repo.Head(testContext, "main")
	if err != nil || head == nil || *head != commit.ID {
		t.Fatalf("head: %v %v", head, err)
	}
	refs, err := rt.Repo.Refs(testContext)
	if err != nil || len(refs) != 1 || refs[0].Name != "main" {
		t.Fatalf("refs: %+v %v", refs, err)
	}

	addDomain(t, rt, "www.example.com")
	if _, err := rt.Repo.Commit(testContext, "main", "second", nil); err != nil {
		t.Fatalf("second commit: %v", err)
	}

	snapshot, err := rt.Repo.DumpJSON(testContext)
	if err != nil || len(snapshot) == 0 {
		t.Fatalf("dump json: %v", err)
	}
	fingerprint, err := rt.Repo.SnapshotFingerprint(testContext)
	if err != nil || fingerprint == "" {
		t.Fatalf("fingerprint: %v", err)
	}

	compact, err := rt.Repo.Dump(testContext, "zstd1")
	if err != nil || len(compact) == 0 {
		t.Fatalf("dump: %v", err)
	}

	rt2 := openRuntime(t)
	consumed, err := rt2.Repo.LoadJSON(testContext, snapshot)
	if err != nil {
		t.Fatalf("load json: %v", err)
	}
	if consumed != uint64(len(snapshot)) {
		t.Fatalf("consumed %d != %d", consumed, len(snapshot))
	}
	if count, _ := rt2.Graph.NodeCount(testContext); count != 2 {
		t.Fatalf("loaded node count: %d", count)
	}
	if fp, _ := rt2.Repo.SnapshotFingerprint(testContext); fp != fingerprint {
		t.Fatalf("fingerprint mismatch after load: %s != %s", fp, fingerprint)
	}

	rt3 := openRuntime(t)
	if _, err := rt3.Repo.Load(testContext, compact, "zstd1"); err != nil {
		t.Fatalf("load: %v", err)
	}
	if count, _ := rt3.Graph.NodeCount(testContext); count != 2 {
		t.Fatalf("loaded compact node count: %d", count)
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
	if cursor.Next(canceled) {
		t.Fatal("canceled cursor unexpectedly yielded a node")
	}
	if !errors.Is(cursor.Err(), context.Canceled) {
		t.Fatalf("expected canceled cursor error, got %v", cursor.Err())
	}
}

func TestCanonicalV03FixtureMatchesGoContract(t *testing.T) {
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
	payload, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "v03_conformance.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if err := json.Unmarshal(payload, &fixture); err != nil {
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
	var ids []string
	for cursor.Next(testContext) {
		ids = append(ids, cursor.Node().ID)
	}
	if err := cursor.Err(); err != nil {
		t.Fatalf("iterate query: %v", err)
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
