package parser

import "fmt"

// ParseError represents a parsing error with context about what went wrong.
type ParseError struct {
	// Diagram is the diagram type being parsed (e.g., "sequence", "state").
	Diagram string

	// Line is the input line that caused the error, if applicable.
	Line string

	// Message describes the error.
	Message string
}

func (e *ParseError) Error() string {
	if e.Line != "" {
		return fmt.Sprintf("%s parser: %s (line: %q)", e.Diagram, e.Message, e.Line)
	}
	return fmt.Sprintf("%s parser: %s", e.Diagram, e.Message)
}
