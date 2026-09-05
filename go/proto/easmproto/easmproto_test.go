package easmproto

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestExtensionMessagesUseOnePublicPackage(t *testing.T) {
	model := &Ip{Ip: "192.0.2.1"}
	packed, err := anypb.New(model)
	if err != nil {
		t.Fatal(err)
	}
	if packed.TypeUrl != "type.googleapis.com/easm.Ip" {
		t.Fatalf("unexpected EASM Any type URL: %q", packed.TypeUrl)
	}
	wire, err := proto.Marshal(model)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Ip
	if err := proto.Unmarshal(wire, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.GetIp() != model.GetIp() {
		t.Fatalf("round-trip changed model: %q", decoded.GetIp())
	}
}
