package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
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
	g := out.Graph
	if g.Kind != ir.Packet {
		t.Fatalf("Kind = %v, want Packet", g.Kind)
	}
	if len(g.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(g.Fields))
	}
	if g.Fields[0].Start != 0 || g.Fields[0].End != 15 {
		t.Errorf("Fields[0] = %d-%d, want 0-15", g.Fields[0].Start, g.Fields[0].End)
	}
	if g.Fields[0].Description != "Source Port" {
		t.Errorf("Fields[0].Description = %q, want \"Source Port\"", g.Fields[0].Description)
	}
	if g.Fields[2].Start != 32 || g.Fields[2].End != 63 {
		t.Errorf("Fields[2] = %d-%d, want 32-63", g.Fields[2].Start, g.Fields[2].End)
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
	g := out.Graph
	if len(g.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(g.Fields))
	}
	if g.Fields[0].Start != 0 || g.Fields[0].End != 15 {
		t.Errorf("Fields[0] = %d-%d, want 0-15", g.Fields[0].Start, g.Fields[0].End)
	}
	if g.Fields[1].Start != 16 || g.Fields[1].End != 31 {
		t.Errorf("Fields[1] = %d-%d, want 16-31", g.Fields[1].Start, g.Fields[1].End)
	}
	if g.Fields[2].Start != 32 || g.Fields[2].End != 63 {
		t.Errorf("Fields[2] = %d-%d, want 32-63", g.Fields[2].Start, g.Fields[2].End)
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
