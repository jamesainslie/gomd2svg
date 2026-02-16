package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParsePacketRangeNotation(t *testing.T) {
	input := `packet
0-15: "Source Port"
16-31: "Destination Port"
32-63: "Sequence Number"`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	graph := out.Graph
	if graph.Kind != ir.Packet {
		t.Fatalf("Kind = %v, want Packet", graph.Kind)
	}
	if len(graph.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(graph.Fields))
	}
	if graph.Fields[0].Start != 0 || graph.Fields[0].End != 15 {
		t.Errorf("Fields[0] = %d-%d, want 0-15", graph.Fields[0].Start, graph.Fields[0].End)
	}
	if graph.Fields[0].Description != "Source Port" {
		t.Errorf("Fields[0].Description = %q, want \"Source Port\"", graph.Fields[0].Description)
	}
	if graph.Fields[2].Start != 32 || graph.Fields[2].End != 63 {
		t.Errorf("Fields[2] = %d-%d, want 32-63", graph.Fields[2].Start, graph.Fields[2].End)
	}
}

func TestParsePacketBitCountNotation(t *testing.T) {
	input := `packet
+16: "Source Port"
+16: "Destination Port"
+32: "Sequence Number"`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	graph := out.Graph
	if len(graph.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(graph.Fields))
	}
	if graph.Fields[0].Start != 0 || graph.Fields[0].End != 15 {
		t.Errorf("Fields[0] = %d-%d, want 0-15", graph.Fields[0].Start, graph.Fields[0].End)
	}
	if graph.Fields[1].Start != 16 || graph.Fields[1].End != 31 {
		t.Errorf("Fields[1] = %d-%d, want 16-31", graph.Fields[1].Start, graph.Fields[1].End)
	}
	if graph.Fields[2].Start != 32 || graph.Fields[2].End != 63 {
		t.Errorf("Fields[2] = %d-%d, want 32-63", graph.Fields[2].Start, graph.Fields[2].End)
	}
}

func TestParsePacketMixedNotation(t *testing.T) {
	input := `packet
0-15: "Source Port"
+16: "Destination Port"
32-63: "Sequence Number"`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(out.Graph.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(out.Graph.Fields))
	}
	// +16 after 0-15 should be 16-31
	f := out.Graph.Fields[1]
	if f.Start != 16 || f.End != 31 {
		t.Errorf("Fields[1] = %d-%d, want 16-31", f.Start, f.End)
	}
}

func TestParsePacketSingleBit(t *testing.T) {
	input := `packet
0-3: "Version"
+1: "Flag"
+1: "Flag2"`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(out.Graph.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(out.Graph.Fields))
	}
	// +1 after 0-3 should be 4-4
	if out.Graph.Fields[1].Start != 4 || out.Graph.Fields[1].End != 4 {
		t.Errorf("Fields[1] = %d-%d, want 4-4", out.Graph.Fields[1].Start, out.Graph.Fields[1].End)
	}
}
