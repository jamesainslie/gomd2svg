package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// buildSeqGraph creates a sequence diagram ir.Graph with the given participants,
// events, and autonumber flag. Participants are created as ParticipantBox by default.
func buildSeqGraph(names []string, events []*ir.SeqEvent) *ir.Graph {
	graph := ir.NewGraph()
	graph.Kind = ir.Sequence
	for _, name := range names {
		graph.Participants = append(graph.Participants, &ir.SeqParticipant{
			ID:   name,
			Kind: ir.ParticipantBox,
		})
	}
	graph.Events = events
	return graph
}

func renderSeqSVG(seqGraph *ir.Graph) string {
	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := layout.ComputeLayout(seqGraph, th, cfg)
	return RenderSVG(lay, th, cfg)
}

func TestRenderSequenceHasLifelines(t *testing.T) {
	graph := buildSeqGraph(
		[]string{"Alice", "Bob"},
		[]*ir.SeqEvent{
			{Kind: ir.EvMessage, Message: &ir.SeqMessage{
				From: "Alice", To: "Bob", Text: "Hi", Kind: ir.MsgSolidArrow,
			}},
		},
	)

	svg := renderSeqSVG(graph)

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
	graph := buildSeqGraph(
		[]string{"Alice", "Bob"},
		[]*ir.SeqEvent{
			{Kind: ir.EvMessage, Message: &ir.SeqMessage{
				From: "Alice", To: "Bob", Text: "msg", Kind: ir.MsgSolid,
			}},
		},
	)

	svg := renderSeqSVG(graph)

	if !strings.Contains(svg, "Alice") {
		t.Error("expected SVG to contain participant label 'Alice'")
	}
	if !strings.Contains(svg, "Bob") {
		t.Error("expected SVG to contain participant label 'Bob'")
	}
}

func TestRenderSequenceHasMessageText(t *testing.T) {
	graph := buildSeqGraph(
		[]string{"Alice", "Bob"},
		[]*ir.SeqEvent{
			{Kind: ir.EvMessage, Message: &ir.SeqMessage{
				From: "Alice", To: "Bob", Text: "Hello World", Kind: ir.MsgSolidArrow,
			}},
		},
	)

	svg := renderSeqSVG(graph)

	if !strings.Contains(svg, "Hello World") {
		t.Error("expected SVG to contain message text 'Hello World'")
	}
}

func TestRenderSequenceHasActivations(t *testing.T) {
	th := theme.Modern()

	graph := buildSeqGraph(
		[]string{"Alice", "Bob"},
		[]*ir.SeqEvent{
			{Kind: ir.EvActivate, Target: "Bob"},
			{Kind: ir.EvMessage, Message: &ir.SeqMessage{
				From: "Alice", To: "Bob", Text: "call", Kind: ir.MsgSolidArrow,
			}},
			{Kind: ir.EvDeactivate, Target: "Bob"},
		},
	)

	svg := renderSeqSVG(graph)

	if !strings.Contains(svg, th.ActivationBackground) {
		t.Errorf("expected SVG to contain activation fill color %s", th.ActivationBackground)
	}
}
