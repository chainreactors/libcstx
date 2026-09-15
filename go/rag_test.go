package cstx

import (
	"errors"
	"testing"

	"github.com/chainreactors/libcstx/go/plugins/easm"
	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// appNode carries free text. A domain node holds only its host, which the
// projection does not turn into a searchable record, so an app is the smallest
// node that produces one.
func appNode(t *testing.T, identifier, title string) *cstxproto.Node {
	t.Helper()
	node, err := easm.App{AppId: identifier, Title: &title}.Node("test")
	if err != nil {
		t.Fatalf("build app node: %v", err)
	}
	id := "app:" + identifier
	node.Id = &id
	return node
}

// ragRuntime builds the smallest graph that projects both a node record and a
// relationship record, which is what makes the index worth streaming.
func ragRuntime(t *testing.T) *CSTX {
	t.Helper()
	rt := openRuntime(t)
	if _, err := rt.Graph.AddNodes(testContext, []*cstxproto.Node{
		appNode(t, "admin", "nginx admin console"),
	}); err != nil {
		t.Fatalf("add app node: %v", err)
	}
	addDomain(t, rt, "rag-target.example")
	if _, err := rt.Graph.AddRelationships(testContext, []*cstxproto.Relationship{
		usesRelationship("app:admin", "domain:rag-target.example"),
	}); err != nil {
		t.Fatalf("add relationship: %v", err)
	}
	return rt
}

func TestRagIndexStreamsProjectedRecords(t *testing.T) {
	rt := ragRuntime(t)

	session, err := rt.Rag.Index(testContext, &cstxproto.RagIndexPlan{
		Commit: "revision-1",
		Mode:   cstxproto.RagIndexMode_RAG_INDEX_FULL,
	})
	if err != nil {
		t.Fatalf("index: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	metadata, err := session.Metadata(testContext)
	if err != nil {
		t.Fatalf("metadata: %v", err)
	}
	if metadata.GetCommit() != "revision-1" {
		t.Fatalf("commit = %q; want revision-1", metadata.GetCommit())
	}
	if metadata.GetOperationId() == "" {
		t.Fatal("projection reports no idempotent operation id")
	}
	if metadata.GetUpsertCount() == 0 {
		t.Fatal("projection produced no records to upsert")
	}

	// The iterator walks exactly the records metadata counted.
	iterator, err := session.Records(testContext)
	if err != nil {
		t.Fatalf("records: %v", err)
	}
	t.Cleanup(func() { _ = iterator.Close() })
	kinds := map[cstxproto.RagRecordKind]int{}
	var streamed uint64
	for {
		record, ok, err := iterator.Next(testContext)
		if err != nil {
			t.Fatalf("next: %v", err)
		}
		if !ok {
			break
		}
		if record.GetId() == "" {
			t.Fatalf("record carries no id: %+v", record)
		}
		kinds[record.GetKind()]++
		streamed++
	}
	if streamed != metadata.GetUpsertCount() {
		t.Fatalf("streamed %d record(s); metadata counted %d", streamed, metadata.GetUpsertCount())
	}
	if kinds[cstxproto.RagRecordKind_RAG_RECORD_NODE] == 0 {
		t.Fatalf("no node records projected: %v", kinds)
	}
	if kinds[cstxproto.RagRecordKind_RAG_RECORD_RELATIONSHIP] == 0 {
		t.Fatalf("no relationship records projected: %v", kinds)
	}

	// Pending is the same projection through a caller-controlled window.
	page, err := session.Pending(testContext, 0, 1)
	if err != nil {
		t.Fatalf("pending: %v", err)
	}
	if len(page.GetRecords()) != 1 {
		t.Fatalf("pending returned %d record(s); want 1", len(page.GetRecords()))
	}

	if _, err := session.Deletes(testContext); err != nil {
		t.Fatalf("deletes: %v", err)
	}
}

func TestRagIndexSessionRejectsUseAfterClose(t *testing.T) {
	rt := ragRuntime(t)
	session, err := rt.Rag.Index(testContext, &cstxproto.RagIndexPlan{
		Commit: "revision-1",
		Mode:   cstxproto.RagIndexMode_RAG_INDEX_FULL,
	})
	if err != nil {
		t.Fatalf("index: %v", err)
	}
	if err := session.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if !session.Closed() {
		t.Fatal("session does not report itself closed")
	}
	// Repeated closes are safe, and every read after one is an error rather
	// than a use of a freed native pointer.
	if err := session.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}
	if _, err := session.Deletes(testContext); err == nil {
		t.Fatal("deletes succeeded on a closed session")
	}
	if _, err := session.Metadata(testContext); err == nil {
		t.Fatal("metadata succeeded on a closed session")
	}
}

func TestRagRetrieveCompletesOnceAgainstIndexedGraph(t *testing.T) {
	rt := ragRuntime(t)
	if _, err := rt.Rag.Index(testContext, &cstxproto.RagIndexPlan{
		Commit: "revision-1",
		Mode:   cstxproto.RagIndexMode_RAG_INDEX_FULL,
	}); err != nil {
		t.Fatalf("index: %v", err)
	}

	retrieval, err := rt.Rag.Retrieve(testContext, &cstxproto.RagQuery{
		Text:  "nginx admin console",
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("retrieve: %v", err)
	}
	t.Cleanup(func() { _ = retrieval.Close() })

	if _, err := retrieval.Requests(testContext); err != nil {
		t.Fatalf("requests: %v", err)
	}

	// A nil results means "no external recall ran", which is the built-in
	// lexical case: that recall never leaves Rust.
	result, err := retrieval.Complete(testContext, nil)
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if len(result.GetNodes()) == 0 {
		t.Fatalf("retrieval matched no nodes: %+v", result)
	}

	// Completing consumes the retrieval.
	if _, err := retrieval.Complete(testContext, nil); err == nil {
		t.Fatal("retrieval completed twice")
	}
}

func TestRagRetrievalIsInvalidatedByGraphMutation(t *testing.T) {
	rt := ragRuntime(t)
	if _, err := rt.Rag.Index(testContext, &cstxproto.RagIndexPlan{
		Commit: "revision-1",
		Mode:   cstxproto.RagIndexMode_RAG_INDEX_FULL,
	}); err != nil {
		t.Fatalf("index: %v", err)
	}

	retrieval, err := rt.Rag.Retrieve(testContext, &cstxproto.RagQuery{Text: "nginx admin console", Limit: 20})
	if err != nil {
		t.Fatalf("retrieve: %v", err)
	}
	t.Cleanup(func() { _ = retrieval.Close() })

	// A retrieval is bound to one graph generation, so a write between
	// suspending and completing invalidates it rather than silently answering
	// from a stale graph.
	addDomain(t, rt, "rag-late.example")

	_, err = retrieval.Complete(testContext, nil)
	if err == nil {
		t.Fatal("retrieval completed against a mutated graph")
	}
	var cstxErr *Error
	if !errors.As(err, &cstxErr) || cstxErr.Code != CodeStaleOperation {
		t.Fatalf("complete error = %v; want CodeStaleOperation", err)
	}
}

func TestRagRejectsNilRequests(t *testing.T) {
	rt := ragRuntime(t)
	if _, err := rt.Rag.Index(testContext, nil); err == nil {
		t.Fatal("index accepted a nil plan")
	}
	if _, err := rt.Rag.Retrieve(testContext, nil); err == nil {
		t.Fatal("retrieve accepted a nil query")
	}
}
