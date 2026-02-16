package render

import (
	"fmt"
	"strconv"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Sequence diagram rendering constants.
const (
	seqBoxBorderRadius    float32 = 4
	seqBoxLabelOffsetX    float32 = 8
	seqBoxLabelOffsetY    float32 = 16
	seqBoxFontScale       float32 = 0.9
	seqFrameBorderRadius  float32 = 4
	seqFrameTabCharWidth  float32 = 0.6
	seqFrameTabPadding    float32 = 16
	seqFrameTabPadY       float32 = 8
	seqFrameLabelOffsetX  float32 = 6
	seqFrameFontScale     float32 = 0.85
	seqSelfBumpWidth      float32 = 40
	seqSelfBumpHeight     float32 = 30
	seqSelfTextOffsetX    float32 = 24
	seqTextAboveOffset    float32 = 6
	seqAutoNumFontScale   float32 = 0.6
	seqAutoNumTextScale   float32 = 0.7
	seqAutoNumBaselineAdj float32 = 0.3
	seqNoteBorderRadius   float32 = 4
	seqNoteLineHeight     float32 = 1.2
	seqNoteTextPadY       float32 = 4
	seqNoteTextPadX       float32 = 8
	seqPartBorderRadius   float32 = 4
	seqKindFontScale      float32 = 0.65
	seqTextBaselineAdj    float32 = 0.35
	seqHeadScale          float32 = 0.15
	seqBodyScale          float32 = 0.3
	seqArmYFraction       float32 = 0.3
	seqArmSpanScale       float32 = 0.25
	seqLegLenScale        float32 = 0.25
	seqDbEllipseScale     float32 = 0.12
	seqPersonHeadRadius   float32 = 12
	seqPersonHeadCY       float32 = 18
	seqPersonBodyWidth    float32 = 20
	seqPersonBodyArc      float32 = 24
	seqPersonTextStart    float32 = 50
	seqPersonBorderRadius float32 = 6
)

// renderSequence renders all sequence diagram elements in visual stacking
// order (back to front): boxes, frames, lifelines, activations, messages,
// notes, then participants.
func renderSequence(builder *svgBuilder, computed *layout.Layout, th *theme.Theme, _ *config.Layout) {
	sd, ok := computed.Diagram.(layout.SequenceData)
	if !ok {
		return
	}
	renderSeqBoxes(builder, &sd, th)
	renderSeqFrames(builder, &sd, th)
	renderSeqLifelines(builder, &sd, th)
	renderSeqActivations(builder, &sd, th)
	renderSeqMessages(builder, &sd, th)
	renderSeqNotes(builder, &sd, th)
	renderSeqParticipants(builder, &sd, th)
}

// renderSeqBoxes renders participant box groups as rounded rectangles.
func renderSeqBoxes(builder *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, box := range sd.Boxes {
		fill := box.Color
		if fill == "" {
			fill = "rgba(0,0,0,0.05)"
		}
		attrs := []string{
			"fill", fill,
			"stroke", th.ClusterBorder,
			"stroke-width", "1",
		}
		if box.Color != "" && !isTransparentColor(box.Color) {
			attrs = append(attrs, "fill-opacity", "0.15")
		}
		builder.rect(box.X, box.Y, box.Width, box.Height, seqBoxBorderRadius, attrs...)
		if box.Label != "" {
			builder.text(box.X+seqBoxLabelOffsetX, box.Y+seqBoxLabelOffsetY, box.Label,
				"fill", th.TextColor,
				"font-size", fmtFloat(th.FontSize*seqBoxFontScale),
				"font-weight", "bold",
			)
		}
	}
}

// renderSeqFrames renders combined fragment frames (loop, alt, opt, etc.).
func renderSeqFrames(builder *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, frame := range sd.Frames {
		fill := "rgba(0,0,0,0.03)"
		if frame.Kind == ir.FrameRect && frame.Color != "" {
			fill = frame.Color
		}

		// Outer frame rect.
		builder.rect(frame.X, frame.Y, frame.Width, frame.Height, seqFrameBorderRadius,
			"fill", fill,
			"stroke", th.ClusterBorder,
			"stroke-width", "1",
		)

		// Label tab in top-left corner.
		kindLabel := frame.Kind.String()
		tabW := float32(len(kindLabel))*th.FontSize*seqFrameTabCharWidth + seqFrameTabPadding
		tabH := th.FontSize + seqFrameTabPadY
		builder.rect(frame.X, frame.Y, tabW, tabH, seqFrameBorderRadius,
			"fill", th.ClusterBorder,
			"stroke", th.ClusterBorder,
			"stroke-width", "1",
		)
		builder.text(frame.X+seqFrameLabelOffsetX, frame.Y+th.FontSize+1, kindLabel,
			"fill", th.LoopTextColor,
			"font-size", fmtFloat(th.FontSize*seqFrameFontScale),
			"font-weight", "bold",
		)

		// Condition/label text after the tab.
		if frame.Label != "" {
			builder.text(frame.X+tabW+seqFrameLabelOffsetX, frame.Y+th.FontSize+1, frame.Label,
				"fill", th.LoopTextColor,
				"font-size", fmtFloat(th.FontSize*seqFrameFontScale),
			)
		}

		// Divider lines for alt/par/critical fragments.
		switch frame.Kind {
		case ir.FrameAlt, ir.FramePar, ir.FrameCritical:
			for _, divY := range frame.Dividers {
				builder.line(frame.X, divY, frame.X+frame.Width, divY,
					"stroke", th.ClusterBorder,
					"stroke-width", "1",
					"stroke-dasharray", "5,5",
				)
			}
		}
	}
}

// renderSeqLifelines renders vertical dashed lines for each participant.
func renderSeqLifelines(builder *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, ll := range sd.Lifelines {
		builder.line(ll.X, ll.TopY, ll.X, ll.BottomY,
			"stroke", th.ActorLineColor,
			"stroke-width", "1",
			"stroke-dasharray", "5,5",
		)
	}
}

// renderSeqActivations renders narrow filled rectangles for activation bars.
func renderSeqActivations(builder *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, act := range sd.Activations {
		builder.rect(act.X, act.TopY, act.Width, act.BottomY-act.TopY, 2,
			"fill", th.ActivationBackground,
			"stroke", th.ActivationBorderColor,
			"stroke-width", "1",
		)
	}
}

// renderSeqMessages renders message arrows between participants.
func renderSeqMessages(builder *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, msg := range sd.Messages {
		isSelf := msg.From == msg.To

		attrs := []string{
			"stroke", th.SignalColor,
			"stroke-width", "1.5",
			"fill", "none",
		}

		// Dotted stroke for dotted message kinds.
		if msg.Kind.IsDotted() {
			attrs = append(attrs, "stroke-dasharray", "5,5")
		}

		// Arrow markers based on message kind.
		switch msg.Kind {
		case ir.MsgSolidArrow, ir.MsgDottedArrow:
			attrs = append(attrs, "marker-end", "url(#arrowhead)")
		case ir.MsgSolidOpen, ir.MsgDottedOpen:
			attrs = append(attrs, "marker-end", "url(#marker-open-arrow)")
		case ir.MsgSolidCross, ir.MsgDottedCross:
			attrs = append(attrs, "marker-end", "url(#marker-cross)")
		case ir.MsgBiSolid, ir.MsgBiDotted:
			attrs = append(attrs, "marker-start", "url(#arrowhead-start)")
			attrs = append(attrs, "marker-end", "url(#arrowhead)")
		}
		// MsgSolid, MsgDotted: no markers (plain line).

		if isSelf {
			// Self-message: draw a right-bump loop path.
			pathData := fmt.Sprintf("M %s,%s L %s,%s L %s,%s L %s,%s",
				fmtFloat(msg.FromX), fmtFloat(msg.Y),
				fmtFloat(msg.FromX+seqSelfBumpWidth), fmtFloat(msg.Y),
				fmtFloat(msg.FromX+seqSelfBumpWidth), fmtFloat(msg.Y+seqSelfBumpHeight),
				fmtFloat(msg.FromX), fmtFloat(msg.Y+seqSelfBumpHeight),
			)
			builder.path(pathData, attrs...)
		} else {
			builder.line(msg.FromX, msg.Y, msg.ToX, msg.Y, attrs...)
		}

		// Message text label above the arrow line.
		if len(msg.Text.Lines) > 0 && msg.Text.Lines[0] != "" {
			var textX float32
			if isSelf {
				textX = msg.FromX + seqSelfTextOffsetX
			} else {
				textX = (msg.FromX + msg.ToX) / 2
			}
			textY := msg.Y - seqTextAboveOffset

			builder.text(textX, textY, msg.Text.Lines[0],
				"text-anchor", "middle",
				"fill", th.SignalTextColor,
				"font-size", fmtFloat(th.FontSize),
			)
		}

		// Autonumber: filled circle with number at the start of the arrow.
		if msg.Number > 0 {
			numR := th.FontSize * seqAutoNumFontScale
			numX := msg.FromX
			numY := msg.Y

			builder.circle(numX, numY, numR,
				"fill", th.SignalColor,
				"stroke", "none",
			)
			builder.text(numX, numY+th.FontSize*seqAutoNumBaselineAdj, strconv.Itoa(msg.Number),
				"text-anchor", "middle",
				"fill", th.SequenceNumberColor,
				"font-size", fmtFloat(th.FontSize*seqAutoNumTextScale),
				"font-weight", "bold",
			)
		}
	}
}

// renderSeqNotes renders note boxes with text content.
func renderSeqNotes(builder *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, note := range sd.Notes {
		builder.rect(note.X, note.Y, note.Width, note.Height, seqNoteBorderRadius,
			"fill", th.NoteBackground,
			"stroke", th.NoteBorderColor,
			"stroke-width", "1",
		)

		// Render text lines inside the note.
		fontSize := note.Text.FontSize
		if fontSize <= 0 {
			fontSize = th.FontSize
		}
		lineH := fontSize * seqNoteLineHeight
		startY := note.Y + lineH + seqNoteTextPadY

		for idx, line := range note.Text.Lines {
			ly := startY + float32(idx)*lineH
			builder.text(note.X+seqNoteTextPadX, ly, line,
				"fill", th.NoteTextColor,
				"font-size", fmtFloat(fontSize),
			)
		}
	}
}

// renderSeqParticipants renders participant headers (at the top of the diagram)
// and footers (at the bottom). Different participant kinds get different shapes.
func renderSeqParticipants(builder *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, participant := range sd.Participants {
		// Render header at top.
		renderSeqParticipantShape(builder, &participant, participant.X, participant.Y, th)

		// Render footer at bottom (mirror of header).
		footerY := sd.DiagramHeight - participant.Height
		renderSeqParticipantShape(builder, &participant, participant.X, footerY, th)
	}
}

// renderSeqParticipantShape renders a single participant shape at the given position.
func renderSeqParticipantShape(builder *svgBuilder, participant *layout.SeqParticipantLayout, cx, topY float32, th *theme.Theme) {
	label := ""
	if len(participant.Label.Lines) > 0 {
		label = participant.Label.Lines[0]
	}

	switch participant.Kind {
	case ir.ActorStickFigure:
		renderStickFigure(builder, cx, topY, participant.Height, label, th)

	case ir.ParticipantDatabase:
		renderDatabaseShape(builder, cx, topY, participant.Width, participant.Height, label, th)

	default:
		// ParticipantBox and all other kinds: rounded rect with label.
		posX := cx - participant.Width/2
		builder.rect(posX, topY, participant.Width, participant.Height, seqPartBorderRadius,
			"fill", th.ActorBackground,
			"stroke", th.ActorBorder,
			"stroke-width", "1.5",
		)

		// For non-standard kinds, add a small kind label above the main label.
		if participant.Kind != ir.ParticipantBox {
			kindStr := participant.Kind.String()
			builder.text(cx, topY+th.FontSize*seqBoxFontScale, "<<"+kindStr+">>",
				"text-anchor", "middle",
				"fill", th.ActorTextColor,
				"font-size", fmtFloat(th.FontSize*seqKindFontScale),
				"font-style", "italic",
			)
			// Main label below the kind annotation.
			builder.text(cx, topY+participant.Height/2+th.FontSize*seqTextBaselineAdj, label,
				"text-anchor", "middle",
				"fill", th.ActorTextColor,
				"font-size", fmtFloat(th.FontSize),
			)
		} else {
			// Standard participant: label centered.
			builder.text(cx, topY+participant.Height/2+th.FontSize*seqTextBaselineAdj, label,
				"text-anchor", "middle",
				"fill", th.ActorTextColor,
				"font-size", fmtFloat(th.FontSize),
			)
		}
	}
}

// renderStickFigure draws a simple stick figure: circle head, body line,
// arms line, leg lines, with a label below.
func renderStickFigure(builder *svgBuilder, cx, topY, height float32, label string, th *theme.Theme) {
	headR := height * seqHeadScale
	headCY := topY + headR + 2
	bodySt := headCY + headR
	bodyLen := height * seqBodyScale
	bodyEnd := bodySt + bodyLen
	armY := bodySt + bodyLen*seqArmYFraction
	armSpan := height * seqArmSpanScale
	legLen := height * seqLegLenScale

	// Head.
	builder.circle(cx, headCY, headR,
		"fill", "none",
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Body.
	builder.line(cx, bodySt, cx, bodyEnd,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Arms.
	builder.line(cx-armSpan, armY, cx+armSpan, armY,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Left leg.
	builder.line(cx, bodyEnd, cx-armSpan, bodyEnd+legLen,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Right leg.
	builder.line(cx, bodyEnd, cx+armSpan, bodyEnd+legLen,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Label below the figure.
	builder.text(cx, topY+height-2, label,
		"text-anchor", "middle",
		"fill", th.ActorTextColor,
		"font-size", fmtFloat(th.FontSize),
	)
}

// renderDatabaseShape draws a cylinder (rect body + ellipse caps) for database participants.
func renderDatabaseShape(builder *svgBuilder, cx, topY, width, height float32, label string, th *theme.Theme) {
	posX := cx - width/2
	ellipseRY := height * seqDbEllipseScale
	bodyTop := topY + ellipseRY
	bodyH := height - 2*ellipseRY

	// Body rect.
	builder.rect(posX, bodyTop, width, bodyH, 0,
		"fill", th.ActorBackground,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Top ellipse cap.
	builder.ellipse(cx, bodyTop, width/2, ellipseRY,
		"fill", th.ActorBackground,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Bottom ellipse cap.
	builder.ellipse(cx, bodyTop+bodyH, width/2, ellipseRY,
		"fill", th.ActorBackground,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Cover the body-top-ellipse overlap with a filled rect (no stroke).
	builder.rect(posX+1, bodyTop, width-2, ellipseRY, 0,
		"fill", th.ActorBackground,
		"stroke", "none",
	)

	// Label centered.
	builder.text(cx, topY+height/2+th.FontSize*seqTextBaselineAdj, label,
		"text-anchor", "middle",
		"fill", th.ActorTextColor,
		"font-size", fmtFloat(th.FontSize),
	)
}
