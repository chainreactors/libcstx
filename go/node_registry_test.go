package cstx

import (
	"encoding/json"
	"testing"
)

type registryTestNode struct {
	nodeHeader
	Value string `json:"value"`
}

func TestRegisteredParserTakesPrecedence(t *testing.T) {
	typ := "test_registry_node"
	if err := RegisterSCOParser(typ, func(data []byte) (SCONode, error) {
		var node registryTestNode
		return &node, json.Unmarshal(data, &node)
	}); err != nil {
		t.Fatal(err)
	}
	node, err := ParseSCONode([]byte(`{"cstx_type":"test_registry_node","cstx_id":"x","value":"ok"}`))
	if err != nil {
		t.Fatal(err)
	}
	if node.CstxType() != typ || node.CstxID() != "x" {
		t.Fatalf("unexpected node: %#v", node)
	}
}

func TestRegisterSCOParserRejectsDuplicate(t *testing.T) {
	typ := "test_registry_duplicate"
	parser := func([]byte) (SCONode, error) { return &registryTestNode{}, nil }
	if err := RegisterSCOParser(typ, parser); err != nil {
		t.Fatal(err)
	}
	if err := RegisterSCOParser(typ, parser); err == nil {
		t.Fatal("expected duplicate registration to fail")
	}
}
