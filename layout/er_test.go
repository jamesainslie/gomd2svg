package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestComputeERLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Er
	graph.Direction = ir.TopDown

	graph.EnsureNode("CUSTOMER", nil, nil)
	graph.EnsureNode("ORDER", nil, nil)
	graph.Entities["CUSTOMER"] = &ir.Entity{
		ID: "CUSTOMER",
		Attributes: []ir.EntityAttribute{
			{Type: "string", Name: "name"},
			{Type: "int", Name: "id", Keys: []ir.AttributeKey{ir.KeyPrimary}},
		},
	}
	graph.Entities["ORDER"] = &ir.Entity{
		ID: "ORDER",
		Attributes: []ir.EntityAttribute{
			{Type: "int", Name: "id", Keys: []ir.AttributeKey{ir.KeyPrimary}},
		},
	}
	graph.Edges = append(graph.Edges, &ir.Edge{From: "CUSTOMER", To: "ORDER"})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Kind != ir.Er {
		t.Errorf("Kind = %v, want Er", lay.Kind)
	}
	if len(lay.Nodes) != 2 {
		t.Errorf("nodes = %d, want 2", len(lay.Nodes))
	}
	cust := lay.Nodes["CUSTOMER"]
	order := lay.Nodes["ORDER"]
	if cust.Height <= order.Height {
		t.Errorf("CUSTOMER height (%f) should be > ORDER height (%f)", cust.Height, order.Height)
	}
	if _, ok := lay.Diagram.(ERData); !ok {
		t.Errorf("Diagram data type = %T, want ERData", lay.Diagram)
	}
}
