// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"fmt"
	"image"

	. "github.com/gizak/termui"
)

// LineChart has two modes: braille(default) and dot.
// A single braille character is a 2x4 grid of dots, so using braille
// gives 2x X resolution and 4x Y resolution over dot mode.
type LineChart struct {
	Block
	Data            [][]float64
	DataLabels      []string
	HorizontalScale int
	LineType        LineType
	DotChar         rune
	LineColors      []Color
	AxesColor       Color // TODO
	MaxVal          float64
	ShowAxes        bool
	DrawDirection   DrawDirection // TODO
}

const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 4
	xAxisLabelsGap    = 2
	yAxisLabelsGap    = 1
)

type LineType int

const (
	BrailleLine LineType = iota
	DotLine
)

type DrawDirection uint

const (
	DrawLeft DrawDirection = iota
	DrawRight
)

func NewLineChart() *LineChart {
	return &LineChart{
		Block:           *NewBlock(),
		LineColors:      Theme.LineChart.Lines,
		AxesColor:       Theme.LineChart.Axes,
		LineType:        BrailleLine,
		DotChar:         DOT,
		Data:            [][]float64{},
		HorizontalScale: 1,
		DrawDirection:   DrawRight,
		ShowAxes:        true,
	}
}

func (self *LineChart) renderBraille(buf *Buffer, drawArea image.Rectangle, maxVal float64) {
	canvas := NewCanvas()
	canvas.Rectangle = drawArea

	for i, line := range self.Data {
		previousHeight := int((line[1] / maxVal) * float64(drawArea.Dy()-1))
		for j, val := range line[1:] {
			height := int((val / maxVal) * float64(drawArea.Dy()-1))
			canvas.Line(
				image.Pt(
					(drawArea.Min.X+(j*self.HorizontalScale))*2,
					(drawArea.Max.Y-previousHeight-1)*4,
				),
				image.Pt(
					(drawArea.Min.X+((j+1)*self.HorizontalScale))*2,
					(drawArea.Max.Y-height-1)*4,
				),
				SelectColor(self.LineColors, i),
			)
			previousHeight = height
		}
	}

	canvas.Draw(buf)
}

func (self *LineChart) renderDot(buf *Buffer, drawArea image.Rectangle, maxVal float64) {
	for i, line := range self.Data {
		for j := 0; j < len(line) && j*self.HorizontalScale < drawArea.Dx(); j++ {
			val := line[j]
			height := int((val / maxVal) * float64(drawArea.Dy()-1))
			buf.SetCell(
				NewCell(self.DotChar, NewStyle(SelectColor(self.LineColors, i))),
				image.Pt(drawArea.Min.X+(j*self.HorizontalScale), drawArea.Max.Y-1-height),
			)
		}
	}
}

func (self *LineChart) plotAxes(buf *Buffer, maxVal float64) {
	// draw origin cell
	buf.SetCell(
		NewCell(BOTTOM_LEFT, NewStyle(ColorWhite)),
		image.Pt(self.Inner.Min.X+yAxisLabelsWidth, self.Inner.Max.Y-xAxisLabelsHeight-1),
	)
	// draw x axis line
	for i := yAxisLabelsWidth + 1; i < self.Inner.Dx(); i++ {
		buf.SetCell(
			NewCell(HORIZONTAL_DASH, NewStyle(ColorWhite)),
			image.Pt(i+self.Inner.Min.X, self.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}
	// draw y axis line
	for i := 0; i < self.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			NewCell(VERTICAL_DASH, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Min.X+yAxisLabelsWidth, i+self.Inner.Min.Y),
		)
	}
	// draw x axis labels
	// draw 0
	buf.SetString(
		"0",
		NewStyle(ColorWhite),
		image.Pt(self.Inner.Min.X+yAxisLabelsWidth, self.Inner.Max.Y-1),
	)
	// draw rest
	for x := self.Inner.Min.X + yAxisLabelsWidth + (xAxisLabelsGap)*self.HorizontalScale + 1; x < self.Inner.Max.X-1; {
		label := fmt.Sprintf(
			"%d",
			(x-(self.Inner.Min.X+yAxisLabelsWidth)-1)/(self.HorizontalScale)+1,
		)
		buf.SetString(
			label,
			NewStyle(ColorWhite),
			image.Pt(x, self.Inner.Max.Y-1),
		)
		x += (len(label) + xAxisLabelsGap) * self.HorizontalScale
	}
	// draw y axis labels
	verticalScale := maxVal / float64(self.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 0; i*(yAxisLabelsGap+1) < self.Inner.Dy()-1; i++ {
		buf.SetString(
			fmt.Sprintf("%.2f", float64(i)*verticalScale*(yAxisLabelsGap+1)),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Min.X, self.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2),
		)
	}
}

func (self *LineChart) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	maxVal := self.MaxVal
	if maxVal == 0 {
		maxVal, _ = GetMaxFloat64From2dSlice(self.Data)
	}

	if self.ShowAxes {
		self.plotAxes(buf, maxVal)
	}

	drawArea := self.Inner
	if self.ShowAxes {
		drawArea = image.Rect(
			self.Inner.Min.X+yAxisLabelsWidth+1, self.Inner.Min.Y,
			self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
		)
	}

	if self.LineType == BrailleLine {
		self.renderBraille(buf, drawArea, maxVal)
	} else {
		self.renderDot(buf, drawArea, maxVal)
	}
}
