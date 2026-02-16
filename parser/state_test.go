package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseStateSimple(t *testing.T) {
	input := `stateDiagram-v2
    [*] --> First
    First --> Second
    Second --> [*]`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph
	if g.Kind != ir.State {
		t.Errorf("Kind = %v, want State", g.Kind)
	}
	if len(g.Edges) != 3 {
		t.Errorf("Edges = %d, want 3", len(g.Edges))
	}

	// Verify start/end node mapping
	if g.Edges[0].From != "__start__" {
		t.Errorf("Edge[0].From = %q, want __start__", g.Edges[0].From)
	}
	if g.Edges[0].To != "First" {
		t.Errorf("Edge[0].To = %q, want First", g.Edges[0].To)
	}
	if g.Edges[2].To != "__end__" {
		t.Errorf("Edge[2].To = %q, want __end__", g.Edges[2].To)
	}
}

func TestParseStateDescription(t *testing.T) {
	input := `stateDiagram-v2
    s1 : This is state s1`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph
	desc, ok := g.StateDescriptions["s1"]
	if !ok {
		t.Fatal("StateDescriptions missing key s1")
	}
	if desc != "This is state s1" {
		t.Errorf("StateDescriptions[s1] = %q, want %q", desc, "This is state s1")
	}
}

func TestParseStateAsKeyword(t *testing.T) {
	input := `stateDiagram-v2
    state "Moving state" as s1`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph
	if _, ok := g.Nodes["s1"]; !ok {
		t.Fatal("expected node s1 to exist")
	}
	desc, ok := g.StateDescriptions["s1"]
	if !ok {
		t.Fatal("StateDescriptions missing key s1")
	}
	if desc != "Moving state" {
		t.Errorf("StateDescriptions[s1] = %q, want %q", desc, "Moving state")
	}
}

func TestParseStateComposite(t *testing.T) {
	input := `stateDiagram-v2
    state First {
        [*] --> fir
        fir --> [*]
    }`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph
	cs, ok := g.CompositeStates["First"]
	if !ok {
		t.Fatal("CompositeStates missing key First")
	}
	if cs.Inner == nil {
		t.Fatal("CompositeState.Inner is nil")
	}
	if len(cs.Inner.Edges) != 2 {
		t.Errorf("Inner edges = %d, want 2", len(cs.Inner.Edges))
	}
}

func TestParseStateChoice(t *testing.T) {
	input := `stateDiagram-v2
    state if_state <<choice>>`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph
	ann, ok := g.StateAnnotations["if_state"]
	if !ok {
		t.Fatal("StateAnnotations missing key if_state")
	}
	if ann != ir.StateChoice {
		t.Errorf("StateAnnotations[if_state] = %v, want StateChoice", ann)
	}
}

func TestParseStateForkJoin(t *testing.T) {
	input := `stateDiagram-v2
    state fork_state <<fork>>
    state join_state <<join>>`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph

	forkAnn, ok := g.StateAnnotations["fork_state"]
	if !ok {
		t.Fatal("StateAnnotations missing key fork_state")
	}
	if forkAnn != ir.StateFork {
		t.Errorf("StateAnnotations[fork_state] = %v, want StateFork", forkAnn)
	}

	joinAnn, ok := g.StateAnnotations["join_state"]
	if !ok {
		t.Fatal("StateAnnotations missing key join_state")
	}
	if joinAnn != ir.StateJoin {
		t.Errorf("StateAnnotations[join_state] = %v, want StateJoin", joinAnn)
	}
}

func TestParseStateConcurrent(t *testing.T) {
	input := `stateDiagram-v2
    state First {
        [*] --> fir
        --
        [*] --> sec
    }`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph
	cs, ok := g.CompositeStates["First"]
	if !ok {
		t.Fatal("CompositeStates missing key First")
	}
	if len(cs.Regions) != 2 {
		t.Errorf("Regions = %d, want 2", len(cs.Regions))
	}
	// Each region should have 1 edge
	for i, region := range cs.Regions {
		if len(region.Edges) != 1 {
			t.Errorf("Region[%d] edges = %d, want 1", i, len(region.Edges))
		}
	}
}

func TestParseStateTransitionLabel(t *testing.T) {
	input := `stateDiagram-v2
    s1 --> s2 : A transition`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph
	if len(g.Edges) != 1 {
		t.Fatalf("Edges = %d, want 1", len(g.Edges))
	}
	e := g.Edges[0]
	if e.From != "s1" {
		t.Errorf("Edge.From = %q, want s1", e.From)
	}
	if e.To != "s2" {
		t.Errorf("Edge.To = %q, want s2", e.To)
	}
	if e.Label == nil {
		t.Fatal("Edge.Label is nil, want 'A transition'")
	}
	if *e.Label != "A transition" {
		t.Errorf("Edge.Label = %q, want %q", *e.Label, "A transition")
	}
}

func TestParseStateDirection(t *testing.T) {
	input := `stateDiagram-v2
    direction LR`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph
	if g.Direction != ir.LeftRight {
		t.Errorf("Direction = %v, want LeftRight", g.Direction)
	}
}

func TestParseStateNote(t *testing.T) {
	input := `stateDiagram-v2
    State1
    note right of State1 : Important info`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph
	if len(g.Notes) != 1 {
		t.Fatalf("Notes = %d, want 1", len(g.Notes))
	}
	n := g.Notes[0]
	if n.Position != "right of" {
		t.Errorf("Note.Position = %q, want %q", n.Position, "right of")
	}
	if n.Target != "State1" {
		t.Errorf("Note.Target = %q, want %q", n.Target, "State1")
	}
	if n.Text != "Important info" {
		t.Errorf("Note.Text = %q, want %q", n.Text, "Important info")
	}
}

func TestParseStateBracketAnnotation(t *testing.T) {
	input := `stateDiagram-v2
    state fork_state [[fork]]`

	out, err := parseState(input)
	if err != nil {
		t.Fatalf("parseState() error: %v", err)
	}

	g := out.Graph
	ann, ok := g.StateAnnotations["fork_state"]
	if !ok {
		t.Fatal("StateAnnotations missing key fork_state")
	}
	if ann != ir.StateFork {
		t.Errorf("StateAnnotations[fork_state] = %v, want StateFork", ann)
	}
}
