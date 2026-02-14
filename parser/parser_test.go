package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestDetectDiagramKind(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ir.DiagramKind
	}{
		{"flowchart LR", "flowchart LR\n  A-->B", ir.Flowchart},
		{"graph TD", "graph TD\n  A-->B", ir.Flowchart},
		{"sequenceDiagram", "sequenceDiagram\n  Alice->>Bob: Hi", ir.Sequence},
		{"classDiagram", "classDiagram\n  A <|-- B", ir.Class},
		{"stateDiagram", "stateDiagram-v2\n  [*] --> A", ir.State},
		{"pie", "pie\n  \"A\" : 10", ir.Pie},
		{"skip comments", "%%{init}%%\nflowchart LR\n  A-->B", ir.Flowchart},
		{"skip empty lines", "\n\n  flowchart TD\n  A-->B", ir.Flowchart},
		{"default flowchart", "A-->B", ir.Flowchart},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectDiagramKind(tt.input)
			if got != tt.want {
				t.Errorf("detectDiagramKind() = %v, want %v", got, tt.want)
			}
		})
	}
}
