// Package textmetrics provides font measurement for text layout.
// It measures text width using system fonts when available,
// falling back to a character-width estimation heuristic.
package textmetrics

import (
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// Measurer measures text dimensions for layout purposes.
type Measurer struct {
	mu         sync.Mutex
	fonts      map[string]*sfnt.Font // fontFamily -> loaded font
	widthCache map[widthKey]float32
}

type widthKey struct {
	text       string
	fontSize   float32
	fontFamily string
}

// New creates a new Measurer.
func New() *Measurer {
	return &Measurer{
		fonts:      make(map[string]*sfnt.Font),
		widthCache: make(map[widthKey]float32),
	}
}

// Width returns the width of text rendered at the given font size and family.
// It returns 0 for empty text.
func (m *Measurer) Width(text string, fontSize float32, fontFamily string) float32 {
	if text == "" {
		return 0
	}

	key := widthKey{text: text, fontSize: fontSize, fontFamily: fontFamily}

	m.mu.Lock()
	if cached, ok := m.widthCache[key]; ok {
		m.mu.Unlock()
		return cached
	}
	m.mu.Unlock()

	w := m.measure(text, fontSize, fontFamily)

	m.mu.Lock()
	m.widthCache[key] = w
	m.mu.Unlock()

	return w
}

// AverageCharWidth returns the average character width for a font family
// and size, measured across the Latin alphabet (upper and lower case).
func (m *Measurer) AverageCharWidth(fontFamily string, fontSize float32) float32 {
	const sample = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	total := m.Width(sample, fontSize, fontFamily)
	return total / float32(len(sample))
}

// measure computes text width, trying system fonts first, then falling back
// to a heuristic estimation.
func (m *Measurer) measure(text string, fontSize float32, fontFamily string) float32 {
	f := m.loadFont(fontFamily)
	if f != nil {
		if w, ok := m.measureWithFont(f, text, fontSize); ok {
			return w
		}
	}
	// Fallback: estimate at 0.6 em per character.
	return fontSize * 0.6 * float32(len([]rune(text)))
}

// measureWithFont uses an sfnt.Font to compute the advance width of text.
func (m *Measurer) measureWithFont(f *sfnt.Font, text string, fontSize float32) (float32, bool) {
	var buf sfnt.Buffer
	ppem := fixed.I(int(fontSize))

	var totalAdvance fixed.Int26_6
	for _, r := range text {
		idx, err := f.GlyphIndex(&buf, r)
		if err != nil || idx == 0 {
			// Glyph not found; fall back to heuristic for entire string.
			return 0, false
		}
		adv, err := f.GlyphAdvance(&buf, idx, ppem, font.HintingNone)
		if err != nil {
			return 0, false
		}
		totalAdvance += adv
	}

	return float32(totalAdvance) / 64.0, true // fixed.Int26_6 has 6 fractional bits
}

// loadFont attempts to find and load a system font matching the family name.
// Results are cached. Returns nil if no suitable font is found.
func (m *Measurer) loadFont(fontFamily string) *sfnt.Font {
	m.mu.Lock()
	defer m.mu.Unlock()

	if f, ok := m.fonts[fontFamily]; ok {
		return f // may be nil if previously failed
	}

	f := m.findSystemFont(fontFamily)
	m.fonts[fontFamily] = f
	return f
}

// findSystemFont searches common macOS font directories for a sans-serif font.
func (m *Measurer) findSystemFont(fontFamily string) *sfnt.Font {
	// Map generic family names to actual font file names.
	var candidates []string
	switch fontFamily {
	case "sans-serif", "sans", "arial", "helvetica":
		candidates = []string{
			"Helvetica.ttc",
			"HelveticaNeue.ttc",
			"ArialHB.ttc",
			"Arial.ttf",
			"Arial Unicode.ttf",
		}
	case "serif", "times", "times new roman":
		candidates = []string{
			"Times.ttc",
			"Times New Roman.ttf",
		}
	case "monospace", "courier", "courier new":
		candidates = []string{
			"Courier.ttc",
			"Courier New.ttf",
		}
	default:
		// Try the family name directly as a file name.
		candidates = []string{
			fontFamily + ".ttc",
			fontFamily + ".ttf",
			fontFamily + ".otf",
		}
	}

	searchDirs := []string{
		"/System/Library/Fonts",
		"/Library/Fonts",
	}

	// Also check user fonts if HOME is set.
	if home, err := os.UserHomeDir(); err == nil {
		searchDirs = append(searchDirs, filepath.Join(home, "Library", "Fonts"))
	}

	for _, dir := range searchDirs {
		for _, name := range candidates {
			path := filepath.Join(dir, name)
			if f := m.tryLoadFontFile(path); f != nil {
				return f
			}
		}
	}

	return nil
}

// tryLoadFontFile attempts to parse a font from a file path.
// For .ttc (TrueType Collection) files, the first font in the collection is used.
func (m *Measurer) tryLoadFontFile(path string) *sfnt.Font {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	// Try as a single font first.
	f, err := sfnt.Parse(data)
	if err == nil {
		return f
	}

	// Try as a font collection (.ttc).
	col, err := sfnt.ParseCollection(data)
	if err != nil {
		return nil
	}

	f, err = col.Font(0)
	if err != nil {
		return nil
	}
	return f
}
