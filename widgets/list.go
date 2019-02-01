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
	Rows             []string
	WrapText         bool
	TextStyle        Style
	SelectedRow      uint
	topRow           uint
	SelectedRowStyle Style
}

func NewList() *List {
	return &List{
		Block:            *NewBlock(),
		TextStyle:        Theme.List.Text,
		SelectedRowStyle: Theme.List.Text,
	}
}

func (self *List) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	point := self.Inner.Min

	if self.SelectedRow >= uint(self.Inner.Max.Y)+self.topRow-2 {
		self.topRow = self.SelectedRow - uint(self.Inner.Max.Y) + 2
	} else if self.SelectedRow < self.topRow {
		self.topRow = self.SelectedRow
	}

	for row := self.topRow; row < uint(len(self.Rows)) && point.Y < self.Inner.Max.Y; row++ {
		cells := ParseText(self.Rows[row], self.TextStyle)
		if self.WrapText {
			cells = WrapCells(cells, uint(self.Inner.Dx()))
		}
		for j := 0; j < len(cells) && point.Y < self.Inner.Max.Y; j++ {
			style := cells[j].Style
			if row == self.SelectedRow {
				style = self.SelectedRowStyle
			}
			if cells[j].Rune == '\n' {
				point = image.Pt(self.Inner.Min.X, point.Y+1)
			} else {
				if point.X+1 == self.Inner.Max.X+1 && len(cells) > self.Inner.Dx() {
					buf.SetCell(NewCell(ELLIPSES, style), point.Add(image.Pt(-1, 0)))
					break
				} else {
					buf.SetCell(NewCell(cells[j].Rune, style), point)
					point = point.Add(image.Pt(1, 0))
				}
			}
		}
		point = image.Pt(self.Inner.Min.X, point.Y+1)
	}

	if self.topRow > 0 {
		buf.SetCell(
			NewCell(UP_ARROW, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Max.X-1, self.Inner.Min.Y),
		)
	}
	if len(self.Rows) > int(self.topRow)+self.Inner.Dy() {
		buf.SetCell(
			NewCell(DOWN_ARROW, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Max.X-1, self.Inner.Max.Y-1),
		)
	}
}

func (self *List) ScrollUp() {
	if self.SelectedRow > 0 {
		self.SelectedRow--
		if self.SelectedRow < self.topRow {
			self.topRow--
		}
	}
}

func (self *List) ScrollDown() {
	if self.SelectedRow < uint(len(self.Rows))-1 {
		self.SelectedRow++
		if self.SelectedRow-self.topRow > uint(self.Inner.Dy()-1) {
			self.topRow++
		}
	}
}

// PageUp scrolls up one whole page.
func (self *List) PageUp() {
	// if on the first 'page'
	if int(self.SelectedRow)-self.Inner.Dy() < 0 {
		// go to the top
		self.topRow = 0
	} else {
		self.topRow = uint(MaxInt(int(self.topRow)-self.Inner.Dy(), 0))
	}
	self.SelectedRow = self.topRow
}

// PageDown scolls down one whole page.
func (self *List) PageDown() {
	// if on last 'page'
	if len(self.Rows)-int(self.topRow) <= self.Inner.Dy() {
		// select last item
		self.SelectedRow = uint(len(self.Rows) - 1)
	} else {
		self.topRow += uint(self.Inner.Dy())
		self.SelectedRow = self.topRow
	}
}
