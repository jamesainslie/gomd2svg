package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestSequenceLayoutParticipantPositions(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Sequence
	g.Participants = []*ir.SeqParticipant{
		{ID: "A", Alias: "Alice"},
		{ID: "B", Alias: "Bob"},
		{ID: "C", Alias: "Charlie"},
	}
	g.Events = []*ir.SeqEvent{
		{Kind: ir.EvMessage, Message: &ir.SeqMessage{From: "A", To: "B", Text: "hello", Kind: ir.MsgSolid}},
	}

	l := ComputeLayout(g, theme.Modern(), config.DefaultLayout())
	sd, ok := l.Diagram.(SequenceData)
	if !ok {
		t.Fatal("expected SequenceData diagram type")
	}

	if len(sd.Participants) != 3 {
		t.Fatalf("expected 3 participants, got %d", len(sd.Participants))
	}

	// X positions must monotonically increase.
	for i := 1; i < len(sd.Participants); i++ {
		if sd.Participants[i].X <= sd.Participants[i-1].X {
			t.Errorf("participant %d X (%f) not greater than participant %d X (%f)",
				i, sd.Participants[i].X, i-1, sd.Participants[i-1].X)
		}
	}

	// Verify all participant IDs are present.
	ids := map[string]bool{}
	for _, p := range sd.Participants {
		ids[p.ID] = true
	}
	for _, want := range []string{"A", "B", "C"} {
		if !ids[want] {
			t.Errorf("missing participant %s", want)
		}
	}
}

func TestSequenceLayoutMessagePositions(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Sequence
	g.Participants = []*ir.SeqParticipant{
		{ID: "A", Alias: "Alice"},
		{ID: "B", Alias: "Bob"},
	}
	g.Events = []*ir.SeqEvent{
		{Kind: ir.EvMessage, Message: &ir.SeqMessage{From: "A", To: "B", Text: "first", Kind: ir.MsgSolid}},
		{Kind: ir.EvMessage, Message: &ir.SeqMessage{From: "B", To: "A", Text: "second", Kind: ir.MsgDotted}},
	}

	l := ComputeLayout(g, theme.Modern(), config.DefaultLayout())
	sd, ok := l.Diagram.(SequenceData)
	if !ok {
		t.Fatal("expected SequenceData diagram type")
	}

	if len(sd.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(sd.Messages))
	}

	if sd.Messages[1].Y <= sd.Messages[0].Y {
		t.Errorf("second message Y (%f) should be greater than first message Y (%f)",
			sd.Messages[1].Y, sd.Messages[0].Y)
	}
}

func TestSequenceLayoutSelfMessage(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Sequence
	g.Participants = []*ir.SeqParticipant{
		{ID: "A", Alias: "Alice"},
	}
	g.Events = []*ir.SeqEvent{
		{Kind: ir.EvMessage, Message: &ir.SeqMessage{From: "A", To: "A", Text: "self", Kind: ir.MsgSolid}},
	}

	l := ComputeLayout(g, theme.Modern(), config.DefaultLayout())
	sd, ok := l.Diagram.(SequenceData)
	if !ok {
		t.Fatal("expected SequenceData diagram type")
	}

	if len(sd.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(sd.Messages))
	}

	msg := sd.Messages[0]
	if msg.ToX <= msg.FromX {
		t.Errorf("self-message ToX (%f) should be greater than FromX (%f)", msg.ToX, msg.FromX)
	}
}

func TestSequenceLayoutActivations(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Sequence
	g.Participants = []*ir.SeqParticipant{
		{ID: "A", Alias: "Alice"},
		{ID: "B", Alias: "Bob"},
	}
	g.Events = []*ir.SeqEvent{
		{Kind: ir.EvActivate, Target: "B"},
		{Kind: ir.EvMessage, Message: &ir.SeqMessage{From: "A", To: "B", Text: "call", Kind: ir.MsgSolid}},
		{Kind: ir.EvDeactivate, Target: "B"},
	}

	l := ComputeLayout(g, theme.Modern(), config.DefaultLayout())
	sd, ok := l.Diagram.(SequenceData)
	if !ok {
		t.Fatal("expected SequenceData diagram type")
	}

	if len(sd.Activations) != 1 {
		t.Fatalf("expected 1 activation, got %d", len(sd.Activations))
	}

	act := sd.Activations[0]
	if act.ParticipantID != "B" {
		t.Errorf("expected activation on B, got %s", act.ParticipantID)
	}
	if act.BottomY <= act.TopY {
		t.Errorf("activation BottomY (%f) should be greater than TopY (%f)", act.BottomY, act.TopY)
	}
}

func TestSequenceLayoutFrames(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Sequence
	g.Participants = []*ir.SeqParticipant{
		{ID: "A", Alias: "Alice"},
		{ID: "B", Alias: "Bob"},
	}
	g.Events = []*ir.SeqEvent{
		{Kind: ir.EvFrameStart, Frame: &ir.SeqFrame{Kind: ir.FrameLoop, Label: "retry"}},
		{Kind: ir.EvMessage, Message: &ir.SeqMessage{From: "A", To: "B", Text: "ping", Kind: ir.MsgSolid}},
		{Kind: ir.EvFrameEnd},
	}

	l := ComputeLayout(g, theme.Modern(), config.DefaultLayout())
	sd, ok := l.Diagram.(SequenceData)
	if !ok {
		t.Fatal("expected SequenceData diagram type")
	}

	if len(sd.Frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(sd.Frames))
	}

	frame := sd.Frames[0]
	if frame.Height <= 0 {
		t.Errorf("frame Height (%f) should be > 0", frame.Height)
	}
	if frame.Kind != ir.FrameLoop {
		t.Errorf("expected FrameLoop, got %v", frame.Kind)
	}
	if frame.Label != "retry" {
		t.Errorf("expected label 'retry', got %q", frame.Label)
	}
}

func TestSequenceLayoutCreatedParticipant(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Sequence
	g.Participants = []*ir.SeqParticipant{
		{ID: "A", Alias: "Alice"},
		{ID: "B", Alias: "Bob", IsCreated: true},
	}
	g.Events = []*ir.SeqEvent{
		{Kind: ir.EvMessage, Message: &ir.SeqMessage{From: "A", To: "A", Text: "init", Kind: ir.MsgSolid}},
		{Kind: ir.EvCreate, Target: "B"},
		{Kind: ir.EvMessage, Message: &ir.SeqMessage{From: "A", To: "B", Text: "new", Kind: ir.MsgSolid}},
	}

	l := ComputeLayout(g, theme.Modern(), config.DefaultLayout())
	sd, ok := l.Diagram.(SequenceData)
	if !ok {
		t.Fatal("expected SequenceData diagram type")
	}

	// Find the created participant B.
	var bLayout *SeqParticipantLayout
	for i := range sd.Participants {
		if sd.Participants[i].ID == "B" {
			bLayout = &sd.Participants[i]
			break
		}
	}
	if bLayout == nil {
		t.Fatal("participant B not found")
	}

	if bLayout.Y <= 0 {
		t.Errorf("created participant B Y (%f) should be > 0", bLayout.Y)
	}
}
