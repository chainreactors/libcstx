package cstx

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
	"google.golang.org/protobuf/types/known/anypb"
)

// One field as the shared fixture's schema document declares it.
type declaredField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Repeated bool   `json:"repeated"`
}

type dynamicFixture struct {
	Document json.RawMessage `json:"document"`
	// Raw, so it can be decoded with UseNumber: one of the fixture's int64
	// values is outside what a float64 can hold, and that is why it is there.
	RawValues  json.RawMessage `json:"values"`
	ExpectedID string          `json:"expected_id"`
	Relation   struct {
		RelationType string `json:"relation_type"`
		TypeURL      string `json:"type_url"`
	} `json:"relation"`
}

// fixtureValues converts the fixture's JSON into the Go types the schema's
// columns hold. The document says which type each field is, so nothing here
// guesses — and the conversion is what proves Go and Python agree about it.
func fixtureValues(t *testing.T, fixture dynamicFixture) NodeValues {
	t.Helper()
	var document struct {
		Nodes map[string]struct {
			Fields []declaredField `json:"fields"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(fixture.Document, &document); err != nil {
		t.Fatalf("document: %v", err)
	}
	declared := map[string]declaredField{}
	for _, field := range document.Nodes["acme_asset"].Fields {
		declared[field.Name] = field
	}

	decoder := json.NewDecoder(strings.NewReader(string(fixture.RawValues)))
	decoder.UseNumber()
	raw := map[string]any{}
	if err := decoder.Decode(&raw); err != nil {
		t.Fatalf("values: %v", err)
	}

	values := NodeValues{}
	for name, value := range raw {
		field, ok := declared[name]
		if !ok {
			t.Fatalf("fixture value %q is not declared by the document", name)
		}
		switch {
		case field.Repeated:
			items := value.([]any)
			list := make([]string, 0, len(items))
			for _, item := range items {
				list = append(list, item.(string))
			}
			values[name] = list
		case field.Type == "bool":
			values[name] = value.(bool)
		case field.Type == "string":
			values[name] = value.(string)
		case field.Type == "double" || field.Type == "float":
			real, err := value.(json.Number).Float64()
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			values[name] = real
		default:
			number, err := value.(json.Number).Int64()
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			values[name] = number
		}
	}
	return values
}

func openDynamicRuntime(t *testing.T, fixture dynamicFixture) *CSTX {
	t.Helper()
	runtime, err := Open(context.Background(), &cstxproto.RuntimeConfig{
		ProjectId:     "go-dyn-conformance",
		PayloadFormat: cstxproto.PayloadFormat_PAYLOAD_FORMAT_VALUE,
	})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { runtime.Close() })
	contract := newExtensionBuilder("acme", "1.0").Schema(string(fixture.Document)).Build()
	if err := runtime.Extensions.Register(context.Background(), contract); err != nil {
		t.Fatalf("register: %v", err)
	}
	return runtime
}

func loadDynamicFixture(t *testing.T) dynamicFixture {
	t.Helper()
	var fixture struct {
		DynamicExtension dynamicFixture `json:"dynamic_extension"`
	}
	if err := json.Unmarshal(conformanceFixture, &fixture); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	return fixture.DynamicExtension
}

// The shared fixture's runtime-declared type, written and read from Go.
//
// The point of the fixture being shared is that Python runs the same document
// and the same values through its own SDK and asserts the same id and the same
// values. Neither language has a generated message type for `acme.Asset`, and
// neither can: the type is declared by a document that ships as test data.
//
// Its fields cover every proto type a schema document may declare, so a column
// that cannot survive the trip fails here rather than the first time a plugin
// uses it.
func TestDynamicExtensionConformance(t *testing.T) {
	ctx := context.Background()
	fixture := loadDynamicFixture(t)
	runtime := openDynamicRuntime(t, fixture)
	values := fixtureValues(t, fixture)

	if _, err := runtime.Graph.AddNodeValues(ctx, "acme_asset", values); err != nil {
		t.Fatalf("add_node_values: %v", err)
	}

	node, err := runtime.Graph.Node(ctx, fixture.ExpectedID)
	if err != nil {
		t.Fatalf("node %q: %v", fixture.ExpectedID, err)
	}
	nodeType, read, err := FieldValues(node)
	if err != nil {
		t.Fatalf("field values: %v", err)
	}
	if nodeType != "acme_asset" {
		t.Fatalf("node_type = %q", nodeType)
	}
	if !reflect.DeepEqual(read, values) {
		t.Fatalf("read back %#v; want %#v", read, values)
	}
}

// A relation type declared by the same document, with no generated code.
//
// A relation carries no payload — the document declares that the type exists
// and which message names it, and the bytes are empty. So building one needs
// the type URL and nothing else, which is the whole reason a relation type
// declared at runtime can work at all.
func TestDynamicExtensionRelationRoundTrips(t *testing.T) {
	ctx := context.Background()
	fixture := loadDynamicFixture(t)
	runtime := openDynamicRuntime(t, fixture)
	values := fixtureValues(t, fixture)

	other := NodeValues{"asset_id": "a-2"}
	for _, payload := range []NodeValues{values, other} {
		if _, err := runtime.Graph.AddNodeValues(ctx, "acme_asset", payload); err != nil {
			t.Fatalf("add_node_values: %v", err)
		}
	}

	edge := &cstxproto.Relationship{
		SourceId: fixture.ExpectedID,
		TargetId: "acme_asset:a-2",
		Sources:  []string{"test"},
		Relation: &anypb.Any{TypeUrl: fixture.Relation.TypeURL},
	}
	if _, err := runtime.Graph.AddRelationships(ctx, []*cstxproto.Relationship{edge}); err != nil {
		t.Fatalf("add_relationships: %v", err)
	}

	cursor, err := runtime.Graph.Relationships(ctx, &cstxproto.RelationshipQuery{})
	if err != nil {
		t.Fatalf("relationships: %v", err)
	}
	defer cursor.Close()
	page, err := cursor.Page(ctx, 10, 1)
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	stored := page.GetRelationships().GetValues()
	if len(stored) != 1 {
		t.Fatalf("relationship count = %d; want 1", len(stored))
	}
	if got := stored[0].GetRelation().GetTypeUrl(); got != fixture.Relation.TypeURL {
		t.Fatalf("relation type_url = %q; want %q", got, fixture.Relation.TypeURL)
	}
}

// conformanceFixture is the whole shared file, as the three languages read it.
type conformanceFile struct {
	Document json.RawMessage `json:"document"`
	Nodes    []struct {
		ID      string            `json:"id"`
		Type    string            `json:"type"`
		Model   map[string]string `json:"model"`
		Sources []string          `json:"sources"`
	} `json:"nodes"`
	Relationships []struct {
		ID           string   `json:"id"`
		SourceID     string   `json:"source_id"`
		TargetID     string   `json:"target_id"`
		RelationType string   `json:"relation_type"`
		Sources      []string `json:"sources"`
	} `json:"relationships"`
	Query    string `json:"query"`
	Expected struct {
		NodeIDs           []string `json:"node_ids"`
		NodeCount         uint64   `json:"node_count"`
		RelationshipCount uint64   `json:"relationship_count"`
	} `json:"expected"`
}

func loadConformanceFixture(t *testing.T) conformanceFile {
	t.Helper()
	var fixture conformanceFile
	if err := json.Unmarshal(conformanceFixture, &fixture); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	return fixture
}

// sortedKeys keeps one map producing one message: Go randomizes map iteration,
// and a node's payload is content, not a bag.
func sortedKeys(values map[string]string) []string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
