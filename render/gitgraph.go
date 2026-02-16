package render

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func renderGitGraph(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	ggd, ok := l.Diagram.(layout.GitGraphData)
	if !ok {
		return
	}

	commitRadius := cfg.GitGraph.CommitRadius

	// Draw branch lines.
	for _, br := range ggd.Branches {
		if br.StartX < br.EndX {
			b.line(br.StartX, br.Y, br.EndX, br.Y,
				"stroke", br.Color,
				"stroke-width", "2",
			)
		}
		// Branch label.
		b.text(br.StartX-10, br.Y, br.Name,
			"text-anchor", "end",
			"dominant-baseline", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize-2),
			"fill", th.TextColor,
		)
	}

	// Draw merge/cherry-pick connections.
	for _, conn := range ggd.Connections {
		dashArray := ""
		if conn.IsCherryPick {
			dashArray = "4,4"
		}
		attrs := []string{
			"stroke", th.LineColor,
			"stroke-width", "1.5",
		}
		if dashArray != "" {
			attrs = append(attrs, "stroke-dasharray", dashArray)
		}
		b.line(conn.FromX, conn.FromY, conn.ToX, conn.ToY, attrs...)
	}

	// Build a map from branch name to color.
	branchColor := make(map[string]string)
	for _, br := range ggd.Branches {
		branchColor[br.Name] = br.Color
	}

	// Draw commits.
	for _, c := range ggd.Commits {
		color := branchColor[c.Branch]
		if color == "" {
			color = th.GitCommitFill
		}

		switch c.Type {
		case ir.GitCommitHighlight:
			b.circle(c.X, c.Y, commitRadius,
				"fill", th.GitHighlightFill,
				"stroke", color,
				"stroke-width", "2",
			)
		case ir.GitCommitReverse:
			// Reverse: filled circle with a cross.
			b.circle(c.X, c.Y, commitRadius,
				"fill", color,
				"stroke", th.GitCommitStroke,
				"stroke-width", "2",
			)
			halfR := commitRadius * 0.6
			b.line(c.X-halfR, c.Y-halfR, c.X+halfR, c.Y+halfR,
				"stroke", th.Background,
				"stroke-width", "2",
			)
			b.line(c.X-halfR, c.Y+halfR, c.X+halfR, c.Y-halfR,
				"stroke", th.Background,
				"stroke-width", "2",
			)
		default:
			b.circle(c.X, c.Y, commitRadius,
				"fill", color,
				"stroke", th.GitCommitStroke,
				"stroke-width", "2",
			)
		}

		// Tag label.
		if c.Tag != "" {
			tagX := c.X
			tagY := c.Y - commitRadius - 4
			b.rect(tagX-20, tagY-10, 40, 14, 3,
				"fill", th.GitTagFill,
				"stroke", th.GitTagBorder,
				"stroke-width", "1",
			)
			b.text(tagX, tagY-1, c.Tag,
				"text-anchor", "middle",
				"dominant-baseline", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(cfg.GitGraph.TagFontSize),
				"fill", th.TextColor,
			)
		}
	}
}
