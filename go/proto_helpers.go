package cstx

import (
	"fmt"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// NewNode constructs the canonical generated Node message around a typed SCO.
func NewNode(entity proto.Message, sources []string, annotations map[string]any, flags ...cstxproto.NodeFlag) (*cstxproto.Node, error) {
	if entity == nil {
		return nil, fmt.Errorf("cstx: node entity is required")
	}
	packed, err := anypb.New(entity)
	if err != nil {
		return nil, fmt.Errorf("cstx: pack node entity: %w", err)
	}
	value := &cstxproto.Node{
		Entity:  packed,
		Sources: append([]string(nil), sources...),
		Flags:   append([]cstxproto.NodeFlag(nil), flags...),
	}
	if annotations != nil {
		value.Annotations, err = structpb.NewStruct(annotations)
		if err != nil {
			return nil, fmt.Errorf("cstx: node annotations: %w", err)
		}
	}
	return value, nil
}

// NewRelationship constructs the canonical generated Relationship message
// around a typed SRO.
func NewRelationship(sourceID, targetID string, relation proto.Message, sources []string, annotations map[string]any) (*cstxproto.Relationship, error) {
	if relation == nil {
		return nil, fmt.Errorf("cstx: relationship relation is required")
	}
	packed, err := anypb.New(relation)
	if err != nil {
		return nil, fmt.Errorf("cstx: pack relationship relation: %w", err)
	}
	value := &cstxproto.Relationship{
		SourceId: sourceID,
		TargetId: targetID,
		Relation: packed,
		Sources:  append([]string(nil), sources...),
	}
	if annotations != nil {
		value.Annotations, err = structpb.NewStruct(annotations)
		if err != nil {
			return nil, fmt.Errorf("cstx: relationship annotations: %w", err)
		}
	}
	return value, nil
}
