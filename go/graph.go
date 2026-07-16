package cstx

/*
#include "cstx_ffi.h"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"unsafe"
)

type Graph struct {
	ptr *C.struct_CstxGraph
}

func New() *Graph {
	return &Graph{ptr: C.cstx_graph_new()}
}

func (g *Graph) Close() {
	if g.ptr != nil {
		C.cstx_graph_free(g.ptr)
		g.ptr = nil
	}
}

// ── Plugin loading ──

func (g *Graph) LoadPlugin(name string) error {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	if C.cstx_graph_load_plugin(g.ptr, cs) != 0 {
		return fmt.Errorf("cstx: plugin %q not found", name)
	}
	return nil
}

func (g *Graph) LoadAllPlugins() int {
	return int(C.cstx_graph_load_all_plugins(g.ptr))
}

func AvailablePlugins() []string {
	cs := C.cstx_available_plugins()
	defer C.cstx_free_string(cs)
	var out []string
	json.Unmarshal([]byte(C.GoString(cs)), &out)
	return out
}

// ── Schema + Rules ──

func (g *Graph) RegisterSchema(nodeType, jsonSchema string, valueField string) error {
	cType := C.CString(nodeType)
	cSchema := C.CString(jsonSchema)
	defer C.free(unsafe.Pointer(cType))
	defer C.free(unsafe.Pointer(cSchema))

	var cVF *C.char
	if valueField != "" {
		cVF = C.CString(valueField)
		defer C.free(unsafe.Pointer(cVF))
	}
	if C.cstx_graph_register_schema(g.ptr, cType, cSchema, cVF) != 0 {
		return fmt.Errorf("cstx: register schema %q failed", nodeType)
	}
	return nil
}

func (g *Graph) AddJoinRule(leftType, rightType, relation, leftKey, rightKey string, predicted bool) error {
	clt := C.CString(leftType)
	crt := C.CString(rightType)
	cr := C.CString(relation)
	clk := C.CString(leftKey)
	crk := C.CString(rightKey)
	defer C.free(unsafe.Pointer(clt))
	defer C.free(unsafe.Pointer(crt))
	defer C.free(unsafe.Pointer(cr))
	defer C.free(unsafe.Pointer(clk))
	defer C.free(unsafe.Pointer(crk))

	pred := C.int(0)
	if predicted {
		pred = 1
	}
	if C.cstx_graph_add_join_rule(g.ptr, clt, crt, cr, clk, crk, pred, nil, nil) != 0 {
		return fmt.Errorf("cstx: add join rule failed")
	}
	return nil
}

// ── Node / Edge CRUD ──

func (g *Graph) AddNode(nodeType, cstxID, payloadJSON string, sources []string, flags uint64) error {
	cType := C.CString(nodeType)
	cID := C.CString(cstxID)
	cPayload := C.CString(payloadJSON)
	srcJSON, _ := json.Marshal(sources)
	cSrc := C.CString(string(srcJSON))
	defer C.free(unsafe.Pointer(cType))
	defer C.free(unsafe.Pointer(cID))
	defer C.free(unsafe.Pointer(cPayload))
	defer C.free(unsafe.Pointer(cSrc))

	if C.cstx_graph_add_node(g.ptr, cType, cID, cPayload, cSrc, C.uint64_t(flags)) != 0 {
		return fmt.Errorf("cstx: add node failed")
	}
	return nil
}

func (g *Graph) AddEdge(sourceID, targetID, relation, dataSource string) error {
	cs := C.CString(sourceID)
	ct := C.CString(targetID)
	cr := C.CString(relation)
	cd := C.CString(dataSource)
	defer C.free(unsafe.Pointer(cs))
	defer C.free(unsafe.Pointer(ct))
	defer C.free(unsafe.Pointer(cr))
	defer C.free(unsafe.Pointer(cd))

	if C.cstx_graph_add_edge(g.ptr, cs, ct, cr, cd) != 0 {
		return fmt.Errorf("cstx: add edge failed")
	}
	return nil
}

func (g *Graph) AddNodesBatch(nodesJSON string) (*BatchResult, error) {
	cj := C.CString(nodesJSON)
	defer C.free(unsafe.Pointer(cj))

	var out *C.char
	var outLen C.size_t
	rc := C.cstx_graph_add_nodes_batch(g.ptr, cj, &out, &outLen)
	if out != nil {
		defer C.cstx_free_string(out)
	}
	if rc != 0 {
		return nil, ffiError("add_nodes_batch", out, outLen)
	}
	var r BatchResult
	json.Unmarshal([]byte(C.GoStringN(out, C.int(outLen))), &r)
	return &r, nil
}

func (g *Graph) AddEdgesBatch(edgesJSON string) (int, error) {
	cj := C.CString(edgesJSON)
	defer C.free(unsafe.Pointer(cj))
	rc := C.cstx_graph_add_edges_batch(g.ptr, cj)
	if rc < 0 {
		return 0, fmt.Errorf("cstx: add_edges_batch failed")
	}
	return int(rc), nil
}

// ── Ingest ──

func (g *Graph) Link(source string, data []byte) (*IngestResult, error) {
	cs := C.CString(source)
	defer C.free(unsafe.Pointer(cs))

	var dp *C.uchar
	if len(data) > 0 {
		dp = (*C.uchar)(unsafe.Pointer(&data[0]))
	}
	var out *C.char
	var outLen C.size_t
	rc := C.cstx_graph_link(g.ptr, cs, dp, C.size_t(len(data)), &out, &outLen)
	if out != nil {
		defer C.cstx_free_string(out)
	}
	if rc != 0 {
		return nil, ffiError("link", out, outLen)
	}
	var r IngestResult
	json.Unmarshal([]byte(C.GoStringN(out, C.int(outLen))), &r)
	return &r, nil
}

func (g *Graph) IngestNative(plugin, artifact string, data []byte) (*IngestResult, error) {
	cp := C.CString(plugin)
	ca := C.CString(artifact)
	defer C.free(unsafe.Pointer(cp))
	defer C.free(unsafe.Pointer(ca))

	var dp *C.uchar
	if len(data) > 0 {
		dp = (*C.uchar)(unsafe.Pointer(&data[0]))
	}
	var out *C.char
	var outLen C.size_t
	rc := C.cstx_graph_ingest_native(g.ptr, cp, ca, dp, C.size_t(len(data)), &out, &outLen)
	if out != nil {
		defer C.cstx_free_string(out)
	}
	if rc != 0 {
		return nil, ffiError("ingest_native", out, outLen)
	}
	var r IngestResult
	json.Unmarshal([]byte(C.GoStringN(out, C.int(outLen))), &r)
	return &r, nil
}

func (g *Graph) HasNativeArtifact(artifact string) bool {
	ca := C.CString(artifact)
	defer C.free(unsafe.Pointer(ca))
	return C.cstx_graph_has_native_artifact(g.ptr, ca) != 0
}

func (g *Graph) IngestJSONL(nodeType string, data []byte, idExpr, dataSource string) (*IngestResult, error) {
	cType := C.CString(nodeType)
	cExpr := C.CString(idExpr)
	cSrc := C.CString(dataSource)
	defer C.free(unsafe.Pointer(cType))
	defer C.free(unsafe.Pointer(cExpr))
	defer C.free(unsafe.Pointer(cSrc))

	var dp *C.uchar
	if len(data) > 0 {
		dp = (*C.uchar)(unsafe.Pointer(&data[0]))
	}
	var out *C.char
	var outLen C.size_t
	rc := C.cstx_graph_ingest_jsonl(g.ptr, cType, dp, C.size_t(len(data)), cExpr, cSrc, &out, &outLen)
	if out != nil {
		defer C.cstx_free_string(out)
	}
	if rc != 0 {
		return nil, ffiError("ingest_jsonl", out, outLen)
	}
	var r IngestResult
	json.Unmarshal([]byte(C.GoStringN(out, C.int(outLen))), &r)
	return &r, nil
}

func (g *Graph) LinkNodes(dataSource string) (string, error) {
	cs := C.CString(dataSource)
	defer C.free(unsafe.Pointer(cs))
	var out *C.char
	var outLen C.size_t
	C.cstx_graph_link_nodes(g.ptr, cs, &out, &outLen)
	if out != nil {
		defer C.cstx_free_string(out)
		return C.GoStringN(out, C.int(outLen)), nil
	}
	return "{}", nil
}

// ── Queries ──

func (g *Graph) Node(nodeID string) (string, bool) {
	cs := C.CString(nodeID)
	defer C.free(unsafe.Pointer(cs))
	out := C.cstx_graph_node(g.ptr, cs)
	if out == nil {
		return "", false
	}
	defer C.cstx_free_string(out)
	return C.GoString(out), true
}

func (g *Graph) Nodes(typeFilter string) string {
	var cf *C.char
	if typeFilter != "" {
		cf = C.CString(typeFilter)
		defer C.free(unsafe.Pointer(cf))
	}
	out := C.cstx_graph_nodes(g.ptr, cf)
	defer C.cstx_free_string(out)
	return C.GoString(out)
}

func (g *Graph) Edges(relation string) string {
	var cr *C.char
	if relation != "" {
		cr = C.CString(relation)
		defer C.free(unsafe.Pointer(cr))
	}
	out := C.cstx_graph_edges(g.ptr, cr)
	defer C.cstx_free_string(out)
	return C.GoString(out)
}

func (g *Graph) NodeTypes() []string {
	out := C.cstx_graph_node_types(g.ptr)
	defer C.cstx_free_string(out)
	var types []string
	json.Unmarshal([]byte(C.GoString(out)), &types)
	return types
}

func (g *Graph) Neighbors(nodeID string) string {
	cs := C.CString(nodeID)
	defer C.free(unsafe.Pointer(cs))
	out := C.cstx_graph_neighbors(g.ptr, cs)
	defer C.cstx_free_string(out)
	return C.GoString(out)
}

func (g *Graph) ContainsNode(nodeID string) bool {
	cs := C.CString(nodeID)
	defer C.free(unsafe.Pointer(cs))
	return C.cstx_graph_contains_node(g.ptr, cs) != 0
}

func (g *Graph) NodeCount() int {
	return int(C.cstx_graph_node_count(g.ptr))
}

func (g *Graph) EdgeCount() int {
	return int(C.cstx_graph_edge_count(g.ptr))
}

func (g *Graph) NodeIDs() []string {
	out := C.cstx_graph_node_ids(g.ptr)
	defer C.cstx_free_string(out)
	var ids []string
	json.Unmarshal([]byte(C.GoString(out)), &ids)
	return ids
}

func (g *Graph) Stats() string {
	out := C.cstx_graph_stats(g.ptr)
	defer C.cstx_free_string(out)
	return C.GoString(out)
}

func (g *Graph) ToJSON() string {
	out := C.cstx_graph_to_json(g.ptr)
	defer C.cstx_free_string(out)
	return C.GoString(out)
}

func (g *Graph) AllNodesJSON() string {
	out := C.cstx_graph_all_nodes_json(g.ptr)
	defer C.cstx_free_string(out)
	return C.GoString(out)
}

func (g *Graph) NodePayload(nodeID string) (string, bool) {
	cs := C.CString(nodeID)
	defer C.free(unsafe.Pointer(cs))
	out := C.cstx_graph_node_payload(g.ptr, cs)
	if out == nil {
		return "", false
	}
	defer C.cstx_free_string(out)
	return C.GoString(out), true
}

// ── Stateless transform (compatible with utils/cstx Parse) ──

func Transform(sourceType string, data []byte) ([]byte, error) {
	cs := C.CString(sourceType)
	defer C.free(unsafe.Pointer(cs))

	var dp *C.uchar
	if len(data) > 0 {
		dp = (*C.uchar)(unsafe.Pointer(&data[0]))
	}
	var out *C.char
	var outLen C.size_t
	rc := C.cstx_transform(cs, dp, C.size_t(len(data)), &out, &outLen)
	if out != nil {
		defer C.cstx_free_string(out)
	}
	if rc != 0 {
		return nil, ffiError("transform", out, outLen)
	}
	return C.GoBytes(unsafe.Pointer(out), C.int(outLen)), nil
}

// ── Utilities ──

func Version() string {
	return C.GoString(C.cstx_version())
}

func SupportedArtifacts() []string {
	out := C.cstx_supported_artifacts()
	defer C.cstx_free_string(out)
	var arts []string
	json.Unmarshal([]byte(C.GoString(out)), &arts)
	return arts
}

func ffiError(op string, out *C.char, outLen C.size_t) error {
	if out != nil {
		return fmt.Errorf("cstx %s: %s", op, C.GoStringN(out, C.int(outLen)))
	}
	return fmt.Errorf("cstx %s failed", op)
}
