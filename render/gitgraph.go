package render

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// GitGraph rendering constants.
const (
	gitBranchLabelOffset float32 = 10
	gitReverseCrossScale float32 = 0.6
	gitTagPadding        float32 = 4
	gitTagHalfWidth      float32 = 20
	gitTagRectWidth      float32 = 40
	gitTagRectHeight     float32 = 14
	gitTagBorderRadius   float32 = 3
	gitTagOffsetY        float32 = 10
)

func renderGitGraph(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	ggd, ok := lay.Diagram.(layout.GitGraphData)
	if !ok {
		return
	}

	commitRadius := cfg.GitGraph.CommitRadius

	// Draw branch lines.
	for _, br := range ggd.Branches {
		if br.StartX < br.EndX {
			builder.line(br.StartX, br.Y, br.EndX, br.Y,
				"stroke", br.Color,
				"stroke-width", "2",
			)
		}
		// Branch label.
		builder.text(br.StartX-gitBranchLabelOffset, br.Y, br.Name,
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
		builder.line(conn.FromX, conn.FromY, conn.ToX, conn.ToY, attrs...)
	}

	// Build a map from branch name to color.
	branchColor := make(map[string]string)
	for _, br := range ggd.Branches {
		branchColor[br.Name] = br.Color
	}

	// Draw commits.
	for _, commit := range ggd.Commits {
		color := branchColor[commit.Branch]
		if color == "" {
			color = th.GitCommitFill
		}

		switch commit.Type {
		case ir.GitCommitHighlight:
			builder.circle(commit.X, commit.Y, commitRadius,
				"fill", th.GitHighlightFill,
				"stroke", color,
				"stroke-width", "2",
			)
		case ir.GitCommitReverse:
			// Reverse: filled circle with a cross.
			builder.circle(commit.X, commit.Y, commitRadius,
				"fill", color,
				"stroke", th.GitCommitStroke,
				"stroke-width", "2",
			)
			halfR := commitRadius * gitReverseCrossScale
			builder.line(commit.X-halfR, commit.Y-halfR, commit.X+halfR, commit.Y+halfR,
				"stroke", th.Background,
				"stroke-width", "2",
			)
			builder.line(commit.X-halfR, commit.Y+halfR, commit.X+halfR, commit.Y-halfR,
				"stroke", th.Background,
				"stroke-width", "2",
			)
		default:
			builder.circle(commit.X, commit.Y, commitRadius,
				"fill", color,
				"stroke", th.GitCommitStroke,
				"stroke-width", "2",
			)
		}

		// Tag label.
		if commit.Tag != "" {
			tagX := commit.X
			tagY := commit.Y - commitRadius - gitTagPadding
			builder.rect(tagX-gitTagHalfWidth, tagY-gitTagOffsetY, gitTagRectWidth, gitTagRectHeight, gitTagBorderRadius,
				"fill", th.GitTagFill,
				"stroke", th.GitTagBorder,
				"stroke-width", "1",
			)
			builder.text(tagX, tagY-1, commit.Tag,
				"text-anchor", "middle",
				"dominant-baseline", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(cfg.GitGraph.TagFontSize),
				"fill", th.TextColor,
			)
		}
	}
}
