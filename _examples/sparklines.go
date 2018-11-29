// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import ui "github.com/gizak/termui"

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	data := []int{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}
	spl0 := ui.NewSparkline()
	spl0.Data = data[3:]
	spl0.Title = "Sparkline 0"
	spl0.LineColor = ui.ColorGreen

	// single
	spls0 := ui.NewSparklines(spl0)
	spls0.Height = 2
	spls0.Width = 20
	spls0.Border = false

	spl1 := ui.NewSparkline()
	spl1.Data = data
	spl1.Title = "Sparkline 1"
	spl1.LineColor = ui.ColorRed

	spl2 := ui.NewSparkline()
	spl2.Data = data[5:]
	spl2.Title = "Sparkline 2"
	spl2.LineColor = ui.ColorMagenta

	// group
	spls1 := ui.NewSparklines(spl0, spl1, spl2)
	spls1.Height = 8
	spls1.Width = 20
	spls1.Y = 3
	spls1.BorderLabel = "Group Sparklines"

	spl3 := ui.NewSparkline()
	spl3.Data = data
	spl3.Title = "Enlarged Sparkline"
	spl3.Height = 8
	spl3.LineColor = ui.ColorYellow

	spls2 := ui.NewSparklines(spl3)
	spls2.Height = 11
	spls2.Width = 30
	spls2.BorderFg = ui.ColorCyan
	spls2.X = 21
	spls2.BorderLabel = "Tweeked Sparkline"

	ui.Render(spls0, spls1, spls2)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
