package cstx

import (
	"context"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// Rag is the retrieval namespace of a CSTX runtime. It owns two suspended
// operations rather than two calls: Index projects the graph into records the
// caller embeds and stores, and Retrieve plans recall the caller executes
// against its own vector store. CSTX never talks to an embedder or an index;
// it decides what to ask for and what to do with the answers.
type Rag struct{ eng engine }

// Index projects the graph into a retained session of records to upsert and
// IDs to delete. Nothing is written anywhere: the caller streams the records
// out, embeds and stores them, and closes the session.
func (r *Rag) Index(ctx context.Context, plan *cstxproto.RagIndexPlan) (*RagIndexSession, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.rag.index", Message: "plan must not be nil"}
	}
	session, err := r.eng.ragIndex(ctx, plan)
	if err != nil {
		return nil, err
	}
	return &RagIndexSession{inner: session}, nil
}

// Retrieve suspends a retrieval against one graph generation and checkpoint.
// Read the plan with Requests, run it against the vector store, and hand the
// answers back through Complete.
func (r *Rag) Retrieve(ctx context.Context, query *cstxproto.RagQuery) (*RagRetrieval, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if query == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.rag.retrieve", Message: "query must not be nil"}
	}
	retrieval, err := r.eng.ragRetrieve(ctx, query)
	if err != nil {
		return nil, err
	}
	return &RagRetrieval{inner: retrieval}, nil
}

// RagIndexSession is a retained deterministic projection. It holds native
// state until closed, so callers should Close it when done; repeated calls are
// safe.
type RagIndexSession struct {
	inner ragIndexSession
	done  bool
}

func (s *RagIndexSession) closedError(operation string) error {
	return &Error{Code: CodeNotInitialized, Operation: operation, Message: "index session is closed"}
}

// Metadata returns the projection's identity and counts: the idempotent
// operation ID, the commit it targets, whether it replaces the whole index,
// and how many records it upserts and deletes.
func (s *RagIndexSession) Metadata(ctx context.Context) (*cstxproto.RagIndexResult, error) {
	if s.done {
		return nil, s.closedError("graph.rag.index.metadata")
	}
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return s.inner.metadata(ctx)
}

// Pending returns one bounded page of projected records. Records streams the
// same records without paging arithmetic; Pending is for a caller that wants
// to control the window itself.
func (s *RagIndexSession) Pending(ctx context.Context, offset, limit int) (*cstxproto.RagRecordPage, error) {
	if s.done {
		return nil, s.closedError("graph.rag.index.pending")
	}
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if offset < 0 || limit < 0 {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.rag.index.pending", Message: "offset and limit must not be negative"}
	}
	return s.inner.pending(ctx, offset, limit)
}

// Deletes returns the record IDs this projection removes from the index.
func (s *RagIndexSession) Deletes(ctx context.Context) ([]string, error) {
	if s.done {
		return nil, s.closedError("graph.rag.index.deletes")
	}
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return s.inner.deletes(ctx)
}

// Records streams every projected record through a bounded native iterator,
// which never materializes the whole projection in memory. The caller must
// close the iterator.
func (s *RagIndexSession) Records(ctx context.Context) (*RagRecordIterator, error) {
	if s.done {
		return nil, s.closedError("graph.rag.index.records")
	}
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	iterator, err := s.inner.records(ctx)
	if err != nil {
		return nil, err
	}
	return &RagRecordIterator{inner: iterator}, nil
}

// Close releases the retained projection; repeated calls are safe.
func (s *RagIndexSession) Close() error {
	if !s.done {
		s.done = true
		s.inner.close()
	}
	return nil
}

// Closed reports whether Close has been called.
func (s *RagIndexSession) Closed() bool { return s.done }

// RagRecordIterator walks the projected records of one index session.
type RagRecordIterator struct {
	inner ragRecordIterator
	done  bool
}

// Next returns the next record. The boolean is false at the end of the
// projection, which is not an error.
func (i *RagRecordIterator) Next(ctx context.Context) (*cstxproto.RagRecord, bool, error) {
	if i.done {
		return nil, false, &Error{Code: CodeInvalidArgument, Operation: "graph.rag.index.records.next", Message: "iterator is closed"}
	}
	if err := contextError(ctx); err != nil {
		return nil, false, err
	}
	return i.inner.next(ctx)
}

// Close releases the iterator early; repeated calls are safe.
func (i *RagRecordIterator) Close() error {
	if !i.done {
		i.done = true
		i.inner.close()
	}
	return nil
}

// Closed reports whether Close has been called.
func (i *RagRecordIterator) Closed() bool { return i.done }

// RagRetrieval is a suspended retrieval bound to one graph generation and
// checkpoint. Complete consumes it.
type RagRetrieval struct {
	inner     ragRetrieval
	done      bool
	completed bool
}

// Requests returns the recall plan: what the caller must look up in its own
// vector store before CSTX can finish the retrieval.
func (r *RagRetrieval) Requests(ctx context.Context) (*cstxproto.RecallPlan, error) {
	if r.done || r.completed {
		return nil, &Error{Code: CodeNotInitialized, Operation: "graph.rag.retrieve.requests", Message: "retrieval is closed or completed"}
	}
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.inner.requests(ctx)
}

// Complete fuses the recall results into a final answer and consumes the
// retrieval. A nil results means the caller ran no external recall, which is
// the normal case when the built-in lexical projection is the only source —
// that recall never leaves Rust. Calling Complete twice reports CodeConflict.
func (r *RagRetrieval) Complete(ctx context.Context, results *cstxproto.RecallResults) (*cstxproto.RagResult, error) {
	if r.done {
		return nil, &Error{Code: CodeNotInitialized, Operation: "graph.rag.complete", Message: "retrieval is closed"}
	}
	if r.completed {
		return nil, &Error{Code: CodeConflict, Operation: "graph.rag.complete", Message: "retrieval has already completed"}
	}
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if results == nil {
		results = &cstxproto.RecallResults{}
	}
	result, err := r.inner.complete(ctx, results)
	if err != nil {
		return nil, err
	}
	r.completed = true
	return result, nil
}

// Close releases the suspended retrieval; repeated calls are safe.
func (r *RagRetrieval) Close() error {
	if !r.done {
		r.done = true
		r.inner.close()
	}
	return nil
}

// Closed reports whether Close has been called.
func (r *RagRetrieval) Closed() bool { return r.done }
