// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import "image"

// Cell is a rune with assigned Fg and Bg
type Cell struct {
	Ch rune
	Fg Attribute
	Bg Attribute
}

// Buffer is a renderable rectangle cell data container.
type Buffer struct {
	Area    image.Rectangle // selected drawing area
	CellMap map[image.Point]Cell
}

// At returns the cell at (x,y).
func (b Buffer) At(x, y int) Cell {
	return b.CellMap[image.Pt(x, y)]
}

// Set assigns a char to (x,y)
func (b Buffer) Set(x, y int, c Cell) {
	b.CellMap[image.Pt(x, y)] = c
}

// Bounds returns the domain for which At can return non-zero color.
func (b Buffer) Bounds() image.Rectangle {
	x0, y0, x1, y1 := 0, 0, 0, 0
	for p := range b.CellMap {
		switch {
		case p.X > x1:
			x1 = p.X
		case p.X < x0:
			x0 = p.X
		case p.Y > y1:
			y1 = p.Y
		case p.Y < y0:
			y0 = p.Y
		}
	}
	return image.Rect(x0, y0, x1, y1)
}

// Align sets drawing area to the buffer's bound
func (b *Buffer) Align() {
	b.Area = b.Bounds()
}

// NewCell returns a new cell
func NewCell(ch rune, fg, bg Attribute) Cell {
	return Cell{ch, fg, bg}
}

// Union squeezes buf into b
func (b Buffer) Union(buf Buffer) {
	for p, c := range buf.CellMap {
		b.Set(p.X, p.Y, c)
	}
}

// Union returns a new Buffer formed by squeezing bufs into one Buffer
func Union(bufs ...Buffer) Buffer {
	buf := NewBuffer()
	for _, b := range bufs {
		buf.Union(b)
	}
	buf.Align()
	return buf
}

// Point for adapting use, will be removed after resolving bridging.
type Point struct {
	X  int
	Y  int
	Ch rune
	Fg Attribute
	Bg Attribute
}

// NewBuffer returns a new Buffer
func NewBuffer() Buffer {
	return Buffer{CellMap: make(map[image.Point]Cell)}
}
