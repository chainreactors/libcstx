package cstx

import (
	"context"
	"fmt"
	"sort"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// NodeValues is one node's payload as field names and Go values.
//
// It is the shape a caller uses when it has no generated message type for the
// node — which is every type an extension declares at runtime, since no code
// generator ran for it. The runtime encodes and decodes it with the schema
// document that extension registered, so nothing here needs a field number.
type NodeValues map[string]any

// AddNodeValues writes one node from field names and values.
//
// The counterpart to AddNodes for callers without generated types. Identity
// comes from the schema document exactly as it does for a generated message,
// so a node written this way lands on the same id as the same content written
// as an Any.
func (g *Graph) AddNodeValues(ctx context.Context, nodeType string, values NodeValues, options ...NodeValueOption) (uint64, error) {
	node, err := ValueNode(nodeType, values, options...)
	if err != nil {
		return 0, err
	}
	return g.AddNodes(ctx, []*cstxproto.Node{node})
}

// NodeValueOption sets one non-payload field on a value-shaped node.
type NodeValueOption func(*cstxproto.Node)

// WithNodeID sets the node's id explicitly, for a type whose schema declares
// its identity computed.
func WithNodeID(id string) NodeValueOption {
	return func(node *cstxproto.Node) { node.Id = &id }
}

// WithNodeSources records which artifacts observed the node.
func WithNodeSources(sources ...string) NodeValueOption {
	return func(node *cstxproto.Node) { node.Sources = sources }
}

// ValueNode builds a value-shaped node without writing it.
func ValueNode(nodeType string, values NodeValues, options ...NodeValueOption) (*cstxproto.Node, error) {
	entity, err := EntityValue(nodeType, values)
	if err != nil {
		return nil, err
	}
	node := &cstxproto.Node{Value: entity}
	for _, option := range options {
		option(node)
	}
	return node, nil
}

// EntityValue converts a field map into the wire payload.
//
// The branch is chosen by the Go type, and the runtime checks it against the
// column the schema declared. A mismatch is an error there rather than a
// silent coercion here: a number stored as text is a field the encoder cannot
// reproduce, which the runtime refuses at registration for the same reason.
func EntityValue(nodeType string, values NodeValues) (*cstxproto.EntityValue, error) {
	fields := make([]*cstxproto.EntityField, 0, len(values))
	for name, value := range values {
		field := &cstxproto.EntityField{Name: name}
		switch typed := value.(type) {
		case nil:
			continue
		case string:
			field.Value = &cstxproto.EntityField_Text{Text: typed}
		case bool:
			field.Value = &cstxproto.EntityField_Flag{Flag: typed}
		case int:
			field.Value = &cstxproto.EntityField_Number{Number: int64(typed)}
		case int32:
			field.Value = &cstxproto.EntityField_Number{Number: int64(typed)}
		case int64:
			field.Value = &cstxproto.EntityField_Number{Number: typed}
		case float64:
			field.Value = &cstxproto.EntityField_Real{Real: typed}
		case []string:
			field.Value = &cstxproto.EntityField_List{List: &cstxproto.StringList{Values: typed}}
		default:
			return nil, fmt.Errorf("cstx: field %q has unsupported type %T", name, value)
		}
		fields = append(fields, field)
	}
	// Sorted so one map produces one message: Go randomizes map iteration, and
	// the payload is content, not a bag.
	sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })
	return &cstxproto.EntityValue{NodeType: nodeType, Fields: fields}, nil
}

// FieldValues reads a value-shaped payload back into a field map.
//
// Returns an error when the node came back as an Any, which means the runtime
// was opened with the entity payload format.
func FieldValues(node *cstxproto.Node) (string, NodeValues, error) {
	entity := node.GetValue()
	if entity == nil {
		return "", nil, fmt.Errorf(
			"cstx: node %q carries no value payload; open the runtime with PayloadFormat_PAYLOAD_FORMAT_VALUE",
			node.GetId(),
		)
	}
	values := make(NodeValues, len(entity.GetFields()))
	for _, field := range entity.GetFields() {
		switch carried := field.GetValue().(type) {
		case *cstxproto.EntityField_Text:
			values[field.GetName()] = carried.Text
		case *cstxproto.EntityField_Number:
			values[field.GetName()] = carried.Number
		case *cstxproto.EntityField_Flag:
			values[field.GetName()] = carried.Flag
		case *cstxproto.EntityField_Real:
			values[field.GetName()] = carried.Real
		case *cstxproto.EntityField_List:
			values[field.GetName()] = carried.List.GetValues()
		}
	}
	return entity.GetNodeType(), values, nil
}
