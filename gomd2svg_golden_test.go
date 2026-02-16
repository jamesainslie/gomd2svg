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
				buf.WriteString("first difference at line ")
				buf.WriteString(strings.TrimLeft(strings.Repeat("0", 3)+string(rune('0'+i%10)), "0"))
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
