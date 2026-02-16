package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseTreemapBasic(t *testing.T) {
	input := `treemap-beta
"Root"
    "Section A"
        "Leaf 1": 30
        "Leaf 2": 50
    "Section B"
        "Leaf 3": 20`

	out, err := parseTreemap(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Treemap {
		t.Fatalf("Kind = %v, want Treemap", g.Kind)
	}
	if g.TreemapRoot == nil {
		t.Fatal("TreemapRoot is nil")
	}
	if g.TreemapRoot.Label != "Root" {
		t.Errorf("root label = %q, want Root", g.TreemapRoot.Label)
	}
	if len(g.TreemapRoot.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(g.TreemapRoot.Children))
	}
	secA := g.TreemapRoot.Children[0]
	if secA.Label != "Section A" {
		t.Errorf("secA label = %q", secA.Label)
	}
	if len(secA.Children) != 2 {
		t.Fatalf("secA children = %d, want 2", len(secA.Children))
	}
	if secA.Children[0].Value != 30 {
		t.Errorf("leaf1 value = %v, want 30", secA.Children[0].Value)
	}
}

func TestParseTreemapFlat(t *testing.T) {
	input := `treemap
"Root"
    "A": 10
    "B": 20
    "C": 30`

	out, err := parseTreemap(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if len(g.TreemapRoot.Children) != 3 {
		t.Fatalf("children = %d, want 3", len(g.TreemapRoot.Children))
	}
	total := g.TreemapRoot.TotalValue()
	if total != 60 {
		t.Errorf("total = %v, want 60", total)
	}
}

func TestParseTreemapTitle(t *testing.T) {
	input := `treemap-beta
    title Budget
"Operations"
    "Salaries": 700
    "Equipment": 200`

	out, err := parseTreemap(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if out.Graph.TreemapTitle != "Budget" {
		t.Errorf("title = %q, want Budget", out.Graph.TreemapTitle)
	}
}

func TestParseTreemapEmpty(t *testing.T) {
	input := `treemap-beta`

	out, err := parseTreemap(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if out.Graph.TreemapRoot != nil {
		t.Error("expected nil root for empty treemap")
	}
}

func TestParseTreemapClass(t *testing.T) {
	input := `treemap
"Root"
    "Leaf": 42 :::highlight`

	out, err := parseTreemap(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	leaf := out.Graph.TreemapRoot.Children[0]
	if leaf.Class != "highlight" {
		t.Errorf("class = %q, want highlight", leaf.Class)
	}
	if leaf.Value != 42 {
		t.Errorf("value = %v, want 42", leaf.Value)
	}
}

func TestParseTreemapNodeLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantLbl string
		wantVal float64
		wantHas bool
		wantCls string
	}{
		{"quoted with value", `"Sales": 100`, "Sales", 100, true, ""},
		{"quoted no value", `"Marketing"`, "Marketing", 0, false, ""},
		{"single quoted", `'R&D': 50`, "R&D", 50, true, ""},
		{"with class", `"Ops": 25 :::critical`, "Ops", 25, true, "critical"},
		{"comma separator", `"IT", 75`, "IT", 75, true, ""},
		{"no quotes", `Bare`, "", 0, false, ""},
		{"empty", ``, "", 0, false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lbl, val, has, cls := parseTreemapNodeLine(tt.line)
			if lbl != tt.wantLbl {
				t.Errorf("label = %q, want %q", lbl, tt.wantLbl)
			}
			if val != tt.wantVal {
				t.Errorf("value = %v, want %v", val, tt.wantVal)
			}
			if has != tt.wantHas {
				t.Errorf("hasValue = %v, want %v", has, tt.wantHas)
			}
			if cls != tt.wantCls {
				t.Errorf("class = %q, want %q", cls, tt.wantCls)
			}
		})
	}
}
