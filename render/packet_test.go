package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderPacketContainsFields(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Packet
	g.Fields = []*ir.PacketField{
		{Start: 0, End: 15, Description: "Source Port"},
		{Start: 16, End: 31, Description: "Dest Port"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "Source Port") {
		t.Error("SVG should contain 'Source Port'")
	}
	if !strings.Contains(svg, "Dest Port") {
		t.Error("SVG should contain 'Dest Port'")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("SVG should contain rect elements for fields")
	}
}

func TestRenderPacketBitNumbers(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Packet
	g.Fields = []*ir.PacketField{
		{Start: 0, End: 31, Description: "Data"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// Should contain bit number "0" and "31"
	if !strings.Contains(svg, ">0<") {
		t.Error("SVG should contain bit number '0'")
	}
	if !strings.Contains(svg, ">31<") {
		t.Error("SVG should contain bit number '31'")
	}
}

func TestRenderPacketValidSVG(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Packet
	g.Fields = []*ir.PacketField{
		{Start: 0, End: 15, Description: "Field"},
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
