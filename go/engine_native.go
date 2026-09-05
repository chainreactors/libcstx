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
	"fmt"
	"runtime"
	"unsafe"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type nativeEngine struct {
	handle *C.CstxHandle
}

func newEngine(config runtimeConfig) (engine, error) {
	payload, err := proto.Marshal(&cstxproto.RuntimeConfig{
		ProjectId:      config.projectID,
		CursorPageSize: uint64(config.cursorPageSize),
		PayloadFormat:  config.payloadFormat,
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
	payload, err := proto.Marshal(&cstxproto.GraphSelection{NodeIds: seedIDs})
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
	payload, err := proto.Marshal(&cstxproto.GraphSelection{NodeIds: nodeIDs})
	if err != nil {
		return 0, err
	}
	return countResult("graph.delete_nodes", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_delete_nodes(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
}

func (e *nativeEngine) graphDeleteRelationships(_ context.Context, relationshipIDs []string) (uint64, error) {
	payload, err := proto.Marshal(&cstxproto.GraphSelection{RelationshipIds: relationshipIDs})
	if err != nil {
		return 0, err
	}
	return countResult("graph.delete_relationships", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_delete_relationships(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
}

// --- C transport helpers -------------------------------------------------

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

func bufferResult(op string, call func(out, errBuf *C.CstxBuffer) C.CstxStatusCode) ([]byte, error) {
	var out, errBuf C.CstxBuffer
	if err := statusError(call(&out, &errBuf), op, &errBuf); err != nil {
		C.cstx_buffer_free(&out)
		return nil, err
	}
	return takeBuffer(&out), nil
}

func textResult(op string, call func(out, errBuf *C.CstxBuffer) C.CstxStatusCode) (string, error) {
	data, err := bufferResult(op, call)
	if err != nil {
		return "", err
	}
	return string(data), nil
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

// --- runtime -------------------------------------------------------------

func (e *nativeEngine) lastChange(_ context.Context) (*cstxproto.GraphChangeSet, error) {
	data, err := bufferResult("cstx.last_change", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_last_change(e.handle, out, errBuf)
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.GraphChangeSet
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, fmt.Errorf("cstx: decode change set protobuf: %w", err)
	}
	return &wire, nil
}

// --- extensions ----------------------------------------------------------

func (e *nativeEngine) extensionRegister(_ context.Context, contract *cstxproto.ExtensionContract) error {
	payload, err := proto.Marshal(contract)
	if err != nil {
		return err
	}
	return statusCall("extensions.register", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_extension_register(e.handle, byteSlice(payload), errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
}

func (e *nativeEngine) extensionExportContract(_ context.Context) (cstxproto.ExtensionContract, error) {
	data, err := bufferResult("extensions.export_contract", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_extension_export_contract(e.handle, out, errBuf)
	})
	if err != nil {
		return cstxproto.ExtensionContract{}, err
	}
	var value cstxproto.ExtensionContract
	if err := proto.Unmarshal(data, &value); err != nil {
		return cstxproto.ExtensionContract{}, fmt.Errorf("cstx: decode extensions.export_contract protobuf: %w", err)
	}
	return value, nil
}

func (e *nativeEngine) extensionEnable(_ context.Context, name string) error {
	return statusCall("extensions.enable", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_extension_enable(e.handle, stringSlice(name), errBuf)
		runtime.KeepAlive(name)
		return rc
	})
}

func (e *nativeEngine) extensionList(_ context.Context) (*cstxproto.ExtensionCatalog, error) {
	data, err := bufferResult("extensions.list", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode { return C.cstx_extension_list(e.handle, out, errBuf) })
	if err != nil {
		return nil, err
	}
	var value cstxproto.ExtensionCatalog
	if err := proto.Unmarshal(data, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func (e *nativeEngine) extensionInfo(_ context.Context, name string) (*cstxproto.ExtensionInfo, error) {
	data, err := bufferResult("extensions.info", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_extension_info(e.handle, stringSlice(name), out, errBuf)
		runtime.KeepAlive(name)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var value cstxproto.ExtensionInfo
	if err := proto.Unmarshal(data, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func (e *nativeEngine) extensionContains(_ context.Context, nodeType string) (bool, error) {
	return boolResult("extensions.contains", func(out *C.uint8_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_extension_contains(e.handle, stringSlice(nodeType), out, errBuf)
		runtime.KeepAlive(nodeType)
		return rc
	})
}

func (e *nativeEngine) extensionSchema(_ context.Context, nodeType string) (cstxproto.NodeType, error) {
	data, err := bufferResult("extensions.schema", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_extension_schema(e.handle, stringSlice(nodeType), out, errBuf)
		runtime.KeepAlive(nodeType)
		return rc
	})
	if err != nil {
		return cstxproto.NodeType{}, err
	}
	var value cstxproto.NodeType
	if err := proto.Unmarshal(data, &value); err != nil {
		return cstxproto.NodeType{}, err
	}
	return value, nil
}

func (e *nativeEngine) extensionSchemas(_ context.Context) (cstxproto.NodeTypeCatalog, error) {
	data, err := bufferResult("extensions.schemas", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_extension_schemas(e.handle, out, errBuf)
	})
	if err != nil {
		return cstxproto.NodeTypeCatalog{}, err
	}
	var value cstxproto.NodeTypeCatalog
	if err := proto.Unmarshal(data, &value); err != nil {
		return cstxproto.NodeTypeCatalog{}, err
	}
	return value, nil
}

func (e *nativeEngine) extensionHasNativeArtifact(_ context.Context, artifact string) (bool, error) {
	return boolResult("extensions.has_native_artifact", func(out *C.uint8_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_extension_has_native_artifact(e.handle, stringSlice(artifact), out, errBuf)
		runtime.KeepAlive(artifact)
		return rc
	})
}

func (e *nativeEngine) extensionAnchorConcepts(_ context.Context) (cstxproto.AnchorConceptCatalog, error) {
	data, err := bufferResult("extensions.anchor_concepts", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_extension_anchor_concepts(e.handle, out, errBuf)
	})
	if err != nil {
		return cstxproto.AnchorConceptCatalog{}, err
	}
	var value cstxproto.AnchorConceptCatalog
	if err := proto.Unmarshal(data, &value); err != nil {
		return cstxproto.AnchorConceptCatalog{}, err
	}
	return value, nil
}

// --- graph ---------------------------------------------------------------

func (e *nativeEngine) graphAddNodes(_ context.Context, nodes []*cstxproto.Node) (uint64, error) {
	return e.graphAddNodesWire(context.Background(), &cstxproto.Graph{Nodes: nodes})
}

func (e *nativeEngine) graphReplaceNodes(_ context.Context, nodes []*cstxproto.Node) (uint64, error) {
	return e.graphReplaceNodesWire(context.Background(), &cstxproto.Graph{Nodes: nodes})
}

func (e *nativeEngine) graphAddRelationships(_ context.Context, relationships []*cstxproto.Relationship) (uint64, error) {
	return e.graphAddRelationshipsWire(context.Background(), &cstxproto.Graph{Relationships: relationships})
}

func (e *nativeEngine) graphIngest(_ context.Context, plugin, artifact string, data []byte) (cstxproto.GraphIngestResult, error) {
	request, err := proto.Marshal(&cstxproto.ParserPayload{
		Plugin:   plugin,
		Artifact: artifact,
		Data:     data,
	})
	if err != nil {
		return cstxproto.GraphIngestResult{}, err
	}
	bytes, err := bufferResult("graph.ingest", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_ingest(e.handle, byteSlice(request), out, errBuf)
		runtime.KeepAlive(request)
		return rc
	})
	var value cstxproto.GraphIngestResult
	if err := proto.Unmarshal(bytes, &value); err != nil {
		return cstxproto.GraphIngestResult{}, err
	}
	return value, nil
}

func (e *nativeEngine) graphNode(_ context.Context, nodeID string) (*cstxproto.Node, error) {
	node, err := e.graphNodeWire(context.Background(), nodeID)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (e *nativeEngine) graphRelationship(_ context.Context, relationshipID string) (*cstxproto.Relationship, error) {
	relationship, err := e.graphRelationshipWire(context.Background(), relationshipID)
	if err != nil {
		return nil, err
	}
	return &relationship, nil
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

func (e *nativeEngine) graphRelationshipCount(_ context.Context) (uint64, error) {
	return countResult("graph.relationship_count", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_graph_relationship_count(e.handle, out, errBuf)
	})
}

func (e *nativeEngine) graphStats(_ context.Context) (*cstxproto.GraphStats, error) {
	data, err := bufferResult("graph.stats", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_graph_stats(e.handle, 0, 0, out, errBuf)
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.GraphStats
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
}

func (e *nativeEngine) graphNodes(_ context.Context, query *cstxproto.NodeQuery) (graphCursor, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, err
	}
	var cursor *C.CstxGraphCursor
	err = statusCall("graph.nodes", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_nodes(e.handle, byteSlice(payload), &cursor, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNativeGraphCursor(cursor), nil
}

func (e *nativeEngine) graphRelationships(_ context.Context, query *cstxproto.RelationshipQuery) (graphCursor, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, err
	}
	var cursor *C.CstxGraphCursor
	err = statusCall("graph.relationships", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_relationships(e.handle, byteSlice(payload), &cursor, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNativeGraphCursor(cursor), nil
}

func (e *nativeEngine) graphNeighbors(_ context.Context, query *cstxproto.NeighborQuery) (graphCursor, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, err
	}
	var cursor *C.CstxGraphCursor
	err = statusCall("graph.neighbors", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_neighbors(e.handle, byteSlice(payload), &cursor, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNativeGraphCursor(cursor), nil
}

func (e *nativeEngine) graphQuery(_ context.Context, query *cstxproto.GraphQuery) (graphCursor, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, err
	}
	var cursor *C.CstxGraphCursor
	err = statusCall("graph.query", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_query(e.handle, byteSlice(payload), &cursor, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return nil, err
	}
	return newNativeGraphCursor(cursor), nil
}

func (e *nativeEngine) graphAnalyze(_ context.Context, algorithm *cstxproto.Algorithm, selection *string) (uint8, bool, graphCursor, error) {
	payload, err := proto.Marshal(algorithm)
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
	return textResult("repo.resolve", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_resolve(e.handle, stringSlice(revision), out, errBuf)
		runtime.KeepAlive(revision)
		return rc
	})
}

func (e *nativeEngine) repoCheckout(_ context.Context, revision string, force bool) (*cstxproto.Commit, error) {
	var nativeForce C.uint8_t
	if force {
		nativeForce = 1
	}
	data, err := bufferResult("repo.checkout", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_checkout(e.handle, stringSlice(revision), nativeForce, out, errBuf)
		runtime.KeepAlive(revision)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.Commit
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
}

func (e *nativeEngine) repoCommit(
	_ context.Context,
	message string,
	refName string,
	expectedHead *string,
	metadata *structpb.Struct,
) (*cstxproto.Commit, error) {
	if metadata == nil {
		metadata = &structpb.Struct{}
	}
	metadataBytes, err := proto.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	data, err := bufferResult("repo.commit", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		var expected C.CstxSlice
		if expectedHead != nil {
			expected = stringSlice(*expectedHead)
		}
		rc := C.cstx_repo_commit(e.handle, stringSlice(message), stringSlice(refName), expected, byteSlice(metadataBytes), out, errBuf)
		runtime.KeepAlive(refName)
		runtime.KeepAlive(message)
		runtime.KeepAlive(expectedHead)
		runtime.KeepAlive(metadataBytes)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.Commit
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
}

func (e *nativeEngine) repoPrepare(
	_ context.Context,
	message string,
	refName string,
	expectedHead *string,
	metadata *structpb.Struct,
	timestamp *int64,
) (*cstxproto.PublicationPlan, error) {
	if metadata == nil {
		metadata = &structpb.Struct{}
	}
	metadataBytes, err := proto.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	var nativeTimestamp C.int64_t
	var hasTimestamp C.uint8_t
	if timestamp != nil {
		nativeTimestamp = C.int64_t(*timestamp)
		hasTimestamp = 1
	}
	data, err := bufferResult("repo.prepare", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		var expected C.CstxSlice
		if expectedHead != nil {
			expected = stringSlice(*expectedHead)
		}
		rc := C.cstx_repo_prepare(e.handle, stringSlice(message), stringSlice(refName), expected, byteSlice(metadataBytes), nativeTimestamp, hasTimestamp, out, errBuf)
		runtime.KeepAlive(message)
		runtime.KeepAlive(refName)
		runtime.KeepAlive(expectedHead)
		runtime.KeepAlive(metadataBytes)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.PublicationPlan
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
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

func (e *nativeEngine) repoSynchronize(_ context.Context, state *cstxproto.RepositoryState) error {
	if state == nil {
		state = &cstxproto.RepositoryState{}
	}
	data, err := proto.Marshal(state)
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

func (e *nativeEngine) repoMissing(_ context.Context, plan *cstxproto.RepositoryObjectPlan) (*cstxproto.ObjectSelection, error) {
	if plan == nil {
		return nil, fmt.Errorf("cstx: repository object plan must not be nil")
	}
	payload, err := proto.Marshal(plan)
	if err != nil {
		return nil, err
	}
	data, err := bufferResult("repo.missing", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_missing(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var value cstxproto.ObjectSelection
	if err := proto.Unmarshal(data, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func (e *nativeEngine) repoReleaseTransientObjects(_ context.Context) error {
	return statusCall("repo.release_transient_objects", func(errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_repo_release_transient_objects(e.handle, errBuf)
	})
}

func (e *nativeEngine) repoDiff(_ context.Context, baseRef, headRef string, limit *uint64, detailValue cstxproto.DiffDetail) (*cstxproto.GraphDiff, error) {
	var nativeLimit C.size_t
	var hasLimit C.uint8_t
	if limit != nil {
		nativeLimit = C.size_t(*limit)
		hasLimit = 1
	}
	detail := "entities"
	if detailValue == cstxproto.DiffDetail_DIFF_DETAIL_COUNTS {
		detail = "counts"
	}
	data, err := bufferResult("repo.diff", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_diff(e.handle, stringSlice(baseRef), stringSlice(headRef), nativeLimit, hasLimit, stringSlice(detail), out, errBuf)
		runtime.KeepAlive(baseRef)
		runtime.KeepAlive(headRef)
		runtime.KeepAlive(detail)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.GraphDiff
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
}

func (e *nativeEngine) repoHead(_ context.Context, refName string) (*string, error) {
	value, err := textResult("repo.head", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_head(e.handle, stringSlice(refName), out, errBuf)
		runtime.KeepAlive(refName)
		return rc
	})
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, nil
	}
	return &value, nil
}

func (e *nativeEngine) repoLog(_ context.Context, revision string, limit int) (*cstxproto.CommitLog, error) {
	data, err := bufferResult("repo.log", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_log(e.handle, stringSlice(revision), C.size_t(limit), out, errBuf)
		runtime.KeepAlive(revision)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.CommitLog
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
}

func (e *nativeEngine) repoHistory(
	_ context.Context,
	entityID string,
	revision string,
	limit *int,
) (*cstxproto.EntityHistory, error) {
	var nativeLimit C.size_t
	var hasLimit C.uint8_t
	if limit != nil {
		nativeLimit = C.size_t(*limit)
		hasLimit = 1
	}
	data, err := bufferResult("repo.history", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_history(e.handle, stringSlice(entityID), stringSlice(revision), nativeLimit, hasLimit, out, errBuf)
		runtime.KeepAlive(entityID)
		runtime.KeepAlive(revision)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.EntityHistory
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
}

func (e *nativeEngine) repoBranch(_ context.Context, name, startPoint string) (string, error) {
	return textResult("repo.branch", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_branch(e.handle, stringSlice(name), stringSlice(startPoint), out, errBuf)
		runtime.KeepAlive(name)
		runtime.KeepAlive(startPoint)
		return rc
	})
}

func (e *nativeEngine) repoMerge(
	_ context.Context,
	source string,
	target string,
	expectedHead *string,
	message *string,
) (*cstxproto.Commit, error) {
	data, err := bufferResult("repo.merge", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
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
	if err != nil {
		return nil, err
	}
	var wire cstxproto.Commit
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
}

func (e *nativeEngine) repoStat(
	_ context.Context,
	revision string,
	excludeMask uint64,
	includeMask uint64,
) (*cstxproto.GraphStats, error) {
	data, err := bufferResult("repo.stat", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_stat(e.handle, stringSlice(revision), C.uint64_t(excludeMask), C.uint64_t(includeMask), out, errBuf)
		runtime.KeepAlive(revision)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.GraphStats
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
}

func (e *nativeEngine) repoDelta(
	_ context.Context,
	revision string,
	startTimestamp *int64,
	endTimestamp *int64,
) (*cstxproto.GraphChangeSummary, error) {
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
	data, err := bufferResult("repo.delta", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_repo_delta(e.handle, stringSlice(revision), start, hasStart, end, hasEnd, out, errBuf)
		runtime.KeepAlive(revision)
		return rc
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.GraphChangeSummary
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
}

// --- unified graph cursor ------------------------------------------------
type nativeGraphCursor struct{ cursor *C.CstxGraphCursor }

func newNativeGraphCursor(cursor *C.CstxGraphCursor) *nativeGraphCursor {
	result := &nativeGraphCursor{cursor: cursor}
	runtime.SetFinalizer(result, (*nativeGraphCursor).close)
	return result
}

func (c *nativeGraphCursor) page(_ context.Context, limit, page int) (*cstxproto.GraphResultPage, error) {
	if c.cursor == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "cursor.page", Message: "cursor is closed"}
	}
	data, err := bufferResult("cursor.page", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		return C.cstx_graph_cursor_page(c.cursor, C.size_t(limit), C.size_t(page), out, errBuf)
	})
	if err != nil {
		return nil, err
	}
	var wire cstxproto.GraphResultPage
	if err := proto.Unmarshal(data, &wire); err != nil {
		return nil, err
	}
	return &wire, nil
}

func (c *nativeGraphCursor) close() {
	if c.cursor != nil {
		C.cstx_graph_cursor_free(c.cursor)
		c.cursor = nil
		runtime.SetFinalizer(c, nil)
	}
}
