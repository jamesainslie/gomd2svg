package ir

import "testing"

func TestArchSideString(t *testing.T) {
	tests := []struct {
		side ArchSide
		want string
	}{
		{ArchLeft, "L"},
		{ArchRight, "R"},
		{ArchTop, "T"},
		{ArchBottom, "B"},
	}
	for _, tt := range tests {
		if got := tt.side.String(); got != tt.want {
			t.Errorf("ArchSide(%d).String() = %q, want %q", tt.side, got, tt.want)
		}
	}
}

func TestArchService(t *testing.T) {
	svc := &ArchService{
		ID:      "web",
		Label:   "Web Server",
		Icon:    "server",
		GroupID: "cloud",
	}
	if svc.ID != "web" {
		t.Errorf("ID = %q, want %q", svc.ID, "web")
	}
	if svc.Label != "Web Server" {
		t.Errorf("Label = %q, want %q", svc.Label, "Web Server")
	}
	if svc.Icon != "server" {
		t.Errorf("Icon = %q, want %q", svc.Icon, "server")
	}
	if svc.GroupID != "cloud" {
		t.Errorf("GroupID = %q, want %q", svc.GroupID, "cloud")
	}
}

func TestArchEdge(t *testing.T) {
	edge := &ArchEdge{
		FromID:     "web",
		FromSide:   ArchRight,
		ToID:       "db",
		ToSide:     ArchLeft,
		ArrowLeft:  false,
		ArrowRight: true,
	}
	if edge.FromID != "web" {
		t.Errorf("FromID = %q, want %q", edge.FromID, "web")
	}
	if edge.FromSide != ArchRight {
		t.Errorf("FromSide = %v, want ArchRight", edge.FromSide)
	}
	if edge.ToID != "db" {
		t.Errorf("ToID = %q, want %q", edge.ToID, "db")
	}
	if edge.ToSide != ArchLeft {
		t.Errorf("ToSide = %v, want ArchLeft", edge.ToSide)
	}
	if edge.ArrowLeft {
		t.Error("ArrowLeft = true, want false")
	}
	if !edge.ArrowRight {
		t.Error("ArrowRight = false, want true")
	}
}
