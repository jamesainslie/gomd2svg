package parser

import (
	"errors"
	"strings"
	"testing"
)

func TestParseErrorFormat(t *testing.T) {
	e := &ParseError{Diagram: "sequence", Message: "bad syntax", Line: "foo bar"}
	got := e.Error()
	if !strings.Contains(got, "sequence parser") {
		t.Errorf("error = %q, want contains 'sequence parser'", got)
	}
	if !strings.Contains(got, "bad syntax") {
		t.Errorf("error = %q, want contains 'bad syntax'", got)
	}
	if !strings.Contains(got, "foo bar") {
		t.Errorf("error = %q, want contains line", got)
	}
}

func TestParseErrorFormatNoLine(t *testing.T) {
	e := &ParseError{Diagram: "state", Message: "unclosed block"}
	got := e.Error()
	if !strings.Contains(got, "state parser") {
		t.Errorf("error = %q, want contains 'state parser'", got)
	}
	if strings.Contains(got, "line:") {
		t.Errorf("error = %q, should not contain 'line:' when Line is empty", got)
	}
}

func TestSequenceInvalidLinksJSON(t *testing.T) {
	input := `sequenceDiagram
participant Alice
links Alice: {invalid json
`
	_, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for invalid JSON in links, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Diagram != "sequence" {
		t.Errorf("diagram = %q, want sequence", pe.Diagram)
	}
	if !strings.Contains(pe.Message, "JSON") {
		t.Errorf("message = %q, want contains 'JSON'", pe.Message)
	}
}

func TestSequenceInvalidPropertiesJSON(t *testing.T) {
	input := `sequenceDiagram
participant Alice
properties Alice: not-json-at-all
`
	_, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for invalid JSON in properties, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
}

func TestSequenceInvalidParticipantJSON(t *testing.T) {
	input := `sequenceDiagram
participant API@{not valid json} as Public
`
	_, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for invalid JSON in participant annotation, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
}

func TestSequenceUnmatchedEnd(t *testing.T) {
	input := `sequenceDiagram
Alice->>Bob: Hello
end
end
`
	_, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for unmatched end, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if !strings.Contains(pe.Message, "end") {
		t.Errorf("message = %q, want contains 'end'", pe.Message)
	}
}

func TestSequenceUnclosedFrame(t *testing.T) {
	input := `sequenceDiagram
Alice->>Bob: Hello
loop Every minute
  Bob->>Alice: Great
`
	_, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for unclosed frame, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if !strings.Contains(pe.Message, "unclosed") {
		t.Errorf("message = %q, want contains 'unclosed'", pe.Message)
	}
}

func TestSequenceUnclosedBox(t *testing.T) {
	input := `sequenceDiagram
box Purple Group
participant Alice
participant Bob
`
	_, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for unclosed box, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
}

func TestSequenceValidInputNoError(t *testing.T) {
	input := `sequenceDiagram
participant Alice
participant Bob
Alice->>Bob: Hello
loop Every minute
  Bob->>Alice: Great
end
`
	_, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStateUnclosedComposite(t *testing.T) {
	input := `stateDiagram-v2
state A {
  B --> C
`
	_, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for unclosed composite state, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Diagram != "state" {
		t.Errorf("diagram = %q, want state", pe.Diagram)
	}
}

func TestFlowchartUnclosedSubgraph(t *testing.T) {
	input := `flowchart TD
subgraph one
  A --> B
`
	_, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for unclosed subgraph, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Diagram != "flowchart" {
		t.Errorf("diagram = %q, want flowchart", pe.Diagram)
	}
}

func TestZenUMLUnclosedBlock(t *testing.T) {
	input := `zenuml
@Starter(Client)
Server.process() {
  Database.query()
`
	_, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for unclosed block, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Diagram != "zenuml" {
		t.Errorf("diagram = %q, want zenuml", pe.Diagram)
	}
}
