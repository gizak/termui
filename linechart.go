// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"fmt"
	"math"
	"sort"
)

// only 16 possible combinations, why bother
var braillePatterns = map[[2]int]rune{
	[2]int{0, 0}: '⣀',
	[2]int{0, 1}: '⡠',
	[2]int{0, 2}: '⡐',
	[2]int{0, 3}: '⡈',

	[2]int{1, 0}: '⢄',
	[2]int{1, 1}: '⠤',
	[2]int{1, 2}: '⠔',
	[2]int{1, 3}: '⠌',

	[2]int{2, 0}: '⢂',
	[2]int{2, 1}: '⠢',
	[2]int{2, 2}: '⠒',
	[2]int{2, 3}: '⠊',

	[2]int{3, 0}: '⢁',
	[2]int{3, 1}: '⠡',
	[2]int{3, 2}: '⠑',
	[2]int{3, 3}: '⠉',
}

var lSingleBraille = [4]rune{'\u2840', '⠄', '⠂', '⠁'}
var rSingleBraille = [4]rune{'\u2880', '⠠', '⠐', '⠈'}

// LineChart has two modes: braille(default) and dot.
// A single braille character is a 2x4 grid of dots, so Using braille
// gives 2x X resolution and 4x Y resolution over dot mode.
/*
  lc := termui.NewLineChart()
  lc.Border.Label = "braille-mode Line Chart"
  lc.Data["name'] = [1.2, 1.3, 1.5, 1.7, 1.5, 1.6, 1.8, 2.0]
  lc.Width = 50
  lc.Height = 12
  lc.AxesColor = termui.ColorWhite
  lc.LineColor = termui.ColorGreen | termui.AttrBold
  // termui.Render(lc)...
*/
type LineChart struct {
	Block
	Data             map[string][]float64
	DataLabels       []string // if unset, the data indices will be used
	Mode             string   // braille | dot
	DotStyle         rune
	LineColor        map[string]Attribute
	defaultLineColor Attribute
	scale            float64 // data span per cell on y-axis
	AxesColor        Attribute
	drawingX         int
	drawingY         int
	axisYHeight      int
	axisXWidth       int
	axisYLabelGap    int
	axisXLabelGap    int
	topValue         float64
	bottomValue      float64
	labelX           [][]rune
	labelY           [][]rune
	labelYSpace      int
	maxY             float64
	minY             float64
	YPadding         float64
	YFloor           float64
	YCeil            float64
}

// NewLineChart returns a new LineChart with current theme.
func NewLineChart() *LineChart {
	return &LineChart{
		Block:            *NewBlock(),
		AxesColor:        ThemeAttr("linechart.axes.fg"),
		defaultLineColor: ThemeAttr("linechart.line.fg"),
		Mode:             "braille",
		DotStyle:         '•',
		Data:             make(map[string][]float64),
		LineColor:        make(map[string]Attribute),
		axisXLabelGap:    2,
		axisYLabelGap:    1,
		bottomValue:      math.Inf(1),
		topValue:         math.Inf(-1),
		YPadding:         0.2,
		YFloor:           math.Inf(-1),
		YCeil:            math.Inf(1),
	}
}

// one cell contains two data points, so capicity is 2x dot mode
func (lc *LineChart) renderBraille() Buffer {
	buf := NewBuffer()

	// return: b -> which cell should the point be in
	//         m -> in the cell, divided into 4 equal height levels, which subcell?
	getPos := func(d float64) (b, m int) {
		cnt4 := int((d-lc.bottomValue)/(lc.scale/4) + 0.5)
		b = cnt4 / 4
		m = cnt4 % 4
		return
	}

	// Sort the series so that overlapping data will overlap the same way each time
	seriesList := make([]string, len(lc.Data))
	i := 0
	for seriesName := range lc.Data {
		seriesList[i] = seriesName
		i++
	}
	sort.Strings(seriesList)

	// plot points
	for _, seriesName := range seriesList {
		seriesData := lc.Data[seriesName]
		if len(seriesData) == 0 {
			continue
		}
		thisLineColor, ok := lc.LineColor[seriesName]
		if !ok {
			thisLineColor = lc.defaultLineColor
		}

		minCell := lc.innerArea.Min.X + lc.labelYSpace
		cellPos := lc.innerArea.Max.X - 1
		for dataPos := len(seriesData) - 1; dataPos >= 0 && cellPos > minCell; {
			b0, m0 := getPos(seriesData[dataPos])
			var b1, m1 int

			if dataPos > 0 {
				b1, m1 = getPos(seriesData[dataPos-1])

				if b0 == b1 {
					c := Cell{
						Ch: braillePatterns[[2]int{m1, m0}],
						Bg: lc.Bg,
						Fg: thisLineColor,
					}
					y := lc.innerArea.Min.Y + lc.innerArea.Dy() - 3 - b0
					buf.Set(cellPos, y, c)
				} else {
					c0 := Cell{
						Ch: rSingleBraille[m0],
						Fg: thisLineColor,
						Bg: lc.Bg,
					}
					y0 := lc.innerArea.Min.Y + lc.innerArea.Dy() - 3 - b0
					buf.Set(cellPos, y0, c0)

					c1 := Cell{
						Ch: lSingleBraille[m1],
						Fg: thisLineColor,
						Bg: lc.Bg,
					}
					y1 := lc.innerArea.Min.Y + lc.innerArea.Dy() - 3 - b1
					buf.Set(cellPos, y1, c1)
				}
			} else {
				c0 := Cell{
					Ch: rSingleBraille[m0],
					Fg: thisLineColor,
					Bg: lc.Bg,
				}
				x0 := cellPos
				y0 := lc.innerArea.Min.Y + lc.innerArea.Dy() - 3 - b0
				buf.Set(x0, y0, c0)
			}
			dataPos -= 2
			cellPos--
		}
	}
	return buf
}

func (lc *LineChart) renderDot() Buffer {
	buf := NewBuffer()
	for seriesName, seriesData := range lc.Data {
		thisLineColor, ok := lc.LineColor[seriesName]
		if !ok {
			thisLineColor = lc.defaultLineColor
		}
		minCell := lc.innerArea.Min.X + lc.labelYSpace
		cellPos := lc.innerArea.Max.X - 1
		for dataPos := len(seriesData) - 1; dataPos >= 0 && cellPos > minCell; {
			c := Cell{
				Ch: lc.DotStyle,
				Fg: thisLineColor,
				Bg: lc.Bg,
			}
			x := cellPos
			y := lc.innerArea.Min.Y + lc.innerArea.Dy() - 3 - int((seriesData[dataPos]-lc.bottomValue)/lc.scale+0.5)
			buf.Set(x, y, c)

			cellPos--
			dataPos--
		}
	}

	return buf
}

func (lc *LineChart) calcLabelX() {
	lc.labelX = [][]rune{}

	for i, l := 0, 0; i < len(lc.DataLabels) && l < lc.axisXWidth; i++ {
		if lc.Mode == "dot" {
			if l >= len(lc.DataLabels) {
				break
			}

			s := str2runes(lc.DataLabels[l])
			w := strWidth(lc.DataLabels[l])
			if l+w <= lc.axisXWidth {
				lc.labelX = append(lc.labelX, s)
			}
			l += w + lc.axisXLabelGap
		} else { // braille
			if 2*l >= len(lc.DataLabels) {
				break
			}

			s := str2runes(lc.DataLabels[2*l])
			w := strWidth(lc.DataLabels[2*l])
			if l+w <= lc.axisXWidth {
				lc.labelX = append(lc.labelX, s)
			}
			l += w + lc.axisXLabelGap

		}
	}
}

func shortenFloatVal(x float64) string {
	s := fmt.Sprintf("%.2f", x)
	if len(s)-3 > 3 {
		s = fmt.Sprintf("%.2e", x)
	}

	if x < 0 {
		s = fmt.Sprintf("%.2f", x)
	}
	return s
}

func (lc *LineChart) calcLabelY() {
	span := lc.topValue - lc.bottomValue
	// where does -2 come from? Without it, we might draw on the top border or past the block
	lc.scale = span / float64(lc.axisYHeight-2)

	n := (1 + lc.axisYHeight) / (lc.axisYLabelGap + 1)
	lc.labelY = make([][]rune, n)
	maxLen := 0
	for i := 0; i < n; i++ {
		s := str2runes(shortenFloatVal(lc.bottomValue + float64(i)*span/float64(n)))
		if len(s) > maxLen {
			maxLen = len(s)
		}
		lc.labelY[i] = s
	}

	lc.labelYSpace = maxLen
}

func (lc *LineChart) calcLayout() {
	for _, seriesData := range lc.Data {
		if seriesData == nil || len(seriesData) == 0 {
			continue
		}
		// set datalabels if not provided
		if lc.DataLabels == nil || len(lc.DataLabels) == 0 {
			lc.DataLabels = make([]string, len(seriesData))
			for i := range seriesData {
				lc.DataLabels[i] = fmt.Sprint(i)
			}
		}

		// lazy increase, to avoid y shaking frequently
		lc.minY = seriesData[0]
		lc.maxY = seriesData[0]

		// valid visible range
		vrange := lc.innerArea.Dx()
		if lc.Mode == "braille" {
			vrange = 2 * lc.innerArea.Dx()
		}
		if vrange > len(seriesData) {
			vrange = len(seriesData)
		}

		for _, v := range seriesData[:vrange] {
			if v > lc.maxY {
				lc.maxY = v
			}
			if v < lc.minY {
				lc.minY = v
			}
		}

		span := lc.maxY - lc.minY

		// allow some padding unless we are beyond the flor/ceil
		if lc.minY <= lc.bottomValue {
			lc.bottomValue = lc.minY - lc.YPadding*span
			if lc.bottomValue < lc.YFloor {
				lc.bottomValue = lc.YFloor
			}
		}

		if lc.maxY >= lc.topValue {
			lc.topValue = lc.maxY + lc.YPadding*span
			if lc.topValue > lc.YCeil {
				lc.topValue = lc.YCeil
			}
		}
	}

	lc.axisYHeight = lc.innerArea.Dy() - 1
	lc.calcLabelY()

	lc.axisXWidth = lc.innerArea.Dx() - 1 - lc.labelYSpace
	lc.calcLabelX()

	lc.drawingX = lc.innerArea.Min.X + 1 + lc.labelYSpace
	lc.drawingY = lc.innerArea.Min.Y
}

func (lc *LineChart) plotAxes() Buffer {
	buf := NewBuffer()

	origY := lc.innerArea.Min.Y + lc.innerArea.Dy() - 2
	origX := lc.innerArea.Min.X + lc.labelYSpace

	buf.Set(origX, origY, Cell{Ch: ORIGIN, Fg: lc.AxesColor, Bg: lc.Bg})

	for x := origX + 1; x < origX+lc.axisXWidth; x++ {
		buf.Set(x, origY, Cell{Ch: HDASH, Fg: lc.AxesColor, Bg: lc.Bg})
	}

	for y := origY - 1; y > origY-lc.axisYHeight; y-- {
		buf.Set(origX, y, Cell{Ch: VDASH, Fg: lc.AxesColor, Bg: lc.Bg})
	}

	// x label
	oft := 0
	for _, rs := range lc.labelX {
		if oft+len(rs) > lc.axisXWidth {
			break
		}
		for j, r := range rs {
			c := Cell{
				Ch: r,
				Fg: lc.AxesColor,
				Bg: lc.Bg,
			}
			x := origX + oft + j
			y := lc.innerArea.Min.Y + lc.innerArea.Dy() - 1
			buf.Set(x, y, c)
		}
		oft += len(rs) + lc.axisXLabelGap
	}

	// y labels
	for i, rs := range lc.labelY {
		for j, r := range rs {
			buf.Set(
				lc.innerArea.Min.X+j,
				origY-i*(lc.axisYLabelGap+1),
				Cell{Ch: r, Fg: lc.AxesColor, Bg: lc.Bg})
		}
	}

	return buf
}

// Buffer implements Bufferer interface.
func (lc *LineChart) Buffer() Buffer {
	buf := lc.Block.Buffer()

	seriesCount := 0
	for _, data := range lc.Data {
		if len(data) > 0 {
			seriesCount++
		}
	}
	if seriesCount == 0 {
		return buf
	}
	lc.calcLayout()
	buf.Merge(lc.plotAxes())

	if lc.Mode == "dot" {
		buf.Merge(lc.renderDot())
	} else {
		buf.Merge(lc.renderBraille())
	}

	return buf
}
