// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import ui "github.com/gizak/termui"

func main() {
	ui.Init()
	defer ui.Close()

	g0 := ui.NewGauge()
	g0.Percent = 40
	g0.Width = 50
	g0.Height = 3
	g0.BorderLabel = "Slim Gauge"
	g0.BarColor = ui.ColorRed
	g0.BorderFg = ui.ColorWhite
	g0.BorderLabelFg = ui.ColorCyan

	gg := ui.NewBlock()
	gg.Width = 50
	gg.Height = 5
	gg.Y = 12
	gg.BorderLabel = "TEST"
	gg.Align()

	g2 := ui.NewGauge()
	g2.Percent = 60
	g2.Width = 50
	g2.Height = 3
	g2.PercentColor = ui.ColorBlue
	g2.Y = 3
	g2.BorderLabel = "Slim Gauge"
	g2.BarColor = ui.ColorYellow
	g2.BorderFg = ui.ColorWhite

	g1 := ui.NewGauge()
	g1.Percent = 30
	g1.Width = 50
	g1.Height = 5
	g1.Y = 6
	g1.BorderLabel = "Big Gauge"
	g1.PercentColor = ui.ColorYellow
	g1.BarColor = ui.ColorGreen
	g1.BorderFg = ui.ColorWhite
	g1.BorderLabelFg = ui.ColorMagenta

	g3 := ui.NewGauge()
	g3.Percent = 50
	g3.Width = 50
	g3.Height = 3
	g3.Y = 11
	g3.BorderLabel = "Gauge with custom label"
	g3.Label = "{{percent}}% (100MBs free)"
	g3.LabelAlign = ui.AlignRight

	g4 := ui.NewGauge()
	g4.Percent = 50
	g4.Width = 50
	g4.Height = 3
	g4.Y = 14
	g4.BorderLabel = "Gauge"
	g4.Label = "Gauge with custom highlighted label"
	g4.PercentColor = ui.ColorYellow
	g4.BarColor = ui.ColorGreen
	g4.PercentColorHighlighted = ui.ColorBlack

	ui.Render(g0, g1, g2, g3, g4)

	ui.Handle("q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()
}
