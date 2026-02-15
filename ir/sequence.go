package ir

// SeqParticipantKind distinguishes participant rendering styles.
type SeqParticipantKind int

const (
	ParticipantBox SeqParticipantKind = iota
	ActorStickFigure
	ParticipantBoundary
	ParticipantControl
	ParticipantEntity
	ParticipantDatabase
	ParticipantCollections
	ParticipantQueue
)

// String returns the lowercase name of the participant kind.
func (k SeqParticipantKind) String() string {
	switch k {
	case ParticipantBox:
		return "participant"
	case ActorStickFigure:
		return "actor"
	case ParticipantBoundary:
		return "boundary"
	case ParticipantControl:
		return "control"
	case ParticipantEntity:
		return "entity"
	case ParticipantDatabase:
		return "database"
	case ParticipantCollections:
		return "collections"
	case ParticipantQueue:
		return "queue"
	default:
		return ""
	}
}

// SeqParticipant represents a participant in a sequence diagram.
type SeqParticipant struct {
	ID          string
	Alias       string
	Kind        SeqParticipantKind
	Links       []SeqLink
	Properties  map[string]string
	IsCreated   bool
	IsDestroyed bool
}

// DisplayName returns the alias if non-empty, otherwise the ID.
func (p *SeqParticipant) DisplayName() string {
	if p.Alias != "" {
		return p.Alias
	}
	return p.ID
}

// SeqLink represents a clickable link on a participant.
type SeqLink struct {
	Label string
	URL   string
}

// SeqMessageKind distinguishes message arrow styles.
type SeqMessageKind int

const (
	MsgSolid       SeqMessageKind = iota // ->
	MsgDotted                            // -->
	MsgSolidArrow                        // ->>
	MsgDottedArrow                       // -->>
	MsgSolidCross                        // -x
	MsgDottedCross                       // --x
	MsgSolidOpen                         // -)
	MsgDottedOpen                        // --)
	MsgBiSolid                           // <<->>
	MsgBiDotted                          // <<-->>
)

// IsDotted returns true if the message uses a dotted line style.
func (k SeqMessageKind) IsDotted() bool {
	switch k {
	case MsgDotted, MsgDottedArrow, MsgDottedCross, MsgDottedOpen, MsgBiDotted:
		return true
	default:
		return false
	}
}

// SeqMessage represents a message between participants.
type SeqMessage struct {
	From             string
	To               string
	Text             string
	Kind             SeqMessageKind
	ActivateTarget   bool
	DeactivateSource bool
}

// SeqEventKind distinguishes the types of events in a sequence diagram.
type SeqEventKind int

const (
	EvMessage SeqEventKind = iota
	EvNote
	EvActivate
	EvDeactivate
	EvFrameStart
	EvFrameMiddle
	EvFrameEnd
	EvCreate
	EvDestroy
)

// SeqEvent represents a single event in the sequence diagram timeline.
type SeqEvent struct {
	Kind    SeqEventKind
	Message *SeqMessage
	Note    *SeqNote
	Frame   *SeqFrame
	Target  string
}

// SeqNotePosition indicates where a note is placed relative to participants.
type SeqNotePosition int

const (
	NoteLeft SeqNotePosition = iota
	NoteRight
	NoteOver
)

// SeqNote represents a note in a sequence diagram.
type SeqNote struct {
	Position     SeqNotePosition
	Participants []string
	Text         string
}

// SeqFrameKind distinguishes the types of combined fragments.
type SeqFrameKind int

const (
	FrameLoop SeqFrameKind = iota
	FrameAlt
	FrameOpt
	FramePar
	FrameCritical
	FrameBreak
	FrameRect
)

// String returns the lowercase name of the frame kind.
func (k SeqFrameKind) String() string {
	switch k {
	case FrameLoop:
		return "loop"
	case FrameAlt:
		return "alt"
	case FrameOpt:
		return "opt"
	case FramePar:
		return "par"
	case FrameCritical:
		return "critical"
	case FrameBreak:
		return "break"
	case FrameRect:
		return "rect"
	default:
		return ""
	}
}

// SeqFrame represents a combined fragment (loop, alt, opt, etc.).
type SeqFrame struct {
	Kind  SeqFrameKind
	Label string
	Color string
}

// SeqBox represents a colored box grouping participants.
type SeqBox struct {
	Label        string
	Color        string
	Participants []string
}
