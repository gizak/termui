// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import "image"

type Border struct {
	Area   image.Rectangle
	Left   bool
	Top    bool
	Right  bool
	Bottom bool
	Fg     Attribute
	Bg     Attribute
}

type Hline struct {
	X   int
	Y   int
	Len int
	Fg  Attribute
	Bg  Attribute
}

type Vline struct {
	X   int
	Y   int
	Len int
	Fg  Attribute
	Bg  Attribute
}

// Buffer draws a horizontal line.
func (l Hline) Buffer() Buffer {
	buf := NewBuffer()
	for i := 0; i < l.Len; i++ {
		buf.Set(l.X+i, l.Y, Cell{HORIZONTAL_LINE, l.Fg, l.Bg})
	}
	buf.Align()
	return buf
}

// Buffer draws a vertical line.
func (l Vline) Buffer() Buffer {
	buf := NewBuffer()
	for i := 0; i < l.Len; i++ {
		buf.Set(l.X, l.Y+i, Cell{VERTICAL_LINE, l.Fg, l.Bg})
	}
	buf.Align()
	return buf
}

// Buffer draws a box border.
func (b Border) Buffer() Buffer {
	buf := NewBuffer()
	if b.Area.Size().X < 2 || b.Area.Size().Y < 2 {
		return buf
	}

	min := b.Area.Min
	max := b.Area.Max

	x0 := min.X
	y0 := min.Y
	x1 := max.X
	y1 := max.Y

	// draw lines
	switch {
	case b.Top:
		buf.Union(Hline{x0, y0, x1 - x0, b.Fg, b.Bg}.Buffer())
		fallthrough
	case b.Bottom:
		buf.Union(Hline{x0, y1, x1 - x0, b.Fg, b.Bg}.Buffer())
		fallthrough
	case b.Left:
		buf.Union(Vline{x0, y0, y1 - y0, b.Fg, b.Bg}.Buffer())
		fallthrough
	case b.Right:
		buf.Union(Vline{x1, y0, y1 - y0, b.Fg, b.Bg}.Buffer())
	}

	// draw corners
	switch {
	case b.Top && b.Left:
		buf.Set(x0, y0, Cell{TOP_LEFT, b.Fg, b.Bg})
		fallthrough
	case b.Top && b.Right:
		buf.Set(x1, y0, Cell{TOP_RIGHT, b.Fg, b.Bg})
		fallthrough
	case b.Bottom && b.Left:
		buf.Set(x0, y1, Cell{BOTTOM_LEFT, b.Fg, b.Bg})
		fallthrough
	case b.Bottom && b.Right:
		buf.Set(x1, y1, Cell{BOTTOM_RIGHT, b.Fg, b.Bg})
	}

	return buf
}

// LabeledBorder defined label upon Border
type LabeledBorder struct {
	Border
	Label      string
	LabelFgClr Attribute
	LabelBgClr Attribute
}

// Buffer draw a box border with label.
func (lb LabeledBorder) Buffer() Buffer {
	border := lb.Border.Buffer()
	maxTxtW := lb.Area.Dx() + 1 - 2
	tx := DTrimTxCls(TextCells(lb.Label, lb.LabelFgClr, lb.LabelBgClr), maxTxtW)

	for i, w := 0, 0; i < len(tx); i++ {
		border.Set(border.Area.Min.X+1+w, border.Area.Min.Y, tx[i])
		w += tx[i].Width()
	}

	return border
}
