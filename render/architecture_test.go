package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderArchitecture(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Architecture

	dbLabel := "Database"
	srvLabel := "Server"
	graph.EnsureNode("db", &dbLabel, nil)
	graph.EnsureNode("srv", &srvLabel, nil)

	graph.ArchServices = []*ir.ArchService{
		{ID: "db", Label: "Database", Icon: "database"},
		{ID: "srv", Label: "Server", Icon: "server"},
	}
	graph.ArchGroups = []*ir.ArchGroup{
		{ID: "api", Label: "API", Icon: "cloud", Children: []string{"db", "srv"}},
	}
	graph.ArchEdges = []*ir.ArchEdge{
		{FromID: "db", FromSide: ir.ArchRight, ToID: "srv", ToSide: ir.ArchLeft, ArrowRight: true},
	}
	graph.Edges = append(graph.Edges, &ir.Edge{From: "db", To: "srv", Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := layout.ComputeLayout(graph, th, cfg)

	// Call renderArchitecture directly to verify our renderer works,
	// independent of svg.go dispatch (which another agent adds).
	var builder svgBuilder
	builder.openTag("svg",
		"xmlns", "http://www.w3.org/2000/svg",
		"width", "400",
		"height", "300",
	)
	renderArchitecture(&builder, lay, th, cfg)
	builder.closeTag("svg")
	svg := builder.String()

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Database") {
		t.Error("missing Database label")
	}
	if !strings.Contains(svg, "Server") {
		t.Error("missing Server label")
	}
	if !strings.Contains(svg, "API") {
		t.Error("missing group label API")
	}
	// Group should have a dashed border.
	if !strings.Contains(svg, "stroke-dasharray") {
		t.Error("missing dashed stroke on group")
	}
	// Should contain at least one rect for groups/services.
	if !strings.Contains(svg, "<rect") {
		t.Error("missing <rect elements")
	}
}

func TestRenderArchitectureJunctions(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Architecture

	aLabel := "ServiceA"
	bLabel := "ServiceB"
	graph.EnsureNode("a", &aLabel, nil)
	graph.EnsureNode("b", &bLabel, nil)
	graph.EnsureNode("j1", nil, nil)

	graph.ArchServices = []*ir.ArchService{
		{ID: "a", Label: "ServiceA"},
		{ID: "b", Label: "ServiceB"},
	}
	graph.ArchJunctions = []*ir.ArchJunction{
		{ID: "j1"},
	}
	graph.ArchEdges = []*ir.ArchEdge{
		{FromID: "a", FromSide: ir.ArchRight, ToID: "j1", ToSide: ir.ArchLeft},
		{FromID: "j1", FromSide: ir.ArchRight, ToID: "b", ToSide: ir.ArchLeft},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := layout.ComputeLayout(graph, th, cfg)

	var builder svgBuilder
	builder.openTag("svg", "xmlns", "http://www.w3.org/2000/svg", "width", "600", "height", "200")
	renderArchitecture(&builder, lay, th, cfg)
	builder.closeTag("svg")
	svg := builder.String()

	// Junction should be rendered as a circle.
	if !strings.Contains(svg, "<circle") {
		t.Error("missing <circle for junction")
	}
	if !strings.Contains(svg, "ServiceA") {
		t.Error("missing ServiceA label")
	}
	if !strings.Contains(svg, "ServiceB") {
		t.Error("missing ServiceB label")
	}
}

func TestRenderArchitectureIcons(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Architecture

	// Single service with each icon type.
	icons := []string{"database", "server", "cloud", "internet", "disk"}
	for _, icon := range icons {
		label := icon
		graph.EnsureNode(icon, &label, nil)
		graph.ArchServices = append(graph.ArchServices, &ir.ArchService{
			ID: icon, Label: label, Icon: icon,
		})
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := layout.ComputeLayout(graph, th, cfg)

	var builder svgBuilder
	builder.openTag("svg", "xmlns", "http://www.w3.org/2000/svg", "width", "800", "height", "200")
	renderArchitecture(&builder, lay, th, cfg)
	builder.closeTag("svg")
	svg := builder.String()

	// Each icon type should produce some SVG shape.
	if !strings.Contains(svg, "<ellipse") {
		t.Error("missing <ellipse for database/cloud icon")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing <circle for internet/disk icon")
	}
}

func TestRenderArchitectureEmpty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Architecture

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := layout.ComputeLayout(graph, th, cfg)

	var builder svgBuilder
	builder.openTag("svg", "xmlns", "http://www.w3.org/2000/svg", "width", "100", "height", "100")
	renderArchitecture(&builder, lay, th, cfg)
	builder.closeTag("svg")
	svg := builder.String()

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	// Empty diagram should not crash.
	if !strings.Contains(svg, "</svg>") {
		t.Error("missing closing </svg> tag")
	}
}

func TestRenderArchitectureEdgeArrows(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Architecture

	aLabel := "A"
	bLabel := "B"
	graph.EnsureNode("a", &aLabel, nil)
	graph.EnsureNode("b", &bLabel, nil)

	graph.ArchServices = []*ir.ArchService{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
	}
	graph.ArchEdges = []*ir.ArchEdge{
		{FromID: "a", FromSide: ir.ArchRight, ToID: "b", ToSide: ir.ArchLeft, ArrowLeft: true, ArrowRight: true},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := layout.ComputeLayout(graph, th, cfg)

	var builder svgBuilder
	builder.openTag("svg", "xmlns", "http://www.w3.org/2000/svg", "width", "400", "height", "200")
	renderArchitecture(&builder, lay, th, cfg)
	builder.closeTag("svg")
	svg := builder.String()

	// Bidirectional arrows should produce both marker-start and marker-end.
	if !strings.Contains(svg, "marker-end") {
		t.Error("missing marker-end for arrow")
	}
	if !strings.Contains(svg, "marker-start") {
		t.Error("missing marker-start for arrow")
	}
}

func TestRenderArchitectureNoDispatch(t *testing.T) {
	// Verify that renderArchitecture handles wrong diagram data gracefully.
	lay := &layout.Layout{
		Kind:    ir.Architecture,
		Nodes:   make(map[string]*layout.NodeLayout),
		Diagram: layout.GraphData{}, // wrong type on purpose
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()

	var builder svgBuilder
	builder.openTag("svg", "xmlns", "http://www.w3.org/2000/svg", "width", "100", "height", "100")
	renderArchitecture(&builder, lay, th, cfg)
	builder.closeTag("svg")
	svg := builder.String()

	// Should produce valid SVG even with wrong data type (early return).
	if !strings.Contains(svg, "</svg>") {
		t.Error("missing closing </svg> tag after wrong data type")
	}
}
