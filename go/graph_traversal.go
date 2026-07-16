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

// NeighborIDs returns IDs of neighbors in the given direction ("in", "out", "both").
func (g *Graph) NeighborIDs(nodeID string, direction string) []string {
	cID := C.CString(nodeID)
	cDir := C.CString(direction)
	defer C.free(unsafe.Pointer(cID))
	defer C.free(unsafe.Pointer(cDir))

	out := C.cstx_graph_neighbor_ids(g.ptr, cID, cDir)
	if out == nil {
		return nil
	}
	defer C.cstx_free_string(out)
	var ids []string
	json.Unmarshal([]byte(C.GoString(out)), &ids)
	return ids
}

// BFS performs breadth-first search from seed. depth=0 means unlimited.
func (g *Graph) BFS(seedID string, depth uint32, reverse bool) []string {
	cSeed := C.CString(seedID)
	defer C.free(unsafe.Pointer(cSeed))

	rev := C.int(0)
	if reverse {
		rev = 1
	}
	out := C.cstx_graph_bfs(g.ptr, cSeed, C.uint(depth), rev)
	if out == nil {
		return nil
	}
	defer C.cstx_free_string(out)
	var ids []string
	json.Unmarshal([]byte(C.GoString(out)), &ids)
	return ids
}

// ShortestPaths returns all shortest paths between two nodes.
func (g *Graph) ShortestPaths(startID string, endID string, maxDepth uint32) [][]string {
	cStart := C.CString(startID)
	cEnd := C.CString(endID)
	defer C.free(unsafe.Pointer(cStart))
	defer C.free(unsafe.Pointer(cEnd))

	out := C.cstx_graph_shortest_paths(g.ptr, cStart, cEnd, C.uint(maxDepth))
	if out == nil {
		return nil
	}
	defer C.cstx_free_string(out)
	var paths [][]string
	json.Unmarshal([]byte(C.GoString(out)), &paths)
	return paths
}

// Degree returns node degree in the given direction.
func (g *Graph) Degree(nodeID string, direction string) int {
	cID := C.CString(nodeID)
	cDir := C.CString(direction)
	defer C.free(unsafe.Pointer(cID))
	defer C.free(unsafe.Pointer(cDir))

	return int(C.cstx_graph_degree(g.ptr, cID, cDir))
}

// SubgraphNodeIDs returns node and edge IDs reachable from seeds within the given depth.
func (g *Graph) SubgraphNodeIDs(seedIDs []string, depth uint32) (nodeIDs []string, edgeIDs []string, err error) {
	seedJSON, _ := json.Marshal(seedIDs)
	cSeeds := C.CString(string(seedJSON))
	defer C.free(unsafe.Pointer(cSeeds))

	var out *C.char
	var outLen C.size_t
	rc := C.cstx_graph_subgraph_node_ids(g.ptr, cSeeds, C.uint(depth), &out, &outLen)
	if out != nil {
		defer C.cstx_free_string(out)
	}
	if rc != 0 {
		return nil, nil, ffiError("subgraph_node_ids", out, outLen)
	}

	var result struct {
		Nodes []string `json:"nodes"`
		Edges []string `json:"edges"`
	}
	json.Unmarshal([]byte(C.GoStringN(out, C.int(outLen))), &result)
	return result.Nodes, result.Edges, nil
}
