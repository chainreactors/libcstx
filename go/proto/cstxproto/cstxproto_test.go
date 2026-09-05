package cstxproto

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestNodeRoundTrip(t *testing.T) {
	extras, err := structpb.NewStruct(map[string]any{"source": "test"})
	if err != nil {
		t.Fatal(err)
	}
	entity := &anypb.Any{TypeUrl: "type.googleapis.com/easm.Ip", Value: []byte{0x01, 0x02}}
	id := "ip:one"
	want := &Node{Id: &id, Entity: entity, Annotations: extras}
	wire, err := proto.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got Node
	if err := proto.Unmarshal(wire, &got); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(want, &got) {
		t.Fatalf("round-trip changed node: %v", &got)
	}
}
