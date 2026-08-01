package cstx

import (
	"context"
	"encoding/json"
	"runtime"
	"sync"
)

// Raw exposes advanced native operations without duplicating Rust DTOs in Go.
// Payloads use the JSON shapes documented by the C ABI.
type Raw struct{ eng rawEngine }

type rawEngine interface {
	rawSchemaRegisterJoinRule(context.Context, []byte) error
	rawSchemaAnchorConcepts(context.Context) ([]byte, error)
	rawGraphIngestNative(context.Context, string, string, []byte) ([]byte, error)
	rawGraphAddNodes(context.Context, []byte) (uint64, error)
	rawGraphAddEdges(context.Context, []byte) (uint64, error)
	rawGraphFindNode(context.Context, string) ([]byte, error)
	rawGraphNodeTypes(context.Context) ([]byte, error)
	rawGraphNodesPage(context.Context, []byte) ([]byte, error)
	rawGraphLink(context.Context, []byte, string) ([]byte, error)
	rawRepositoryIndexGraph(context.Context) ([]byte, error)
	rawRAGIndex(context.Context, []byte) (rawRAGIndexSession, error)
	rawRAGRetrieve(context.Context, []byte) (rawRAGRetrieval, error)
}

type rawRAGIndexSession interface {
	metadata(context.Context) ([]byte, error)
	pending(context.Context, int, int) ([]byte, error)
	deletes(context.Context) ([]byte, error)
	close()
}

type rawRAGRetrieval interface {
	requests(context.Context) ([]byte, error)
	complete(context.Context, []byte) ([]byte, error)
	close()
}

func rawContext(ctx context.Context) error { return contextError(ctx) }

// RegisterJoinRule registers one linker rule encoded as JSON.
func (r *Raw) RegisterJoinRule(ctx context.Context, request json.RawMessage) error {
	if err := rawContext(ctx); err != nil {
		return err
	}
	return r.eng.rawSchemaRegisterJoinRule(ctx, request)
}

// AnchorConcepts returns native concept tuples as JSON.
func (r *Raw) AnchorConcepts(ctx context.Context) (json.RawMessage, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	return r.eng.rawSchemaAnchorConcepts(ctx)
}

// IngestNative invokes a linked native plugin and returns mutation details as JSON.
func (r *Raw) IngestNative(ctx context.Context, plugin, artifact string, data []byte) (json.RawMessage, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	return r.eng.rawGraphIngestNative(ctx, plugin, artifact, data)
}

// AddNodes submits the canonical JSON node array without Go DTO conversion.
func (r *Raw) AddNodes(ctx context.Context, nodes json.RawMessage) (uint64, error) {
	if err := rawContext(ctx); err != nil {
		return 0, err
	}
	return r.eng.rawGraphAddNodes(ctx, nodes)
}

// AddEdges submits the canonical JSON edge array without Go DTO conversion.
func (r *Raw) AddEdges(ctx context.Context, edges json.RawMessage) (uint64, error) {
	if err := rawContext(ctx); err != nil {
		return 0, err
	}
	return r.eng.rawGraphAddEdges(ctx, edges)
}

// FindNode returns a node or JSON null.
func (r *Raw) FindNode(ctx context.Context, identifier string) (json.RawMessage, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	return r.eng.rawGraphFindNode(ctx, identifier)
}

// NodeTypes returns registered graph node types as JSON.
func (r *Raw) NodeTypes(ctx context.Context) (json.RawMessage, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	return r.eng.rawGraphNodeTypes(ctx)
}

// NodesPage executes the raw node-page request and returns its JSON response.
func (r *Raw) NodesPage(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	return r.eng.rawGraphNodesPage(ctx, request)
}

// Link runs native linker rules for the selected node IDs and returns the
// linker result as JSON.
func (r *Raw) Link(ctx context.Context, nodeIDs json.RawMessage, dataSource string) (json.RawMessage, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	return r.eng.rawGraphLink(ctx, nodeIDs, dataSource)
}

// IndexGraph returns canonical CAS entries and objects as JSON.
func (r *Raw) IndexGraph(ctx context.Context) (json.RawMessage, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	return r.eng.rawRepositoryIndexGraph(ctx)
}

// RAGIndex opens an opaque native projection session from a JSON request.
func (r *Raw) RAGIndex(ctx context.Context, request json.RawMessage) (*RawRAGIndexSession, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	inner, err := r.eng.rawRAGIndex(ctx, request)
	if err != nil {
		return nil, err
	}
	session := &RawRAGIndexSession{inner: inner}
	runtime.SetFinalizer(session, (*RawRAGIndexSession).finalize)
	return session, nil
}

// RAGRetrieve opens an opaque native retrieval from a JSON query.
func (r *Raw) RAGRetrieve(ctx context.Context, query json.RawMessage) (*RawRAGRetrieval, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	inner, err := r.eng.rawRAGRetrieve(ctx, query)
	if err != nil {
		return nil, err
	}
	retrieval := &RawRAGRetrieval{inner: inner}
	runtime.SetFinalizer(retrieval, (*RawRAGRetrieval).finalize)
	return retrieval, nil
}

// RawRAGIndexSession is an opaque native index-session handle.
type RawRAGIndexSession struct {
	mu     sync.Mutex
	inner  rawRAGIndexSession
	closed bool
}

func (s *RawRAGIndexSession) use(ctx context.Context, call func(rawRAGIndexSession) ([]byte, error)) (json.RawMessage, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, &Error{Code: CodeNotInitialized, Operation: "rag.index", Message: "RAG index session is closed"}
	}
	return call(s.inner)
}

// Metadata returns operation, commit, mode, and counts as JSON.
func (s *RawRAGIndexSession) Metadata(ctx context.Context) (json.RawMessage, error) {
	return s.use(ctx, func(inner rawRAGIndexSession) ([]byte, error) { return inner.metadata(ctx) })
}

// Pending returns one page of projected records as JSON.
func (s *RawRAGIndexSession) Pending(ctx context.Context, offset, limit int) (json.RawMessage, error) {
	return s.use(ctx, func(inner rawRAGIndexSession) ([]byte, error) { return inner.pending(ctx, offset, limit) })
}

// Deletes returns projected deletion IDs as JSON.
func (s *RawRAGIndexSession) Deletes(ctx context.Context) (json.RawMessage, error) {
	return s.use(ctx, func(inner rawRAGIndexSession) ([]byte, error) { return inner.deletes(ctx) })
}

// Close releases the native index session. Repeated calls are safe.
func (s *RawRAGIndexSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.closed {
		s.inner.close()
		s.closed = true
		runtime.SetFinalizer(s, nil)
	}
	return nil
}
func (s *RawRAGIndexSession) finalize() { _ = s.Close() }

// RawRAGRetrieval is an opaque native retrieval handle that completes once.
type RawRAGRetrieval struct {
	mu                sync.Mutex
	inner             rawRAGRetrieval
	closed, completed bool
}

// Requests returns external recall requests as JSON.
func (r *RawRAGRetrieval) Requests(ctx context.Context) (json.RawMessage, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed || r.completed {
		return nil, &Error{Code: CodeNotInitialized, Operation: "rag.retrieve.requests", Message: "RAG retrieval is closed or completed"}
	}
	return r.inner.requests(ctx)
}

// Complete submits recall batches as JSON and returns the Rust-computed result.
func (r *RawRAGRetrieval) Complete(ctx context.Context, batches json.RawMessage) (json.RawMessage, error) {
	if err := rawContext(ctx); err != nil {
		return nil, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil, &Error{Code: CodeNotInitialized, Operation: "rag.retrieve.complete", Message: "RAG retrieval is closed"}
	}
	if r.completed {
		return nil, &Error{Code: CodeConflict, Operation: "rag.retrieve.complete", Message: "RAG retrieval is already completed"}
	}
	result, err := r.inner.complete(ctx, batches)
	if err == nil {
		r.completed = true
		runtime.SetFinalizer(r, nil)
	}
	return result, err
}

// Close releases an incomplete retrieval. Repeated calls are safe.
func (r *RawRAGRetrieval) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.closed && !r.completed {
		r.inner.close()
	}
	r.closed = true
	runtime.SetFinalizer(r, nil)
	return nil
}
func (r *RawRAGRetrieval) finalize() { _ = r.Close() }
