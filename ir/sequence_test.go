package ir

import "testing"

func TestSeqParticipantKindString(t *testing.T) {
	tests := []struct {
		kind SeqParticipantKind
		want string
	}{
		{ParticipantBox, "participant"},
		{ActorStickFigure, "actor"},
		{ParticipantBoundary, "boundary"},
		{ParticipantControl, "control"},
		{ParticipantEntity, "entity"},
		{ParticipantDatabase, "database"},
		{ParticipantCollections, "collections"},
		{ParticipantQueue, "queue"},
		{SeqParticipantKind(99), ""},
	}
	for _, tt := range tests {
		got := tt.kind.String()
		if got != tt.want {
			t.Errorf("SeqParticipantKind(%d).String() = %q, want %q", int(tt.kind), got, tt.want)
		}
	}
}

func TestSeqParticipantDisplayName(t *testing.T) {
	tests := []struct {
		name string
		p    SeqParticipant
		want string
	}{
		{
			name: "with alias",
			p:    SeqParticipant{ID: "A", Alias: "Alice"},
			want: "Alice",
		},
		{
			name: "without alias",
			p:    SeqParticipant{ID: "Bob"},
			want: "Bob",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.DisplayName()
			if got != tt.want {
				t.Errorf("DisplayName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSeqMessageKindIsDotted(t *testing.T) {
	tests := []struct {
		kind SeqMessageKind
		want bool
	}{
		{MsgSolid, false},
		{MsgDotted, true},
		{MsgSolidArrow, false},
		{MsgDottedArrow, true},
		{MsgSolidCross, false},
		{MsgDottedCross, true},
		{MsgSolidOpen, false},
		{MsgDottedOpen, true},
		{MsgBiSolid, false},
		{MsgBiDotted, true},
	}
	for _, tt := range tests {
		got := tt.kind.IsDotted()
		if got != tt.want {
			t.Errorf("SeqMessageKind(%d).IsDotted() = %v, want %v", int(tt.kind), got, tt.want)
		}
	}
}

func TestSeqFrameKindString(t *testing.T) {
	tests := []struct {
		kind SeqFrameKind
		want string
	}{
		{FrameLoop, "loop"},
		{FrameAlt, "alt"},
		{FrameOpt, "opt"},
		{FramePar, "par"},
		{FrameCritical, "critical"},
		{FrameBreak, "break"},
		{FrameRect, "rect"},
		{SeqFrameKind(99), ""},
	}
	for _, tt := range tests {
		got := tt.kind.String()
		if got != tt.want {
			t.Errorf("SeqFrameKind(%d).String() = %q, want %q", int(tt.kind), got, tt.want)
		}
	}
}

func TestGraphSequenceFields(t *testing.T) {
	g := NewGraph()
	if g.Participants != nil {
		t.Error("Participants should be nil on new graph")
	}
	if g.Events != nil {
		t.Error("Events should be nil on new graph")
	}
	if g.Boxes != nil {
		t.Error("Boxes should be nil on new graph")
	}
	if g.Autonumber {
		t.Error("Autonumber should be false on new graph")
	}
}
