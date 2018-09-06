// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	ui "github.com/gizak/termui"
)

func main() {
	ui.Init()
	defer ui.Close()

	p := ui.NewPar(":PRESS q TO QUIT DEMO")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = "Text Box"
	p.BorderFg = ui.ColorCyan
	p.Handle("/timer/1s", func(e ui.Event) {
		cnt := e.Data.(ui.EvtTimer)
		if cnt.Count%2 == 0 {
			p.TextFgColor = ui.ColorRed
		} else {
			p.TextFgColor = ui.ColorWhite
		}
	})

	listData := []string{"[0] gizak/termui", "[1] editbox.go", "[2] interrupt.go", "[3] keyboard.go", "[4] output.go", "[5] random_out.go", "[6] dashboard.go", "[7] nsf/termbox-go"}

	l := ui.NewList()
	l.Items = listData
	l.ItemFgColor = ui.ColorYellow
	l.BorderLabel = "List"
	l.Height = 7
	l.Width = 25
	l.Y = 4

	g := ui.NewGauge()
	g.Percent = 50
	g.Width = 50
	g.Height = 3
	g.Y = 11
	g.BorderLabel = "Gauge"
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan

	sparklineData := []int{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}

	sl := ui.Sparkline{}
	sl.Height = 1
	sl.Title = "srv 0:"
	sl.Data = sparklineData
	sl.LineColor = ui.ColorCyan
	sl.TitleColor = ui.ColorWhite

	sl2 := ui.Sparkline{}
	sl2.Height = 1
	sl2.Title = "srv 1:"
	sl2.Data = sparklineData
	sl2.TitleColor = ui.ColorWhite
	sl2.LineColor = ui.ColorRed

	sls := ui.NewSparklines(sl, sl2)
	sls.Width = 25
	sls.Height = 7
	sls.BorderLabel = "Sparkline"
	sls.Y = 4
	sls.X = 25

	sinData := (func() []float64 {
		n := 220
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5)
		}
		return ps
	})()

	lc := ui.NewLineChart()
	lc.BorderLabel = "dot-mode Line Chart"
	lc.Data["default"] = sinData
	lc.Width = 50
	lc.Height = 11
	lc.X = 0
	lc.Y = 14
	lc.AxesColor = ui.ColorWhite
	lc.LineColor["default"] = ui.ColorRed | ui.AttrBold
	lc.Mode = "dot"

	barchartData := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}

	bc := ui.NewBarChart()
	bc.BorderLabel = "Bar Chart"
	bc.Width = 26
	bc.Height = 10
	bc.X = 51
	bc.Y = 0
	bc.DataLabels = []string{"S0", "S1", "S2", "S3", "S4", "S5"}
	bc.BarColor = ui.ColorGreen
	bc.NumColor = ui.ColorBlack

	lc2 := ui.NewLineChart()
	lc2.BorderLabel = "braille-mode Line Chart"
	lc2.Data["default"] = sinData
	lc2.Width = 26
	lc2.Height = 11
	lc2.X = 51
	lc2.Y = 14
	lc2.AxesColor = ui.ColorWhite
	lc2.LineColor["default"] = ui.ColorYellow | ui.AttrBold

	p2 := ui.NewPar("Hey!\nI am a borderless block!")
	p2.Border = false
	p2.Width = 26
	p2.Height = 2
	p2.TextFgColor = ui.ColorMagenta
	p2.X = 52
	p2.Y = 11

	draw := func(count int) {
		g.Percent = count % 101
		l.Items = listData[count%9:]
		sls.Lines[0].Data = sparklineData[:30+count%50]
		sls.Lines[1].Data = sparklineData[:35+count%50]
		lc.Data["default"] = sinData[count/2%220:]
		lc2.Data["default"] = sinData[2*count%220:]
		bc.Data = barchartData[count/2%10:]

		ui.Render(p, l, g, sls, lc, bc, lc2, p2)
	}

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		draw(int(t.Count))
	})

	ui.Loop()
}
