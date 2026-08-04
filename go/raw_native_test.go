package cstx

import (
	"encoding/json"
	"testing"
)

func TestRawRAGCanDisableLexicalRecall(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "example.com")
	payload := json.RawMessage(`{
		"text":"example domain",
		"limit":5,
		"filters":{"node_types":[],"relation_types":[],"exclude_flags":0,"include_flags":0},
		"policy":{"rrf_k":60,"candidate_multiplier":4,"damping":0.85,"propagation_iterations":20,"max_path_depth":4,"epsilon":0.000001,"communities":true,"use_lexical":false},
		"context_budget":null
	}`)
	retrieval, err := rt.Raw.RAGRetrieve(testContext, payload)
	if err != nil {
		t.Fatal(err)
	}
	defer retrieval.Close()
	resultJSON, err := retrieval.Complete(testContext, json.RawMessage(`[]`))
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		Nodes      []json.RawMessage `json:"nodes"`
		Edges      []json.RawMessage `json:"edges"`
		Extensions []string          `json:"extensions"`
	}
	if err := json.Unmarshal(resultJSON, &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Nodes) != 0 || len(result.Edges) != 0 || len(result.Extensions) != 0 {
		t.Fatalf("vector-only retrieval without external batches must be empty: %s", resultJSON)
	}
}

func TestRawLinkAndTypedSubgraph(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "example.com")
	addDomain(t, rt, "www.example.com")
	if _, err := rt.Graph.AddEdges(testContext, []Edge{relatedEdge("domain:www.example.com", "domain:example.com")}); err != nil {
		t.Fatal(err)
	}
	if _, err := rt.Raw.Link(testContext, json.RawMessage(`["domain:example.com"]`), "test"); err != nil {
		t.Fatal(err)
	}
	derived, err := rt.Graph.Subgraph(testContext, []string{"domain:www.example.com"}, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer derived.Close()
	if count, err := derived.Graph.NodeCount(testContext); err != nil || count != 2 {
		t.Fatalf("derived node count=%d err=%v", count, err)
	}
}

func TestRepositoryIndexGraphKeepsObjectsRaw(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "example.com")
	index, err := rt.Repo.IndexGraph(testContext)
	if err != nil {
		t.Fatal(err)
	}
	hash := index.Entries["domain"]["domain:example.com"]
	if hash == "" {
		t.Fatalf("missing domain entry: %#v", index.Entries)
	}
	object := index.Objects[hash]
	if len(object) == 0 || !json.Valid(object) {
		t.Fatalf("invalid raw CAS object: %q", object)
	}
}
