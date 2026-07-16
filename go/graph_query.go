package cstx

/*
#include "cstx_ffi.h"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"unsafe"
)

// QueryDSL executes a DSL expression and returns results as JSON.
func (g *Graph) QueryDSL(expression string, limit int, offset int) (string, error) {
	cExpr := C.CString(expression)
	defer C.free(unsafe.Pointer(cExpr))

	var out *C.char
	var outLen C.uintptr_t
	rc := C.cstx_graph_query_dsl(g.ptr, cExpr, C.intptr_t(limit), C.uintptr_t(offset), &out, &outLen)
	if out != nil {
		defer C.cstx_free_string(out)
	}
	if rc != 0 {
		return "", ffiError("query_dsl", out, outLen)
	}
	return C.GoStringN(out, C.int(outLen)), nil
}

// QueryNodeIDs executes a DSL expression and returns matching node IDs.
func (g *Graph) QueryNodeIDs(expression string, limit int, offset int) ([]string, error) {
	cExpr := C.CString(expression)
	defer C.free(unsafe.Pointer(cExpr))

	out := C.cstx_graph_query_node_ids(g.ptr, cExpr, C.intptr_t(limit), C.uintptr_t(offset))
	if out == nil {
		return nil, nil
	}
	defer C.cstx_free_string(out)
	var ids []string
	json.Unmarshal([]byte(C.GoString(out)), &ids)
	return ids, nil
}

// IsPathExpression returns whether the expression is a path/traversal query.
func (g *Graph) IsPathExpression(expression string) bool {
	cExpr := C.CString(expression)
	defer C.free(unsafe.Pointer(cExpr))
	return C.cstx_graph_is_path_expression(g.ptr, cExpr) != 0
}
