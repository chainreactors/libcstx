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

func (e *nativeEngine) graphParseWire(_ context.Context, payload *cstxproto.ParserPayload) (cstxproto.Graph, uint64, error) {
	if payload == nil {
		return cstxproto.Graph{}, 0, &Error{Code: CodeInvalidArgument, Operation: "graph.parse", Message: "payload must not be nil"}
	}
	encoded, err := proto.Marshal(payload)
	if err != nil {
		return cstxproto.Graph{}, 0, err
	}
	var records C.uint64_t
	var out, errBuf C.CstxBuffer
	if err := statusError(C.cstx_graph_parse(e.handle, byteSlice(encoded), &records, &out, &errBuf), "graph.parse", &errBuf); err != nil {
		C.cstx_buffer_free(&out)
		return cstxproto.Graph{}, 0, err
	}
	runtime.KeepAlive(encoded)
	data := takeBuffer(&out)
	var graph cstxproto.Graph
	if err := proto.Unmarshal(data, &graph); err != nil {
		return cstxproto.Graph{}, 0, fmt.Errorf("cstx: decode graph protobuf: %w", err)
	}
	return graph, uint64(records), nil
}

func (e *nativeEngine) graphLinkWire(_ context.Context, nodeIDs []string, dataSource string) (cstxproto.GraphLinkResult, error) {
	selection, err := proto.Marshal(&cstxproto.GraphSelection{NodeIds: nodeIDs})
	if err != nil {
		return cstxproto.GraphLinkResult{}, err
	}
	data, err := bufferResult("graph.link", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_link(e.handle, byteSlice(selection), stringSlice(dataSource), out, errBuf)
		runtime.KeepAlive(selection)
		runtime.KeepAlive(dataSource)
		return rc
	})
	if err != nil {
		return cstxproto.GraphLinkResult{}, err
	}
	var result cstxproto.GraphLinkResult
	if err := proto.Unmarshal(data, &result); err != nil {
		return cstxproto.GraphLinkResult{}, fmt.Errorf("cstx: decode graph link result protobuf: %w", err)
	}
	return result, nil
}

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

func (e *nativeEngine) graphAddRelationshipWire(_ context.Context, relationship *cstxproto.Relationship) (cstxproto.Relationship, error) {
	if relationship == nil {
		return cstxproto.Relationship{}, &Error{Code: CodeInvalidArgument, Operation: "graph.add_relationship", Message: "relationship must not be nil"}
	}
	payload, err := proto.Marshal(relationship)
	if err != nil {
		return cstxproto.Relationship{}, err
	}
	data, err := bufferResult("graph.add_relationship", func(out, errBuf *C.CstxBuffer) C.CstxStatusCode {
		rc := C.cstx_graph_add_relationship(e.handle, byteSlice(payload), out, errBuf)
		runtime.KeepAlive(payload)
		return rc
	})
	if err != nil {
		return cstxproto.Relationship{}, err
	}
	var stored cstxproto.Relationship
	if err := proto.Unmarshal(data, &stored); err != nil {
		return cstxproto.Relationship{}, fmt.Errorf("cstx: decode relationship protobuf: %w", err)
	}
	return stored, nil
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
