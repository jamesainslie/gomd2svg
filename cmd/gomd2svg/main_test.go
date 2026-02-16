package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderFile(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"render", "../../testdata/fixtures/flowchart-simple.mmd"}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "<svg") {
		t.Error("expected SVG output")
	}
}

func TestRenderStdin(t *testing.T) {
	stdin := strings.NewReader("flowchart LR\n  A-->B")
	var stdout, stderr bytes.Buffer
	err := run([]string{"render"}, stdin, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "<svg") {
		t.Error("expected SVG output")
	}
}

func TestRenderOutputFile(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.svg")
	var stdout, stderr bytes.Buffer
	err := run([]string{"render", "-o", out, "../../testdata/fixtures/flowchart-simple.mmd"}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "<svg") {
		t.Error("expected SVG in output file")
	}
}

func TestRenderThemeFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"render", "-theme", "dark", "../../testdata/fixtures/flowchart-simple.mmd"}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "#1A1A2E") {
		t.Error("expected dark background")
	}
}

func TestRenderTiming(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"render", "-timing", "../../testdata/fixtures/flowchart-simple.mmd"}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr.String(), "total:") {
		t.Error("expected timing on stderr")
	}
}

func TestThemes(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"themes"}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"modern", "default", "dark", "forest", "neutral"} {
		if !strings.Contains(stdout.String(), name) {
			t.Errorf("missing theme %q", name)
		}
	}
}

func TestVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"version"}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "gomd2svg") {
		t.Error("expected version output")
	}
}

func TestEmptyInput(t *testing.T) {
	stdin := strings.NewReader("")
	var stdout, stderr bytes.Buffer
	err := run([]string{"render"}, stdin, &stdout, &stderr)
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"bogus"}, nil, &stdout, &stderr)
	if err == nil {
		t.Error("expected error for unknown command")
	}
}

func TestHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"help"}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr.String(), "render") {
		t.Error("expected help text")
	}
}

func TestNoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run(nil, nil, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr.String(), "render") {
		t.Error("expected usage text")
	}
}
