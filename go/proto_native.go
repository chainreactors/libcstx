package cstx

// The native adapter transports only generated protobuf messages across the
// C ABI. It intentionally contains no SDK-owned graph or repository models.

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

func (e *nativeEngine) graphAddNodesWire(_ context.Context, graph *cstxproto.Graph) (uint64, error) {
	payload, err := proto.Marshal(graph)
	if err != nil {
		return 0, err
	}
	return countResult("graph.add_nodes", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_add_nodes(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
}

func (e *nativeEngine) graphReplaceNodesWire(_ context.Context, graph *cstxproto.Graph) (uint64, error) {
	payload, err := proto.Marshal(graph)
	if err != nil {
		return 0, err
	}
	return countResult("graph.replace_nodes", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_replace_nodes(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
}

func (e *nativeEngine) graphAddRelationshipsWire(_ context.Context, graph *cstxproto.Graph) (uint64, error) {
	payload, err := proto.Marshal(graph)
	if err != nil {
		return 0, err
	}
	return countResult("graph.add_relationships", func(out *C.uint64_t, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_add_relationships(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
}

func (e *nativeEngine) graphNodeWire(_ context.Context, nodeID string) (cstxproto.Node, error) {
	data, err := bufferResult("graph.node", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_node(e.handle, stringSlice(nodeID), out, errBuf)
		runtime.KeepAlive(nodeID)
		return rc
	})
	if err != nil {
		return cstxproto.Node{}, err
	}
	var node cstxproto.Node
	if err := proto.Unmarshal(data, &node); err != nil {
		return cstxproto.Node{}, fmt.Errorf("cstx: decode node protobuf: %w", err)
	}
	return node, nil
}

func (e *nativeEngine) graphRelationshipWire(_ context.Context, relationshipID string) (cstxproto.Relationship, error) {
	data, err := bufferResult("graph.relationship", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_relationship(e.handle, stringSlice(relationshipID), out, errBuf)
		runtime.KeepAlive(relationshipID)
		return rc
	})
	if err != nil {
		return cstxproto.Relationship{}, err
	}
	var relationship cstxproto.Relationship
	if err := proto.Unmarshal(data, &relationship); err != nil {
		return cstxproto.Relationship{}, fmt.Errorf("cstx: decode relationship protobuf: %w", err)
	}
	return relationship, nil
}
