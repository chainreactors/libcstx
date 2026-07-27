package cstx

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

// The transport contract is owned by cstx.core.elements (Node/Edge pydantic
// models) and mirrored by codegen/generate_transport_ts.py. These tests pin
// the Go structs to the same field sets so the three languages cannot drift.
func jsonFieldNames(t *testing.T, value any) []string {
	t.Helper()
	typ := reflect.TypeOf(value)
	var names []string
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		names = append(names, strings.Split(tag, ",")[0])
	}
	sort.Strings(names)
	return names
}

func TestNodeMatchesTransportContract(t *testing.T) {
	got := jsonFieldNames(t, Node{})
	want := []string{"extras", "id", "model", "sources", "type", "value"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("node fields drifted from transport contract: got %v want %v", got, want)
	}
}

func TestEdgeMatchesTransportContract(t *testing.T) {
	got := jsonFieldNames(t, Edge{})
	want := []string{"attrs", "id", "relation_type", "source_id", "sources", "target_id"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("edge fields drifted from transport contract: got %v want %v", got, want)
	}
}
