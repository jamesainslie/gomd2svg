package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseSankeyBasic(t *testing.T) {
	input := `sankey-beta

Solar,Grid,60
Wind,Grid,290
Grid,Industry,340
Grid,Homes,114`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	if graph.Kind != ir.Sankey {
		t.Fatalf("Kind = %v, want Sankey", graph.Kind)
	}
	if len(graph.SankeyLinks) != 4 {
		t.Fatalf("links len = %d, want 4", len(graph.SankeyLinks))
	}
	if graph.SankeyLinks[0].Source != "Solar" {
		t.Errorf("link[0] source = %q, want Solar", graph.SankeyLinks[0].Source)
	}
	if graph.SankeyLinks[0].Value != 60 {
		t.Errorf("link[0] value = %v, want 60", graph.SankeyLinks[0].Value)
	}
}

func TestParseSankeyQuotedNames(t *testing.T) {
	input := `sankey

"Source A","Target, B",100.5
"Source ""C""",Target D,200`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	if len(graph.SankeyLinks) != 2 {
		t.Fatalf("links len = %d, want 2", len(graph.SankeyLinks))
	}
	if graph.SankeyLinks[0].Target != "Target, B" {
		t.Errorf("link[0] target = %q, want %q", graph.SankeyLinks[0].Target, "Target, B")
	}
	if graph.SankeyLinks[0].Value != 100.5 {
		t.Errorf("link[0] value = %v, want 100.5", graph.SankeyLinks[0].Value)
	}
	if graph.SankeyLinks[1].Source != `Source "C"` {
		t.Errorf("link[1] source = %q, want %q", graph.SankeyLinks[1].Source, `Source "C"`)
	}
}

func TestParseSankeyMinimal(t *testing.T) {
	input := `sankey-beta
A,B,10`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.SankeyLinks) != 1 {
		t.Fatalf("links len = %d, want 1", len(out.Graph.SankeyLinks))
	}
}
