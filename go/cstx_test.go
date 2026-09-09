package cstx

import (
	"testing"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

func TestRuntimeConfigNormalize(t *testing.T) {
	cfg := normalizeRuntimeConfig(&cstxproto.RuntimeConfig{})
	if cfg.projectID != "default" {
		t.Fatalf("project id: %q", cfg.projectID)
	}
	if cfg.cursorPageSize != DefaultCursorPageSize {
		t.Fatalf("page size: %d", cfg.cursorPageSize)
	}
}

func TestParseErrorRoundTrip(t *testing.T) {
	raw := []byte("missing")
	cerr := parseError(raw, CodeNotFound)
	if cerr.Code != CodeNotFound || cerr.Message != "missing" || !IsCode(cerr, CodeNotFound) {
		t.Fatalf("unexpected error: %+v", cerr)
	}
}

func TestAffected(t *testing.T) {
	change := &cstxproto.GraphChangeSet{AddedNodeIds: []string{"a"}, UpdatedRelationshipIds: []string{"b", "c"}}
	if got := Affected(change); got != 3 {
		t.Fatalf("affected: %d", got)
	}
}
