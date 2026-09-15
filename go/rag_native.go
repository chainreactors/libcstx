package cstx

// RAG crosses the ABI as three retained native objects rather than as plain
// calls: an index session holding a projection, a bounded iterator over that
// projection's records, and a suspended retrieval. Each owns a C pointer, so
// each gets the same explicit-close-plus-finalizer treatment as a graph cursor.

/*
#include "cstx_ffi.h"
*/
import "C"

import (
	"context"
	"fmt"
	"runtime"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
	"google.golang.org/protobuf/proto"
)

func (e *nativeEngine) ragIndex(_ context.Context, plan *cstxproto.RagIndexPlan) (ragIndexSession, error) {
	payload, err := proto.Marshal(plan)
	if err != nil {
		return nil, err
	}
	var session *C.CstxRagIndexSession
	err = statusCall("graph.rag.index", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_rag_index(e.handle, byteSlice(payload), &session, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNativeRagIndexSession(session), nil
}

func (e *nativeEngine) ragRetrieve(_ context.Context, query *cstxproto.RagQuery) (ragRetrieval, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, err
	}
	var retrieval *C.CstxRagRetrieval
	err = statusCall("graph.rag.retrieve", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_rag_retrieve(e.handle, byteSlice(payload), &retrieval, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNativeRagRetrieval(retrieval), nil
}

// --- index session -------------------------------------------------------

type nativeRagIndexSession struct{ session *C.CstxRagIndexSession }

func newNativeRagIndexSession(session *C.CstxRagIndexSession) *nativeRagIndexSession {
	result := &nativeRagIndexSession{session: session}
	runtime.SetFinalizer(result, (*nativeRagIndexSession).close)
	return result
}

func (s *nativeRagIndexSession) closedError(operation string) error {
	return &Error{Code: CodeNotInitialized, Operation: operation, Message: "index session is closed"}
}

func (s *nativeRagIndexSession) metadata(_ context.Context) (*cstxproto.RagIndexResult, error) {
	if s.session == nil {
		return nil, s.closedError("graph.rag.index.metadata")
	}
	data, err := bufferResult("graph.rag.index.metadata", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_rag_index_session_metadata(s.session, out, errBuf)
	})
	if err != nil {
		return nil, err
	}
	var result cstxproto.RagIndexResult
	if err := proto.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("cstx: decode rag index result protobuf: %w", err)
	}
	return &result, nil
}

func (s *nativeRagIndexSession) pending(_ context.Context, offset, limit int) (*cstxproto.RagRecordPage, error) {
	if s.session == nil {
		return nil, s.closedError("graph.rag.index.pending")
	}
	data, err := bufferResult("graph.rag.index.pending", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_rag_index_session_pending(s.session, C.size_t(offset), C.size_t(limit), out, errBuf)
	})
	if err != nil {
		return nil, err
	}
	var page cstxproto.RagRecordPage
	if err := proto.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("cstx: decode rag record page protobuf: %w", err)
	}
	return &page, nil
}

func (s *nativeRagIndexSession) deletes(_ context.Context) ([]string, error) {
	if s.session == nil {
		return nil, s.closedError("graph.rag.index.deletes")
	}
	data, err := bufferResult("graph.rag.index.deletes", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_rag_index_session_deletes(s.session, out, errBuf)
	})
	if err != nil {
		return nil, err
	}
	var selection cstxproto.ObjectSelection
	if err := proto.Unmarshal(data, &selection); err != nil {
		return nil, fmt.Errorf("cstx: decode object selection protobuf: %w", err)
	}
	return selection.GetObjectIds(), nil
}

func (s *nativeRagIndexSession) records(_ context.Context) (ragRecordIterator, error) {
	if s.session == nil {
		return nil, s.closedError("graph.rag.index.records")
	}
	var iterator *C.CstxRagRecordIterator
	if err := statusCall("graph.rag.index.records", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_rag_index_session_records(s.session, &iterator, errBuf)
	}); err != nil {
		return nil, err
	}
	return newNativeRagRecordIterator(iterator), nil
}

func (s *nativeRagIndexSession) close() {
	if s.session != nil {
		C.cstx_rag_index_session_close(s.session)
		C.cstx_rag_index_session_free(s.session)
		s.session = nil
		runtime.SetFinalizer(s, nil)
	}
}

// --- record iterator -----------------------------------------------------

type nativeRagRecordIterator struct{ iterator *C.CstxRagRecordIterator }

func newNativeRagRecordIterator(iterator *C.CstxRagRecordIterator) *nativeRagRecordIterator {
	result := &nativeRagRecordIterator{iterator: iterator}
	runtime.SetFinalizer(result, (*nativeRagRecordIterator).close)
	return result
}

func (i *nativeRagRecordIterator) next(_ context.Context) (*cstxproto.RagRecord, bool, error) {
	if i.iterator == nil {
		return nil, false, &Error{Code: CodeInvalidArgument, Operation: "graph.rag.index.records.next", Message: "iterator is closed"}
	}
	var out, errBuf C.CstxBuffer
	var hasValue C.uint8_t
	if err := statusError(
		C.cstx_rag_record_iterator_next(i.iterator, &out, &hasValue, &errBuf),
		"graph.rag.index.records.next",
		&errBuf,
	); err != nil {
		C.cstx_buffer_free(&out)
		return nil, false, err
	}
	data := takeBuffer(&out)
	if hasValue == 0 {
		return nil, false, nil
	}
	var record cstxproto.RagRecord
	if err := proto.Unmarshal(data, &record); err != nil {
		return nil, false, fmt.Errorf("cstx: decode rag record protobuf: %w", err)
	}
	return &record, true, nil
}

func (i *nativeRagRecordIterator) close() {
	if i.iterator != nil {
		C.cstx_rag_record_iterator_close(i.iterator)
		C.cstx_rag_record_iterator_free(i.iterator)
		i.iterator = nil
		runtime.SetFinalizer(i, nil)
	}
}

// --- retrieval -----------------------------------------------------------

type nativeRagRetrieval struct{ retrieval *C.CstxRagRetrieval }

func newNativeRagRetrieval(retrieval *C.CstxRagRetrieval) *nativeRagRetrieval {
	result := &nativeRagRetrieval{retrieval: retrieval}
	runtime.SetFinalizer(result, (*nativeRagRetrieval).close)
	return result
}

func (r *nativeRagRetrieval) requests(_ context.Context) (*cstxproto.RecallPlan, error) {
	if r.retrieval == nil {
		return nil, &Error{Code: CodeNotInitialized, Operation: "graph.rag.retrieve.requests", Message: "retrieval is closed"}
	}
	data, err := bufferResult("graph.rag.retrieve.requests", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_rag_retrieval_requests(r.retrieval, out, errBuf)
	})
	if err != nil {
		return nil, err
	}
	var plan cstxproto.RecallPlan
	if err := proto.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("cstx: decode recall plan protobuf: %w", err)
	}
	return &plan, nil
}

func (r *nativeRagRetrieval) complete(_ context.Context, results *cstxproto.RecallResults) (*cstxproto.RagResult, error) {
	if r.retrieval == nil {
		return nil, &Error{Code: CodeNotInitialized, Operation: "graph.rag.complete", Message: "retrieval is closed"}
	}
	payload, err := proto.Marshal(results)
	if err != nil {
		return nil, err
	}
	data, err := bufferResult("graph.rag.complete", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_rag_retrieval_complete(r.retrieval, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var result cstxproto.RagResult
	if err := proto.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("cstx: decode rag result protobuf: %w", err)
	}
	return &result, nil
}

func (r *nativeRagRetrieval) close() {
	if r.retrieval != nil {
		C.cstx_rag_retrieval_close(r.retrieval)
		C.cstx_rag_retrieval_free(r.retrieval)
		r.retrieval = nil
		runtime.SetFinalizer(r, nil)
	}
}
