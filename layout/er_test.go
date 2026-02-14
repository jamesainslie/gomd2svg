package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestComputeERLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Er
	g.Direction = ir.TopDown

	g.EnsureNode("CUSTOMER", nil, nil)
	g.EnsureNode("ORDER", nil, nil)
	g.Entities["CUSTOMER"] = &ir.Entity{
		ID: "CUSTOMER",
		Attributes: []ir.EntityAttribute{
			{Type: "string", Name: "name"},
			{Type: "int", Name: "id", Keys: []ir.AttributeKey{ir.KeyPrimary}},
		},
	}
	g.Entities["ORDER"] = &ir.Entity{
		ID: "ORDER",
		Attributes: []ir.EntityAttribute{
			{Type: "int", Name: "id", Keys: []ir.AttributeKey{ir.KeyPrimary}},
		},
	}
	g.Edges = append(g.Edges, &ir.Edge{From: "CUSTOMER", To: "ORDER"})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Er {
		t.Errorf("Kind = %v, want Er", l.Kind)
	}
	if len(l.Nodes) != 2 {
		t.Errorf("nodes = %d, want 2", len(l.Nodes))
	}
	cust := l.Nodes["CUSTOMER"]
	order := l.Nodes["ORDER"]
	if cust.Height <= order.Height {
		t.Errorf("CUSTOMER height (%f) should be > ORDER height (%f)", cust.Height, order.Height)
	}
	if _, ok := l.Diagram.(ERData); !ok {
		t.Errorf("Diagram data type = %T, want ERData", l.Diagram)
	}
}
