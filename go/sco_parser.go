package cstx

import (
	"fmt"
	"strings"
	"sync"
)

// SCOParser decodes one registered SCO type. Built-in and extension SCOs use
// the same ParseSCONode entry point, so consumers do not need a second parser
// or a Cairn-owned type switch when libcstx gains a new node type.
type SCOParser func([]byte) (SCONode, error)

var scoParsers = struct {
	sync.RWMutex
	items map[string]SCOParser
}{items: make(map[string]SCOParser)}

// RegisterSCOParser adds a typed parser for an extension SCO. Registration is
// process-wide because parser implementations are part of the libcstx Go
// binding, while schema validation remains scoped to each CSTX runtime.
// Registering the same type twice is rejected to avoid order-dependent decode
// behavior across plugins.
func RegisterSCOParser(nodeType string, parser SCOParser) error {
	nodeType = strings.TrimSpace(nodeType)
	if nodeType == "" {
		return fmt.Errorf("cstx: SCO parser type is empty")
	}
	if parser == nil {
		return fmt.Errorf("cstx: SCO parser %q is nil", nodeType)
	}
	scoParsers.Lock()
	defer scoParsers.Unlock()
	if _, exists := scoParsers.items[nodeType]; exists {
		return fmt.Errorf("cstx: SCO parser %q is already registered", nodeType)
	}
	scoParsers.items[nodeType] = parser
	return nil
}

func registeredSCOParser(nodeType string) SCOParser {
	scoParsers.RLock()
	defer scoParsers.RUnlock()
	return scoParsers.items[nodeType]
}
