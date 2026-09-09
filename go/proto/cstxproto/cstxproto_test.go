package cstxproto

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestRetiredPayloadSurfaceIsAbsent(t *testing.T) {
	if File_cstx_proto.Enums().ByName(protoreflect.Name("PayloadFormat")) != nil {
		t.Fatal("PayloadFormat is still present")
	}
	for message, field := range map[protoreflect.MessageDescriptor]protoreflect.Name{
		(&RuntimeConfig{}).ProtoReflect().Descriptor(): "payload_format",
		(&Node{}).ProtoReflect().Descriptor():          "entity",
		(&Relationship{}).ProtoReflect().Descriptor():  "relation",
	} {
		if message.Fields().ByName(field) != nil {
			t.Fatalf("%s.%s is still present", message.Name(), field)
		}
	}
}

func TestNodeRoundTrip(t *testing.T) {
	extras, err := structpb.NewStruct(map[string]any{"source": "test"})
	if err != nil {
		t.Fatal(err)
	}
	payload := &EntityValue{
		NodeType: "ip",
		Fields:   []*EntityField{{Name: "ip", Value: &EntityField_Text{Text: "1.1.1.1"}}},
	}
	id := "ip:one"
	want := &Node{Id: &id, Value: payload, Annotations: extras}
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
