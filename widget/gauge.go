// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widget

import "github.com/gizak/termui"
import "strconv"

// Gauge is a progress bar like widget.
// A simple example:
/*
  g := termui.NewGauge()
  g.Percent = 40
  g.Width = 50
  g.Height = 3
  g.Border.Label = "Slim Gauge"
  g.BarColor = termui.ColorRed
  g.PercentColor = termui.ColorBlue
*/
type Gauge struct {
	termui.Block
	Percent      int
	BarColor     termui.Attribute
	PercentColor termui.Attribute
}

// NewGauge return a new gauge with current theme.
func NewGauge() *Gauge {
	g := &Gauge{
		Block:        *termui.NewBlock(),
		PercentColor: termui.Theme().GaugePercent,
		BarColor:     termui.Theme().GaugeBar}
	g.Width = 12
	g.Height = 3
	return g
}

// Buffer implements Bufferer interface.
func (g *Gauge) Buffer() termui.Buffer {
	buf := g.Block.Buffer()

	inner := g.InnerBounds()
	w := g.Percent * (inner.Dx() + 1) / 100
	s := strconv.Itoa(g.Percent) + "%"
	tx := termui.TextCells(s, g.PercentColor, g.Bg)

	prx := inner.Min.X + (inner.Dx()+1)/2 - 1
	pry := inner.Min.Y + (inner.Dy()+1)/2

	// plot bar
	for i := 0; i <= inner.Dy(); i++ {
		for j := 0; j < w; j++ {
			c := termui.Cell{' ', g.BarColor, g.BarColor}
			buf.Set(inner.Min.X+j, inner.Min.Y+i, c)
		}
	}

	// plot percentage
	for i, v := range tx {
		if w > (inner.Dx()+1)/2-1+i {
			v.Bg = g.BarColor
		}
		buf.Set(prx+i, pry, v)
	}
	return buf
}
