package layout

import (
	"sort"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeJourneyLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	jcfg := cfg.Journey

	// Collect unique actors
	actorSet := make(map[string]bool)
	for _, t := range g.JourneyTasks {
		for _, a := range t.Actors {
			actorSet[a] = true
		}
	}
	var actorNames []string
	for a := range actorSet {
		actorNames = append(actorNames, a)
	}
	sort.Strings(actorNames)
	var actors []JourneyActorLayout
	for idx, a := range actorNames {
		actors = append(actors, JourneyActorLayout{Name: a, ColorIndex: idx})
	}

	// Title height
	var titleH float32
	if g.JourneyTitle != "" {
		titleH = 30
	}

	trackY := jcfg.PaddingY + titleH + 10
	trackH := jcfg.TrackHeight

	// Build section layouts
	curX := jcfg.PaddingX
	var sections []JourneySectionLayout

	if len(g.JourneySections) == 0 {
		// No sections — lay out all tasks in a single implicit section
		var tasks []JourneyTaskLayout
		for _, t := range g.JourneyTasks {
			tw := jcfg.TaskWidth
			labelW := measurer.Width(t.Name, 14, th.FontFamily)
			if labelW+20 > tw {
				tw = labelW + 20
			}
			// Score 5 = top, score 1 = bottom
			scoreRatio := float32(t.Score-1) / 4.0
			taskY := trackY + trackH*(1-scoreRatio) - jcfg.TaskHeight/2
			tasks = append(tasks, JourneyTaskLayout{
				Label:  t.Name,
				Score:  t.Score,
				X:      curX + tw/2,
				Y:      taskY + jcfg.TaskHeight/2,
				Width:  tw,
				Height: jcfg.TaskHeight,
			})
			curX += tw + jcfg.TaskSpacing
		}
		if len(tasks) > 0 {
			secW := curX - jcfg.PaddingX - jcfg.TaskSpacing
			sections = append(sections, JourneySectionLayout{
				Label:  "",
				X:      jcfg.PaddingX,
				Y:      trackY,
				Width:  secW,
				Height: trackH,
				Tasks:  tasks,
			})
		}
	} else {
		for si, sec := range g.JourneySections {
			secStartX := curX
			var tasks []JourneyTaskLayout
			for _, ti := range sec.Tasks {
				if ti >= len(g.JourneyTasks) {
					continue
				}
				t := g.JourneyTasks[ti]
				tw := jcfg.TaskWidth
				labelW := measurer.Width(t.Name, 14, th.FontFamily)
				if labelW+20 > tw {
					tw = labelW + 20
				}
				scoreRatio := float32(t.Score-1) / 4.0
				taskY := trackY + trackH*(1-scoreRatio) - jcfg.TaskHeight/2
				tasks = append(tasks, JourneyTaskLayout{
					Label:  t.Name,
					Score:  t.Score,
					X:      curX + tw/2,
					Y:      taskY + jcfg.TaskHeight/2,
					Width:  tw,
					Height: jcfg.TaskHeight,
				})
				curX += tw + jcfg.TaskSpacing
			}
			secW := curX - secStartX
			if len(tasks) > 0 {
				secW -= jcfg.TaskSpacing // remove trailing spacing
			}
			if secW < 0 {
				secW = 0
			}

			color := ""
			if len(th.JourneySectionColors) > 0 {
				color = th.JourneySectionColors[si%len(th.JourneySectionColors)]
			}

			sections = append(sections, JourneySectionLayout{
				Label:  sec.Name,
				X:      secStartX,
				Y:      trackY,
				Width:  secW,
				Height: trackH,
				Color:  color,
				Tasks:  tasks,
			})
			curX += jcfg.SectionGap
		}
	}

	totalW := curX + jcfg.PaddingX
	actorLegendH := float32(0)
	if len(actors) > 0 {
		actorLegendH = 30
	}
	totalH := trackY + trackH + actorLegendH + jcfg.PaddingY

	return &Layout{
		Kind:   g.Kind,
		Nodes:  make(map[string]*NodeLayout),
		Width:  totalW,
		Height: totalH,
		Diagram: JourneyData{
			Sections: sections,
			Title:    g.JourneyTitle,
			Actors:   actors,
			TrackY:   trackY,
			TrackH:   trackH,
		},
	}
}
