package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParsePieBasic(t *testing.T) {
	input := `pie
    title Pets adopted by volunteers
    "Dogs" : 386
    "Cats" : 85
    "Rats" : 15`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if out.Graph.Kind != ir.Pie {
		t.Errorf("Kind = %v, want Pie", out.Graph.Kind)
	}
	if out.Graph.PieTitle != "Pets adopted by volunteers" {
		t.Errorf("PieTitle = %q", out.Graph.PieTitle)
	}
	if out.Graph.PieShowData {
		t.Error("PieShowData = true, want false")
	}
	if len(out.Graph.PieSlices) != 3 {
		t.Fatalf("PieSlices = %d, want 3", len(out.Graph.PieSlices))
	}
	if out.Graph.PieSlices[0].Label != "Dogs" || out.Graph.PieSlices[0].Value != 386 {
		t.Errorf("slice[0] = %+v", out.Graph.PieSlices[0])
	}
}

func TestParsePieShowData(t *testing.T) {
	input := `pie showData
    title Budget
    "Engineering" : 45.50
    "Sales" : 25.25`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if !out.Graph.PieShowData {
		t.Error("PieShowData = false, want true")
	}
	if len(out.Graph.PieSlices) != 2 {
		t.Fatalf("PieSlices = %d, want 2", len(out.Graph.PieSlices))
	}
	if out.Graph.PieSlices[0].Value != 45.50 {
		t.Errorf("slice[0].Value = %f, want 45.5", out.Graph.PieSlices[0].Value)
	}
}

func TestParsePieComments(t *testing.T) {
	input := `pie
    %% This is a comment
    "A" : 10
    "B" : 20 %% trailing comment`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.PieSlices) != 2 {
		t.Fatalf("PieSlices = %d, want 2", len(out.Graph.PieSlices))
	}
}

func TestParsePieNoTitle(t *testing.T) {
	input := `pie
    "X" : 50
    "Y" : 50`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if out.Graph.PieTitle != "" {
		t.Errorf("PieTitle = %q, want empty", out.Graph.PieTitle)
	}
	if len(out.Graph.PieSlices) != 2 {
		t.Fatalf("PieSlices = %d, want 2", len(out.Graph.PieSlices))
	}
}
