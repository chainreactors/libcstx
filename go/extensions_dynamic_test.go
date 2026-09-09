package cstx

import (
	"context"
	"encoding/json"
	"testing"
)

// What a Go consumer can ask about a type it declared at runtime.
//
// No generated message type exists for `acme.Asset` and none can, because the
// type is declared by a document that ships as test data — so registering it
// is the whole story, and the runtime must answer for it exactly as it does
// for a built-in. Writing and reading one is `TestDynamicExtensionConformance`
// next door, against the same document; this covers the registration surface
// that test does not touch.
func TestDynamicExtensionRegistration(t *testing.T) {
	ctx := context.Background()
	fixture := loadDynamicFixture(t)
	runtime := openDynamicRuntime(t, fixture)

	ok, err := runtime.Extensions.Contains(ctx, "acme_asset")
	if err != nil || !ok {
		t.Fatalf("contains acme_asset = %v, %v; want true", ok, err)
	}

	var document struct {
		Nodes map[string]struct {
			Message string `json:"message"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(fixture.Document, &document); err != nil {
		t.Fatalf("document: %v", err)
	}

	schema, err := runtime.Extensions.Schema(ctx, "acme_asset")
	if err != nil {
		t.Fatalf("schema: %v", err)
	}
	want := "type.googleapis.com/" + document.Nodes["acme_asset"].Message
	if schema.GetTypeUrl() != want {
		t.Fatalf("type_url = %q; want %q", schema.GetTypeUrl(), want)
	}

	// A type the document never declared is not registered by accident.
	if ok, err := runtime.Extensions.Contains(ctx, "acme_absent"); err != nil || ok {
		t.Fatalf("contains acme_absent = %v, %v; want false", ok, err)
	}
	if _, err := runtime.Extensions.Schema(ctx, "acme_absent"); err == nil {
		t.Fatal("schema of an undeclared type should fail")
	}
}
