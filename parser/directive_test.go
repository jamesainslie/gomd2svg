package parser

import (
	"strings"
	"testing"
)

func TestExtractDirective_ThemeName(t *testing.T) {
	input := "%%{init: {\"theme\": \"dark\"}}%%\nflowchart LR\n  A-->B"
	dir, rest := extractDirective(input)
	if dir.Theme != "dark" {
		t.Errorf("Theme = %q, want dark", dir.Theme)
	}
	if strings.Contains(rest, "%%{init") {
		t.Error("directive not stripped from rest")
	}
}

func TestExtractDirective_None(t *testing.T) {
	input := "flowchart LR\n  A-->B"
	dir, rest := extractDirective(input)
	if dir.Theme != "" {
		t.Errorf("Theme = %q, want empty", dir.Theme)
	}
	if rest != input {
		t.Error("input should be unchanged")
	}
}

func TestExtractDirective_FontOverride(t *testing.T) {
	input := "%%{init: {\"theme\": \"forest\", \"themeVariables\": {\"fontFamily\": \"Fira Code\"}}}%%\nflowchart LR\n  A-->B"
	dir, _ := extractDirective(input)
	if dir.Theme != "forest" {
		t.Errorf("Theme = %q, want forest", dir.Theme)
	}
	if dir.ThemeVariables.FontFamily != "Fira Code" {
		t.Errorf("FontFamily = %q, want Fira Code", dir.ThemeVariables.FontFamily)
	}
}

func TestExtractDirective_SingleQuotes(t *testing.T) {
	input := "%%{init: {'theme': 'dark'}}%%\nflowchart LR\n  A-->B"
	dir, _ := extractDirective(input)
	if dir.Theme != "dark" {
		t.Errorf("Theme = %q, want dark", dir.Theme)
	}
}

func TestExtractDirective_PreservesRest(t *testing.T) {
	input := "%%{init: {\"theme\": \"dark\"}}%%\nflowchart LR\n  A-->B"
	_, rest := extractDirective(input)
	if !strings.HasPrefix(rest, "flowchart") {
		t.Errorf("rest should start with flowchart, got %q", rest[:20])
	}
}
