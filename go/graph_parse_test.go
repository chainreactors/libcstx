package cstx

import (
	"testing"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// The native accelerator is only half of an ingest: Parse returns a batch the
// caller feeds through the same AddNodes/Link path a parser in any language
// uses. This walks that whole path, because either half alone passing does not
// show the batch is actually writable.
func TestParseAddNodesLink(t *testing.T) {
	rt := openRuntime(t)
	ctx := testContext

	data := []byte(`{"ip":"127.0.0.1","port":"443","protocol":"https","host":"example.com","uri":"/","title":"hello"}` + "\n")
	batch, records, err := rt.Graph.Parse(ctx, &cstxproto.ParserPayload{
		Plugin:   "easm",
		Artifact: "gogo",
		Data:     data,
	})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if records != 1 {
		t.Fatalf("records = %d; want 1", records)
	}
	// One gogo record yields an ip, a cidr, a port and an app.
	if len(batch.Nodes) != 4 {
		t.Fatalf("parsed %d nodes; want 4", len(batch.Nodes))
	}

	affected, err := rt.Graph.AddNodes(ctx, batch.Nodes)
	if err != nil {
		t.Fatalf("add nodes: %v", err)
	}
	if affected != 4 {
		t.Fatalf("affected = %d; want 4", affected)
	}

	ids := make([]string, 0, len(batch.Nodes))
	for _, node := range batch.Nodes {
		ids = append(ids, node.GetId())
	}
	link, err := rt.Graph.Link(ctx, ids, "graph_parse_test")
	if err != nil {
		t.Fatalf("link: %v", err)
	}
	if len(link.RelationshipIds) == 0 {
		t.Fatalf("link created no relationships: %+v", link)
	}

	node, err := rt.Graph.Node(ctx, ids[0])
	if err != nil {
		t.Fatalf("node: %v", err)
	}
	if node.GetValue().GetNodeType() != "ip" {
		t.Fatalf("node type = %q; want ip", node.GetValue().GetNodeType())
	}
}

// Parse is the accelerator, not the write: an unenabled extension has no
// registered parser to dispatch to, and the batch is not committed either way.
func TestParseWithoutEnabledExtension(t *testing.T) {
	rt, err := Open(testContext, &cstxproto.RuntimeConfig{ProjectId: "sdk-go-parse-disabled"})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = rt.Close() })

	_, _, err = rt.Graph.Parse(testContext, &cstxproto.ParserPayload{
		Plugin:   "easm",
		Artifact: "gogo",
		Data:     []byte(`{"ip":"127.0.0.1","port":"443"}` + "\n"),
	})
	if err == nil {
		t.Fatal("parse without an enabled extension reported success")
	}
	count, err := rt.Graph.NodeCount(testContext)
	if err != nil {
		t.Fatalf("node count: %v", err)
	}
	if count != 0 {
		t.Fatalf("node count = %d; want 0", count)
	}
}
