// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import "image"

// Block is a base struct for all other upper level widgets,
// consider it as css: display:block.
// Normally you do not need to create it manually.
type Block struct {
	Area          image.Rectangle
	innerArea     image.Rectangle
	X             int
	Y             int
	Border        LabeledBorder
	IsDisplay     bool
	HasBorder     bool
	Bg            Attribute
	Width         int
	Height        int
	PaddingTop    int
	PaddingBottom int
	PaddingLeft   int
	PaddingRight  int
}

// NewBlock returns a *Block which inherits styles from current theme.
func NewBlock() *Block {
	d := Block{}
	d.IsDisplay = true
	d.HasBorder = theme.HasBorder
	d.Border.Left = true
	d.Border.Right = true
	d.Border.Top = true
	d.Border.Bottom = true
	d.Border.Bg = theme.BorderBg
	d.Border.Fg = theme.BorderFg
	d.Border.LabelBgClr = theme.BorderLabelTextBg
	d.Border.LabelFgClr = theme.BorderLabelTextFg
	d.Bg = theme.BlockBg
	d.Width = 2
	d.Height = 2
	return &d
}

// Align computes box model
func (d *Block) Align() {
	d.Area.Min.X = d.X
	d.Area.Min.Y = d.Y
	d.Area.Max.X = d.X + d.Width - 1
	d.Area.Max.Y = d.Y + d.Height - 1

	d.innerArea.Min.X = d.X + d.PaddingLeft
	d.innerArea.Min.Y = d.Y + d.PaddingTop
	d.innerArea.Max.X = d.Area.Max.X - d.PaddingRight
	d.innerArea.Max.Y = d.Area.Max.Y - d.PaddingBottom

	d.Border.Area = d.Area

	if d.HasBorder {
		switch {
		case d.Border.Left:
			d.innerArea.Min.X++
			fallthrough
		case d.Border.Right:
			d.innerArea.Max.X--
			fallthrough
		case d.Border.Top:
			d.innerArea.Min.Y++
			fallthrough
		case d.Border.Bottom:
			d.innerArea.Max.Y--
		}
	}
}

// InnerBounds returns the internal bounds of the block after aligning and
// calculating the padding and border, if any.
func (d *Block) InnerBounds() image.Rectangle {
	d.Align()
	return d.innerArea
}

// Buffer implements Bufferer interface.
// Draw background and border (if any).
func (d *Block) Buffer() Buffer {
	d.Align()

	buf := NewBuffer()
	buf.Area = d.Area
	if !d.IsDisplay {
		return buf
	}

	// render border
	if d.HasBorder {
		buf.Union(d.Border.Buffer())
	}

	// render background
	for p := range buf.CellMap {
		if p.In(d.innerArea) {
			buf.CellMap[p] = Cell{' ', ColorDefault, d.Bg}
		}
	}
	return buf
}

// GetHeight implements GridBufferer.
// It returns current height of the block.
func (d Block) GetHeight() int {
	return d.Height
}

// SetX implements GridBufferer interface, which sets block's x position.
func (d *Block) SetX(x int) {
	d.X = x
}

// SetY implements GridBufferer interface, it sets y position for block.
func (d *Block) SetY(y int) {
	d.Y = y
}

// SetWidth implements GridBuffer interface, it sets block's width.
func (d *Block) SetWidth(w int) {
	d.Width = w
}
