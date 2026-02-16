package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestPacketLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Packet
	graph.Fields = []*ir.PacketField{
		{Start: 0, End: 15, Description: "Source Port"},
		{Start: 16, End: 31, Description: "Dest Port"},
		{Start: 32, End: 63, Description: "Sequence Number"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	pd, ok := lay.Diagram.(PacketData)
	if !ok {
		t.Fatalf("Diagram type = %T, want PacketData", lay.Diagram)
	}

	// 32 bits per row: first row has 2 fields (0-15, 16-31), second row has 1 (32-63)
	if len(pd.Rows) != 2 {
		t.Fatalf("len(Rows) = %d, want 2", len(pd.Rows))
	}
	if len(pd.Rows[0].Fields) != 2 {
		t.Errorf("len(Rows[0].Fields) = %d, want 2", len(pd.Rows[0].Fields))
	}
	if len(pd.Rows[1].Fields) != 1 {
		t.Errorf("len(Rows[1].Fields) = %d, want 1", len(pd.Rows[1].Fields))
	}

	// Field widths should be proportional to bit count
	f0 := pd.Rows[0].Fields[0]
	f1 := pd.Rows[0].Fields[1]
	if f0.Width != f1.Width {
		t.Errorf("16-bit fields should have equal width: %v vs %v", f0.Width, f1.Width)
	}

	// 32-bit field should span full row width
	f2 := pd.Rows[1].Fields[0]
	expectedFullWidth := float32(cfg.Packet.BitsPerRow) * cfg.Packet.BitWidth
	if f2.Width != expectedFullWidth {
		t.Errorf("32-bit field width = %v, want %v", f2.Width, expectedFullWidth)
	}

	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("dimensions = %v x %v, want positive", lay.Width, lay.Height)
	}
}

func TestPacketLayoutSingleBit(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Packet
	graph.Fields = []*ir.PacketField{
		{Start: 0, End: 0, Description: "Flag"},
		{Start: 1, End: 31, Description: "Rest"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	pd, ok := lay.Diagram.(PacketData)
	if !ok {
		t.Fatal("expected PacketData")
	}
	if len(pd.Rows) != 1 {
		t.Fatalf("len(Rows) = %d, want 1", len(pd.Rows))
	}
	if pd.Rows[0].Fields[0].Width != cfg.Packet.BitWidth {
		t.Errorf("1-bit field width = %v, want %v", pd.Rows[0].Fields[0].Width, cfg.Packet.BitWidth)
	}
}
