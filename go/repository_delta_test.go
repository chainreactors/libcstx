package cstx

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestRepositoryDeltaAPI(t *testing.T) {
	ctx := context.Background()
	rt, err := Open(ctx, Config{ProjectID: t.Name()})
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()
	if err := rt.Schemas.Register(ctx, "fact", map[string]any{
		"type": "object", "additionalProperties": true,
	}, ""); err != nil {
		t.Fatal(err)
	}
	node := Node{ID: "fact:1", Type: "fact", Model: map[string]any{}, Sources: []string{}, Extras: map[string]any{"text": "one"}}
	if _, err := rt.Graph.AddNodes(ctx, []Node{node}); err != nil {
		t.Fatal(err)
	}
	first, err := rt.Repo.Commit(ctx, "main", "first", nil)
	if err != nil {
		t.Fatal(err)
	}
	objects, err := rt.Repo.ExportTreeObjects(ctx, first.Tree)
	if err != nil || len(objects) == 0 {
		t.Fatalf("export tree objects: %s %v", objects, err)
	}
	entries, err := rt.Repo.TreeEntries(ctx, first.Tree, nil)
	if err != nil || entries["fact"]["fact:1"] == "" {
		t.Fatalf("tree entries: %#v %v", entries, err)
	}
	found, err := rt.Repo.FindTreeEntry(ctx, first.Tree, "fact:1")
	if err != nil || found == nil || *found != entries["fact"]["fact:1"] {
		t.Fatalf("find tree entry: %v %v", found, err)
	}

	node.Extras["text"] = "two"
	if _, err := rt.Graph.AddNodes(ctx, []Node{node}); err != nil {
		t.Fatal(err)
	}
	indexedJSON, err := rt.Raw.IndexGraph(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var indexed struct {
		Entries map[string]map[string]string `json:"entries"`
	}
	if err := json.Unmarshal(indexedJSON, &indexed); err != nil {
		t.Fatal(err)
	}
	t.Logf("indexed=%s", indexedJSON)
	if indexed.Entries["fact"]["fact:1"] == "" {
		t.Fatalf("index graph: %s", indexedJSON)
	}
	second, err := rt.Repo.CommitDelta(ctx, CommitDeltaRequest{
		DeltaEntries: indexed.Entries, RefName: "main", Message: "second", CreatedAt: time.Now().Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("first=%+v second=%+v", first, second)
	if second.Parent != first.ID {
		t.Fatalf("parent = %q, want %q", second.Parent, first.ID)
	}
	diff, err := rt.Repo.DiffTreeEntries(ctx, first.Tree, second.Tree)
	if err != nil || diff["fact"]["fact:1"][0] == nil || diff["fact"]["fact:1"][1] == nil {
		t.Fatalf("diff tree entries: %#v %v", diff, err)
	}
}

func TestCommitDeltaAcceptsExternallyIndexedHashes(t *testing.T) {
	ctx := context.Background()
	indexer, err := Open(ctx, Config{ProjectID: t.Name() + ":index"})
	if err != nil {
		t.Fatal(err)
	}
	defer indexer.Close()
	if err := indexer.Schemas.Register(ctx, "fact", map[string]any{"type": "object", "additionalProperties": true}, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := indexer.Graph.AddNodes(ctx, []Node{{ID: "fact:1", Type: "fact", Model: map[string]any{"created_at": 1, "updated_at": 1}, Sources: []string{}, Extras: map[string]any{"text": "one"}}}); err != nil {
		t.Fatal(err)
	}
	indexedJSON, err := indexer.Raw.IndexGraph(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var indexed struct {
		Entries map[string]map[string]string `json:"entries"`
	}
	if err := json.Unmarshal(indexedJSON, &indexed); err != nil {
		t.Fatal(err)
	}

	repository, err := Open(ctx, Config{ProjectID: t.Name() + ":repo"})
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	commit, err := repository.Repo.CommitDelta(ctx, CommitDeltaRequest{DeltaEntries: indexed.Entries, RefName: "main", Message: "external", CreatedAt: 1})
	if err != nil {
		t.Fatal(err)
	}
	entries, err := repository.Repo.TreeEntries(ctx, commit.Tree, nil)
	if err != nil {
		t.Fatal(err)
	}
	if entries["fact"]["fact:1"] == "" {
		t.Fatalf("entries=%#v", entries)
	}
}
