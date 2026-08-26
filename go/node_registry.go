package cstx

import (
	"encoding/json"
	"fmt"
	"sync"
)

// SCOParser decodes one native CSTX node type. Parsers are registered by the
// package that owns the type, keeping the graph engine extensible.
type SCOParser func([]byte) (SCONode, error)

var (
	parserMu sync.RWMutex
	parsers  = map[string]SCOParser{}
)

func RegisterSCOParser(nodeType string, parser SCOParser) error {
	if nodeType == "" || parser == nil {
		return fmt.Errorf("cstx: node type and parser are required")
	}
	parserMu.Lock()
	defer parserMu.Unlock()
	if _, exists := parsers[nodeType]; exists {
		return fmt.Errorf("cstx: parser already registered for %q", nodeType)
	}
	parsers[nodeType] = parser
	return nil
}

func parseRegisteredSCONode(data []byte) (SCONode, bool, error) {
	var header struct {
		Type string `json:"cstx_type"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, false, err
	}
	parserMu.RLock()
	parser := parsers[header.Type]
	parserMu.RUnlock()
	if parser == nil {
		return nil, false, nil
	}
	node, err := parser(data)
	return node, true, err
}
