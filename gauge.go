// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

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
	Block
	Percent      int
	BarColor     Attribute
	PercentColor Attribute
}

// NewGauge return a new gauge with current theme.
func NewGauge() *Gauge {
	g := &Gauge{
		Block:        *NewBlock(),
		PercentColor: theme.GaugePercent,
		BarColor:     theme.GaugeBar}
	g.Width = 12
	g.Height = 5
	return g
}

// Buffer implements Bufferer interface.
func (g *Gauge) Buffer() []Point {
	ps := g.Block.Buffer()

	w := g.Percent * g.innerWidth / 100
	s := strconv.Itoa(g.Percent) + "%"
	rs := str2runes(s)

	prx := g.innerX + g.innerWidth/2 - 1
	pry := g.innerY + g.innerHeight/2

	// plot bar
	for i := 0; i < g.innerHeight; i++ {
		for j := 0; j < w; j++ {
			p := Point{}
			p.X = g.innerX + j
			p.Y = g.innerY + i
			p.Ch = ' '
			p.Bg = g.BarColor
			if p.Bg == ColorDefault {
				p.Bg |= AttrReverse
			}
			ps = append(ps, p)
		}
	}

	// plot percentage
	for i, v := range rs {
		p := Point{}
		p.X = prx + i
		p.Y = pry
		p.Ch = v
		p.Fg = g.PercentColor
		if w > g.innerWidth/2-1+i {
			p.Bg = g.BarColor
			if p.Bg == ColorDefault {
				p.Bg |= AttrReverse
			}

		} else {
			p.Bg = g.Block.BgColor
		}
		ps = append(ps, p)
	}
	return ps
}
