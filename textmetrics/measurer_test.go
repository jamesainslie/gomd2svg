package textmetrics

import "testing"

func TestMeasureEmpty(t *testing.T) {
	m := New()
	w := m.Width("", 14, "sans-serif")
	if w != 0 {
		t.Errorf("Width of empty = %f, want 0", w)
	}
}

func TestMeasureNonEmpty(t *testing.T) {
	m := New()
	w := m.Width("Hello", 14, "sans-serif")
	if w <= 0 {
		t.Errorf("Width of 'Hello' = %f, want > 0", w)
	}
}

func TestMeasureLongerIsWider(t *testing.T) {
	m := New()
	short := m.Width("Hi", 14, "sans-serif")
	long := m.Width("Hello World", 14, "sans-serif")
	if long <= short {
		t.Errorf("long (%f) should be > short (%f)", long, short)
	}
}

func TestMeasureLargerFontIsWider(t *testing.T) {
	m := New()
	small := m.Width("Hello", 10, "sans-serif")
	big := m.Width("Hello", 20, "sans-serif")
	if big <= small {
		t.Errorf("big (%f) should be > small (%f)", big, small)
	}
}

func TestAverageCharWidth(t *testing.T) {
	m := New()
	w := m.AverageCharWidth("sans-serif", 14)
	if w <= 0 {
		t.Errorf("AverageCharWidth = %f, want > 0", w)
	}
	if w > 14 {
		t.Errorf("AverageCharWidth = %f, unexpectedly large", w)
	}
}
