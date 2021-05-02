// Copyright 2021 Szymon Błaszczyński <museyoucoulduse@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"image"

	. "github.com/gizak/termui/v3"
)

type Input struct {
	Block
	Text      string
	TextStyle Style
	WrapText  bool
	focusing  bool
}

// Create new Input widget. Somehow similar to Paragraph.
func NewInput() *Input {
	return &Input{
		Block:     *NewBlock(),
		TextStyle: Theme.Input.Text,
		WrapText:  true,
		focusing:  false,
		Text:      "",
	}
}

func (self *Input) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	cells := ParseStyles(self.Text, self.TextStyle)
	if self.WrapText {
		cells = WrapCells(cells, uint(self.Inner.Dx()))
	}

	rows := SplitCells(cells, '\n')

	for y, row := range rows {
		if y+self.Inner.Min.Y >= self.Inner.Max.Y {
			break
		}
		row = TrimCells(row, self.Inner.Dx())
		for _, cx := range BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(self.Inner.Min))
		}
	}
}

// Focus on input field and start typing.
// Best used with go routine 'go inputField.Focus()'
func (i *Input) Focus() {
	i.focusing = true
Loop:
	for {
		events := PollEvents()
		e := <-events
		if e.Type == KeyboardEvent {
			switch e.ID {
			case "<Backspace>":
				if len(i.Text) > 0 {
					i.Text = i.Text[:len(i.Text)-1]
					Render(i)
				}
			case "<Escape>", "C-c", "<Enter>":
				i.focusing = false
				break Loop
			case "<Space>":
				i.Text = i.Text + " "
				Render(i)
			case "<Tab>":
				i.Text = i.Text + "\t"
				Render(i)
			default:
				i.Text = i.Text + e.ID
				Render(i)
			}
		}
	}
}
