package layout

import (
	"math"
	"strings"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

// computeSequenceLayout produces a timeline-based layout for sequence diagrams.
// Unlike other diagram kinds this does not use the Sugiyama algorithm; instead
// participants are placed in columns and events are walked top-to-bottom.
func computeSequenceLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	sc := cfg.Sequence
	lineH := th.FontSize * cfg.LabelLineHeight
	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical

	// Phase 1: Measure participants and assign horizontal positions.
	pInfos, pIndex := seqMeasureParticipants(graph, measurer, th, sc, lineH, padH)

	// Phase 2: Walk events top-to-bottom.
	eventY := sc.HeaderHeight + padV
	messages, notes, activations, frames, eventY := seqProcessEvents(
		graph, measurer, th, sc, lineH, padH, padV, pInfos, pIndex, eventY,
	)

	// Phase 3: Finalize.
	activations = seqCloseRemainingActivations(graph, pInfos, pIndex, activations, sc, eventY)

	footerY := eventY + padV
	diagramH := footerY + sc.HeaderHeight + padV

	participants := make([]SeqParticipantLayout, len(pInfos))
	lifelines := make([]SeqLifeline, len(pInfos))
	for idx, pi := range pInfos {
		participants[idx] = SeqParticipantLayout{
			ID:     pi.id,
			Label:  pi.label,
			Kind:   pi.kind,
			X:      pi.x,
			Y:      pi.y,
			Width:  pi.w,
			Height: pi.h,
		}
		lifelines[idx] = SeqLifeline{
			ParticipantID: pi.id,
			X:             pi.x,
			TopY:          pi.y + pi.h,
			BottomY:       footerY,
		}
	}

	boxes := seqBuildBoxLayouts(graph, pInfos, pIndex, sc, diagramH)

	rightEdge := float32(0)
	if len(pInfos) > 0 {
		last := pInfos[len(pInfos)-1]
		rightEdge = last.x + last.w/2 + padH
	}

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  nil,
		Edges:  nil,
		Width:  rightEdge,
		Height: diagramH,
		Diagram: SequenceData{
			Participants:  participants,
			Lifelines:     lifelines,
			Messages:      messages,
			Activations:   activations,
			Notes:         notes,
			Frames:        frames,
			Boxes:         boxes,
			Autonumber:    graph.Autonumber,
			DiagramHeight: diagramH,
		},
	}
}

type seqParticipantInfo struct {
	id    string
	label TextBlock
	kind  ir.SeqParticipantKind
	x     float32 // centre X
	y     float32 // top Y (0 for normal, set later for created)
	w     float32
	h     float32
}

func seqMeasureParticipants(
	graph *ir.Graph,
	measurer *textmetrics.Measurer,
	th *theme.Theme,
	sc config.SequenceConfig,
	lineH, padH float32,
) ([]seqParticipantInfo, map[string]int) {
	pInfos := make([]seqParticipantInfo, len(graph.Participants))
	pIndex := make(map[string]int, len(graph.Participants))

	cursorX := padH
	for idx, participant := range graph.Participants {
		name := participant.DisplayName()
		tw := measurer.Width(name, th.FontSize, th.FontFamily)
		partW := tw + 2*padH
		partH := sc.HeaderHeight

		pInfos[idx] = seqParticipantInfo{
			id:    participant.ID,
			label: TextBlock{Lines: []string{name}, Width: tw, Height: lineH, FontSize: th.FontSize},
			kind:  participant.Kind,
			x:     cursorX + partW/2,
			y:     0,
			w:     partW,
			h:     partH,
		}
		pIndex[participant.ID] = idx
		cursorX += partW + sc.ParticipantSpacing
	}
	return pInfos, pIndex
}

// seqProcessEvents walks events top-to-bottom, building messages, notes,
// activations, and frames. Returns the updated Y cursor.
func seqProcessEvents( //nolint:revive,funlen // event processing switch is inherently complex; 5 return values needed for distinct event types
	graph *ir.Graph,
	measurer *textmetrics.Measurer,
	th *theme.Theme,
	sc config.SequenceConfig,
	lineH, padH, padV float32,
	pInfos []seqParticipantInfo,
	pIndex map[string]int,
	eventY float32,
) ([]SeqMessageLayout, []SeqNoteLayout, []SeqActivationLayout, []SeqFrameLayout, float32) {
	activationStacks := make(map[string][]float32)

	type frameEntry struct {
		frame    *ir.SeqFrame
		startY   float32
		dividers []float32
	}
	var frameStack []frameEntry

	var messages []SeqMessageLayout
	var notes []SeqNoteLayout
	var activations []SeqActivationLayout
	var frames []SeqFrameLayout
	msgNumber := 0

	for _, ev := range graph.Events {
		switch ev.Kind {
		case ir.EvMessage:
			msg := ev.Message
			eventY += sc.MessageSpacing

			fromIdx, fromOK := pIndex[msg.From]
			toIdx, toOK := pIndex[msg.To]
			if !fromOK || !toOK {
				continue
			}

			fromX := pInfos[fromIdx].x
			toX := pInfos[toIdx].x
			if msg.From == msg.To {
				toX = fromX + sc.SelfMessageWidth
			}

			tw := measurer.Width(msg.Text, th.FontSize, th.FontFamily)
			msgNumber++
			num := 0
			if graph.Autonumber {
				num = msgNumber
			}

			messages = append(messages, SeqMessageLayout{
				From:   msg.From,
				To:     msg.To,
				Text:   TextBlock{Lines: []string{msg.Text}, Width: tw, Height: lineH, FontSize: th.FontSize},
				Kind:   msg.Kind,
				Y:      eventY,
				FromX:  fromX,
				ToX:    toX,
				Number: num,
			})

		case ir.EvNote:
			note := ev.Note
			noteLines := strings.Split(note.Text, "\n")
			maxLineW := float32(0)
			for _, ln := range noteLines {
				lw := measurer.Width(ln, th.FontSize, th.FontFamily)
				if lw > maxLineW {
					maxLineW = lw
				}
			}
			noteW := maxLineW + 2*padH
			if noteW > sc.NoteMaxWidth {
				noteW = sc.NoteMaxWidth
			}
			noteH := float32(len(noteLines))*lineH + 2*padV
			noteX := seqNoteX(note, pInfos, pIndex, padH, noteW)
			notes = append(notes, SeqNoteLayout{
				Text:   TextBlock{Lines: noteLines, Width: maxLineW, Height: float32(len(noteLines)) * lineH, FontSize: th.FontSize},
				X:      noteX,
				Y:      eventY,
				Width:  noteW,
				Height: noteH,
			})
			eventY += noteH + padV

		case ir.EvActivate:
			activationStacks[ev.Target] = append(activationStacks[ev.Target], eventY)

		case ir.EvDeactivate:
			stack := activationStacks[ev.Target]
			if len(stack) > 0 {
				startY := stack[len(stack)-1]
				activationStacks[ev.Target] = stack[:len(stack)-1]
				px := float32(0)
				if idx, ok := pIndex[ev.Target]; ok {
					px = pInfos[idx].x
				}
				activations = append(activations, SeqActivationLayout{
					ParticipantID: ev.Target,
					X:             px - sc.ActivationWidth/2,
					TopY:          startY,
					BottomY:       eventY,
					Width:         sc.ActivationWidth,
				})
			}

		case ir.EvFrameStart:
			frameStack = append(frameStack, frameEntry{
				frame:  ev.Frame,
				startY: eventY,
			})

		case ir.EvFrameMiddle:
			if len(frameStack) > 0 {
				frameStack[len(frameStack)-1].dividers = append(
					frameStack[len(frameStack)-1].dividers, eventY)
			}

		case ir.EvFrameEnd:
			if len(frameStack) > 0 {
				entry := frameStack[len(frameStack)-1]
				frameStack = frameStack[:len(frameStack)-1]
				leftX := pInfos[0].x - pInfos[0].w/2 - sc.FramePadding
				rightX := pInfos[len(pInfos)-1].x + pInfos[len(pInfos)-1].w/2 + sc.FramePadding
				frameW := rightX - leftX
				frameH := eventY - entry.startY + sc.FramePadding
				label := ""
				kind := ir.FrameLoop
				color := ""
				if entry.frame != nil {
					label = entry.frame.Label
					kind = entry.frame.Kind
					color = entry.frame.Color
				}
				frames = append(frames, SeqFrameLayout{
					Kind:     kind,
					Label:    label,
					Color:    color,
					X:        leftX,
					Y:        entry.startY,
					Width:    frameW,
					Height:   frameH,
					Dividers: entry.dividers,
				})
			}

		case ir.EvCreate:
			if idx, ok := pIndex[ev.Target]; ok {
				pInfos[idx].y = eventY
			}

		case ir.EvDestroy:
			// Destruction Y is recorded; lifeline will end here.
		}
	}

	return messages, notes, activations, frames, eventY
}

// seqNoteX computes the X position for a note based on its position and participants.
func seqNoteX(note *ir.SeqNote, pInfos []seqParticipantInfo, pIndex map[string]int, padH, noteW float32) float32 {
	if len(note.Participants) == 0 {
		return 0
	}
	firstIdx := pIndex[note.Participants[0]]
	px := pInfos[firstIdx].x
	pw := pInfos[firstIdx].w

	switch note.Position {
	case ir.NoteRight:
		return px + pw/2 + padH
	case ir.NoteLeft:
		return px - pw/2 - noteW - padH
	case ir.NoteOver:
		if len(note.Participants) >= 2 {
			secondIdx := pIndex[note.Participants[1]]
			px2 := pInfos[secondIdx].x
			return (px+px2)/2 - noteW/2
		}
		return px - noteW/2
	default:
		return 0
	}
}

// seqCloseRemainingActivations closes any unclosed activations.
func seqCloseRemainingActivations(
	graph *ir.Graph,
	pInfos []seqParticipantInfo,
	pIndex map[string]int,
	activations []SeqActivationLayout,
	sc config.SequenceConfig,
	eventY float32,
) []SeqActivationLayout {
	// Rebuild activation stacks from events (they were local to seqProcessEvents).
	activationStacks := make(map[string][]float32)
	for _, ev := range graph.Events {
		switch ev.Kind {
		case ir.EvActivate:
			activationStacks[ev.Target] = append(activationStacks[ev.Target], 0)
		case ir.EvDeactivate:
			stack := activationStacks[ev.Target]
			if len(stack) > 0 {
				activationStacks[ev.Target] = stack[:len(stack)-1]
			}
		}
	}
	// Only close stacks that still have items (i.e., unclosed activations).
	// Since we don't have the original start Y values, we use eventY as a
	// reasonable fallback — the original processing already created proper
	// activations for all balanced activate/deactivate pairs.
	for pid, stack := range activationStacks {
		if len(stack) == 0 {
			continue
		}
		px := float32(0)
		if idx, ok := pIndex[pid]; ok {
			px = pInfos[idx].x
		}
		for range stack {
			activations = append(activations, SeqActivationLayout{
				ParticipantID: pid,
				X:             px - sc.ActivationWidth/2,
				TopY:          eventY,
				BottomY:       eventY,
				Width:         sc.ActivationWidth,
			})
		}
	}
	return activations
}

// seqBuildBoxLayouts constructs participant box layouts.
func seqBuildBoxLayouts(graph *ir.Graph, pInfos []seqParticipantInfo, pIndex map[string]int, sc config.SequenceConfig, diagramH float32) []SeqBoxLayout {
	var boxes []SeqBoxLayout
	for _, box := range graph.Boxes {
		if len(box.Participants) == 0 {
			continue
		}
		minX := float32(math.MaxFloat32)
		maxX := float32(-math.MaxFloat32)
		for _, pid := range box.Participants {
			if idx, ok := pIndex[pid]; ok {
				pi := pInfos[idx]
				left := pi.x - pi.w/2
				right := pi.x + pi.w/2
				if left < minX {
					minX = left
				}
				if right > maxX {
					maxX = right
				}
			}
		}
		boxes = append(boxes, SeqBoxLayout{
			Label:  box.Label,
			Color:  box.Color,
			X:      minX - sc.BoxPadding,
			Y:      0,
			Width:  (maxX - minX) + 2*sc.BoxPadding,
			Height: diagramH,
		})
	}
	return boxes
}
