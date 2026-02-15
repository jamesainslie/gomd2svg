package ir

import "testing"

func TestPacketFieldRange(t *testing.T) {
	f := &PacketField{Start: 0, End: 15, Description: "Source Port"}
	if f.BitWidth() != 16 {
		t.Errorf("BitWidth() = %d, want 16", f.BitWidth())
	}
}

func TestGraphPacketFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Packet
	g.Fields = []*PacketField{
		{Start: 0, End: 15, Description: "Source Port"},
		{Start: 16, End: 31, Description: "Dest Port"},
	}
	if len(g.Fields) != 2 {
		t.Fatalf("len(Fields) = %d, want 2", len(g.Fields))
	}
}
