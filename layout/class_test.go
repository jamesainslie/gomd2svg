package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestComputeClassLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Class
	graph.Direction = ir.TopDown

	graph.EnsureNode("Animal", nil, nil)
	graph.EnsureNode("Dog", nil, nil)
	graph.Members["Animal"] = &ir.ClassMembers{
		Attributes: []ir.ClassMember{
			{Name: "name", Type: "String", Visibility: ir.VisPublic},
		},
		Methods: []ir.ClassMember{
			{Name: "speak", IsMethod: true, Visibility: ir.VisPublic, Type: "void"},
		},
	}
	graph.Edges = append(graph.Edges, &ir.Edge{From: "Dog", To: "Animal", Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Kind != ir.Class {
		t.Errorf("Kind = %v, want Class", lay.Kind)
	}
	if len(lay.Nodes) != 2 {
		t.Errorf("nodes = %d, want 2", len(lay.Nodes))
	}
	if len(lay.Edges) != 1 {
		t.Errorf("edges = %d, want 1", len(lay.Edges))
	}
	animal := lay.Nodes["Animal"]
	dog := lay.Nodes["Dog"]
	if animal.Height <= dog.Height {
		t.Errorf("Animal height (%f) should be > Dog height (%f)", animal.Height, dog.Height)
	}
	if _, ok := lay.Diagram.(ClassData); !ok {
		t.Errorf("Diagram data type = %T, want ClassData", lay.Diagram)
	}
}
