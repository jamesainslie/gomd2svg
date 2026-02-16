package layout

import (
	"sort"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Journey layout constants.
const (
	journeyTrackGap      float32 = 10
	journeyLabelFontSize float32 = 14
	journeyLabelPad      float32 = 20
	journeyScoreRange    float32 = 4.0
)

func computeJourneyLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	jcfg := cfg.Journey

	// Collect unique actors
	actorSet := make(map[string]bool)
	for _, task := range graph.JourneyTasks {
		for _, actor := range task.Actors {
			actorSet[actor] = true
		}
	}
	actorNames := make([]string, 0, len(actorSet))
	for actor := range actorSet {
		actorNames = append(actorNames, actor)
	}
	sort.Strings(actorNames)
	actors := make([]JourneyActorLayout, 0, len(actorNames))
	for idx, actor := range actorNames {
		actors = append(actors, JourneyActorLayout{Name: actor, ColorIndex: idx})
	}

	// Title height
	var titleH float32
	if graph.JourneyTitle != "" {
		titleH = 30
	}

	trackY := jcfg.PaddingY + titleH + journeyTrackGap
	trackH := jcfg.TrackHeight

	// Build section layouts
	curX := jcfg.PaddingX
	var sections []JourneySectionLayout

	if len(graph.JourneySections) == 0 {
		// No sections -- lay out all tasks in a single implicit section
		tasks := make([]JourneyTaskLayout, 0, len(graph.JourneyTasks))
		for _, task := range graph.JourneyTasks {
			tw := jcfg.TaskWidth
			labelW := measurer.Width(task.Name, journeyLabelFontSize, th.FontFamily)
			if labelW+journeyLabelPad > tw {
				tw = labelW + journeyLabelPad
			}
			// Score 5 = top, score 1 = bottom
			scoreRatio := float32(task.Score-1) / journeyScoreRange
			taskY := trackY + trackH*(1-scoreRatio) - jcfg.TaskHeight/2
			tasks = append(tasks, JourneyTaskLayout{
				Label:  task.Name,
				Score:  task.Score,
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
		for si, sec := range graph.JourneySections {
			secStartX := curX
			tasks := make([]JourneyTaskLayout, 0, len(sec.Tasks))
			for _, ti := range sec.Tasks {
				if ti >= len(graph.JourneyTasks) {
					continue
				}
				task := graph.JourneyTasks[ti]
				tw := jcfg.TaskWidth
				labelW := measurer.Width(task.Name, journeyLabelFontSize, th.FontFamily)
				if labelW+journeyLabelPad > tw {
					tw = labelW + journeyLabelPad
				}
				scoreRatio := float32(task.Score-1) / journeyScoreRange
				taskY := trackY + trackH*(1-scoreRatio) - jcfg.TaskHeight/2
				tasks = append(tasks, JourneyTaskLayout{
					Label:  task.Name,
					Score:  task.Score,
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
		Kind:   graph.Kind,
		Nodes:  make(map[string]*NodeLayout),
		Width:  totalW,
		Height: totalH,
		Diagram: JourneyData{
			Sections: sections,
			Title:    graph.JourneyTitle,
			Actors:   actors,
			TrackY:   trackY,
			TrackH:   trackH,
		},
	}
}
