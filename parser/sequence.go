package parser

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

// seqActorKeyword is the keyword used to declare actor participants.
const seqActorKeyword = "actor"

// arrowEntry maps a literal arrow token to its SeqMessageKind.
// Ordered longest-first so scanning finds the right match.
//
//nolint:gochecknoglobals // package-level lookup table is idiomatic for constant sets.
var arrowTable = []struct {
	token string
	kind  ir.SeqMessageKind
}{
	{"<<-->>", ir.MsgBiDotted},
	{"<<->>", ir.MsgBiSolid},
	{"-->>", ir.MsgDottedArrow},
	{"->>", ir.MsgSolidArrow},
	{"-->", ir.MsgDotted},
	{"--x", ir.MsgDottedCross},
	{"--)", ir.MsgDottedOpen},
	{"->", ir.MsgSolid},
	{"-x", ir.MsgSolidCross},
	{"-)", ir.MsgSolidOpen},
}

var (
	// participantRe matches "participant A as Alice" or "actor B as Bob".
	participantRe = regexp.MustCompile(
		`^(participant|actor)\s+(\S+?)(?:\s+as\s+(.+))?$`,
	)

	// jsonParticipantRe matches JSON-typed participant annotations.
	jsonParticipantRe = regexp.MustCompile(
		`^(participant|actor)\s+(\S+?)@\{(.+?)\}(?:\s+as\s+(.+))?$`,
	)

	// noteRe matches Note right of Alice: text.
	noteRe = regexp.MustCompile(
		`(?i)^Note\s+(right\s+of|left\s+of|over)\s+([^:]+):\s*(.+)$`,
	)

	// activateRe and deactivateRe match activate/deactivate commands.
	activateRe   = regexp.MustCompile(`^activate\s+(\S+)$`)
	deactivateRe = regexp.MustCompile(`^deactivate\s+(\S+)$`)

	// linkRe matches link Alice: Label @ URL.
	linkRe = regexp.MustCompile(`^link\s+(\S+)\s*:\s*(.+?)\s*@\s*(.+)$`)

	// linksRe matches links Alice: {"Label": "URL", ...}.
	linksRe = regexp.MustCompile(`^links\s+(\S+)\s*:\s*(.+)$`)

	// propertiesRe matches properties Alice: {"key": "value", ...}.
	propertiesRe = regexp.MustCompile(`^properties\s+(\S+)\s*:\s*(.+)$`)

	// createRe matches create participant Carl OR create actor Carl as C.
	createRe = regexp.MustCompile(`^create\s+(participant|actor)\s+(\S+?)(?:\s+as\s+(.+))?$`)

	// destroyRe matches destroy Carl.
	destroyRe = regexp.MustCompile(`^destroy\s+(\S+)$`)

	// rectRe matches rect rgb(...) or rect rgba(...).
	rectRe = regexp.MustCompile(`(?i)^rect\s+(.+)$`)

	// boxRe matches box Purple Team A.
	boxRe = regexp.MustCompile(`(?i)^box\s+(.*)$`)
)

// parseSequence parses a Mermaid sequence diagram into a Graph.
//
//nolint:gocognit,funlen,maintidx // sequence parsing has inherent complexity from 15+ distinct line types.
func parseSequence(input string) (*ParseOutput, error) {
	graph := ir.NewGraph()
	graph.Kind = ir.Sequence

	lines := preprocessInput(input)

	// Track participant IDs for ordering / implicit creation.
	participantIndex := map[string]int{}

	// Box tracking: when inBox is true, new participants are added to the box.
	var currentBox *ir.SeqBox
	inBox := false

	// Track frame/box nesting depth for structural validation.
	frameDepth := 0

	ensureParticipant := func(id string) {
		if _, exists := participantIndex[id]; exists {
			return
		}
		participantIndex[id] = len(graph.Participants)
		graph.Participants = append(graph.Participants, &ir.SeqParticipant{
			ID:   id,
			Kind: ir.ParticipantBox,
		})
		if inBox && currentBox != nil {
			currentBox.Participants = append(currentBox.Participants, id)
		}
	}

	findParticipant := func(id string) *ir.SeqParticipant {
		if idx, ok := participantIndex[id]; ok {
			return graph.Participants[idx]
		}
		return nil
	}

	for _, line := range lines {
		lower := strings.ToLower(line)

		// Skip header line.
		if strings.HasPrefix(lower, "sequencediagram") {
			continue
		}

		// autonumber
		if lower == "autonumber" {
			graph.Autonumber = true
			continue
		}

		// end — closes frame or box
		if lower == "end" {
			switch {
			case inBox:
				graph.Boxes = append(graph.Boxes, currentBox)
				currentBox = nil
				inBox = false
			case frameDepth > 0:
				frameDepth--
				graph.Events = append(graph.Events, &ir.SeqEvent{Kind: ir.EvFrameEnd})
			default:
				return nil, &ParseError{
					Diagram: "sequence",
					Line:    line,
					Message: "unexpected \"end\" without matching frame or box",
				}
			}
			continue
		}

		// JSON-typed participant: participant API@{ "type": "boundary" } as Public API
		if match := jsonParticipantRe.FindStringSubmatch(line); match != nil {
			baseKind := strings.ToLower(match[1])
			id := match[2]
			jsonBody := strings.TrimSpace(match[3])
			alias := strings.TrimSpace(match[4])

			ensureParticipant(id)
			participant := findParticipant(id)
			if alias != "" {
				participant.Alias = alias
			}

			// Determine kind from base keyword.
			if baseKind == seqActorKeyword {
				participant.Kind = ir.ActorStickFigure
			}

			// Parse JSON to override kind if "type" is present.
			jsonStr := "{" + jsonBody + "}"
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
				return nil, &ParseError{
					Diagram: "sequence",
					Line:    line,
					Message: "invalid JSON in participant annotation: " + err.Error(),
				}
			}
			if typeStr, ok := parsed["type"].(string); ok {
				participant.Kind = seqKindFromString(typeStr)
			}
			continue
		}

		// participant / actor
		if match := participantRe.FindStringSubmatch(line); match != nil {
			kind := strings.ToLower(match[1])
			id := match[2]
			alias := strings.TrimSpace(match[3])

			ensureParticipant(id)
			participant := findParticipant(id)
			if alias != "" {
				participant.Alias = alias
			}
			if kind == seqActorKeyword {
				participant.Kind = ir.ActorStickFigure
			}
			continue
		}

		// create participant / create actor
		if match := createRe.FindStringSubmatch(line); match != nil {
			kind := strings.ToLower(match[1])
			id := match[2]
			alias := strings.TrimSpace(match[3])

			ensureParticipant(id)
			participant := findParticipant(id)
			participant.IsCreated = true
			if alias != "" {
				participant.Alias = alias
			}
			if kind == seqActorKeyword {
				participant.Kind = ir.ActorStickFigure
			}
			graph.Events = append(graph.Events, &ir.SeqEvent{
				Kind:   ir.EvCreate,
				Target: id,
			})
			continue
		}

		// destroy
		if match := destroyRe.FindStringSubmatch(line); match != nil {
			id := match[1]
			ensureParticipant(id)
			participant := findParticipant(id)
			participant.IsDestroyed = true
			graph.Events = append(graph.Events, &ir.SeqEvent{
				Kind:   ir.EvDestroy,
				Target: id,
			})
			continue
		}

		// activate / deactivate
		if match := activateRe.FindStringSubmatch(line); match != nil {
			target := match[1]
			ensureParticipant(target)
			graph.Events = append(graph.Events, &ir.SeqEvent{
				Kind:   ir.EvActivate,
				Target: target,
			})
			continue
		}
		if match := deactivateRe.FindStringSubmatch(line); match != nil {
			target := match[1]
			ensureParticipant(target)
			graph.Events = append(graph.Events, &ir.SeqEvent{
				Kind:   ir.EvDeactivate,
				Target: target,
			})
			continue
		}

		// Notes
		if match := noteRe.FindStringSubmatch(line); match != nil {
			posToken := strings.ToLower(strings.TrimSpace(match[1]))
			participantsPart := strings.TrimSpace(match[2])
			text := replaceBR(strings.TrimSpace(match[3]))

			var pos ir.SeqNotePosition
			var participants []string

			switch {
			case strings.HasPrefix(posToken, "right"):
				pos = ir.NoteRight
				participants = []string{strings.TrimSpace(participantsPart)}
			case strings.HasPrefix(posToken, "left"):
				pos = ir.NoteLeft
				participants = []string{strings.TrimSpace(participantsPart)}
			case posToken == "over":
				pos = ir.NoteOver
				for _, part := range strings.Split(participantsPart, ",") {
					trimmed := strings.TrimSpace(part)
					if trimmed != "" {
						participants = append(participants, trimmed)
					}
				}
			}

			for _, pid := range participants {
				ensureParticipant(pid)
			}

			graph.Events = append(graph.Events, &ir.SeqEvent{
				Kind: ir.EvNote,
				Note: &ir.SeqNote{
					Position:     pos,
					Participants: participants,
					Text:         text,
				},
			})
			continue
		}

		// Frames: loop, alt, else, opt, par, and, critical, option, break
		if handled, isStart := parseFrameLine(lower, line, graph); handled {
			if isStart {
				frameDepth++
			}
			continue
		}

		// rect rgb(...)
		if match := rectRe.FindStringSubmatch(line); match != nil {
			color := strings.TrimSpace(match[1])
			frameDepth++
			graph.Events = append(graph.Events, &ir.SeqEvent{
				Kind: ir.EvFrameStart,
				Frame: &ir.SeqFrame{
					Kind:  ir.FrameRect,
					Color: color,
				},
			})
			continue
		}

		// box
		if match := boxRe.FindStringSubmatch(line); match != nil {
			label, color := parseBoxLabel(strings.TrimSpace(match[1]))
			currentBox = &ir.SeqBox{
				Label: label,
				Color: color,
			}
			inBox = true
			continue
		}

		// link
		if match := linkRe.FindStringSubmatch(line); match != nil {
			id := match[1]
			label := strings.TrimSpace(match[2])
			url := strings.TrimSpace(match[3])
			ensureParticipant(id)
			participant := findParticipant(id)
			participant.Links = append(participant.Links, ir.SeqLink{Label: label, URL: url})
			continue
		}

		// links (JSON)
		if match := linksRe.FindStringSubmatch(line); match != nil {
			id := match[1]
			body := strings.TrimSpace(match[2])
			ensureParticipant(id)
			participant := findParticipant(id)
			var parsed map[string]string
			if err := json.Unmarshal([]byte(body), &parsed); err != nil {
				return nil, &ParseError{
					Diagram: "sequence",
					Line:    line,
					Message: "invalid JSON in links: " + err.Error(),
				}
			}
			for lbl, url := range parsed {
				participant.Links = append(participant.Links, ir.SeqLink{Label: lbl, URL: url})
			}
			continue
		}

		// properties (JSON)
		if match := propertiesRe.FindStringSubmatch(line); match != nil {
			id := match[1]
			body := strings.TrimSpace(match[2])
			ensureParticipant(id)
			participant := findParticipant(id)
			var parsed map[string]string
			if err := json.Unmarshal([]byte(body), &parsed); err != nil {
				return nil, &ParseError{
					Diagram: "sequence",
					Line:    line,
					Message: "invalid JSON in properties: " + err.Error(),
				}
			}
			if participant.Properties == nil {
				participant.Properties = make(map[string]string)
			}
			for k, v := range parsed {
				participant.Properties[k] = v
			}
			continue
		}

		// Message: From->>To: text
		if from, to, text, kind, activateTarget, deactivateSource, ok := parseSeqMessage(line); ok {
			ensureParticipant(from)
			ensureParticipant(to)

			msg := &ir.SeqMessage{
				From:             from,
				To:               to,
				Text:             replaceBR(text),
				Kind:             kind,
				ActivateTarget:   activateTarget,
				DeactivateSource: deactivateSource,
			}

			graph.Events = append(graph.Events, &ir.SeqEvent{
				Kind:    ir.EvMessage,
				Message: msg,
			})

			if activateTarget {
				graph.Events = append(graph.Events, &ir.SeqEvent{
					Kind:   ir.EvActivate,
					Target: to,
				})
			}
			if deactivateSource {
				graph.Events = append(graph.Events, &ir.SeqEvent{
					Kind:   ir.EvDeactivate,
					Target: from,
				})
			}
			continue
		}
	}

	// Validate nesting.
	if frameDepth > 0 {
		return nil, &ParseError{
			Diagram: "sequence",
			Message: "unclosed frame (missing \"end\")",
		}
	}
	if inBox {
		return nil, &ParseError{
			Diagram: "sequence",
			Message: "unclosed box (missing \"end\")",
		}
	}

	return &ParseOutput{Graph: graph}, nil
}

// parseSeqMessage scans a line for a sequence message arrow using
// longest-match-first scanning. Returns the parsed components.
//
//nolint:nonamedreturns,revive // named returns clarify multi-value return; 7 return values needed for backward compatibility.
func parseSeqMessage(line string) (from, to, text string, kind ir.SeqMessageKind, activateTarget, deactivateSource bool, ok bool) {
	for _, entry := range arrowTable {
		idx := strings.Index(line, entry.token)
		if idx < 0 {
			continue
		}

		from = strings.TrimSpace(line[:idx])
		if from == "" {
			continue
		}

		rest := line[idx+len(entry.token):]

		// Check for activation shorthand: + or - immediately after the arrow.
		activateTarget = false
		deactivateSource = false
		if len(rest) > 0 && rest[0] == '+' {
			activateTarget = true
			rest = rest[1:]
		} else if len(rest) > 0 && rest[0] == '-' {
			deactivateSource = true
			rest = rest[1:]
		}

		// Split on colon for message text.
		colonIdx := strings.Index(rest, ":")
		if colonIdx < 0 {
			// No colon — to is rest, no text.
			to = strings.TrimSpace(rest)
			text = ""
		} else {
			to = strings.TrimSpace(rest[:colonIdx])
			text = strings.TrimSpace(rest[colonIdx+1:])
		}

		if to == "" {
			continue
		}

		return from, to, text, entry.kind, activateTarget, deactivateSource, true
	}
	return "", "", "", 0, false, false, false
}

// parseFrameLine checks if a line is a frame keyword (loop, alt, etc.) and
// appends the appropriate event. Returns (handled, isFrameStart).
func parseFrameLine(lower, original string, graph *ir.Graph) (bool, bool) {
	type frameMatch struct {
		prefix string
		kind   ir.SeqFrameKind
		event  ir.SeqEventKind
	}

	frames := []frameMatch{
		{"loop ", ir.FrameLoop, ir.EvFrameStart},
		{"alt ", ir.FrameAlt, ir.EvFrameStart},
		{"opt ", ir.FrameOpt, ir.EvFrameStart},
		{"par ", ir.FramePar, ir.EvFrameStart},
		{"critical ", ir.FrameCritical, ir.EvFrameStart},
		{"break ", ir.FrameBreak, ir.EvFrameStart},
	}

	for _, frame := range frames {
		if strings.HasPrefix(lower, frame.prefix) {
			label := strings.TrimSpace(original[len(frame.prefix):])
			graph.Events = append(graph.Events, &ir.SeqEvent{
				Kind: frame.event,
				Frame: &ir.SeqFrame{
					Kind:  frame.kind,
					Label: label,
				},
			})
			return true, true // handled, isStart
		}
	}

	// Middle frame keywords.
	type middleMatch struct {
		prefix string
	}
	middles := []middleMatch{
		{"else "},
		{"and "},
		{"option "},
	}

	for _, mid := range middles {
		if strings.HasPrefix(lower, mid.prefix) {
			label := strings.TrimSpace(original[len(mid.prefix):])
			graph.Events = append(graph.Events, &ir.SeqEvent{
				Kind: ir.EvFrameMiddle,
				Frame: &ir.SeqFrame{
					Label: label,
				},
			})
			return true, false // handled, not a start
		}
	}

	// Handle bare keywords without labels.
	bareFrames := map[string]ir.SeqFrameKind{
		"loop":     ir.FrameLoop,
		"alt":      ir.FrameAlt,
		"opt":      ir.FrameOpt,
		"par":      ir.FramePar,
		"critical": ir.FrameCritical,
		"break":    ir.FrameBreak,
	}
	if kind, found := bareFrames[lower]; found {
		graph.Events = append(graph.Events, &ir.SeqEvent{
			Kind: ir.EvFrameStart,
			Frame: &ir.SeqFrame{
				Kind: kind,
			},
		})
		return true, true // handled, isStart
	}

	bareMiddles := []string{"else", "and", "option"}
	for _, bm := range bareMiddles {
		if lower == bm {
			graph.Events = append(graph.Events, &ir.SeqEvent{
				Kind:  ir.EvFrameMiddle,
				Frame: &ir.SeqFrame{},
			})
			return true, false // handled, not a start
		}
	}

	return false, false
}

// parseBoxLabel extracts color and label from a box declaration.
// Examples:
//
//	"Purple Team A"         -> label="Team A", color="Purple"
//	"rgb(33,66,99) Group"   -> label="Group", color="rgb(33,66,99)"
//	"transparent Name"      -> label="Name", color="transparent"
//	"Team A"                -> label="Team A", color=""
//
//nolint:nonamedreturns // named returns clarify the multi-value return.
func parseBoxLabel(input string) (label, color string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", ""
	}

	// Check for rgb/rgba/hsl(...) color.
	if strings.HasPrefix(strings.ToLower(input), "rgb") ||
		strings.HasPrefix(strings.ToLower(input), "hsl") {
		// Find the closing paren.
		parenIdx := strings.Index(input, ")")
		if parenIdx >= 0 {
			color = strings.TrimSpace(input[:parenIdx+1])
			label = strings.TrimSpace(input[parenIdx+1:])
			return label, color
		}
	}

	// Check for known CSS color keywords or "transparent".
	knownColors := []string{
		"transparent", "red", "blue", "green", "yellow", "orange", "purple",
		"pink", "cyan", "magenta", "white", "black", "gray", "grey",
		"lightblue", "lightgreen", "lightyellow", "darkblue", "darkgreen",
		"darkred", "aqua", "lime", "navy", "teal", "olive", "maroon",
		"silver", "fuchsia",
	}

	lowerInput := strings.ToLower(input)
	for _, colorName := range knownColors {
		if strings.HasPrefix(lowerInput, colorName) {
			// Check that the color word is followed by a space or end of string.
			rest := input[len(colorName):]
			if rest == "" {
				// Color only, no label.
				return "", input
			}
			if rest[0] == ' ' {
				return strings.TrimSpace(rest), input[:len(colorName)]
			}
		}
	}

	// Check for hex color (#xxx or #xxxxxx).
	if strings.HasPrefix(input, "#") {
		parts := strings.SplitN(input, " ", 2)
		if len(parts) == 2 {
			return strings.TrimSpace(parts[1]), parts[0]
		}
		return "", input
	}

	// No recognized color; the whole thing is the label.
	return input, ""
}

// seqKindFromString converts a type string to SeqParticipantKind.
func seqKindFromString(typeStr string) ir.SeqParticipantKind {
	switch strings.ToLower(typeStr) {
	case "participant":
		return ir.ParticipantBox
	case seqActorKeyword:
		return ir.ActorStickFigure
	case "boundary":
		return ir.ParticipantBoundary
	case "control":
		return ir.ParticipantControl
	case "entity":
		return ir.ParticipantEntity
	case "database":
		return ir.ParticipantDatabase
	case "collections":
		return ir.ParticipantCollections
	case "queue":
		return ir.ParticipantQueue
	default:
		return ir.ParticipantBox
	}
}

// replaceBR replaces <br/> and <br> tags with newlines.
func replaceBR(text string) string {
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br>", "\n")
	return text
}
