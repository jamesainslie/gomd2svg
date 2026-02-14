package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestComputeClassLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Class
	g.Direction = ir.TopDown

	g.EnsureNode("Animal", nil, nil)
	g.EnsureNode("Dog", nil, nil)
	g.Members["Animal"] = &ir.ClassMembers{
		Attributes: []ir.ClassMember{
			{Name: "name", Type: "String", Visibility: ir.VisPublic},
		},
		Methods: []ir.ClassMember{
			{Name: "speak", IsMethod: true, Visibility: ir.VisPublic, Type: "void"},
		},
	}
	g.Edges = append(g.Edges, &ir.Edge{From: "Dog", To: "Animal", Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Class {
		t.Errorf("Kind = %v, want Class", l.Kind)
	}
	if len(l.Nodes) != 2 {
		t.Errorf("nodes = %d, want 2", len(l.Nodes))
	}
	if len(l.Edges) != 1 {
		t.Errorf("edges = %d, want 1", len(l.Edges))
	}
	animal := l.Nodes["Animal"]
	dog := l.Nodes["Dog"]
	if animal.Height <= dog.Height {
		t.Errorf("Animal height (%f) should be > Dog height (%f)", animal.Height, dog.Height)
	}
	if _, ok := l.Diagram.(ClassData); !ok {
		t.Errorf("Diagram data type = %T, want ClassData", l.Diagram)
	}
}
