// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"image"

	. "github.com/gizak/termui"
)

type List struct {
	Block
	Rows      []string
	Wrap      bool
	TextStyle Style
}

func NewList() *List {
	return &List{
		Block:     *NewBlock(),
		TextStyle: Theme.List.Text,
	}
}

func (self *List) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	point := self.Inner.Min

	for row := 0; row < len(self.Rows) && point.Y < self.Inner.Max.Y; row++ {
		cells := ParseText(self.Rows[row], self.TextStyle)
		if self.Wrap {
			cells = WrapCells(cells, uint(self.Inner.Dx()))
		}
		for j := 0; j < len(cells) && point.Y < self.Inner.Max.Y; j++ {
			if cells[j].Rune == '\n' {
				point = image.Pt(self.Inner.Min.X, point.Y+1)
			} else {
				if point.X+1 == self.Inner.Max.X+1 && len(cells) > self.Inner.Dx() {
					buf.SetCell(NewCell(ELLIPSES, cells[j].Style), point.Add(image.Pt(-1, 0)))
					break
				} else {
					buf.SetCell(cells[j], point)
					point = point.Add(image.Pt(1, 0))
				}
			}
		}
		point = image.Pt(self.Inner.Min.X, point.Y+1)
	}
}
