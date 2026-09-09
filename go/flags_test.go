package cstx

import (
	"encoding/json"
	"testing"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// The registry reads declarations, it does not carry them: a flag any
// extension declares is reachable on the same terms as a built-in one, which
// is the whole reason the wire stopped carrying a closed enum.
func TestFlagRegistryReadsExtensionDeclarations(t *testing.T) {
	document, err := json.Marshal(map[string]any{
		"schema_version": 1,
		"extension":      "acme",
		"flags": map[string]any{
			"honeypot":    map[string]any{"bit": 0, "default_exclude": true},
			"quarantined": map[string]any{"bit": 40, "default_exclude": true},
			"reviewed":    map[string]any{"bit": 41},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	contract := &cstxproto.ExtensionContract{
		Extensions: map[string]*cstxproto.ExtensionDefinition{
			"acme": {Name: "acme", Schema: string(document)},
		},
	}
	registry, err := Flags(contract)
	if err != nil {
		t.Fatal(err)
	}

	// A bit the retired enum could not express at all.
	if bit, ok := registry.Bit("quarantined"); !ok || bit != 40 {
		t.Fatalf("quarantined: got bit %d ok=%v, want 40 true", bit, ok)
	}
	if mask := registry.Mask("quarantined"); mask != 1<<40 {
		t.Fatalf("quarantined mask: got %#x, want %#x", mask, uint64(1)<<40)
	}
	if mask := registry.Mask("nosuchflag"); mask != 0 {
		t.Fatalf("undeclared flag must mask to 0, got %#x", mask)
	}

	wantAll := uint64(1)<<0 | uint64(1)<<40 | uint64(1)<<41
	if got := registry.AllMask(); got != wantAll {
		t.Fatalf("AllMask: got %#x, want %#x", got, wantAll)
	}
	// Advice, not enforcement -- `reviewed` declares no default_exclude.
	wantExclude := uint64(1)<<0 | uint64(1)<<40
	if got := registry.DefaultExcludeMask(); got != wantExclude {
		t.Fatalf("DefaultExcludeMask: got %#x, want %#x", got, wantExclude)
	}

	names := registry.Names()
	want := []string{"honeypot", "quarantined", "reviewed"} // lowest bit first
	if len(names) != len(want) {
		t.Fatalf("Names: got %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("Names: got %v, want %v", names, want)
		}
	}
}

func TestFlagRegistryToleratesAnEmptyContract(t *testing.T) {
	registry, err := Flags(nil)
	if err != nil {
		t.Fatal(err)
	}
	if registry.AllMask() != 0 || len(registry.Names()) != 0 {
		t.Fatal("an empty contract must declare no flags")
	}
}
