// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"fmt"
	"image"

	rw "github.com/mattn/go-runewidth" // 必須匯入，因為 Draw 中使用了 rw.StringWidth

	. "github.com/gizak/termui/v3"
)

type HorizontalStackedBar struct {
	Block
	BarColors    []Color
	LabelStyles  []Style
	NumStyles    []Style // only Fg and Modifier are used
	NumFormatter func(float64) string
	Data         [][]float64
	Labels       []string
	BarHeight    int
	BarGap       int
	MaxVal       float64
}

func NewHorizontalStackedBar() *HorizontalStackedBar {
	return &HorizontalStackedBar{
		Block:        *NewBlock(),
		BarColors:    Theme.StackedBarChart.Bars,
		LabelStyles:  Theme.StackedBarChart.Labels,
		NumStyles:    Theme.StackedBarChart.Nums,
		NumFormatter: func(n float64) string { return fmt.Sprint(n) },
		BarGap:       2,
		BarHeight:    3,
	}
}

func (self *HorizontalStackedBar) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	maxVal := self.MaxVal
	if maxVal == 0 {
		for _, data := range self.Data {
			maxVal = MaxFloat64(maxVal, SumFloat64Slice(data))
		}
	}

	barYCoordinate := self.Inner.Min.Y+1

	labelWidth := 10

	for i, bar := range self.Data {
		stackX := labelWidth

		for j, data := range bar {
			width := int((data / maxVal) * float64(self.Inner.Dx()-labelWidth-1))

			for y := barYCoordinate; y < MinInt(barYCoordinate+self.BarHeight, self.Inner.Max.Y); y++ {
				for x := self.Inner.Min.X + stackX; x < MinInt(self.Inner.Min.X+stackX+width, self.Inner.Max.X); x++ {
					c := NewCell(' ', NewStyle(ColorClear, SelectColor(self.BarColors, j)))
					buf.SetCell(c, image.Pt(x, y))
				}
			}

			textX := self.Inner.Min.X + stackX + width/2 - rw.StringWidth(self.NumFormatter(data))/2
			textY := barYCoordinate+self.BarHeight
			buf.SetString(
				self.NumFormatter(data),
				NewStyle(
					SelectStyle(self.NumStyles, j+1).Fg,
					ColorClear,
					SelectStyle(self.NumStyles, j+1).Modifier,
				),
				image.Pt(textX, textY),
			)

			stackX += width
		}

		if i < len(self.Labels) {
			labelY := barYCoordinate+self.BarHeight/2
			buf.SetString(
				TrimString(self.Labels[i], labelWidth-2),
				SelectStyle(self.LabelStyles, i),
				image.Pt(self.Inner.Min.X+1, labelY),
			)
		}

		barYCoordinate += (self.BarHeight + self.BarGap)
	}
}
