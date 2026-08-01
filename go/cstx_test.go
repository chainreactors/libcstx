package cstx

import (
	"encoding/json"
	"testing"
)

func TestConfigNormalize(t *testing.T) {
	cfg := Config{}.normalize()
	if cfg.ProjectID != "default" {
		t.Fatalf("project id: %q", cfg.ProjectID)
	}
	if cfg.CursorPageSize != DefaultCursorPageSize {
		t.Fatalf("page size: %d", cfg.CursorPageSize)
	}
}

func TestNodeMarshalKeepsEmptyLists(t *testing.T) {
	data, err := json.Marshal(Node{ID: "n1", Type: "Domain", Value: "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"sources", "model", "extras"} {
		if decoded[field] == nil {
			t.Fatalf("field %s must not be null: %s", field, data)
		}
	}
}

func TestEdgeFilterNullIDs(t *testing.T) {
	data, err := json.Marshal(EdgeFilter{})
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"source_id", "target_id"} {
		value, present := decoded[field]
		if !present || value != nil {
			t.Fatalf("field %s must be explicit null: %s", field, data)
		}
	}
	for _, field := range []string{"relations", "sources"} {
		if decoded[field] == nil {
			t.Fatalf("field %s must be []: %s", field, data)
		}
	}
}

func TestCollectionOptionsNormalize(t *testing.T) {
	options := CollectionOptions{}.normalize()
	if options.PageSize != DefaultCursorPageSize || options.Order != OrderUnspecified {
		t.Fatalf("unexpected defaults: %+v", options)
	}
}

func TestRefUnmarshalTuple(t *testing.T) {
	var refs []Ref
	if err := json.Unmarshal([]byte(`[["main","abc123"],["dev","def456"]]`), &refs); err != nil {
		t.Fatal(err)
	}
	if len(refs) != 2 || refs[0].Name != "main" || refs[0].Head != "abc123" {
		t.Fatalf("unexpected refs: %+v", refs)
	}
}

func TestParseErrorRoundTrip(t *testing.T) {
	raw := []byte(`{"code":"NOT_FOUND","operation":"graph.node","item_index":null,"field":"node_id","message":"missing","expected":null,"actual":null}`)
	cerr := parseError(raw, CodeInternal)
	if cerr.Code != CodeNotFound || cerr.Operation != "graph.node" || cerr.Field != "node_id" {
		t.Fatalf("unexpected error: %+v", cerr)
	}
	if !IsCode(cerr, CodeNotFound) {
		t.Fatal("IsCode mismatch")
	}

	fallback := parseError([]byte("plain failure"), CodeIO)
	if fallback.Code != CodeIO || fallback.Message != "plain failure" {
		t.Fatalf("unexpected fallback: %+v", fallback)
	}
}

func TestChangeSetAffected(t *testing.T) {
	change := ChangeSet{AddedNodeIDs: []string{"a"}, UpdatedEdgeIDs: []string{"b", "c"}}
	if change.Affected() != 3 {
		t.Fatalf("affected: %d", change.Affected())
	}
}
