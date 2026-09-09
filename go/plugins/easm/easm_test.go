package easm

import (
	"testing"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// The generated layer is a projection of the schema document, not a second wire
// format. These tests pin that: it produces the same `EntityValue` a caller
// with no generated code would build by hand, and it reads one back without
// consulting anything but the field names.
//
// Nothing here imports a generated protobuf extension package. Typed access to
// a node type no longer requires a generated protobuf message for it, which is
// the whole point: a
// type declared at runtime has none in any language and never will.

func ptr[T any](value T) *T { return &value }

func TestEntityValueNamesTheSchemaFields(t *testing.T) {
	value, err := Subdomain{
		Host:  "a.example.com",
		IsTld: ptr(false),
		Ttl:   ptr(int64(300)),
		A:     []string{"1.1.1.1", "2.2.2.2"},
	}.EntityValue()
	if err != nil {
		t.Fatal(err)
	}
	if value.GetNodeType() != "subdomain" {
		t.Fatalf("node type = %q, want %q", value.GetNodeType(), "subdomain")
	}

	got := map[string]*cstxproto.EntityField{}
	for _, field := range value.GetFields() {
		got[field.GetName()] = field
	}
	if len(got) != 4 {
		t.Fatalf("field count = %d, want 4 (an unset optional is absent, not zero)", len(got))
	}
	if text := got["host"].GetText(); text != "a.example.com" {
		t.Errorf("host = %q", text)
	}
	// `false` is a value the producer sent, not an absence. A bool that arrives
	// as the flag branch is how the two stay distinguishable.
	if _, ok := got["is_tld"].GetValue().(*cstxproto.EntityField_Flag); !ok {
		t.Errorf("is_tld took the wrong oneof branch: %T", got["is_tld"].GetValue())
	}
	if number := got["ttl"].GetNumber(); number != 300 {
		t.Errorf("ttl = %d", number)
	}
	if values := got["a"].GetList().GetValues(); len(values) != 2 {
		t.Errorf("a = %v", values)
	}
}

func TestFieldsAreOrderedByName(t *testing.T) {
	// One value produces one message: the order is decided when the code is
	// generated, so nothing sorts at run time and two encodings of the same
	// content are byte-identical.
	value, err := Subdomain{Host: "b.example.com", Ttl: ptr(int64(1)), A: []string{"9.9.9.9"}}.EntityValue()
	if err != nil {
		t.Fatal(err)
	}
	previous := ""
	for _, field := range value.GetFields() {
		if field.GetName() <= previous {
			t.Fatalf("fields are not in ascending name order: %q after %q", field.GetName(), previous)
		}
		previous = field.GetName()
	}
}

func TestRoundTripThroughEntityValue(t *testing.T) {
	original := Subdomain{
		Host:  "c.example.com",
		IsTld: ptr(true),
		Ttl:   ptr(int64(60)),
		Cname: []string{"cdn.example.net"},
		Extra: map[string]any{"observed_by": "subfinder"},
	}
	value, err := original.EntityValue()
	if err != nil {
		t.Fatal(err)
	}
	restored, err := SubdomainFrom(value)
	if err != nil {
		t.Fatal(err)
	}
	if restored.Host != original.Host {
		t.Errorf("host = %q, want %q", restored.Host, original.Host)
	}
	if restored.IsTld == nil || *restored.IsTld != true {
		t.Errorf("is_tld = %v", restored.IsTld)
	}
	if restored.Ttl == nil || *restored.Ttl != 60 {
		t.Errorf("ttl = %v", restored.Ttl)
	}
	if len(restored.Cname) != 1 || restored.Cname[0] != "cdn.example.net" {
		t.Errorf("cname = %v", restored.Cname)
	}
	// The declared bag is text on the wire and a document in the caller's hands.
	if restored.Extra["observed_by"] != "subfinder" {
		t.Errorf("extra = %v", restored.Extra)
	}
	// An optional the producer never set stays absent rather than becoming zero.
	if restored.Resolver != nil {
		t.Errorf("resolver = %v, want nil", restored.Resolver)
	}
}

func TestDecodeDispatchesOnTheDeclaredNodeType(t *testing.T) {
	value, err := Ip{Ip: "192.0.2.1", Cdn: ptr(true)}.EntityValue()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := Decode(value)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.CstxType() != "ip" {
		t.Fatalf("decoded as %q", decoded.CstxType())
	}
	if _, ok := decoded.(Ip); !ok {
		t.Fatalf("decoded to %T, want Ip", decoded)
	}
}

func TestDecodeRefusesATypeThisBuildHasNoneFor(t *testing.T) {
	// A type declared at runtime has no generated struct here, and saying so is
	// the correct answer — the untyped `NodeValues` path is what reads it.
	_, err := Decode(&cstxproto.EntityValue{NodeType: "acme_asset"})
	if err == nil {
		t.Fatal("expected an error for a node type with no generated struct")
	}
}

func TestFromRefusesAPayloadOfAnotherType(t *testing.T) {
	value, err := Ip{Ip: "192.0.2.2"}.EntityValue()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := SubdomainFrom(value); err == nil {
		t.Fatal("expected SubdomainFrom to refuse an ip payload")
	}
}

func TestNodeLeavesIdentityToTheRuntime(t *testing.T) {
	// Identity is the schema document's rule; minting one here would be a second
	// place deciding what a node is called.
	node, err := Subdomain{Host: "d.example.com"}.Node("subfinder")
	if err != nil {
		t.Fatal(err)
	}
	if node.Id != nil {
		t.Errorf("id = %v, want unset", node.Id)
	}
	if node.GetValue().GetNodeType() != "subdomain" {
		t.Errorf("payload node type = %q", node.GetValue().GetNodeType())
	}
	if len(node.GetSources()) != 1 || node.GetSources()[0] != "subfinder" {
		t.Errorf("sources = %v", node.GetSources())
	}
}
