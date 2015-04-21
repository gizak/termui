// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import "github.com/gizak/termui"
import "github.com/gizak/termui/widget"

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	termui.UseTheme("helloworld")

	g0 := widget.NewGauge()
	g0.Percent = 40
	g0.Width = 50
	g0.Height = 3
	g0.Border.Label = "Slim Gauge"
	g0.BarColor = termui.ColorRed
	g0.Border.Fg = termui.ColorWhite
	g0.Border.LabelFgClr = termui.ColorCyan

	gg := termui.NewBlock()
	gg.Width = 50
	gg.Height = 5
	gg.Y = 12
	gg.Border.Label = "TEST"
	gg.Align()

	g2 := widget.NewGauge()
	g2.Percent = 60
	g2.Width = 50
	g2.Height = 3
	g2.PercentColor = termui.ColorBlue
	g2.Y = 3
	g2.Border.Label = "Slim Gauge"
	g2.BarColor = termui.ColorYellow
	g2.Border.Fg = termui.ColorWhite

	g1 := widget.NewGauge()
	g1.Percent = 30
	g1.Width = 50
	g1.Height = 5
	g1.Y = 6
	g1.Border.Label = "Big Gauge"
	g1.PercentColor = termui.ColorYellow
	g1.BarColor = termui.ColorGreen
	g1.Border.Fg = termui.ColorWhite
	g1.Border.LabelFgClr = termui.ColorMagenta

	termui.Render(g0, g1, g2, gg)

	<-termui.EventCh()
}
