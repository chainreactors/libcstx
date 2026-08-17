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
	"encoding/hex"
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

func (e *nativeEngine) graphDeleteNodes(_ context.Context, nodeIDs []string) (uint64, error) {
	payload, err := marshalInput("graph.delete_nodes", nodeIDs)
	if err != nil {
		return 0, err
	}
	return countResult("graph.delete_nodes", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_delete_nodes(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
}

func (e *nativeEngine) graphDeleteEdges(_ context.Context, edgeIDs []string) (uint64, error) {
	payload, err := marshalInput("graph.delete_edges", edgeIDs)
	if err != nil {
		return 0, err
	}
	return countResult("graph.delete_edges", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_delete_edges(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
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

func boolByte(value bool) uint8 {
	if value {
		return 1
	}
	return 0
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

func (e *nativeEngine) schemaImport(_ context.Context, contract SchemaContract) error {
	payload, err := marshalInput("schemas.import_schema", contract)
	if err != nil {
		return err
	}
	return statusCall("schemas.import_schema", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_schema_import_schema(e.handle, byteSlice(payload), errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
}

func (e *nativeEngine) schemaExport(_ context.Context) (SchemaContract, error) {
	var contract SchemaContract
	err := jsonResult("schemas.export_schema", &contract, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_schema_export_schema_json(e.handle, out, errBuf)
	})
	return contract, err
}

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

func (e *nativeEngine) schemaRegisterJoinRule(_ context.Context, rule JoinRuleSpec) error {
	payload, err := marshalInput("schemas.register_join_rule", rule)
	if err != nil {
		return err
	}
	return statusCall("schemas.register_join_rule", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_schema_register_join_rule(e.handle, byteSlice(payload), errBuf)
		runtime.KeepAlive(payload)
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

func (e *nativeEngine) schemaHasNativeArtifact(_ context.Context, artifact string) (bool, error) {
	return boolResult("schemas.has_native_artifact", func(out *C.uint8_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_schema_has_native_artifact(e.handle, stringSlice(artifact), out, errBuf)
		runtime.KeepAlive(artifact)
		return rc
	})
}

func (e *nativeEngine) schemaAnchorConcepts(_ context.Context) ([]AnchorConcept, error) {
	var concepts []AnchorConcept
	err := jsonResult("schemas.anchor_concepts", &concepts, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_schema_anchor_concepts_json(e.handle, out, errBuf)
	})
	return concepts, err
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

func (e *nativeEngine) graphReplaceNodes(_ context.Context, nodes []Node) (uint64, error) {
	payload := marshal(nodes)
	return countResult("graph.replace_nodes", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_replace_nodes(e.handle, byteSlice(payload), out, errBuf)
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
		return C.cstx_graph_stats(e.handle, 0, 0, out, errBuf)
	})
	return stats, err
}

func (e *nativeEngine) graphNodes(_ context.Context, filter NodeFilter, options CollectionOptions) (graphCursor, error) {
	filterJSON := marshal(filter)
	optionsJSON := marshal(options)
	var cursor *C.CstxGraphCursor
	err := statusCall("graph.nodes", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_nodes(e.handle, byteSlice(filterJSON), byteSlice(optionsJSON), &cursor, errBuf)
		runtime.KeepAlive(filterJSON)
		runtime.KeepAlive(optionsJSON)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNativeGraphCursor(cursor), nil
}

func (e *nativeEngine) graphEdges(_ context.Context, filter EdgeFilter, options CollectionOptions) (graphCursor, error) {
	filterJSON := marshal(filter)
	optionsJSON := marshal(options)
	var cursor *C.CstxGraphCursor
	err := statusCall("graph.edges", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_edges(e.handle, byteSlice(filterJSON), byteSlice(optionsJSON), &cursor, errBuf)
		runtime.KeepAlive(filterJSON)
		runtime.KeepAlive(optionsJSON)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNativeGraphCursor(cursor), nil
}

func (e *nativeEngine) graphNeighbors(_ context.Context, nodeID, direction string, options CollectionOptions) (graphCursor, error) {
	optionsJSON := marshal(options)
	var cursor *C.CstxGraphCursor
	err := statusCall("graph.neighbors", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_neighbors(e.handle, stringSlice(nodeID), stringSlice(direction), byteSlice(optionsJSON), &cursor, errBuf)
		runtime.KeepAlive(nodeID)
		runtime.KeepAlive(direction)
		runtime.KeepAlive(optionsJSON)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNativeGraphCursor(cursor), nil
}

func (e *nativeEngine) graphQuery(_ context.Context, expression string, options QueryOptions) (graphCursor, error) {
	optionsJSON := marshal(options)
	var cursor *C.CstxGraphCursor
	err := statusCall("graph.query", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_query(e.handle, stringSlice(expression), byteSlice(optionsJSON), &cursor, errBuf)
		runtime.KeepAlive(expression)
		runtime.KeepAlive(optionsJSON)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNativeGraphCursor(cursor), nil
}

func (e *nativeEngine) graphAnalyze(_ context.Context, algorithm any, selection *string) (uint8, bool, graphCursor, error) {
	payload, err := marshalInput("graph.analyze", algorithm)
	if err != nil {
		return 0, false, nil, err
	}
	var kind C.uint8_t
	var boolean C.uint8_t
	var cursor *C.CstxGraphCursor
	var selectionSlice C.CstxSlice
	if selection != nil {
		selectionSlice = stringSlice(*selection)
	}
	err = statusCall("graph.analyze", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_analyze(
			e.handle,
			byteSlice(payload),
			selectionSlice,
			&kind,
			&boolean,
			&cursor,
			errBuf,
		)
		runtime.KeepAlive(payload)
		runtime.KeepAlive(selection)
		return rc
	})
	if err != nil {
		return 0, false, nil, err
	}
	if cursor == nil {
		return uint8(kind), boolean != 0, nil, nil
	}
	return uint8(kind), boolean != 0, newNativeGraphCursor(cursor), nil
}

// --- repository ----------------------------------------------------------

func (e *nativeEngine) repoResolve(_ context.Context, revision string) (string, error) {
	var resolved string
	err := jsonResult("repo.resolve", &resolved, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_resolve(e.handle, stringSlice(revision), out, errBuf)
		runtime.KeepAlive(revision)
		return rc
	})
	return resolved, err
}

func (e *nativeEngine) repoCheckout(_ context.Context, revision string, force bool) (Commit, error) {
	var commit Commit
	var nativeForce C.uint8_t
	if force {
		nativeForce = 1
	}
	err := jsonResult("repo.checkout", &commit, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_checkout(e.handle, stringSlice(revision), nativeForce, out, errBuf)
		runtime.KeepAlive(revision)
		return rc
	})
	return commit, err
}

func (e *nativeEngine) repoCommit(
	_ context.Context,
	message string,
	refName string,
	expectedHead *string,
	metadata any,
) (Commit, error) {
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
		var expected C.CstxSlice
		if expectedHead != nil {
			expected = stringSlice(*expectedHead)
		}
		rc := C.cstx_repo_commit(e.handle, stringSlice(message), stringSlice(refName), expected, byteSlice(metadataJSON), out, errBuf)
		runtime.KeepAlive(refName)
		runtime.KeepAlive(message)
		runtime.KeepAlive(expectedHead)
		runtime.KeepAlive(metadataJSON)
		return rc
	})
	return commit, err
}

type preparedObjectWire struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Envelope string `json:"envelope"`
}

type preparedCommitWire struct {
	Commit    Commit               `json:"commit"`
	IndexRoot string               `json:"index_root"`
	Objects   []preparedObjectWire `json:"objects"`
}

func (e *nativeEngine) repoPrepare(
	_ context.Context,
	message string,
	refName string,
	expectedHead *string,
	metadata any,
	timestamp *int64,
) (PreparedCommit, error) {
	metadataJSON, err := marshalInput("repo.prepare", metadata)
	if err != nil {
		return PreparedCommit{}, err
	}
	var wire preparedCommitWire
	var nativeTimestamp C.int64_t
	var hasTimestamp C.uint8_t
	if timestamp != nil {
		nativeTimestamp = C.int64_t(*timestamp)
		hasTimestamp = 1
	}
	err = jsonResult("repo.prepare", &wire, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		var expected C.CstxSlice
		if expectedHead != nil {
			expected = stringSlice(*expectedHead)
		}
		rc := C.cstx_repo_prepare(e.handle, stringSlice(message), stringSlice(refName), expected, byteSlice(metadataJSON), nativeTimestamp, hasTimestamp, out, errBuf)
		runtime.KeepAlive(message)
		runtime.KeepAlive(refName)
		runtime.KeepAlive(expectedHead)
		runtime.KeepAlive(metadataJSON)
		return rc
	})
	if err != nil {
		return PreparedCommit{}, err
	}
	prepared := PreparedCommit{Commit: wire.Commit, IndexRoot: wire.IndexRoot, Objects: make([]PreparedObject, len(wire.Objects))}
	for i, object := range wire.Objects {
		envelope, err := hex.DecodeString(object.Envelope)
		if err != nil {
			return PreparedCommit{}, &Error{Code: CodeCorruptData, Operation: "repo.prepare", Message: "invalid object envelope: " + err.Error()}
		}
		prepared.Objects[i] = PreparedObject{ID: object.ID, Kind: object.Kind, Envelope: envelope}
	}
	return prepared, nil
}

func (e *nativeEngine) repoAccept(_ context.Context, commit string) error {
	return statusCall("repo.accept", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_accept(e.handle, stringSlice(commit), errBuf)
		runtime.KeepAlive(commit)
		return rc
	})
}

func (e *nativeEngine) repoDiscard(_ context.Context) error {
	return statusCall("repo.discard", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_repo_discard(e.handle, errBuf)
	})
}

func (e *nativeEngine) repoSynchronize(_ context.Context, state RepositorySync) error {
	type objectWire struct {
		ID       string `json:"id"`
		Envelope string `json:"envelope"`
	}
	type refWire struct {
		Name   string  `json:"name"`
		Commit *string `json:"commit"`
	}
	type indexWire struct {
		Commit    string `json:"commit"`
		IndexRoot string `json:"index_root"`
	}
	payload := struct {
		Objects []objectWire `json:"objects"`
		Refs    []refWire    `json:"refs"`
		Indexes []indexWire  `json:"indexes"`
	}{
		Objects: make([]objectWire, len(state.Objects)),
		Refs:    make([]refWire, len(state.Refs)),
		Indexes: make([]indexWire, len(state.Indexes)),
	}
	for i, object := range state.Objects {
		payload.Objects[i] = objectWire{ID: object.ID, Envelope: hex.EncodeToString(object.Envelope)}
	}
	for i, ref := range state.Refs {
		payload.Refs[i] = refWire{Name: ref.Name, Commit: ref.Commit}
	}
	for i, index := range state.Indexes {
		payload.Indexes[i] = indexWire{Commit: index.Commit, IndexRoot: index.IndexRoot}
	}
	data, err := marshalInput("repo.synchronize", payload)
	if err != nil {
		return err
	}
	return statusCall("repo.synchronize", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_synchronize(e.handle, byteSlice(data), errBuf)
		runtime.KeepAlive(data)
		return rc
	})
}

func (e *nativeEngine) repoContains(_ context.Context, object string) (bool, error) {
	return boolResult("repo.contains", func(out *C.uint8_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_contains(e.handle, stringSlice(object), out, errBuf)
		runtime.KeepAlive(object)
		return rc
	})
}

func (e *nativeEngine) repoMissingTree(_ context.Context, commit string) ([]string, error) {
	var ids []string
	err := jsonResult("repo.missing_tree", &ids, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_missing_tree(e.handle, stringSlice(commit), out, errBuf)
		runtime.KeepAlive(commit)
		return rc
	})
	return ids, err
}

func (e *nativeEngine) repoObjectClosure(_ context.Context, commit string) ([]string, error) {
	var ids []string
	err := jsonResult("repo.object_closure", &ids, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_object_closure(e.handle, stringSlice(commit), out, errBuf)
		runtime.KeepAlive(commit)
		return rc
	})
	return ids, err
}

func (e *nativeEngine) repoMissingPrepare(_ context.Context, commit string) ([]string, error) {
	var ids []string
	err := jsonResult("repo.missing_prepare", &ids, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_missing_prepare(e.handle, stringSlice(commit), out, errBuf)
		runtime.KeepAlive(commit)
		return rc
	})
	return ids, err
}

func (e *nativeEngine) repoMissingHistory(_ context.Context, commit, entity string) ([]string, error) {
	var ids []string
	err := jsonResult("repo.missing_history", &ids, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_missing_history(e.handle, stringSlice(commit), stringSlice(entity), out, errBuf)
		runtime.KeepAlive(commit)
		runtime.KeepAlive(entity)
		return rc
	})
	return ids, err
}

func (e *nativeEngine) repoMissingStat(_ context.Context, commit string) ([]string, error) {
	var ids []string
	err := jsonResult("repo.missing_stat", &ids, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_missing_stat(e.handle, stringSlice(commit), out, errBuf)
		runtime.KeepAlive(commit)
		return rc
	})
	return ids, err
}

func (e *nativeEngine) repoMissingCommits(_ context.Context, commit string, limit int) ([]string, error) {
	var ids []string
	err := jsonResult("repo.missing_commits", &ids, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_missing_commits(e.handle, stringSlice(commit), C.size_t(limit), out, errBuf)
		runtime.KeepAlive(commit)
		return rc
	})
	return ids, err
}

func (e *nativeEngine) repoMissingDiff(_ context.Context, base, head string, detail DiffDetail) ([]string, error) {
	var ids []string
	nativeDetail := string(detail)
	err := jsonResult("repo.missing_diff", &ids, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_missing_diff(e.handle, stringSlice(base), stringSlice(head), stringSlice(nativeDetail), out, errBuf)
		runtime.KeepAlive(base)
		runtime.KeepAlive(head)
		runtime.KeepAlive(nativeDetail)
		return rc
	})
	return ids, err
}

func (e *nativeEngine) repoMissingDelta(_ context.Context, commit string, start, end *int64) ([]string, error) {
	var ids []string
	var nativeStart, nativeEnd C.int64_t
	var hasStart, hasEnd C.uint8_t
	if start != nil {
		nativeStart = C.int64_t(*start)
		hasStart = 1
	}
	if end != nil {
		nativeEnd = C.int64_t(*end)
		hasEnd = 1
	}
	err := jsonResult("repo.missing_delta", &ids, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_missing_delta(e.handle, stringSlice(commit), nativeStart, hasStart, nativeEnd, hasEnd, out, errBuf)
		runtime.KeepAlive(commit)
		return rc
	})
	return ids, err
}

// repoMissingMerge takes an empty target to mean "merge into the current head",
// which optionalStringSlice turns into the runtime's None.
func (e *nativeEngine) repoMissingMerge(_ context.Context, source, target string) ([]string, error) {
	var ids []string
	err := jsonResult("repo.missing_merge", &ids, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_missing_merge(e.handle, stringSlice(source), optionalStringSlice(target), out, errBuf)
		runtime.KeepAlive(source)
		runtime.KeepAlive(target)
		return rc
	})
	return ids, err
}

func (e *nativeEngine) repoReleaseTransientObjects(_ context.Context) error {
	return statusCall("repo.release_transient_objects", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_repo_release_transient_objects(e.handle, errBuf)
	})
}

func (e *nativeEngine) repoDiff(_ context.Context, baseRef, headRef string, options DiffOptions) (GraphDiff, error) {
	var diff GraphDiff
	var nativeLimit C.size_t
	var hasLimit C.uint8_t
	if options.Limit != nil {
		nativeLimit = C.size_t(*options.Limit)
		hasLimit = 1
	}
	detail := string(options.detail())
	err := jsonResult("repo.diff", &diff, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_diff(e.handle, stringSlice(baseRef), stringSlice(headRef), nativeLimit, hasLimit, stringSlice(detail), out, errBuf)
		runtime.KeepAlive(baseRef)
		runtime.KeepAlive(headRef)
		runtime.KeepAlive(detail)
		return rc
	})
	return diff, err
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

func (e *nativeEngine) repoLog(_ context.Context, revision string, limit int) ([]map[string]any, error) {
	var commits []map[string]any
	err := jsonResult("repo.log", &commits, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_log(e.handle, stringSlice(revision), C.size_t(limit), out, errBuf)
		runtime.KeepAlive(revision)
		return rc
	})
	return commits, err
}

func (e *nativeEngine) repoHistory(
	_ context.Context,
	entityID string,
	revision string,
	limit *int,
) ([]map[string]any, error) {
	var entries []map[string]any
	var nativeLimit C.size_t
	var hasLimit C.uint8_t
	if limit != nil {
		nativeLimit = C.size_t(*limit)
		hasLimit = 1
	}
	err := jsonResult("repo.history", &entries, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_history(e.handle, stringSlice(entityID), stringSlice(revision), nativeLimit, hasLimit, out, errBuf)
		runtime.KeepAlive(entityID)
		runtime.KeepAlive(revision)
		return rc
	})
	return entries, err
}

func (e *nativeEngine) repoBranch(_ context.Context, name, startPoint string) (string, error) {
	var commit string
	err := jsonResult("repo.branch", &commit, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_branch(e.handle, stringSlice(name), stringSlice(startPoint), out, errBuf)
		runtime.KeepAlive(name)
		runtime.KeepAlive(startPoint)
		return rc
	})
	return commit, err
}

func (e *nativeEngine) repoMerge(
	_ context.Context,
	source string,
	target string,
	expectedHead *string,
	message *string,
) (Commit, error) {
	var commit Commit
	err := jsonResult("repo.merge", &commit, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		var expected, commitMessage C.CstxSlice
		if expectedHead != nil {
			expected = stringSlice(*expectedHead)
		}
		if message != nil {
			commitMessage = stringSlice(*message)
		}
		rc := C.cstx_repo_merge(e.handle, stringSlice(source), stringSlice(target), expected, commitMessage, out, errBuf)
		runtime.KeepAlive(source)
		runtime.KeepAlive(target)
		runtime.KeepAlive(expectedHead)
		runtime.KeepAlive(message)
		return rc
	})
	return commit, err
}

func (e *nativeEngine) repoStat(
	_ context.Context,
	revision string,
	excludeMask uint64,
	includeMask uint64,
) (GraphStats, error) {
	var value GraphStats
	err := jsonResult("repo.stat", &value, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_stat(e.handle, stringSlice(revision), C.uint64_t(excludeMask), C.uint64_t(includeMask), out, errBuf)
		runtime.KeepAlive(revision)
		return rc
	})
	return value, err
}

func (e *nativeEngine) repoDelta(
	_ context.Context,
	revision string,
	startTimestamp *int64,
	endTimestamp *int64,
) (Delta, error) {
	var value Delta
	var start, end C.int64_t
	var hasStart, hasEnd C.uint8_t
	if startTimestamp != nil {
		start = C.int64_t(*startTimestamp)
		hasStart = 1
	}
	if endTimestamp != nil {
		end = C.int64_t(*endTimestamp)
		hasEnd = 1
	}
	err := jsonResult("repo.delta", &value, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_delta(e.handle, stringSlice(revision), start, hasStart, end, hasEnd, out, errBuf)
		runtime.KeepAlive(revision)
		return rc
	})
	return value, err
}

// --- unified graph cursor ------------------------------------------------
type nativeGraphCursor struct{ cursor *C.CstxGraphCursor }

func newNativeGraphCursor(cursor *C.CstxGraphCursor) *nativeGraphCursor {
	result := &nativeGraphCursor{cursor: cursor}
	runtime.SetFinalizer(result, (*nativeGraphCursor).close)
	return result
}

func (c *nativeGraphCursor) page(_ context.Context, limit, page int) (CursorPage, error) {
	if c.cursor == nil {
		return CursorPage{}, &Error{Code: CodeInvalidArgument, Operation: "cursor.page", Message: "cursor is closed"}
	}
	var result CursorPage
	err := jsonResult("cursor.page", &result, func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_graph_cursor_page(c.cursor, C.size_t(limit), C.size_t(page), out, errBuf)
	})
	return result, err
}

func (c *nativeGraphCursor) close() {
	if c.cursor != nil {
		C.cstx_graph_cursor_free(c.cursor)
		c.cursor = nil
		runtime.SetFinalizer(c, nil)
	}
}
