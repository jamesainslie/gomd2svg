package parser

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Directive holds parsed %%{init: ...}%% values.
type Directive struct {
	Theme          string         `json:"theme"`
	ThemeVariables ThemeVariables `json:"themeVariables"`
}

// ThemeVariables holds theme field overrides from directives.
type ThemeVariables struct {
	FontFamily   string `json:"fontFamily"`
	Background   string `json:"background"`
	PrimaryColor string `json:"primaryColor"`
	LineColor    string `json:"lineColor"`
	TextColor    string `json:"textColor"`
}

var directiveRe = regexp.MustCompile(`(?m)^\s*%%\{init:\s*(.*?)\}%%\s*$`)

// extractDirective finds and removes a %%{init: ...}%% directive from input.
// Returns the parsed directive and the input with the directive line removed.
func extractDirective(input string) (Directive, string) {
	var dir Directive
	loc := directiveRe.FindStringSubmatchIndex(input)
	if loc == nil {
		return dir, input
	}

	jsonStr := input[loc[2]:loc[3]]
	// Normalize single quotes to double quotes for JSON compatibility.
	jsonStr = strings.ReplaceAll(jsonStr, "'", "\"")
	// Wrap in braces if needed.
	if !strings.HasPrefix(strings.TrimSpace(jsonStr), "{") {
		jsonStr = "{" + jsonStr + "}"
	}

	_ = json.Unmarshal([]byte(jsonStr), &dir)

	// Remove the directive line from input.
	rest := input[:loc[0]] + input[loc[1]:]
	rest = strings.TrimLeft(rest, "\n")
	return dir, rest
}
