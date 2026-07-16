package cstx

/*
#include "cstx_ffi.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// EdgesFiltered returns edges matching optional filters. Empty string means no filter.
func (g *Graph) EdgesFiltered(sourceID, targetID, relation string) string {
	var cSrc, cTgt, cRel *C.char
	if sourceID != "" {
		cSrc = C.CString(sourceID)
		defer C.free(unsafe.Pointer(cSrc))
	}
	if targetID != "" {
		cTgt = C.CString(targetID)
		defer C.free(unsafe.Pointer(cTgt))
	}
	if relation != "" {
		cRel = C.CString(relation)
		defer C.free(unsafe.Pointer(cRel))
	}

	out := C.cstx_graph_edges_filtered(g.ptr, cSrc, cTgt, cRel)
	if out == nil {
		return "[]"
	}
	defer C.cstx_free_string(out)
	return C.GoString(out)
}

// ExportSnapshotJSON exports a complete graph snapshot as JSON.
func (g *Graph) ExportSnapshotJSON() string {
	out := C.cstx_graph_export_snapshot_json(g.ptr)
	if out == nil {
		return "{}"
	}
	defer C.cstx_free_string(out)
	return C.GoString(out)
}

// UpdateNodeFlags modifies flags on a node. setTo=-1 means don't set absolute value.
func (g *Graph) UpdateNodeFlags(nodeID string, add, remove uint64, setTo int64) bool {
	cID := C.CString(nodeID)
	defer C.free(unsafe.Pointer(cID))
	return C.cstx_graph_update_node_flags(g.ptr, cID, C.uint64_t(add), C.uint64_t(remove), C.int64_t(setTo)) == 0
}

// FlagsAllMask returns the mask of all defined flags.
func FlagsAllMask() uint64 {
	return uint64(C.cstx_flags_all_mask())
}

// FlagsDefaultExcludeMask returns the default exclude mask for queries.
func FlagsDefaultExcludeMask() uint64 {
	return uint64(C.cstx_flags_default_exclude_mask())
}
