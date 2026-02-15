package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// buildSeqGraph creates a sequence diagram ir.Graph with the given participants,
// events, and autonumber flag. Participants are created as ParticipantBox by default.
func buildSeqGraph(names []string, events []*ir.SeqEvent, autonumber bool) *ir.Graph {
	g := ir.NewGraph()
	g.Kind = ir.Sequence
	for _, name := range names {
		g.Participants = append(g.Participants, &ir.SeqParticipant{
			ID:   name,
			Kind: ir.ParticipantBox,
		})
	}
	g.Events = events
	g.Autonumber = autonumber
	return g
}

func renderSeqSVG(g *ir.Graph) string {
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	return RenderSVG(l, th, cfg)
}

func TestRenderSequenceHasLifelines(t *testing.T) {
	g := buildSeqGraph(
		[]string{"Alice", "Bob"},
		[]*ir.SeqEvent{
			{Kind: ir.EvMessage, Message: &ir.SeqMessage{
				From: "Alice", To: "Bob", Text: "Hi", Kind: ir.MsgSolidArrow,
			}},
		},
		false,
	)

	svg := renderSeqSVG(g)

	// Lifelines are dashed vertical lines.
	if !strings.Contains(svg, "stroke-dasharray") {
		t.Error("expected SVG to contain stroke-dasharray for lifeline dashes")
	}
	// There should be two lifelines (one per participant).
	count := strings.Count(svg, `stroke-dasharray="5,5"`)
	if count < 2 {
		t.Errorf("expected at least 2 dashed lines for lifelines, got %d", count)
	}
}

func TestRenderSequenceHasParticipantLabels(t *testing.T) {
	g := buildSeqGraph(
		[]string{"Alice", "Bob"},
		[]*ir.SeqEvent{
			{Kind: ir.EvMessage, Message: &ir.SeqMessage{
				From: "Alice", To: "Bob", Text: "msg", Kind: ir.MsgSolid,
			}},
		},
		false,
	)

	svg := renderSeqSVG(g)

	if !strings.Contains(svg, "Alice") {
		t.Error("expected SVG to contain participant label 'Alice'")
	}
	if !strings.Contains(svg, "Bob") {
		t.Error("expected SVG to contain participant label 'Bob'")
	}
}

func TestRenderSequenceHasMessageText(t *testing.T) {
	g := buildSeqGraph(
		[]string{"Alice", "Bob"},
		[]*ir.SeqEvent{
			{Kind: ir.EvMessage, Message: &ir.SeqMessage{
				From: "Alice", To: "Bob", Text: "Hello World", Kind: ir.MsgSolidArrow,
			}},
		},
		false,
	)

	svg := renderSeqSVG(g)

	if !strings.Contains(svg, "Hello World") {
		t.Error("expected SVG to contain message text 'Hello World'")
	}
}

func TestRenderSequenceHasActivations(t *testing.T) {
	th := theme.Modern()

	g := buildSeqGraph(
		[]string{"Alice", "Bob"},
		[]*ir.SeqEvent{
			{Kind: ir.EvActivate, Target: "Bob"},
			{Kind: ir.EvMessage, Message: &ir.SeqMessage{
				From: "Alice", To: "Bob", Text: "call", Kind: ir.MsgSolidArrow,
			}},
			{Kind: ir.EvDeactivate, Target: "Bob"},
		},
		false,
	)

	svg := renderSeqSVG(g)

	if !strings.Contains(svg, th.ActivationBackground) {
		t.Errorf("expected SVG to contain activation fill color %s", th.ActivationBackground)
	}
}
