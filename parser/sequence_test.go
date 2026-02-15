package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestSequenceParticipants(t *testing.T) {
	input := `sequenceDiagram
    participant A as Alice
    actor B as Bob
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Sequence {
		t.Fatalf("expected Sequence, got %v", g.Kind)
	}
	if len(g.Participants) != 2 {
		t.Fatalf("expected 2 participants, got %d", len(g.Participants))
	}

	alice := g.Participants[0]
	if alice.ID != "A" {
		t.Errorf("expected ID 'A', got %q", alice.ID)
	}
	if alice.Alias != "Alice" {
		t.Errorf("expected Alias 'Alice', got %q", alice.Alias)
	}
	if alice.Kind != ir.ParticipantBox {
		t.Errorf("expected ParticipantBox, got %v", alice.Kind)
	}

	bob := g.Participants[1]
	if bob.ID != "B" {
		t.Errorf("expected ID 'B', got %q", bob.ID)
	}
	if bob.Alias != "Bob" {
		t.Errorf("expected Alias 'Bob', got %q", bob.Alias)
	}
	if bob.Kind != ir.ActorStickFigure {
		t.Errorf("expected ActorStickFigure, got %v", bob.Kind)
	}
}

func TestSequenceAllMessageTypes(t *testing.T) {
	input := `sequenceDiagram
    participant A
    participant B
    A->B: solid
    A-->B: dotted
    A->>B: solid arrow
    A-->>B: dotted arrow
    A-xB: solid cross
    A--xB: dotted cross
    A-)B: solid open
    A--)B: dotted open
    A<<->>B: bi solid
    A<<-->>B: bi dotted
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := out.Graph

	expectedKinds := []ir.SeqMessageKind{
		ir.MsgSolid,
		ir.MsgDotted,
		ir.MsgSolidArrow,
		ir.MsgDottedArrow,
		ir.MsgSolidCross,
		ir.MsgDottedCross,
		ir.MsgSolidOpen,
		ir.MsgDottedOpen,
		ir.MsgBiSolid,
		ir.MsgBiDotted,
	}

	expectedTexts := []string{
		"solid",
		"dotted",
		"solid arrow",
		"dotted arrow",
		"solid cross",
		"dotted cross",
		"solid open",
		"dotted open",
		"bi solid",
		"bi dotted",
	}

	// Count message events.
	var msgEvents []*ir.SeqEvent
	for _, ev := range g.Events {
		if ev.Kind == ir.EvMessage {
			msgEvents = append(msgEvents, ev)
		}
	}

	if len(msgEvents) != len(expectedKinds) {
		t.Fatalf("expected %d message events, got %d", len(expectedKinds), len(msgEvents))
	}

	for i, ev := range msgEvents {
		if ev.Message.Kind != expectedKinds[i] {
			t.Errorf("message %d: expected kind %v, got %v", i, expectedKinds[i], ev.Message.Kind)
		}
		if ev.Message.Text != expectedTexts[i] {
			t.Errorf("message %d: expected text %q, got %q", i, expectedTexts[i], ev.Message.Text)
		}
		if ev.Message.From != "A" {
			t.Errorf("message %d: expected From 'A', got %q", i, ev.Message.From)
		}
		if ev.Message.To != "B" {
			t.Errorf("message %d: expected To 'B', got %q", i, ev.Message.To)
		}
	}
}

func TestSequenceActivationShorthand(t *testing.T) {
	input := `sequenceDiagram
    participant Alice
    participant Bob
    Alice->>+Bob: Hello
    Bob-->>-Alice: Done
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := out.Graph

	// Expect: EvMessage(Alice->Bob) + EvActivate(Bob) + EvMessage(Bob->Alice) + EvDeactivate(Bob)
	if len(g.Events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(g.Events))
	}

	// First event: message
	if g.Events[0].Kind != ir.EvMessage {
		t.Errorf("event 0: expected EvMessage, got %v", g.Events[0].Kind)
	}
	if !g.Events[0].Message.ActivateTarget {
		t.Error("event 0: expected ActivateTarget=true")
	}

	// Second event: activate Bob
	if g.Events[1].Kind != ir.EvActivate {
		t.Errorf("event 1: expected EvActivate, got %v", g.Events[1].Kind)
	}
	if g.Events[1].Target != "Bob" {
		t.Errorf("event 1: expected target 'Bob', got %q", g.Events[1].Target)
	}

	// Third event: message
	if g.Events[2].Kind != ir.EvMessage {
		t.Errorf("event 2: expected EvMessage, got %v", g.Events[2].Kind)
	}
	if !g.Events[2].Message.DeactivateSource {
		t.Error("event 2: expected DeactivateSource=true")
	}

	// Fourth event: deactivate Bob
	if g.Events[3].Kind != ir.EvDeactivate {
		t.Errorf("event 3: expected EvDeactivate, got %v", g.Events[3].Kind)
	}
	if g.Events[3].Target != "Bob" {
		t.Errorf("event 3: expected target 'Bob', got %q", g.Events[3].Target)
	}
}

func TestSequenceNotes(t *testing.T) {
	input := `sequenceDiagram
    participant Alice
    participant Bob
    Note right of Alice: Hello there
    Note left of Bob: Goodbye
    Note over Alice,Bob: spanning text
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := out.Graph

	var noteEvents []*ir.SeqEvent
	for _, ev := range g.Events {
		if ev.Kind == ir.EvNote {
			noteEvents = append(noteEvents, ev)
		}
	}

	if len(noteEvents) != 3 {
		t.Fatalf("expected 3 note events, got %d", len(noteEvents))
	}

	// right of Alice
	if noteEvents[0].Note.Position != ir.NoteRight {
		t.Errorf("note 0: expected NoteRight, got %v", noteEvents[0].Note.Position)
	}
	if len(noteEvents[0].Note.Participants) != 1 || noteEvents[0].Note.Participants[0] != "Alice" {
		t.Errorf("note 0: expected participants [Alice], got %v", noteEvents[0].Note.Participants)
	}
	if noteEvents[0].Note.Text != "Hello there" {
		t.Errorf("note 0: expected text 'Hello there', got %q", noteEvents[0].Note.Text)
	}

	// left of Bob
	if noteEvents[1].Note.Position != ir.NoteLeft {
		t.Errorf("note 1: expected NoteLeft, got %v", noteEvents[1].Note.Position)
	}
	if len(noteEvents[1].Note.Participants) != 1 || noteEvents[1].Note.Participants[0] != "Bob" {
		t.Errorf("note 1: expected participants [Bob], got %v", noteEvents[1].Note.Participants)
	}

	// over Alice,Bob
	if noteEvents[2].Note.Position != ir.NoteOver {
		t.Errorf("note 2: expected NoteOver, got %v", noteEvents[2].Note.Position)
	}
	if len(noteEvents[2].Note.Participants) != 2 {
		t.Fatalf("note 2: expected 2 participants, got %d", len(noteEvents[2].Note.Participants))
	}
	if noteEvents[2].Note.Participants[0] != "Alice" || noteEvents[2].Note.Participants[1] != "Bob" {
		t.Errorf("note 2: expected participants [Alice, Bob], got %v", noteEvents[2].Note.Participants)
	}
	if noteEvents[2].Note.Text != "spanning text" {
		t.Errorf("note 2: expected text 'spanning text', got %q", noteEvents[2].Note.Text)
	}
}

func TestSequenceFrames(t *testing.T) {
	input := `sequenceDiagram
    participant A
    participant B
    loop Every minute
        A->>B: ping
    end
    alt Success
        B->>A: ok
    else Failure
        B->>A: error
    end
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := out.Graph

	// Expected events:
	// 1. EvFrameStart(loop, "Every minute")
	// 2. EvMessage
	// 3. EvFrameEnd
	// 4. EvFrameStart(alt, "Success")
	// 5. EvMessage
	// 6. EvFrameMiddle(else, "Failure")
	// 7. EvMessage
	// 8. EvFrameEnd
	if len(g.Events) != 8 {
		t.Fatalf("expected 8 events, got %d", len(g.Events))
	}

	// loop frame start
	if g.Events[0].Kind != ir.EvFrameStart {
		t.Errorf("event 0: expected EvFrameStart, got %v", g.Events[0].Kind)
	}
	if g.Events[0].Frame.Kind != ir.FrameLoop {
		t.Errorf("event 0: expected FrameLoop, got %v", g.Events[0].Frame.Kind)
	}
	if g.Events[0].Frame.Label != "Every minute" {
		t.Errorf("event 0: expected label 'Every minute', got %q", g.Events[0].Frame.Label)
	}

	// loop end
	if g.Events[2].Kind != ir.EvFrameEnd {
		t.Errorf("event 2: expected EvFrameEnd, got %v", g.Events[2].Kind)
	}

	// alt frame start
	if g.Events[3].Kind != ir.EvFrameStart {
		t.Errorf("event 3: expected EvFrameStart, got %v", g.Events[3].Kind)
	}
	if g.Events[3].Frame.Kind != ir.FrameAlt {
		t.Errorf("event 3: expected FrameAlt, got %v", g.Events[3].Frame.Kind)
	}

	// else middle
	if g.Events[5].Kind != ir.EvFrameMiddle {
		t.Errorf("event 5: expected EvFrameMiddle, got %v", g.Events[5].Kind)
	}
	if g.Events[5].Frame.Label != "Failure" {
		t.Errorf("event 5: expected label 'Failure', got %q", g.Events[5].Frame.Label)
	}

	// alt end
	if g.Events[7].Kind != ir.EvFrameEnd {
		t.Errorf("event 7: expected EvFrameEnd, got %v", g.Events[7].Kind)
	}
}

func TestSequenceBoxes(t *testing.T) {
	input := `sequenceDiagram
    box Purple Team A
        participant A as Alice
        participant B as Bob
    end
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := out.Graph

	if len(g.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(g.Boxes))
	}

	box := g.Boxes[0]
	if box.Label != "Team A" {
		t.Errorf("expected label 'Team A', got %q", box.Label)
	}
	if box.Color != "Purple" {
		t.Errorf("expected color 'Purple', got %q", box.Color)
	}
	if len(box.Participants) != 2 {
		t.Fatalf("expected 2 participants in box, got %d", len(box.Participants))
	}
	if box.Participants[0] != "A" || box.Participants[1] != "B" {
		t.Errorf("expected participants [A, B], got %v", box.Participants)
	}
}

func TestSequenceAutonumber(t *testing.T) {
	input := `sequenceDiagram
    autonumber
    participant A
    participant B
    A->>B: hello
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Graph.Autonumber {
		t.Error("expected Autonumber=true")
	}
}

func TestSequenceImplicitParticipants(t *testing.T) {
	input := `sequenceDiagram
    Alice->>Bob: Hello
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := out.Graph

	if len(g.Participants) != 2 {
		t.Fatalf("expected 2 implicit participants, got %d", len(g.Participants))
	}
	if g.Participants[0].ID != "Alice" {
		t.Errorf("expected first participant 'Alice', got %q", g.Participants[0].ID)
	}
	if g.Participants[1].ID != "Bob" {
		t.Errorf("expected second participant 'Bob', got %q", g.Participants[1].ID)
	}

	// Implicit participants should have no alias.
	if g.Participants[0].Alias != "" {
		t.Errorf("expected empty alias for Alice, got %q", g.Participants[0].Alias)
	}
}

func TestSequenceCreateDestroy(t *testing.T) {
	input := `sequenceDiagram
    participant Alice
    create participant Carl
    Alice->>Carl: Hi Carl
    destroy Carl
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := out.Graph

	// Find Carl.
	var carl *ir.SeqParticipant
	for _, p := range g.Participants {
		if p.ID == "Carl" {
			carl = p
			break
		}
	}
	if carl == nil {
		t.Fatal("expected Carl to be in participants")
	}
	if !carl.IsCreated {
		t.Error("expected Carl.IsCreated=true")
	}
	if !carl.IsDestroyed {
		t.Error("expected Carl.IsDestroyed=true")
	}

	// Verify EvCreate and EvDestroy events exist.
	var hasCreate, hasDestroy bool
	for _, ev := range g.Events {
		if ev.Kind == ir.EvCreate && ev.Target == "Carl" {
			hasCreate = true
		}
		if ev.Kind == ir.EvDestroy && ev.Target == "Carl" {
			hasDestroy = true
		}
	}
	if !hasCreate {
		t.Error("expected EvCreate event for Carl")
	}
	if !hasDestroy {
		t.Error("expected EvDestroy event for Carl")
	}
}

func TestSequenceLinks(t *testing.T) {
	input := `sequenceDiagram
    participant Alice
    link Alice: Dashboard @ https://example.com/dashboard
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := out.Graph

	if len(g.Participants) != 1 {
		t.Fatalf("expected 1 participant, got %d", len(g.Participants))
	}
	alice := g.Participants[0]
	if len(alice.Links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(alice.Links))
	}
	if alice.Links[0].Label != "Dashboard" {
		t.Errorf("expected link label 'Dashboard', got %q", alice.Links[0].Label)
	}
	if alice.Links[0].URL != "https://example.com/dashboard" {
		t.Errorf("expected link URL 'https://example.com/dashboard', got %q", alice.Links[0].URL)
	}
}

func TestSequenceLineBreaks(t *testing.T) {
	input := `sequenceDiagram
    participant A
    participant B
    A->>B: Hello<br/>World
    Note right of A: Line1<br>Line2
`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := out.Graph

	var msgEvents []*ir.SeqEvent
	var noteEvents []*ir.SeqEvent
	for _, ev := range g.Events {
		if ev.Kind == ir.EvMessage {
			msgEvents = append(msgEvents, ev)
		}
		if ev.Kind == ir.EvNote {
			noteEvents = append(noteEvents, ev)
		}
	}

	if len(msgEvents) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgEvents))
	}
	if msgEvents[0].Message.Text != "Hello\nWorld" {
		t.Errorf("expected 'Hello\\nWorld', got %q", msgEvents[0].Message.Text)
	}

	if len(noteEvents) != 1 {
		t.Fatalf("expected 1 note, got %d", len(noteEvents))
	}
	if noteEvents[0].Note.Text != "Line1\nLine2" {
		t.Errorf("expected 'Line1\\nLine2', got %q", noteEvents[0].Note.Text)
	}
}
