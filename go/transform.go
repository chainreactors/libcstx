//go:build cstx_native

package cstx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// Parse converts tool output into typed SCO nodes (stateless transform).
// tool: artifact name ("gogo", "spray", "zombie", "neutron", etc.)
// input: []byte (raw JSONL), a struct, or a slice of structs.
func Parse(tool string, input any) ([]SCONode, error) {
	data, err := toJSONL(input)
	if err != nil {
		return nil, err
	}
	raw, err := Transform(tool, data)
	if err != nil {
		return nil, err
	}
	return parseSCONodes(raw)
}

func parseSCONodes(data []byte) ([]SCONode, error) {
	var rawNodes []json.RawMessage
	if err := json.Unmarshal(data, &rawNodes); err != nil {
		return nil, fmt.Errorf("cstx parse: %w", err)
	}
	nodes := make([]SCONode, 0, len(rawNodes))
	for _, r := range rawNodes {
		n, err := ParseSCONode(r)
		if err != nil {
			continue
		}
		if n != nil {
			nodes = append(nodes, n)
		}
	}
	return nodes, nil
}

func toJSONL(input any) ([]byte, error) {
	if b, ok := input.([]byte); ok {
		return b, nil
	}
	rv := reflect.ValueOf(input)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice {
		b, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		return append(b, '\n'), nil
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	for i := 0; i < rv.Len(); i++ {
		if err := enc.Encode(rv.Index(i).Interface()); err != nil {
			return nil, fmt.Errorf("cstx marshal [%d]: %w", i, err)
		}
	}
	return buf.Bytes(), nil
}
