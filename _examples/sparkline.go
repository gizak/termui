// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	data := []float64{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}

	sl0 := widgets.NewSparkline()
	sl0.Data = data[3:]
	sl0.LineColor = ui.ColorGreen

	// single
	slg0 := widgets.NewSparklineGroup(sl0)
	slg0.Title = "Sparkline 0"
	slg0.SetRect(0, 0, 20, 10)

	sl1 := widgets.NewSparkline()
	sl1.Title = "Sparkline 1"
	sl1.Data = data
	sl1.LineColor = ui.ColorRed

	sl2 := widgets.NewSparkline()
	sl2.Title = "Sparkline 2"
	sl2.Data = data[5:]
	sl2.LineColor = ui.ColorMagenta

	slg1 := widgets.NewSparklineGroup(sl0, sl1, sl2)
	slg1.Title = "Group Sparklines"
	slg1.SetRect(0, 10, 25, 25)

	sl3 := widgets.NewSparkline()
	sl3.Title = "Enlarged Sparkline"
	sl3.Data = data
	sl3.LineColor = ui.ColorYellow

	slg2 := widgets.NewSparklineGroup(sl3)
	slg2.Title = "Tweeked Sparkline"
	slg2.SetRect(20, 0, 50, 10)
	slg2.BorderStyle.Fg = ui.ColorCyan

	ui.Render(slg0, slg1, slg2)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
