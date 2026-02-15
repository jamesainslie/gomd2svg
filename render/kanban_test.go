package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderKanbanContainsColumns(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Kanban
	g.Columns = []*ir.KanbanColumn{
		{ID: "todo", Label: "Todo", Cards: []*ir.KanbanCard{
			{ID: "t1", Label: "Task One"},
		}},
		{ID: "done", Label: "Done"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "Todo") {
		t.Error("SVG should contain column label 'Todo'")
	}
	if !strings.Contains(svg, "Done") {
		t.Error("SVG should contain column label 'Done'")
	}
	if !strings.Contains(svg, "Task One") {
		t.Error("SVG should contain card label 'Task One'")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("SVG should contain rect elements for cards")
	}
}

func TestRenderKanbanValidSVG(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Kanban
	g.Columns = []*ir.KanbanColumn{
		{ID: "col", Label: "Column", Cards: []*ir.KanbanCard{
			{ID: "c1", Label: "Card"},
		}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.HasPrefix(svg, "<svg") {
		t.Error("SVG should start with <svg")
	}
	if !strings.HasSuffix(strings.TrimSpace(svg), "</svg>") {
		t.Error("SVG should end with </svg>")
	}
}
