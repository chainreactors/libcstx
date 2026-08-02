package cstx

/*
#include "cstx_ffi.h"
*/
import "C"

import (
	"context"
	"runtime"
)

func (e *nativeEngine) rawSchemaRegisterJoinRule(_ context.Context, request []byte) error {
	return statusCall("schemas.register_join_rule", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_schema_register_join_rule(e.handle, byteSlice(request), errBuf)
		runtime.KeepAlive(request)
		return rc
	})
}

func (e *nativeEngine) rawSchemaAnchorConcepts(_ context.Context) ([]byte, error) {
	return bufferResult("schemas.anchor_concepts", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_schema_anchor_concepts_json(e.handle, out, errBuf)
	})
}

func (e *nativeEngine) rawGraphIngestNative(_ context.Context, plugin, artifact string, data []byte) ([]byte, error) {
	return bufferResult("graph.ingest_native", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_ingest_native_json(e.handle, stringSlice(plugin), stringSlice(artifact), byteSlice(data), out, errBuf)
		runtime.KeepAlive(plugin)
		runtime.KeepAlive(artifact)
		runtime.KeepAlive(data)
		return rc
	})
}

func (e *nativeEngine) rawGraphAddNodes(_ context.Context, nodes []byte) (uint64, error) {
	return countResult("graph.add_nodes_json", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_add_nodes_json(e.handle, byteSlice(nodes), out, errBuf)
		runtime.KeepAlive(nodes)
		return rc
	})
}

func (e *nativeEngine) rawGraphAddEdges(_ context.Context, edges []byte) (uint64, error) {
	return countResult("graph.add_edges_json", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_add_edges_json(e.handle, byteSlice(edges), out, errBuf)
		runtime.KeepAlive(edges)
		return rc
	})
}

func (e *nativeEngine) rawGraphFindNode(_ context.Context, identifier string) ([]byte, error) {
	return bufferResult("graph.find_node", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_find_node_json(e.handle, stringSlice(identifier), out, errBuf)
		runtime.KeepAlive(identifier)
		return rc
	})
}

func (e *nativeEngine) rawGraphNodeTypes(_ context.Context) ([]byte, error) {
	return bufferResult("graph.node_types", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_graph_node_types_json(e.handle, out, errBuf)
	})
}

func (e *nativeEngine) rawGraphNodesPage(_ context.Context, request []byte) ([]byte, error) {
	return bufferResult("graph.nodes_page", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_nodes_page_json(e.handle, byteSlice(request), out, errBuf)
		runtime.KeepAlive(request)
		return rc
	})
}

func (e *nativeEngine) rawGraphLink(_ context.Context, nodeIDs []byte, dataSource string) ([]byte, error) {
	return bufferResult("graph.link", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_link_json(e.handle, byteSlice(nodeIDs), stringSlice(dataSource), out, errBuf)
		runtime.KeepAlive(nodeIDs)
		runtime.KeepAlive(dataSource)
		return rc
	})
}

func (e *nativeEngine) rawRepositoryIndexGraph(_ context.Context) ([]byte, error) {
	return bufferResult("repo.index_graph", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_repo_index_graph_json(e.handle, out, errBuf)
	})
}

type nativeRawRAGIndexSession struct{ handle *C.CstxRagIndexSession }

func (e *nativeEngine) rawRAGIndex(_ context.Context, request []byte) (rawRAGIndexSession, error) {
	var handle *C.CstxRagIndexSession
	err := statusCall("rag.index", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_rag_index(e.handle, byteSlice(request), &handle, errBuf)
		runtime.KeepAlive(request)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return &nativeRawRAGIndexSession{handle: handle}, nil
}

func (s *nativeRawRAGIndexSession) metadata(_ context.Context) ([]byte, error) {
	return bufferResult("rag.index.metadata", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_rag_index_session_metadata_json(s.handle, out, errBuf)
	})
}

func (s *nativeRawRAGIndexSession) pending(_ context.Context, offset, limit int) ([]byte, error) {
	if offset < 0 || limit < 0 {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "rag.index.pending", Message: "offset and limit must be non-negative"}
	}
	return bufferResult("rag.index.pending", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_rag_index_session_pending_json(s.handle, C.size_t(offset), C.size_t(limit), out, errBuf)
	})
}

func (s *nativeRawRAGIndexSession) deletes(_ context.Context) ([]byte, error) {
	return bufferResult("rag.index.deletes", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_rag_index_session_deletes_json(s.handle, out, errBuf)
	})
}

func (s *nativeRawRAGIndexSession) close() {
	if s.handle != nil {
		C.cstx_rag_index_session_close(s.handle)
		C.cstx_rag_index_session_free(s.handle)
		s.handle = nil
	}
}

type nativeRawRAGRetrieval struct{ handle *C.CstxRagRetrieval }

func (e *nativeEngine) rawRAGRetrieve(_ context.Context, query []byte) (rawRAGRetrieval, error) {
	var handle *C.CstxRagRetrieval
	err := statusCall("rag.retrieve", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_rag_retrieve(e.handle, byteSlice(query), &handle, errBuf)
		runtime.KeepAlive(query)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return &nativeRawRAGRetrieval{handle: handle}, nil
}

func (r *nativeRawRAGRetrieval) requests(_ context.Context) ([]byte, error) {
	return bufferResult("rag.retrieve.requests", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_rag_retrieval_requests_json(r.handle, out, errBuf)
	})
}

func (r *nativeRawRAGRetrieval) complete(_ context.Context, batches []byte) ([]byte, error) {
	result, err := bufferResult("rag.retrieve.complete", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_rag_retrieval_complete_json(r.handle, byteSlice(batches), out, errBuf)
		runtime.KeepAlive(batches)
		return rc
	})
	if err == nil {
		C.cstx_rag_retrieval_free(r.handle)
		r.handle = nil
	}
	return result, err
}

func (r *nativeRawRAGRetrieval) close() {
	if r.handle != nil {
		C.cstx_rag_retrieval_close(r.handle)
		C.cstx_rag_retrieval_free(r.handle)
		r.handle = nil
	}
}
