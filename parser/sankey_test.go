package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
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
	g := out.Graph
	if g.Kind != ir.Sankey {
		t.Fatalf("Kind = %v, want Sankey", g.Kind)
	}
	if len(g.SankeyLinks) != 4 {
		t.Fatalf("links len = %d, want 4", len(g.SankeyLinks))
	}
	if g.SankeyLinks[0].Source != "Solar" {
		t.Errorf("link[0] source = %q, want Solar", g.SankeyLinks[0].Source)
	}
	if g.SankeyLinks[0].Value != 60 {
		t.Errorf("link[0] value = %v, want 60", g.SankeyLinks[0].Value)
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
	g := out.Graph
	if len(g.SankeyLinks) != 2 {
		t.Fatalf("links len = %d, want 2", len(g.SankeyLinks))
	}
	if g.SankeyLinks[0].Target != "Target, B" {
		t.Errorf("link[0] target = %q, want %q", g.SankeyLinks[0].Target, "Target, B")
	}
	if g.SankeyLinks[0].Value != 100.5 {
		t.Errorf("link[0] value = %v, want 100.5", g.SankeyLinks[0].Value)
	}
	if g.SankeyLinks[1].Source != `Source "C"` {
		t.Errorf("link[1] source = %q, want %q", g.SankeyLinks[1].Source, `Source "C"`)
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
