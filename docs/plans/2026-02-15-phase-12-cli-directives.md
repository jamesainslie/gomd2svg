# Phase 12: CLI Tool & Directive Support

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make gomd2svg usable as a standalone command-line tool and support inline `%%{init: ...}%%` directives for per-diagram theme/config overrides.

**Architecture:** CLI binary at `cmd/gomd2svg/` using stdlib `flag` package (no external deps). Directive parsing in `parser/` extracts `%%{init:}%%` before diagram detection. CLI flags override directives which override defaults.

**Tech Stack:** Go stdlib (`flag`, `encoding/json`, `os`, `io`, `path/filepath`)

---

## Precedence

```
CLI flags > %%{init:}%% directives > library defaults
```

---

### Task 1: Directive parsing in parser

**Files:**
- Create: `parser/directive.go`
- Create: `parser/directive_test.go`
- Modify: `parser/parser.go` (wire directive extraction into Parse)

**Step 1: Write the failing tests**

```go
// parser/directive_test.go
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
    input := `%%{init: {"theme": "forest", "themeVariables": {"fontFamily": "Fira Code"}}}%%
flowchart LR
  A-->B`
    dir, _ := extractDirective(input)
    if dir.Theme != "forest" {
        t.Errorf("Theme = %q, want forest", dir.Theme)
    }
    if dir.ThemeVariables.FontFamily != "Fira Code" {
        t.Errorf("FontFamily = %q, want Fira Code", dir.ThemeVariables.FontFamily)
    }
}

func TestExtractDirective_SingleQuotes(t *testing.T) {
    // mermaid.js uses single quotes; we normalize to double for JSON
    input := "%%{init: {'theme': 'dark'}}%%\nflowchart LR\n  A-->B"
    dir, _ := extractDirective(input)
    if dir.Theme != "dark" {
        t.Errorf("Theme = %q, want dark", dir.Theme)
    }
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./parser/ -run TestExtractDirective -v`
Expected: FAIL — `extractDirective` undefined

**Step 3: Implement directive extraction**

```go
// parser/directive.go
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
    FontFamily string `json:"fontFamily"`
    Background string `json:"background"`
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
    // Wrap in braces if the regex captured without them.
    if !strings.HasPrefix(strings.TrimSpace(jsonStr), "{") {
        jsonStr = "{" + jsonStr + "}"
    }

    _ = json.Unmarshal([]byte(jsonStr), &dir)

    // Remove the directive line from input.
    rest := input[:loc[0]] + input[loc[1]:]
    rest = strings.TrimLeft(rest, "\n")
    return dir, rest
}
```

**Step 4: Wire into Parse()**

In `parser/parser.go`, extract directive before `detectDiagramKind`:
```go
func Parse(input string) (*ParseOutput, error) {
    dir, cleaned := extractDirective(input)
    kind := detectDiagramKind(cleaned)
    // ... existing parsing with cleaned input ...
    return &ParseOutput{
        Graph:     g,
        Directive: dir,
    }, nil
}
```

Add `Directive` field to `ParseOutput`.

**Step 5: Run tests**

Run: `go test ./parser/ -v`
Expected: All pass including new directive tests

**Step 6: Commit**

```bash
git add parser/directive.go parser/directive_test.go parser/parser.go
git commit -m "feat(parser): add %%{init:}%% directive extraction"
```

---

### Task 2: Apply directives in the render pipeline

**Files:**
- Modify: `mermaid.go` (apply directive theme/overrides)
- Modify: `options.go` (add directive application helper)
- Create: `mermaid_directive_test.go`

**Step 1: Write the failing tests**

```go
// mermaid_directive_test.go
func TestDirectiveThemeOverride(t *testing.T) {
    input := "%%{init: {\"theme\": \"dark\"}}%%\nflowchart LR\n  A-->B"
    svg, err := Render(input)
    if err != nil {
        t.Fatal(err)
    }
    // Dark theme has #1A1A2E background
    if !strings.Contains(svg, "#1A1A2E") {
        t.Error("expected dark background color in SVG")
    }
}

func TestDirectiveOverriddenByCLI(t *testing.T) {
    input := "%%{init: {\"theme\": \"dark\"}}%%\nflowchart LR\n  A-->B"
    svg, err := RenderWithOptions(input, Options{ThemeName: "forest"})
    if err != nil {
        t.Fatal(err)
    }
    // CLI flag (forest) should win over directive (dark)
    if strings.Contains(svg, "#1A1A2E") {
        t.Error("CLI theme should override directive theme")
    }
}
```

**Step 2: Run tests to verify they fail**

Run: `go test -run TestDirective -v`
Expected: FAIL — dark background not present (directive ignored)

**Step 3: Apply directives in pipeline**

In `mermaid.go`, after parsing:
```go
func RenderWithOptions(input string, opts Options) (string, error) {
    po, err := parser.Parse(input)
    if err != nil { return "", err }

    th := opts.resolveTheme(po.Directive)
    // ... layout and render with th ...
}
```

In `options.go`:
```go
func (o Options) resolveTheme(dir parser.Directive) *theme.Theme {
    // CLI ThemeName takes highest precedence.
    if o.ThemeName != "" {
        return theme.ByName(o.ThemeName)
    }
    // Directive theme is second.
    if dir.Theme != "" {
        th := theme.ByName(dir.Theme)
        // Apply themeVariables as overrides.
        if dir.ThemeVariables != (parser.ThemeVariables{}) {
            ov := theme.Overrides{}
            if dir.ThemeVariables.FontFamily != "" {
                ov.FontFamily = &dir.ThemeVariables.FontFamily
            }
            if dir.ThemeVariables.Background != "" {
                ov.Background = &dir.ThemeVariables.Background
            }
            if dir.ThemeVariables.PrimaryColor != "" {
                ov.PrimaryColor = &dir.ThemeVariables.PrimaryColor
            }
            if dir.ThemeVariables.LineColor != "" {
                ov.LineColor = &dir.ThemeVariables.LineColor
            }
            if dir.ThemeVariables.TextColor != "" {
                ov.TextColor = &dir.ThemeVariables.TextColor
            }
            th = theme.WithOverrides(th, ov)
        }
        return th
    }
    // Explicit Theme object.
    if o.Theme != nil {
        return o.Theme
    }
    return theme.Modern()
}
```

**Step 4: Run tests**

Run: `go test -v -run TestDirective`
Expected: PASS

**Step 5: Commit**

```bash
git add mermaid.go options.go mermaid_directive_test.go
git commit -m "feat: apply %%{init:}%% directives in render pipeline"
```

---

### Task 3: CLI binary — core render command

**Files:**
- Create: `cmd/gomd2svg/main.go`

**Step 1: Implement CLI**

```go
// cmd/gomd2svg/main.go
package main

import (
    "flag"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"

    mermaid "github.com/jamesainslie/gomd2svg"
    "github.com/jamesainslie/gomd2svg/theme"
)

var version = "dev"

func main() {
    if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
    if len(args) == 0 {
        return usage(stderr)
    }

    switch args[0] {
    case "render":
        return runRender(args[1:], stdin, stdout, stderr)
    case "themes":
        return runThemes(stdout)
    case "version":
        fmt.Fprintf(stdout, "gomd2svg %s\n", version)
        return nil
    case "help", "-h", "--help":
        return usage(stderr)
    default:
        return fmt.Errorf("unknown command: %s", args[0])
    }
}

func runRender(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
    fs := flag.NewFlagSet("render", flag.ContinueOnError)
    fs.SetOutput(stderr)
    output := fs.String("o", "", "output file (default: stdout)")
    themeName := fs.String("theme", "", "theme name (modern|default|dark|forest|neutral)")
    timing := fs.Bool("timing", false, "print timing info to stderr")
    if err := fs.Parse(args); err != nil {
        return err
    }

    // Read input: file arg or stdin.
    var input []byte
    var err error
    if fs.NArg() > 0 {
        input, err = os.ReadFile(fs.Arg(0))
        if err != nil {
            return err
        }
    } else {
        input, err = io.ReadAll(stdin)
        if err != nil {
            return err
        }
    }

    if len(strings.TrimSpace(string(input))) == 0 {
        return fmt.Errorf("empty input")
    }

    opts := mermaid.Options{}
    if *themeName != "" {
        opts.ThemeName = *themeName
    }

    if *timing {
        result, err := mermaid.RenderWithTiming(string(input), opts)
        if err != nil {
            return err
        }
        fmt.Fprintf(stderr, "parse: %dus  layout: %dus  render: %dus  total: %.1fms\n",
            result.ParseUs, result.LayoutUs, result.RenderUs, result.TotalMs())
        return writeOutput(*output, result.SVG, stdout)
    }

    svg, err := mermaid.RenderWithOptions(string(input), opts)
    if err != nil {
        return err
    }
    return writeOutput(*output, svg, stdout)
}

func writeOutput(path, svg string, stdout io.Writer) error {
    if path == "" {
        _, err := io.WriteString(stdout, svg)
        return err
    }
    dir := filepath.Dir(path)
    if dir != "." && dir != "" {
        os.MkdirAll(dir, 0o755)
    }
    return os.WriteFile(path, []byte(svg), 0o644)
}

func runThemes(w io.Writer) error {
    for _, name := range theme.Names() {
        fmt.Fprintln(w, name)
    }
    return nil
}

func usage(w io.Writer) error {
    fmt.Fprintln(w, `Usage: gomd2svg <command> [options]

Commands:
  render [file]   Render a .mmd file to SVG
  themes          List available themes
  version         Print version

Render options:
  -o <file>       Output file (default: stdout)
  -theme <name>   Theme: modern, default, dark, forest, neutral
  -timing         Print timing info to stderr

Examples:
  gomd2svg render diagram.mmd -o diagram.svg
  gomd2svg render -theme dark diagram.mmd > out.svg
  cat diagram.mmd | gomd2svg render > out.svg
  gomd2svg render -theme forest -timing diagram.mmd -o out.svg`)
    return nil
}
```

**Step 2: Build and test manually**

```bash
go build -o /tmp/gomd2svg ./cmd/gomd2svg/
/tmp/gomd2svg version
/tmp/gomd2svg themes
/tmp/gomd2svg render testdata/fixtures/flowchart-simple.mmd -o /tmp/test.svg
cat testdata/fixtures/pie-basic.mmd | /tmp/gomd2svg render --theme dark > /tmp/pie-dark.svg
```

**Step 3: Commit**

```bash
git add cmd/gomd2svg/
git commit -m "feat: add gomd2svg CLI tool"
```

---

### Task 4: CLI tests

**Files:**
- Create: `cmd/gomd2svg/main_test.go`

**Step 1: Write tests**

```go
// cmd/gomd2svg/main_test.go
func TestRenderFile(t *testing.T) {
    // render a fixture file to stdout
    var stdout, stderr bytes.Buffer
    err := run([]string{"render", "../../testdata/fixtures/flowchart-simple.mmd"}, nil, &stdout, &stderr)
    if err != nil { t.Fatal(err) }
    if !strings.Contains(stdout.String(), "<svg") {
        t.Error("expected SVG output")
    }
}

func TestRenderStdin(t *testing.T) {
    stdin := strings.NewReader("flowchart LR\n  A-->B")
    var stdout, stderr bytes.Buffer
    err := run([]string{"render"}, stdin, &stdout, &stderr)
    if err != nil { t.Fatal(err) }
    if !strings.Contains(stdout.String(), "<svg") {
        t.Error("expected SVG output")
    }
}

func TestRenderOutputFile(t *testing.T) {
    out := filepath.Join(t.TempDir(), "out.svg")
    var stdout, stderr bytes.Buffer
    err := run([]string{"render", "-o", out, "../../testdata/fixtures/flowchart-simple.mmd"}, nil, &stdout, &stderr)
    if err != nil { t.Fatal(err) }
    data, _ := os.ReadFile(out)
    if !strings.Contains(string(data), "<svg") {
        t.Error("expected SVG in output file")
    }
}

func TestRenderThemeFlag(t *testing.T) {
    var stdout, stderr bytes.Buffer
    err := run([]string{"render", "-theme", "dark", "../../testdata/fixtures/flowchart-simple.mmd"}, nil, &stdout, &stderr)
    if err != nil { t.Fatal(err) }
    if !strings.Contains(stdout.String(), "#1A1A2E") {
        t.Error("expected dark background")
    }
}

func TestRenderTiming(t *testing.T) {
    var stdout, stderr bytes.Buffer
    err := run([]string{"render", "-timing", "../../testdata/fixtures/flowchart-simple.mmd"}, nil, &stdout, &stderr)
    if err != nil { t.Fatal(err) }
    if !strings.Contains(stderr.String(), "total:") {
        t.Error("expected timing on stderr")
    }
}

func TestThemes(t *testing.T) {
    var stdout, stderr bytes.Buffer
    err := run([]string{"themes"}, nil, &stdout, &stderr)
    if err != nil { t.Fatal(err) }
    for _, name := range []string{"modern", "default", "dark", "forest", "neutral"} {
        if !strings.Contains(stdout.String(), name) {
            t.Errorf("missing theme %q", name)
        }
    }
}

func TestVersion(t *testing.T) {
    var stdout, stderr bytes.Buffer
    err := run([]string{"version"}, nil, &stdout, &stderr)
    if err != nil { t.Fatal(err) }
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
```

**Step 2: Run tests**

Run: `go test ./cmd/gomd2svg/ -v`
Expected: All pass

**Step 3: Commit**

```bash
git add cmd/gomd2svg/main_test.go
git commit -m "test: add CLI test suite"
```

---

### Task 5: SVG accessibility (title + aria-label)

**Files:**
- Modify: `render/svg.go` (add title element and aria attributes)
- Modify: `render/svg.go` test or create `render/accessibility_test.go`

**Step 1: Write the failing test**

```go
// render/accessibility_test.go
func TestSVGHasAriaLabel(t *testing.T) {
    g := ir.NewGraph()
    g.Kind = ir.Flowchart
    g.Title = "My Flowchart"
    g.EnsureNode("A").Label.Lines = []string{"Start"}
    th := theme.Modern()
    cfg := config.DefaultLayout()
    l := layout.ComputeLayout(g, th, cfg)
    svg := RenderSVG(l, th, cfg)
    if !strings.Contains(svg, `role="img"`) {
        t.Error("missing role=img")
    }
    if !strings.Contains(svg, `aria-label="My Flowchart"`) {
        t.Error("missing aria-label")
    }
    if !strings.Contains(svg, "<title>My Flowchart</title>") {
        t.Error("missing <title> element")
    }
}

func TestSVGNoTitleNoAria(t *testing.T) {
    g := ir.NewGraph()
    g.Kind = ir.Flowchart
    g.EnsureNode("A").Label.Lines = []string{"Start"}
    th := theme.Modern()
    cfg := config.DefaultLayout()
    l := layout.ComputeLayout(g, th, cfg)
    svg := RenderSVG(l, th, cfg)
    // Without title, should still have role=img but use diagram kind as label
    if !strings.Contains(svg, `role="img"`) {
        t.Error("missing role=img")
    }
    if !strings.Contains(svg, `aria-label="Flowchart diagram"`) {
        t.Error("missing fallback aria-label")
    }
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./render/ -run TestSVGHas -v`
Expected: FAIL

**Step 3: Implement accessibility**

In `render/svg.go`, pass title from layout:
```go
// Add Title field to Layout
// In svg.go:
ariaLabel := l.Title
if ariaLabel == "" {
    ariaLabel = l.Kind.String() + " diagram"
}

b.openTag("svg",
    "xmlns", "http://www.w3.org/2000/svg",
    "width", fmtFloat(width),
    "height", fmtFloat(height),
    "viewBox", "0 0 "+fmtFloat(width)+" "+fmtFloat(height),
    "font-family", th.FontFamily,
    "role", "img",
    "aria-label", ariaLabel,
)

if l.Title != "" {
    b.openTag("title")
    b.content(l.Title)
    b.closeTag("title")
}
```

Add `Title` field to `layout.Layout`, populated from `g.Title` in each layout function.

**Step 4: Run tests**

Run: `go test ./render/ -v`
Expected: All pass

**Step 5: Commit**

```bash
git add render/svg.go render/accessibility_test.go layout/layout.go layout/types.go ir/graph.go
git commit -m "feat(render): add SVG accessibility (role, aria-label, title)"
```

---

### Task 6: Wire title through parsers and layout

**Files:**
- Modify: `ir/graph.go` (ensure Title field exists — check if already there)
- Modify: `layout/layout.go` (copy g.Title to Layout.Title)
- Modify: `layout/types.go` (add Title to Layout)

**Step 1: Add Title to Layout struct**

Check if `ir.Graph` already has a Title field (it likely does for pie/gantt). Add `Title string` to `layout.Layout`. In `ComputeLayout()`, copy `g.Title` to `l.Title` before returning.

**Step 2: Run all tests**

Run: `go test -race ./...`
Expected: All pass

**Step 3: Commit**

```bash
git add layout/layout.go layout/types.go ir/graph.go
git commit -m "feat(layout): propagate diagram title to Layout for accessibility"
```

---

### Task 7: Final validation and cleanup

**Files:**
- Remove: `cmd/testall/` (replaced by CLI tool)
- Modify: tests as needed

**Step 1: Run full test suite**

```bash
go test -race ./...
go vet ./...
gofmt -l .
```

**Step 2: Build CLI and test all fixtures**

```bash
go build -o /tmp/gomd2svg ./cmd/gomd2svg/
for f in testdata/fixtures/*.mmd; do
    /tmp/gomd2svg render "$f" > /dev/null && echo "PASS $f" || echo "FAIL $f"
done
```

**Step 3: Clean up testall**

```bash
rm -rf cmd/testall/
```

**Step 4: Commit**

```bash
git add -A
git commit -m "chore: remove cmd/testall (replaced by CLI tool)"
```
