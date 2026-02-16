package gomd2svg

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
	t.Parallel()
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
	maxIdx := len(wantLines)
	if len(gotLines) > maxIdx {
		maxIdx = len(gotLines)
	}

	for lineIdx := 0; lineIdx < maxIdx && shown < maxLines; lineIdx++ {
		var wantLine, gotLine string
		if lineIdx < len(wantLines) {
			wantLine = wantLines[lineIdx]
		}
		if lineIdx < len(gotLines) {
			gotLine = gotLines[lineIdx]
		}
		if wantLine != gotLine {
			if shown == 0 {
				buf.WriteString("first difference at line ")
				buf.WriteString(strings.TrimLeft(strings.Repeat("0", 3)+string(rune('0'+lineIdx%10)), "0"))
				buf.WriteByte('\n')
			}
			buf.WriteString("  want: ")
			if len(wantLine) > 120 {
				buf.WriteString(wantLine[:120])
				buf.WriteString("...")
			} else {
				buf.WriteString(wantLine)
			}
			buf.WriteByte('\n')
			buf.WriteString("  got:  ")
			if len(gotLine) > 120 {
				buf.WriteString(gotLine[:120])
				buf.WriteString("...")
			} else {
				buf.WriteString(gotLine)
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
