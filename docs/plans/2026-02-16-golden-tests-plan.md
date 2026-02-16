# Golden Tests Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add golden snapshot tests that compare rendered SVG output byte-for-byte against committed baselines for all 50 fixtures across 5 themes (250 golden files).

**Architecture:** A single `mermaid_golden_test.go` at the root package iterates `.mmd` fixtures, renders each with all themes via `RenderWithOptions()`, and compares against `testdata/golden/<fixture>-<theme>.svg`. A `-update` flag regenerates golden files. Five parallel agents generate the initial baseline with anomaly detection.

**Tech Stack:** Go stdlib (`testing`, `flag`, `os`, `path/filepath`, `strings`), existing `mermaid.RenderWithOptions()` and `theme.Names()` APIs.

---

### Task 1: Write the golden test harness

**Files:**
- Create: `mermaid_golden_test.go`

**Step 1: Write the test file**

Create `mermaid_golden_test.go` with:

```go
package mermaid

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/theme"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestGolden(t *testing.T) {
	fixtures, err := filepath.Glob("testdata/fixtures/*.mmd")
	if err != nil {
		t.Fatal(err)
	}
	if len(fixtures) == 0 {
		t.Fatal("no fixtures found in testdata/fixtures/")
	}

	themes := theme.Names()

	for _, fixture := range fixtures {
		base := strings.TrimSuffix(filepath.Base(fixture), ".mmd")
		input, err := os.ReadFile(fixture)
		if err != nil {
			t.Fatal(err)
		}

		for _, themeName := range themes {
			themeName := themeName
			name := base + "-" + themeName
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				svg, err := RenderWithOptions(string(input), Options{ThemeName: themeName})
				if err != nil {
					t.Fatalf("RenderWithOptions(%s, %s): %v", base, themeName, err)
				}

				goldenPath := filepath.Join("testdata", "golden", name+".svg")

				if *updateGolden {
					if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(goldenPath, []byte(svg), 0o644); err != nil {
						t.Fatal(err)
					}
					return
				}

				expected, err := os.ReadFile(goldenPath)
				if err != nil {
					t.Fatalf("golden file missing (run with -update to create): %v", err)
				}

				if svg != string(expected) {
					// Show a truncated diff to aid debugging.
					diff := diffSnippet(string(expected), svg, 20)
					t.Errorf("golden mismatch for %s\n%s\nRun: go test -run TestGolden -update", name, diff)
				}
			})
		}
	}
}

// diffSnippet returns the first maxLines differing lines between two strings.
func diffSnippet(want, got string, maxLines int) string {
	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(got, "\n")

	var buf strings.Builder
	shown := 0
	max := len(wantLines)
	if len(gotLines) > max {
		max = len(gotLines)
	}

	for i := 0; i < max && shown < maxLines; i++ {
		var w, g string
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if w != g {
			if shown == 0 {
				buf.WriteString("first difference at line " + strings.Repeat(" ", 0))
				buf.WriteString(string(rune('0'+i/100)) + string(rune('0'+(i/10)%10)) + string(rune('0'+i%10)))
				buf.WriteByte('\n')
			}
			buf.WriteString("  want: ")
			if len(w) > 120 {
				buf.WriteString(w[:120])
				buf.WriteString("...")
			} else {
				buf.WriteString(w)
			}
			buf.WriteByte('\n')
			buf.WriteString("  got:  ")
			if len(g) > 120 {
				buf.WriteString(g[:120])
				buf.WriteString("...")
			} else {
				buf.WriteString(g)
			}
			buf.WriteByte('\n')
			shown++
		}
	}

	if shown == 0 {
		return "(files differ in length only)"
	}
	return buf.String()
}
```

**Step 2: Verify it compiles**

Run: `go build ./...`
Expected: no errors.

**Step 3: Run the test without golden files to confirm it fails**

Run: `go test -run TestGolden -count=1`
Expected: FAIL — "golden file missing (run with -update to create)"

**Step 4: Commit the harness**

```bash
git add mermaid_golden_test.go
git commit -m "test: add golden test harness with -update flag"
```

---

### Task 2: Create golden directory

**Files:**
- Create: `testdata/golden/.gitkeep`

**Step 1: Create the directory**

```bash
mkdir -p testdata/golden
touch testdata/golden/.gitkeep
```

**Step 2: Commit**

```bash
git add testdata/golden/.gitkeep
git commit -m "test: add testdata/golden directory for golden SVG files"
```

---

### Task 3: Generate golden files — Agent A (Flowchart, Class, State, ER)

**Fixtures:** `flowchart-simple`, `flowchart-shapes`, `flowchart-labels`, `class-simple`, `class-relationships`, `state-simple`, `state-composite`, `er-simple`, `er-attributes` (9 fixtures, 45 golden files)

**Step 1: Generate golden files**

Run: `go test -run "TestGolden/(flowchart|class|state|er)" -update -count=1`
Expected: PASS, 45 `.svg` files created in `testdata/golden/`.

**Step 2: Verify file count**

Run: `ls testdata/golden/{flowchart,class,state,er}*.svg | wc -l`
Expected: 45

**Step 3: Run anomaly detection on each generated file**

For every generated `.svg` file, verify:

1. **Valid SVG structure**: file starts with `<svg` and ends with `</svg>`
2. **Has `viewBox`**: contains the string `viewBox=`
3. **Non-trivial size**: file is > 200 bytes
4. **Accessibility**: contains `role="img"` and `aria-label=`
5. **Font-family**: contains `font-family=`
6. **Diagram-specific elements**:
   - Flowchart files: contain `<rect` (node shapes)
   - Class files: contain `<rect` and `<line` (class boxes and dividers)
   - State files: contain `<rect` or `<circle` (states and start/end markers)
   - ER files: contain `<rect` (entity boxes)
7. **Theme differentiation**: for each fixture, all 5 `.svg` files are distinct (no two are identical)

If any check fails, stop and report the anomaly. Do NOT proceed.

**Step 4: Run golden test to confirm match**

Run: `go test -run "TestGolden/(flowchart|class|state|er)" -count=1`
Expected: PASS — all 45 subtests pass.

---

### Task 4: Generate golden files — Agent B (Sequence, Kanban, Packet)

**Fixtures:** `sequence-simple`, `sequence-activations`, `sequence-frames`, `sequence-full`, `kanban-basic`, `kanban-metadata`, `packet-tcp`, `packet-bitcount` (8 fixtures, 40 golden files)

**Step 1: Generate golden files**

Run: `go test -run "TestGolden/(sequence|kanban|packet)" -update -count=1`
Expected: PASS, 40 `.svg` files created in `testdata/golden/`.

**Step 2: Verify file count**

Run: `ls testdata/golden/{sequence,kanban,packet}*.svg | wc -l`
Expected: 40

**Step 3: Run anomaly detection**

Same checklist as Task 3, with diagram-specific elements:
- Sequence files: contain `<line` (lifelines/messages) and `<rect` (participant boxes)
- Kanban files: contain `<rect` (columns and cards)
- Packet files: contain `<rect` (bit field cells)

Theme differentiation: for each fixture, all 5 `.svg` files are distinct.

If any check fails, stop and report the anomaly.

**Step 4: Run golden test to confirm match**

Run: `go test -run "TestGolden/(sequence|kanban|packet)" -count=1`
Expected: PASS — all 40 subtests pass.

---

### Task 5: Generate golden files — Agent C (Pie, Quadrant, Timeline, Gantt, GitGraph)

**Fixtures:** `pie-basic`, `pie-showdata`, `quadrant-minimal`, `quadrant-campaigns`, `timeline-basic`, `timeline-sections`, `gantt-basic`, `gantt-dependencies`, `gitgraph-basic`, `gitgraph-branches` (10 fixtures, 50 golden files)

**Step 1: Generate golden files**

Run: `go test -run "TestGolden/(pie|quadrant|timeline|gantt|gitgraph)" -update -count=1`
Expected: PASS, 50 `.svg` files created in `testdata/golden/`.

**Step 2: Verify file count**

Run: `ls testdata/golden/{pie,quadrant,timeline,gantt,gitgraph}*.svg | wc -l`
Expected: 50

**Step 3: Run anomaly detection**

Same base checklist, with diagram-specific elements:
- Pie files: contain `<path` (arc slices)
- Quadrant files: contain `<rect` (quadrant grid) and `<circle` (data points)
- Timeline files: contain `<rect` (event boxes)
- Gantt files: contain `<rect` (task bars)
- GitGraph files: contain `<circle` (commit dots)

Theme differentiation: for each fixture, all 5 `.svg` files are distinct.

If any check fails, stop and report the anomaly.

**Step 4: Run golden test to confirm match**

Run: `go test -run "TestGolden/(pie|quadrant|timeline|gantt|gitgraph)" -count=1`
Expected: PASS — all 50 subtests pass.

---

### Task 6: Generate golden files — Agent D (XYChart, Radar, Mindmap, Sankey, Treemap)

**Fixtures:** `xychart-basic`, `xychart-horizontal`, `radar-basic`, `radar-polygon`, `mindmap-basic`, `mindmap-shapes`, `sankey-basic`, `sankey-energy`, `treemap-basic`, `treemap-nested` (10 fixtures, 50 golden files)

**Step 1: Generate golden files**

Run: `go test -run "TestGolden/(xychart|radar|mindmap|sankey|treemap)" -update -count=1`
Expected: PASS, 50 `.svg` files created in `testdata/golden/`.

**Step 2: Verify file count**

Run: `ls testdata/golden/{xychart,radar,mindmap,sankey,treemap}*.svg | wc -l`
Expected: 50

**Step 3: Run anomaly detection**

Same base checklist, with diagram-specific elements:
- XYChart files: contain `<rect` (bars) or `<path` (line series)
- Radar files: contain `<polygon` (data overlay) or `<path`
- Mindmap files: contain `<path` (curved branches)
- Sankey files: contain `<path` (flow curves) and `<rect` (nodes)
- Treemap files: contain `<rect` (nested rectangles)

Theme differentiation: for each fixture, all 5 `.svg` files are distinct.

If any check fails, stop and report the anomaly.

**Step 4: Run golden test to confirm match**

Run: `go test -run "TestGolden/(xychart|radar|mindmap|sankey|treemap)" -count=1`
Expected: PASS — all 50 subtests pass.

---

### Task 7: Generate golden files — Agent E (Requirement, Block, C4, Journey, Architecture, ZenUML)

**Fixtures:** `requirement-basic`, `requirement-multiple`, `block-grid`, `block-edges`, `c4-context`, `c4-container`, `journey-basic`, `journey-multiactor`, `architecture-basic`, `architecture-groups`, `zenuml-basic`, `zenuml-controlflow`, `zenuml-trycatch` (13 fixtures, 65 golden files)

**Step 1: Generate golden files**

Run: `go test -run "TestGolden/(requirement|block|c4|journey|architecture|zenuml)" -update -count=1`
Expected: PASS, 65 `.svg` files created in `testdata/golden/`.

**Step 2: Verify file count**

Run: `ls testdata/golden/{requirement,block,c4,journey,architecture,zenuml}*.svg | wc -l`
Expected: 65

**Step 3: Run anomaly detection**

Same base checklist, with diagram-specific elements:
- Requirement files: contain `<rect` (requirement/element boxes)
- Block files: contain `<rect` (block containers)
- C4 files: contain `<rect` (C4 boxes)
- Journey files: contain `<rect` (task bars) and `<circle` (score dots)
- Architecture files: contain `<rect` (service boxes)
- ZenUML files: contain `<line` (lifelines) and `<rect` (participant boxes)

Theme differentiation: for each fixture, all 5 `.svg` files are distinct.

If any check fails, stop and report the anomaly.

**Step 4: Run golden test to confirm match**

Run: `go test -run "TestGolden/(requirement|block|c4|journey|architecture|zenuml)" -count=1`
Expected: PASS — all 65 subtests pass.

---

### Task 8: Full verification and commit

**Step 1: Verify total golden file count**

Run: `ls testdata/golden/*.svg | wc -l`
Expected: 250

**Step 2: Run full golden test suite**

Run: `go test -run TestGolden -count=1 -v 2>&1 | tail -5`
Expected: PASS — 250 subtests pass.

**Step 3: Run full project test suite**

Run: `go test ./... -count=1`
Expected: all 9 packages pass.

**Step 4: Remove .gitkeep**

```bash
rm testdata/golden/.gitkeep
```

**Step 5: Commit all golden files**

```bash
git add mermaid_golden_test.go testdata/golden/
git commit -m "test: add 250 golden SVG files for all fixtures and themes"
```

---

## Agent Dispatch Strategy

Tasks 3-7 run as **parallel agents** (one per task). Each agent:
1. Runs `go test -run TestGolden/<pattern> -update` to generate golden files
2. Runs the anomaly detection checklist (shared verification rules)
3. Runs `go test -run TestGolden/<pattern>` to confirm byte-for-byte match
4. Reports success or anomaly

Tasks 1-2 run first (scaffold). Task 8 runs last (verification).

## Anomaly Detection Checklist (shared across all agents)

For every generated `.svg` file:

| Check | How | Fail condition |
|-------|-----|----------------|
| Valid SVG | `<svg` present and `</svg>` present | Missing either tag |
| viewBox | contains `viewBox=` | Missing |
| Non-trivial | file size > 200 bytes | Too small |
| Accessibility | contains `role="img"` and `aria-label=` | Missing either |
| Font-family | contains `font-family=` | Missing |
| Diagram elements | diagram-specific element check (see per-task) | Missing expected elements |
| Theme differentiation | for each fixture, all 5 theme SVGs are distinct | Any two identical |

**If any check fails:** Stop. Report the fixture name, theme, and which check failed. Do NOT write the golden file for that combination.
