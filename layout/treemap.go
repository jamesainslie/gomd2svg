package layout

import (
	"math"
	"sort"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

// treemapColorCount is the number of distinct colors for treemap rects.
const treemapColorCount = 8

// treemapItem is an internal type pairing a TreemapNode with its computed
// total value and original index (used for colour assignment).
type treemapItem struct {
	node  *ir.TreemapNode
	value float64
	idx   int
}

// treemapSquarifiedRect associates a treemapItem with the rectangle assigned
// to it by the squarified algorithm.
type treemapSquarifiedRect struct {
	item treemapItem
	x, y float32
	w, h float32
}

// computeTreemapLayout builds a squarified treemap layout from the IR graph.
func computeTreemapLayout(graph *ir.Graph, _ *theme.Theme, cfg *config.Layout) *Layout {
	padX := cfg.Treemap.PaddingX
	padY := cfg.Treemap.PaddingY
	chartW := cfg.Treemap.ChartWidth
	chartH := cfg.Treemap.ChartHeight

	if graph.TreemapRoot == nil || graph.TreemapRoot.TotalValue() == 0 {
		return &Layout{
			Kind:    graph.Kind,
			Nodes:   map[string]*NodeLayout{},
			Width:   chartW + padX*2,
			Height:  chartH + padY*2,
			Diagram: TreemapData{Title: graph.TreemapTitle},
		}
	}

	var rects []TreemapRectLayout

	// If root has children, lay them out; otherwise treat root as single rect.
	if len(graph.TreemapRoot.Children) > 0 {
		rects = treemapLayoutChildren(
			graph.TreemapRoot.Children,
			padX, padY, chartW, chartH,
			cfg.Treemap.Padding, cfg.Treemap.HeaderHeight,
			0, 0,
		)
	} else {
		rects = []TreemapRectLayout{{
			Label:      graph.TreemapRoot.Label,
			Value:      graph.TreemapRoot.Value,
			X:          padX,
			Y:          padY,
			Width:      chartW,
			Height:     chartH,
			Depth:      0,
			ColorIndex: 0,
		}}
	}

	return &Layout{
		Kind:    graph.Kind,
		Nodes:   map[string]*NodeLayout{},
		Width:   chartW + padX*2,
		Height:  chartH + padY*2,
		Diagram: TreemapData{Rects: rects, Title: graph.TreemapTitle},
	}
}

// treemapLayoutChildren lays out children within a given rectangle using the
// squarified treemap algorithm (Bruls-Huizing-van Wijk). Section nodes get a
// header band and their children are recursively laid out beneath it.
func treemapLayoutChildren( //nolint:revive // argument-limit: recursive layout requires rectangle + styling + depth params
	children []*ir.TreemapNode,
	rectX, rectY, width, height float32,
	padding, headerH float32,
	depth, colorStart int,
) []TreemapRectLayout {
	var rects []TreemapRectLayout

	// Build items with positive values, sorted descending by value.
	var items []treemapItem
	for i, c := range children {
		v := c.TotalValue()
		if v > 0 {
			items = append(items, treemapItem{node: c, value: v, idx: i})
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].value > items[j].value })

	if len(items) == 0 {
		return rects
	}

	totalValue := float64(0)
	for _, it := range items {
		totalValue += it.value
	}

	squarifiedRects := treemapSquarify(items, rectX, rectY, width, height, totalValue)

	for _, sr := range squarifiedRects {
		it := sr.item
		colorIdx := (colorStart + it.idx) % treemapColorCount

		if it.node.IsLeaf() {
			rects = append(rects, TreemapRectLayout{
				Label:      it.node.Label,
				Value:      it.node.Value,
				X:          sr.x + padding/2,
				Y:          sr.y + padding/2,
				Width:      sr.w - padding,
				Height:     sr.h - padding,
				Depth:      depth,
				ColorIndex: colorIdx,
			})
		} else {
			// Section node: header rect, then recurse into children below it.
			rects = append(rects, TreemapRectLayout{
				Label:      it.node.Label,
				X:          sr.x + padding/2,
				Y:          sr.y + padding/2,
				Width:      sr.w - padding,
				Height:     sr.h - padding,
				Depth:      depth,
				IsSection:  true,
				ColorIndex: colorIdx,
			})
			innerX := sr.x + padding
			innerY := sr.y + padding + headerH
			innerW := sr.w - padding*2
			innerH := sr.h - padding*2 - headerH
			if innerW > 0 && innerH > 0 {
				childRects := treemapLayoutChildren(
					it.node.Children,
					innerX, innerY, innerW, innerH,
					padding, headerH,
					depth+1, colorIdx,
				)
				rects = append(rects, childRects...)
			}
		}
	}

	return rects
}

// treemapSquarify implements the squarified treemap algorithm. Items must be
// sorted in descending value order. The algorithm greedily assigns items to
// rows, choosing the layout direction (horizontal or vertical) based on the
// shorter side of the remaining rectangle and adding items to the current row
// as long as doing so improves (or maintains) the worst aspect ratio.
func treemapSquarify(items []treemapItem, rectX, rectY, width, height float32, totalValue float64) []treemapSquarifiedRect {
	if len(items) == 0 {
		return nil
	}

	var result []treemapSquarifiedRect
	remaining := make([]treemapItem, len(items))
	copy(remaining, items)

	rx, ry, rw, rh := rectX, rectY, width, height
	remainingValue := totalValue

	for len(remaining) > 0 {
		row := []treemapItem{remaining[0]}
		remaining = remaining[1:]
		rowValue := row[0].value

		// Try adding more items to the row while aspect ratio improves.
		for len(remaining) > 0 {
			testRow := make([]treemapItem, len(row)+1)
			copy(testRow, row)
			testRow[len(row)] = remaining[0]
			testValue := rowValue + remaining[0].value

			if treemapWorstAspect(testRow, testValue, rw, rh, remainingValue) >
				treemapWorstAspect(row, rowValue, rw, rh, remainingValue) {
				break
			}
			row = testRow
			rowValue = testValue
			remaining = remaining[1:]
		}

		// Lay out the finalised row within the remaining rectangle.
		fraction := float32(rowValue / remainingValue)
		if rw >= rh {
			// Lay row out as a vertical strip on the left side.
			rowW := rw * fraction
			iy := ry
			for _, it := range row {
				itemFrac := float32(it.value / rowValue)
				itemH := rh * itemFrac
				result = append(result, treemapSquarifiedRect{
					item: it, x: rx, y: iy, w: rowW, h: itemH,
				})
				iy += itemH
			}
			rx += rowW
			rw -= rowW
		} else {
			// Lay row out as a horizontal strip on the top side.
			rowH := rh * fraction
			ix := rx
			for _, it := range row {
				itemFrac := float32(it.value / rowValue)
				itemW := rw * itemFrac
				result = append(result, treemapSquarifiedRect{
					item: it, x: ix, y: ry, w: itemW, h: rowH,
				})
				ix += itemW
			}
			ry += rowH
			rh -= rowH
		}

		remainingValue -= rowValue
		if remainingValue <= 0 {
			break
		}
	}

	return result
}

// treemapWorstAspect computes the worst aspect ratio among items in a
// candidate row. The row occupies a fraction of the remaining rectangle
// determined by rowValue/totalValue. Lower values are better (1.0 is a
// perfect square).
func treemapWorstAspect(row []treemapItem, rowValue float64, width, height float32, totalValue float64) float64 {
	if totalValue == 0 || len(row) == 0 {
		return math.MaxFloat64
	}

	fraction := rowValue / totalValue
	var side, otherSide float32
	if width >= height {
		side = width * float32(fraction)
		otherSide = height
	} else {
		side = height * float32(fraction)
		otherSide = width
	}

	if side == 0 || otherSide == 0 {
		return math.MaxFloat64
	}

	worst := float64(0)
	for _, it := range row {
		itemFrac := it.value / rowValue
		itemSize := float64(otherSide) * itemFrac
		if itemSize == 0 {
			return math.MaxFloat64
		}
		aspect := float64(side) / itemSize
		if aspect < 1 {
			aspect = 1 / aspect
		}
		if aspect > worst {
			worst = aspect
		}
	}
	return worst
}
