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
func computeSequenceLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	sc := cfg.Sequence
	lineH := th.FontSize * cfg.LabelLineHeight
	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical

	// ----------------------------------------------------------------
	// Phase 1: Measure participants and assign horizontal positions.
	// ----------------------------------------------------------------

	type participantInfo struct {
		id    string
		label TextBlock
		kind  ir.SeqParticipantKind
		x     float32 // centre X
		y     float32 // top Y (0 for normal, set later for created)
		w     float32
		h     float32
	}

	pInfos := make([]participantInfo, len(g.Participants))
	pIndex := make(map[string]int, len(g.Participants)) // id -> index into pInfos

	cursorX := padH
	for i, p := range g.Participants {
		name := p.DisplayName()
		tw := measurer.Width(name, th.FontSize, th.FontFamily)
		w := tw + 2*padH
		h := sc.HeaderHeight

		pInfos[i] = participantInfo{
			id:    p.ID,
			label: TextBlock{Lines: []string{name}, Width: tw, Height: lineH, FontSize: th.FontSize},
			kind:  p.Kind,
			x:     cursorX + w/2, // centre
			y:     0,
			w:     w,
			h:     h,
		}
		pIndex[p.ID] = i

		cursorX += w + sc.ParticipantSpacing
	}

	// ----------------------------------------------------------------
	// Phase 2: Walk events top-to-bottom.
	// ----------------------------------------------------------------

	y := sc.HeaderHeight + padV

	// Activation stack per participant (stores start Y values).
	activationStacks := make(map[string][]float32)

	// Frame stack for nested combined fragments.
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

	for _, ev := range g.Events {
		switch ev.Kind {
		case ir.EvMessage:
			msg := ev.Message
			y += sc.MessageSpacing

			fromIdx, fromOK := pIndex[msg.From]
			toIdx, toOK := pIndex[msg.To]
			if !fromOK || !toOK {
				continue
			}

			fromX := pInfos[fromIdx].x
			toX := pInfos[toIdx].x

			// Self-message: bump to the right and return.
			if msg.From == msg.To {
				toX = fromX + sc.SelfMessageWidth
			}

			tw := measurer.Width(msg.Text, th.FontSize, th.FontFamily)

			msgNumber++
			num := 0
			if g.Autonumber {
				num = msgNumber
			}

			messages = append(messages, SeqMessageLayout{
				From:   msg.From,
				To:     msg.To,
				Text:   TextBlock{Lines: []string{msg.Text}, Width: tw, Height: lineH, FontSize: th.FontSize},
				Kind:   msg.Kind,
				Y:      y,
				FromX:  fromX,
				ToX:    toX,
				Number: num,
			})

		case ir.EvNote:
			note := ev.Note
			lines := strings.Split(note.Text, "\n")

			maxLineW := float32(0)
			for _, ln := range lines {
				lw := measurer.Width(ln, th.FontSize, th.FontFamily)
				if lw > maxLineW {
					maxLineW = lw
				}
			}

			noteW := maxLineW + 2*padH
			if noteW > sc.NoteMaxWidth {
				noteW = sc.NoteMaxWidth
			}
			noteH := float32(len(lines))*lineH + 2*padV

			var noteX float32
			if len(note.Participants) > 0 {
				firstIdx := pIndex[note.Participants[0]]
				px := pInfos[firstIdx].x
				pw := pInfos[firstIdx].w

				switch note.Position {
				case ir.NoteRight:
					noteX = px + pw/2 + padH
				case ir.NoteLeft:
					noteX = px - pw/2 - noteW - padH
				case ir.NoteOver:
					if len(note.Participants) >= 2 {
						secondIdx := pIndex[note.Participants[1]]
						px2 := pInfos[secondIdx].x
						noteX = (px+px2)/2 - noteW/2
					} else {
						noteX = px - noteW/2
					}
				}
			}

			notes = append(notes, SeqNoteLayout{
				Text:   TextBlock{Lines: lines, Width: maxLineW, Height: float32(len(lines)) * lineH, FontSize: th.FontSize},
				X:      noteX,
				Y:      y,
				Width:  noteW,
				Height: noteH,
			})

			y += noteH + padV

		case ir.EvActivate:
			activationStacks[ev.Target] = append(activationStacks[ev.Target], y)

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
					BottomY:       y,
					Width:         sc.ActivationWidth,
				})
			}

		case ir.EvFrameStart:
			frameStack = append(frameStack, frameEntry{
				frame:  ev.Frame,
				startY: y,
			})

		case ir.EvFrameMiddle:
			if len(frameStack) > 0 {
				frameStack[len(frameStack)-1].dividers = append(
					frameStack[len(frameStack)-1].dividers, y)
			}

		case ir.EvFrameEnd:
			if len(frameStack) > 0 {
				entry := frameStack[len(frameStack)-1]
				frameStack = frameStack[:len(frameStack)-1]

				// Frame spans the full participant range with padding.
				leftX := pInfos[0].x - pInfos[0].w/2 - sc.FramePadding
				rightX := pInfos[len(pInfos)-1].x + pInfos[len(pInfos)-1].w/2 + sc.FramePadding
				frameW := rightX - leftX
				frameH := y - entry.startY + sc.FramePadding

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
				pInfos[idx].y = y
			}

		case ir.EvDestroy:
			// Destruction Y is recorded; lifeline will end here.
			// We store it via the participant info y value as negative
			// signal, but it's simpler to use a separate map.
			// For now we just note it (lifeline bottom will be set in Phase 3).
		}
	}

	// ----------------------------------------------------------------
	// Phase 3: Finalize.
	// ----------------------------------------------------------------

	// Close any remaining activations.
	for pid, stack := range activationStacks {
		px := float32(0)
		if idx, ok := pIndex[pid]; ok {
			px = pInfos[idx].x
		}
		for _, startY := range stack {
			activations = append(activations, SeqActivationLayout{
				ParticipantID: pid,
				X:             px - sc.ActivationWidth/2,
				TopY:          startY,
				BottomY:       y,
				Width:         sc.ActivationWidth,
			})
		}
	}

	// Footer height matches header height.
	footerY := y + padV
	diagramH := footerY + sc.HeaderHeight + padV

	// Build participant layouts and lifelines.
	participants := make([]SeqParticipantLayout, len(pInfos))
	lifelines := make([]SeqLifeline, len(pInfos))

	for i, pi := range pInfos {
		participants[i] = SeqParticipantLayout{
			ID:     pi.id,
			Label:  pi.label,
			Kind:   pi.kind,
			X:      pi.x,
			Y:      pi.y,
			Width:  pi.w,
			Height: pi.h,
		}

		topY := pi.y + pi.h // lifeline starts below header
		bottomY := footerY  // lifeline ends at footer

		lifelines[i] = SeqLifeline{
			ParticipantID: pi.id,
			X:             pi.x,
			TopY:          topY,
			BottomY:       bottomY,
		}
	}

	// Build box layouts.
	var boxes []SeqBoxLayout
	for _, box := range g.Boxes {
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

	// Compute bounding box.
	rightEdge := float32(0)
	if len(pInfos) > 0 {
		last := pInfos[len(pInfos)-1]
		rightEdge = last.x + last.w/2 + padH
	}

	return &Layout{
		Kind:   g.Kind,
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
			Autonumber:    g.Autonumber,
			DiagramHeight: diagramH,
		},
	}
}
