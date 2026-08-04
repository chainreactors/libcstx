package cstx

/*
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib/linux_amd64 -lcstx_ffi -lm -ldl -lpthread
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/lib/linux_arm64 -lcstx_ffi -lm -ldl -lpthread
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin_amd64 -lcstx_ffi -lm -framework Security -framework CoreFoundation
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -lcstx_ffi -lm -framework Security -framework CoreFoundation
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/lib/windows_amd64 -lcstx_ffi -lws2_32 -luserenv -lbcrypt -lntdll -ladvapi32
#include "cstx_ffi.h"
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"encoding/json"
	"runtime"
	"unsafe"
)

type nativeEngine struct {
	handle *C.CstxHandle
}

func newEngine(config Config) (engine, error) {
	payload, err := json.Marshal(map[string]any{
		"project_id":       config.ProjectID,
		"cursor_page_size": config.CursorPageSize,
	})
	if err != nil {
		return nil, err
	}
	var handle *C.CstxHandle
	err = statusCall("cstx.open", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_open(byteSlice(payload), &handle, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return nil, err
	}
	eng := &nativeEngine{handle: handle}
	runtime.SetFinalizer(eng, (*nativeEngine).finalize)
	return eng, nil
}

func (e *nativeEngine) finalize() { _ = e.close() }

func (e *nativeEngine) close() error {
	if e.handle != nil {
		C.cstx_free(e.handle)
		e.handle = nil
		runtime.SetFinalizer(e, nil)
	}
	return nil
}

func (e *nativeEngine) graphSubgraph(_ context.Context, seedIDs []string, depth uint32) (engine, error) {
	payload, err := json.Marshal(seedIDs)
	if err != nil {
		return nil, err
	}
	var handle *C.CstxHandle
	err = statusCall("graph.subgraph", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_subgraph(e.handle, byteSlice(payload), C.uint32_t(depth), &handle, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return nil, err
	}
	derived := &nativeEngine{handle: handle}
	runtime.SetFinalizer(derived, (*nativeEngine).finalize)
	return derived, nil
}

// --- C transport helpers -------------------------------------------------

var emptySliceByte byte

func byteSlice(value []byte) C.CstxSlice {
	if len(value) == 0 {
		return C.CstxSlice{}
	}
	return C.CstxSlice{data: (*C.uint8_t)(unsafe.Pointer(&value[0])), len: C.size_t(len(value))}
}

func stringSlice(value string) C.CstxSlice {
	if len(value) == 0 {
		return C.CstxSlice{}
	}
	return C.CstxSlice{data: (*C.uint8_t)(unsafe.Pointer(unsafe.StringData(value))), len: C.size_t(len(value))}
}

// optionalStringSlice maps "" to a null slice (Rust None). A non-empty value
// is borrowed for the duration of the call.
func optionalStringSlice(value string) C.CstxSlice {
	if value == "" {
		return C.CstxSlice{}
	}
	return stringSlice(value)
}

func statusCodeFallback(code C.CstxStatusCode) Code {
	switch code {
	case C.CSTX_INVALID_ARGUMENT:
		return CodeInvalidArgument
	case C.CSTX_NOT_FOUND:
		return CodeNotFound
	case C.CSTX_CONFLICT:
		return CodeConflict
	case C.CSTX_NOT_INITIALIZED:
		return CodeNotInitialized
	case C.CSTX_UNSUPPORTED:
		return CodeUnsupported
	case C.CSTX_VALIDATION:
		return CodeValidation
	case C.CSTX_PARSE:
		return CodeParse
	case C.CSTX_IO:
		return CodeIO
	case C.CSTX_CORRUPT_DATA:
		return CodeCorruptData
	case C.CSTX_CURSOR_INVALIDATED:
		return CodeCursorInvalidated
	case C.CSTX_STALE_OPERATION:
		return CodeStaleOperation
	case C.CSTX_STALE_RECALL:
		return CodeStaleRecall
	default:
		return CodeInternal
	}
}

func statusError(rc C.CstxStatusCode, op string, errBuf *C.CstxBuffer) error {
	defer C.cstx_buffer_free(errBuf)
	if rc == C.CSTX_OK {
		return nil
	}
	cerr := parseError(takeBuffer(errBuf), statusCodeFallback(rc))
	if cerr.Operation == "" {
		cerr.Operation = op
	}
	return cerr
}

func takeBuffer(buffer *C.CstxBuffer) []byte {
	defer C.cstx_buffer_free(buffer)
	if buffer.data == nil || buffer.len == 0 {
		return []byte{}
	}
	source := unsafe.Slice((*byte)(unsafe.Pointer(buffer.data)), int(buffer.len))
	result := make([]byte, len(source))
	copy(result, source)
	return result
}

func statusCall(op string, call func(errBuf *C.CstxBuffer) C.CstxStatusCode) error {
	var errBuf C.CstxBuffer
	return statusError(call(&errBuf), op, &errBuf)
}

// jsonResult runs a call whose success output is a JSON buffer and decodes it
// into result.
func jsonResult(op string, result any, call func(out, errBuf *C.CstxBuffer) C.CstxStatusCode) error {
	var out, errBuf C.CstxBuffer
	if err := statusError(call(&out, &errBuf), op, &errBuf); err != nil {
		C.cstx_buffer_free(&out)
		return err
	}
	return json.Unmarshal(takeBuffer(&out), result)
}

func bufferResult(op string, call func(out, errBuf *C.CstxBuffer) C.CstxStatusCode) ([]byte, error) {
	var out, errBuf C.CstxBuffer
	if err := statusError(call(&out, &errBuf), op, &errBuf); err != nil {
		C.cstx_buffer_free(&out)
		return nil, err
	}
	return takeBuffer(&out), nil
}

func countResult(op string, call func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode) (uint64, error) {
	var out C.uint64_t
	var errBuf C.CstxBuffer
	if err := statusError(call(&out, &errBuf), op, &errBuf); err != nil {
		return 0, err
	}
	return uint64(out), nil
}

func boolResult(op string, call func(out *C.uint8_t, errBuf *C.CstxBuffer) C.CstxStatusCode) (bool, error) {
	var out C.uint8_t
	var errBuf C.CstxBuffer
	if err := statusError(call(&out, &errBuf), op, &errBuf); err != nil {
		return false, err
	}
	return out != 0, nil
}

func marshal(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		// All marshaled types are internal contracts; a failure here is a bug.
		panic("cstx: marshal transport value: " + err.Error())
	}
	return data
}

func marshalInput(op string, value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: op, Message: err.Error()}
	}
	return data, nil
}

// --- runtime -------------------------------------------------------------

func (e *nativeEngine) lastChange(_ context.Context) (ChangeSet, error) {
	var change ChangeSet
	err := jsonResult("cstx.last_change", &change, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_last_change_json(e.handle, out, errBuf)
	})
	return change, err
}

// --- schemas -------------------------------------------------------------

func (e *nativeEngine) schemaRegister(_ context.Context, nodeType string, schema map[string]any, valueField string) error {
	payload, err := marshalInput("schemas.register", schema)
	if err != nil {
		return err
	}
	return statusCall("schemas.register", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_schema_register(e.handle, stringSlice(nodeType), byteSlice(payload), optionalStringSlice(valueField), errBuf)
		runtime.KeepAlive(nodeType)
		runtime.KeepAlive(payload)
		runtime.KeepAlive(valueField)
		return rc
	})
}

func (e *nativeEngine) schemaContains(_ context.Context, nodeType string) (bool, error) {
	return boolResult("schemas.contains", func(out *C.uint8_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_schema_contains(e.handle, stringSlice(nodeType), out, errBuf)
		runtime.KeepAlive(nodeType)
		return rc
	})
}

func (e *nativeEngine) schemaGet(_ context.Context, nodeType string) (map[string]any, error) {
	var schema map[string]any
	err := jsonResult("schemas.get", &schema, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_schema_get_json(e.handle, stringSlice(nodeType), out, errBuf)
		runtime.KeepAlive(nodeType)
		return rc
	})
	return schema, err
}

func (e *nativeEngine) schemaList(_ context.Context) ([]map[string]any, error) {
	var schemas []map[string]any
	err := jsonResult("schemas.list", &schemas, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_schema_list_json(e.handle, out, errBuf)
	})
	return schemas, err
}

func (e *nativeEngine) schemaLoadPlugin(_ context.Context, name string) error {
	return statusCall("schemas.load_plugin", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_schema_load_plugin(e.handle, stringSlice(name), errBuf)
		runtime.KeepAlive(name)
		return rc
	})
}

func (e *nativeEngine) schemaLoadAllPlugins(_ context.Context) error {
	return statusCall("schemas.load_all_plugins", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_schema_load_all_plugins(e.handle, errBuf)
	})
}

func (e *nativeEngine) schemaAvailablePlugins(_ context.Context) ([]string, error) {
	var plugins []string
	err := jsonResult("schemas.available_plugins", &plugins, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_schema_available_plugins_json(e.handle, out, errBuf)
	})
	return plugins, err
}

func (e *nativeEngine) schemaPluginArtifacts(_ context.Context, name string) ([]string, error) {
	var artifacts []string
	err := jsonResult("schemas.plugin_artifacts", &artifacts, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_schema_plugin_artifacts_json(e.handle, stringSlice(name), out, errBuf)
		runtime.KeepAlive(name)
		return rc
	})
	return artifacts, err
}

// --- graph ---------------------------------------------------------------

func (e *nativeEngine) graphAddNodes(_ context.Context, nodes []Node) (uint64, error) {
	payload := marshal(nodes)
	return countResult("graph.add_nodes", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_add_nodes(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
}

func (e *nativeEngine) graphAddEdges(_ context.Context, edges []Edge) (uint64, error) {
	payload := marshal(edges)
	return countResult("graph.add_edges", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_add_edges(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
}

func (e *nativeEngine) graphIngest(_ context.Context, source string, data []byte) (uint64, error) {
	return countResult("graph.ingest", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_ingest(e.handle, stringSlice(source), byteSlice(data), out, errBuf)
		runtime.KeepAlive(source)
		runtime.KeepAlive(data)
		return rc
	})
}

func (e *nativeEngine) graphNode(_ context.Context, nodeID string) (Node, error) {
	var node Node
	err := jsonResult("graph.node", &node, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_node(e.handle, stringSlice(nodeID), out, errBuf)
		runtime.KeepAlive(nodeID)
		return rc
	})
	return node, err
}

func (e *nativeEngine) graphContains(_ context.Context, nodeID string) (bool, error) {
	return boolResult("graph.contains", func(out *C.uint8_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_contains(e.handle, stringSlice(nodeID), out, errBuf)
		runtime.KeepAlive(nodeID)
		return rc
	})
}

func (e *nativeEngine) graphNodeCount(_ context.Context) (uint64, error) {
	return countResult("graph.node_count", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_graph_node_count(e.handle, out, errBuf)
	})
}

func (e *nativeEngine) graphEdgeCount(_ context.Context) (uint64, error) {
	return countResult("graph.edge_count", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_graph_edge_count(e.handle, out, errBuf)
	})
}

func (e *nativeEngine) graphStats(_ context.Context) (GraphStats, error) {
	var stats GraphStats
	err := jsonResult("graph.stats", &stats, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_graph_stats(e.handle, out, errBuf)
	})
	return stats, err
}

func (e *nativeEngine) graphNodes(_ context.Context, filter NodeFilter, options CollectionOptions) (nodeIter, error) {
	filterJSON := marshal(filter)
	optionsJSON := marshal(options)
	var iter *C.CstxNodeIterator
	err := statusCall("graph.nodes", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_nodes(e.handle, byteSlice(filterJSON), byteSlice(optionsJSON), &iter, errBuf)
		runtime.KeepAlive(filterJSON)
		runtime.KeepAlive(optionsJSON)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNodeIter(iter), nil
}

func (e *nativeEngine) graphEdges(_ context.Context, filter EdgeFilter, options CollectionOptions) (edgeIter, error) {
	filterJSON := marshal(filter)
	optionsJSON := marshal(options)
	var iter *C.CstxEdgeIterator
	err := statusCall("graph.edges", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_edges(e.handle, byteSlice(filterJSON), byteSlice(optionsJSON), &iter, errBuf)
		runtime.KeepAlive(filterJSON)
		runtime.KeepAlive(optionsJSON)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newEdgeIter(iter), nil
}

func (e *nativeEngine) graphNeighbors(_ context.Context, nodeID, direction string, options CollectionOptions) (nodeIter, error) {
	optionsJSON := marshal(options)
	var iter *C.CstxNodeIterator
	err := statusCall("graph.neighbors", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_neighbors(e.handle, stringSlice(nodeID), stringSlice(direction), byteSlice(optionsJSON), &iter, errBuf)
		runtime.KeepAlive(nodeID)
		runtime.KeepAlive(direction)
		runtime.KeepAlive(optionsJSON)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNodeIter(iter), nil
}

func (e *nativeEngine) graphQuery(_ context.Context, expression string, options QueryOptions) (nodeIter, error) {
	optionsJSON := marshal(options)
	var iter *C.CstxNodeIterator
	err := statusCall("graph.query", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_query(e.handle, stringSlice(expression), byteSlice(optionsJSON), &iter, errBuf)
		runtime.KeepAlive(expression)
		runtime.KeepAlive(optionsJSON)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNodeIter(iter), nil
}

func (e *nativeEngine) graphDifference(_ context.Context, other engine, nodeType string) (engine, error) {
	right, ok := other.(*nativeEngine)
	if !ok || right == nil || right.handle == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "graph.difference", Message: "other runtime is invalid"}
	}
	var output *C.CstxHandle
	err := statusCall("graph.difference", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_difference(e.handle, right.handle, stringSlice(nodeType), &output, errBuf)
		runtime.KeepAlive(nodeType)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return &nativeEngine{handle: output}, nil
}

// --- repository ----------------------------------------------------------

func (e *nativeEngine) repoCommit(_ context.Context, refName, message string, metadata any) (Commit, error) {
	var metadataJSON []byte
	if metadata != nil {
		var err error
		metadataJSON, err = marshalInput("repo.commit", metadata)
		if err != nil {
			return Commit{}, err
		}
	}
	var commit Commit
	err := jsonResult("repo.commit", &commit, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_commit(e.handle, stringSlice(refName), stringSlice(message), byteSlice(metadataJSON), out, errBuf)
		runtime.KeepAlive(refName)
		runtime.KeepAlive(message)
		runtime.KeepAlive(metadataJSON)
		return rc
	})
	return commit, err
}

func (e *nativeEngine) repoDiff(_ context.Context, baseRef, headRef string) (GraphDiff, error) {
	var diff GraphDiff
	err := jsonResult("repo.diff", &diff, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_diff(e.handle, stringSlice(baseRef), stringSlice(headRef), out, errBuf)
		runtime.KeepAlive(baseRef)
		runtime.KeepAlive(headRef)
		return rc
	})
	return diff, err
}

func (e *nativeEngine) repoDump(_ context.Context, compression string) ([]byte, error) {
	var out, errBuf C.CstxBuffer
	rc := C.cstx_repo_dump(e.handle, stringSlice(compression), &out, &errBuf)
	runtime.KeepAlive(compression)
	if err := statusError(rc, "repo.dump", &errBuf); err != nil {
		C.cstx_buffer_free(&out)
		return nil, err
	}
	return takeBuffer(&out), nil
}

func (e *nativeEngine) repoLoad(_ context.Context, data []byte, compression string) (uint64, error) {
	return countResult("repo.load", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_load(e.handle, byteSlice(data), stringSlice(compression), out, errBuf)
		runtime.KeepAlive(data)
		runtime.KeepAlive(compression)
		return rc
	})
}

func (e *nativeEngine) repoDumpJSON(_ context.Context) ([]byte, error) {
	var out, errBuf C.CstxBuffer
	if err := statusError(C.cstx_repo_dump_json(e.handle, &out, &errBuf), "repo.dump_json", &errBuf); err != nil {
		C.cstx_buffer_free(&out)
		return nil, err
	}
	return takeBuffer(&out), nil
}

func (e *nativeEngine) repoLoadJSON(_ context.Context, data []byte) (uint64, error) {
	return countResult("repo.load_json", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_load_json(e.handle, byteSlice(data), out, errBuf)
		runtime.KeepAlive(data)
		return rc
	})
}

func (e *nativeEngine) repoSnapshotFingerprint(_ context.Context) (string, error) {
	var fingerprint string
	err := jsonResult("repo.snapshot_fingerprint", &fingerprint, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_repo_snapshot_fingerprint(e.handle, out, errBuf)
	})
	return fingerprint, err
}

func (e *nativeEngine) repoHead(_ context.Context, refName string) (*string, error) {
	var head *string
	err := jsonResult("repo.head", &head, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_head(e.handle, stringSlice(refName), out, errBuf)
		runtime.KeepAlive(refName)
		return rc
	})
	return head, err
}

func (e *nativeEngine) repoRefs(_ context.Context) ([]Ref, error) {
	var refs []Ref
	err := jsonResult("repo.refs", &refs, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_repo_refs(e.handle, out, errBuf)
	})
	return refs, err
}

func (e *nativeEngine) repoCommitDelta(_ context.Context, request CommitDeltaRequest) (Commit, error) {
	parent, err := e.repoHead(context.Background(), request.RefName)
	if err != nil {
		return Commit{}, err
	}
	requestJSON := marshal(request)
	var result struct {
		CommitHash string `json:"commit_hash"`
		TreeHash   string `json:"tree_hash"`
		Stats      any    `json:"stats"`
	}
	err = jsonResult("repo.commit_delta", &result, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_commit_delta_json(e.handle, byteSlice(requestJSON), out, errBuf)
		runtime.KeepAlive(requestJSON)
		return rc
	})
	if err != nil {
		return Commit{}, err
	}
	parentHash := ""
	if parent != nil {
		parentHash = *parent
	}
	return Commit{
		ID: result.CommitHash, Tree: result.TreeHash, Parent: parentHash,
		Message: request.Message, Metadata: request.Metadata, Stats: result.Stats,
		CreatedAt: request.CreatedAt,
	}, nil
}

func (e *nativeEngine) repoSetRef(_ context.Context, refName, commitHash string) error {
	return statusCall("repo.set_ref", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_set_ref(e.handle, stringSlice(refName), stringSlice(commitHash), errBuf)
		runtime.KeepAlive(refName)
		runtime.KeepAlive(commitHash)
		return rc
	})
}

func (e *nativeEngine) repoImportCommit(_ context.Context, hash string, data json.RawMessage) error {
	return statusCall("repo.import_commit", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_import_commit(e.handle, stringSlice(hash), byteSlice(data), errBuf)
		runtime.KeepAlive(hash)
		runtime.KeepAlive(data)
		return rc
	})
}

func (e *nativeEngine) repoExportCommit(_ context.Context, hash string) (json.RawMessage, error) {
	return bufferResult("repo.export_commit", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_export_commit_json(e.handle, stringSlice(hash), out, errBuf)
		runtime.KeepAlive(hash)
		return rc
	})
}

func (e *nativeEngine) repoExportTreeObjects(_ context.Context, rootHash string) (json.RawMessage, error) {
	return bufferResult("repo.export_tree_objects", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_export_tree_objects_json(e.handle, stringSlice(rootHash), out, errBuf)
		runtime.KeepAlive(rootHash)
		return rc
	})
}

func (e *nativeEngine) repoImportTreeObjects(_ context.Context, objects json.RawMessage) error {
	return statusCall("repo.import_tree_objects", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_import_tree_objects(e.handle, byteSlice(objects), errBuf)
		runtime.KeepAlive(objects)
		return rc
	})
}

func (e *nativeEngine) repoTreeEntries(_ context.Context, rootHash string, types []string) (map[string]map[string]string, error) {
	typesJSON := marshal(orEmpty(types))
	var entries map[string]map[string]string
	err := jsonResult("repo.tree_entries", &entries, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_tree_entries_json(e.handle, stringSlice(rootHash), byteSlice(typesJSON), out, errBuf)
		runtime.KeepAlive(rootHash)
		runtime.KeepAlive(typesJSON)
		return rc
	})
	return entries, err
}

func (e *nativeEngine) repoFindTreeEntry(_ context.Context, rootHash, elementID string) (*string, error) {
	var hash *string
	err := jsonResult("repo.find_tree_entry", &hash, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_find_tree_entry_json(e.handle, stringSlice(rootHash), stringSlice(elementID), out, errBuf)
		runtime.KeepAlive(rootHash)
		runtime.KeepAlive(elementID)
		return rc
	})
	return hash, err
}

func (e *nativeEngine) repoDiffTreeEntries(_ context.Context, baseRoot, headRoot string) (TreeEntryDiff, error) {
	var diff TreeEntryDiff
	err := jsonResult("repo.diff_tree_entries", &diff, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_diff_tree_entries_json(e.handle, stringSlice(baseRoot), stringSlice(headRoot), out, errBuf)
		runtime.KeepAlive(baseRoot)
		runtime.KeepAlive(headRoot)
		return rc
	})
	return diff, err
}

func (e *nativeEngine) repoTreeRootStats(_ context.Context, rootHash string) (map[string]int64, error) {
	var stats map[string]int64
	err := jsonResult("repo.tree_root_stats", &stats, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_tree_root_stats_json(e.handle, stringSlice(rootHash), out, errBuf)
		runtime.KeepAlive(rootHash)
		return rc
	})
	return stats, err
}

// --- iterators -----------------------------------------------------------
type nativeNodeIter struct{ iter *C.CstxNodeIterator }

func newNodeIter(iter *C.CstxNodeIterator) *nativeNodeIter {
	it := &nativeNodeIter{iter: iter}
	runtime.SetFinalizer(it, (*nativeNodeIter).close)
	return it
}

func (i *nativeNodeIter) next(_ context.Context) (Node, bool, error) {
	var node Node
	if i.iter == nil {
		return node, false, nil
	}
	var out, errBuf C.CstxBuffer
	var hasValue C.uint8_t
	if err := statusError(C.cstx_node_iterator_next(i.iter, &out, &hasValue, &errBuf), "cursor.next", &errBuf); err != nil {
		C.cstx_buffer_free(&out)
		return node, false, err
	}
	if hasValue == 0 {
		C.cstx_buffer_free(&out)
		return node, false, nil
	}
	if err := json.Unmarshal(takeBuffer(&out), &node); err != nil {
		return node, false, err
	}
	return node, true, nil
}

func (i *nativeNodeIter) close() {
	if i.iter != nil {
		C.cstx_node_iterator_free(i.iter)
		i.iter = nil
		runtime.SetFinalizer(i, nil)
	}
}

type nativeEdgeIter struct{ iter *C.CstxEdgeIterator }

func newEdgeIter(iter *C.CstxEdgeIterator) *nativeEdgeIter {
	it := &nativeEdgeIter{iter: iter}
	runtime.SetFinalizer(it, (*nativeEdgeIter).close)
	return it
}

func (i *nativeEdgeIter) next(_ context.Context) (Edge, bool, error) {
	var edge Edge
	if i.iter == nil {
		return edge, false, nil
	}
	var out, errBuf C.CstxBuffer
	var hasValue C.uint8_t
	if err := statusError(C.cstx_edge_iterator_next(i.iter, &out, &hasValue, &errBuf), "cursor.next", &errBuf); err != nil {
		C.cstx_buffer_free(&out)
		return edge, false, err
	}
	if hasValue == 0 {
		C.cstx_buffer_free(&out)
		return edge, false, nil
	}
	if err := json.Unmarshal(takeBuffer(&out), &edge); err != nil {
		return edge, false, err
	}
	return edge, true, nil
}

func (i *nativeEdgeIter) close() {
	if i.iter != nil {
		C.cstx_edge_iterator_free(i.iter)
		i.iter = nil
		runtime.SetFinalizer(i, nil)
	}
}
