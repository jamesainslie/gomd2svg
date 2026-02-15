package render

import (
	"fmt"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// renderSequence renders all sequence diagram elements in visual stacking
// order (back to front): boxes, frames, lifelines, activations, messages,
// notes, then participants.
func renderSequence(b *svgBuilder, l *layout.Layout, th *theme.Theme, _ *config.Layout) {
	sd, ok := l.Diagram.(layout.SequenceData)
	if !ok {
		return
	}
	renderSeqBoxes(b, &sd, th)
	renderSeqFrames(b, &sd, th)
	renderSeqLifelines(b, &sd, th)
	renderSeqActivations(b, &sd, th)
	renderSeqMessages(b, &sd, th)
	renderSeqNotes(b, &sd, th)
	renderSeqParticipants(b, &sd, th)
}

// renderSeqBoxes renders participant box groups as rounded rectangles.
func renderSeqBoxes(b *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, box := range sd.Boxes {
		fill := box.Color
		if fill == "" {
			fill = "rgba(0,0,0,0.05)"
		}
		b.rect(box.X, box.Y, box.Width, box.Height, 4,
			"fill", fill,
			"stroke", th.ClusterBorder,
			"stroke-width", "1",
		)
		if box.Label != "" {
			b.text(box.X+8, box.Y+16, box.Label,
				"fill", th.TextColor,
				"font-size", fmtFloat(th.FontSize*0.9),
				"font-weight", "bold",
			)
		}
	}
}

// renderSeqFrames renders combined fragment frames (loop, alt, opt, etc.).
func renderSeqFrames(b *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, frame := range sd.Frames {
		fill := "rgba(0,0,0,0.03)"
		if frame.Kind == ir.FrameRect && frame.Color != "" {
			fill = frame.Color
		}

		// Outer frame rect.
		b.rect(frame.X, frame.Y, frame.Width, frame.Height, 4,
			"fill", fill,
			"stroke", th.ClusterBorder,
			"stroke-width", "1",
		)

		// Label tab in top-left corner.
		kindLabel := frame.Kind.String()
		tabW := float32(len(kindLabel))*th.FontSize*0.6 + 16
		tabH := th.FontSize + 8
		b.rect(frame.X, frame.Y, tabW, tabH, 4,
			"fill", th.ClusterBorder,
			"stroke", th.ClusterBorder,
			"stroke-width", "1",
		)
		b.text(frame.X+6, frame.Y+th.FontSize+1, kindLabel,
			"fill", th.LoopTextColor,
			"font-size", fmtFloat(th.FontSize*0.85),
			"font-weight", "bold",
		)

		// Condition/label text after the tab.
		if frame.Label != "" {
			b.text(frame.X+tabW+6, frame.Y+th.FontSize+1, frame.Label,
				"fill", th.LoopTextColor,
				"font-size", fmtFloat(th.FontSize*0.85),
			)
		}

		// Divider lines for alt/par/critical fragments.
		switch frame.Kind {
		case ir.FrameAlt, ir.FramePar, ir.FrameCritical:
			for _, divY := range frame.Dividers {
				b.line(frame.X, divY, frame.X+frame.Width, divY,
					"stroke", th.ClusterBorder,
					"stroke-width", "1",
					"stroke-dasharray", "5,5",
				)
			}
		}
	}
}

// renderSeqLifelines renders vertical dashed lines for each participant.
func renderSeqLifelines(b *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, ll := range sd.Lifelines {
		b.line(ll.X, ll.TopY, ll.X, ll.BottomY,
			"stroke", th.ActorLineColor,
			"stroke-width", "1",
			"stroke-dasharray", "5,5",
		)
	}
}

// renderSeqActivations renders narrow filled rectangles for activation bars.
func renderSeqActivations(b *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, act := range sd.Activations {
		b.rect(act.X, act.TopY, act.Width, act.BottomY-act.TopY, 2,
			"fill", th.ActivationBackground,
			"stroke", th.ActivationBorderColor,
			"stroke-width", "1",
		)
	}
}

// renderSeqMessages renders message arrows between participants.
func renderSeqMessages(b *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
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
			bumpW := float32(40)
			bumpH := float32(30)
			d := fmt.Sprintf("M %s,%s L %s,%s L %s,%s L %s,%s",
				fmtFloat(msg.FromX), fmtFloat(msg.Y),
				fmtFloat(msg.FromX+bumpW), fmtFloat(msg.Y),
				fmtFloat(msg.FromX+bumpW), fmtFloat(msg.Y+bumpH),
				fmtFloat(msg.FromX), fmtFloat(msg.Y+bumpH),
			)
			b.path(d, attrs...)
		} else {
			b.line(msg.FromX, msg.Y, msg.ToX, msg.Y, attrs...)
		}

		// Message text label above the arrow line.
		if len(msg.Text.Lines) > 0 && msg.Text.Lines[0] != "" {
			var textX float32
			if isSelf {
				textX = msg.FromX + 24
			} else {
				textX = (msg.FromX + msg.ToX) / 2
			}
			textY := msg.Y - 6

			b.text(textX, textY, msg.Text.Lines[0],
				"text-anchor", "middle",
				"fill", th.SignalTextColor,
				"font-size", fmtFloat(th.FontSize),
			)
		}

		// Autonumber: filled circle with number at the start of the arrow.
		if msg.Number > 0 {
			numR := th.FontSize * 0.6
			numX := msg.FromX
			numY := msg.Y

			b.circle(numX, numY, numR,
				"fill", th.SignalColor,
				"stroke", "none",
			)
			b.text(numX, numY+th.FontSize*0.3, fmt.Sprintf("%d", msg.Number),
				"text-anchor", "middle",
				"fill", th.SequenceNumberColor,
				"font-size", fmtFloat(th.FontSize*0.7),
				"font-weight", "bold",
			)
		}
	}
}

// renderSeqNotes renders note boxes with text content.
func renderSeqNotes(b *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, note := range sd.Notes {
		b.rect(note.X, note.Y, note.Width, note.Height, 4,
			"fill", th.NoteBackground,
			"stroke", th.NoteBorderColor,
			"stroke-width", "1",
		)

		// Render text lines inside the note.
		fontSize := note.Text.FontSize
		if fontSize <= 0 {
			fontSize = th.FontSize
		}
		lineH := fontSize * 1.2
		startY := note.Y + lineH + 4

		for i, line := range note.Text.Lines {
			ly := startY + float32(i)*lineH
			b.text(note.X+8, ly, line,
				"fill", th.NoteTextColor,
				"font-size", fmtFloat(fontSize),
			)
		}
	}
}

// renderSeqParticipants renders participant headers (at the top of the diagram)
// and footers (at the bottom). Different participant kinds get different shapes.
func renderSeqParticipants(b *svgBuilder, sd *layout.SequenceData, th *theme.Theme) {
	for _, p := range sd.Participants {
		// Render header at top.
		renderSeqParticipantShape(b, &p, p.X, p.Y, th)

		// Render footer at bottom (mirror of header).
		footerY := sd.DiagramHeight - p.Height
		renderSeqParticipantShape(b, &p, p.X, footerY, th)
	}
}

// renderSeqParticipantShape renders a single participant shape at the given position.
func renderSeqParticipantShape(b *svgBuilder, p *layout.SeqParticipantLayout, cx, topY float32, th *theme.Theme) {
	label := ""
	if len(p.Label.Lines) > 0 {
		label = p.Label.Lines[0]
	}

	switch p.Kind {
	case ir.ActorStickFigure:
		renderStickFigure(b, cx, topY, p.Height, label, th)

	case ir.ParticipantDatabase:
		renderDatabaseShape(b, cx, topY, p.Width, p.Height, label, th)

	default:
		// ParticipantBox and all other kinds: rounded rect with label.
		x := cx - p.Width/2
		b.rect(x, topY, p.Width, p.Height, 4,
			"fill", th.ActorBackground,
			"stroke", th.ActorBorder,
			"stroke-width", "1.5",
		)

		// For non-standard kinds, add a small kind label above the main label.
		if p.Kind != ir.ParticipantBox {
			kindStr := p.Kind.String()
			b.text(cx, topY+th.FontSize*0.9, "<<"+kindStr+">>",
				"text-anchor", "middle",
				"fill", th.ActorTextColor,
				"font-size", fmtFloat(th.FontSize*0.65),
				"font-style", "italic",
			)
			// Main label below the kind annotation.
			b.text(cx, topY+p.Height/2+th.FontSize*0.35, label,
				"text-anchor", "middle",
				"fill", th.ActorTextColor,
				"font-size", fmtFloat(th.FontSize),
			)
		} else {
			// Standard participant: label centered.
			b.text(cx, topY+p.Height/2+th.FontSize*0.35, label,
				"text-anchor", "middle",
				"fill", th.ActorTextColor,
				"font-size", fmtFloat(th.FontSize),
			)
		}
	}
}

// renderStickFigure draws a simple stick figure: circle head, body line,
// arms line, leg lines, with a label below.
func renderStickFigure(b *svgBuilder, cx, topY, height float32, label string, th *theme.Theme) {
	headR := height * 0.15
	headCY := topY + headR + 2
	bodySt := headCY + headR
	bodyLen := height * 0.3
	bodyEnd := bodySt + bodyLen
	armY := bodySt + bodyLen*0.3
	armSpan := height * 0.25
	legLen := height * 0.25

	// Head.
	b.circle(cx, headCY, headR,
		"fill", "none",
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Body.
	b.line(cx, bodySt, cx, bodyEnd,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Arms.
	b.line(cx-armSpan, armY, cx+armSpan, armY,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Left leg.
	b.line(cx, bodyEnd, cx-armSpan, bodyEnd+legLen,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Right leg.
	b.line(cx, bodyEnd, cx+armSpan, bodyEnd+legLen,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Label below the figure.
	b.text(cx, topY+height-2, label,
		"text-anchor", "middle",
		"fill", th.ActorTextColor,
		"font-size", fmtFloat(th.FontSize),
	)
}

// renderDatabaseShape draws a cylinder (rect body + ellipse caps) for database participants.
func renderDatabaseShape(b *svgBuilder, cx, topY, width, height float32, label string, th *theme.Theme) {
	x := cx - width/2
	ellipseRY := height * 0.12
	bodyTop := topY + ellipseRY
	bodyH := height - 2*ellipseRY

	// Body rect.
	b.rect(x, bodyTop, width, bodyH, 0,
		"fill", th.ActorBackground,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Top ellipse cap.
	b.ellipse(cx, bodyTop, width/2, ellipseRY,
		"fill", th.ActorBackground,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Bottom ellipse cap.
	b.ellipse(cx, bodyTop+bodyH, width/2, ellipseRY,
		"fill", th.ActorBackground,
		"stroke", th.ActorBorder,
		"stroke-width", "1.5",
	)

	// Cover the body-top-ellipse overlap with a filled rect (no stroke).
	b.rect(x+1, bodyTop, width-2, ellipseRY, 0,
		"fill", th.ActorBackground,
		"stroke", "none",
	)

	// Label centered.
	b.text(cx, topY+height/2+th.FontSize*0.35, label,
		"text-anchor", "middle",
		"fill", th.ActorTextColor,
		"font-size", fmtFloat(th.FontSize),
	)
}
