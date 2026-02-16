package ir

import "testing"

func TestPacketFieldRange(t *testing.T) {
	f := &PacketField{Start: 0, End: 15, Description: "Source Port"}
	if f.BitWidth() != 16 {
		t.Errorf("BitWidth() = %d, want 16", f.BitWidth())
	}
}

func TestGraphPacketFields(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Packet
	graph.Fields = []*PacketField{
		{Start: 0, End: 15, Description: "Source Port"},
		{Start: 16, End: 31, Description: "Dest Port"},
	}
	if len(graph.Fields) != 2 {
		t.Fatalf("len(Fields) = %d, want 2", len(graph.Fields))
	}
}
