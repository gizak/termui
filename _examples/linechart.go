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

	sinps := (func() map[string][]float64 {
		n := 220
		ps := make(map[string][]float64)
		ps["first"] = make([]float64, n)
		ps["second"] = make([]float64, n)
		for i := 0; i < n; i++ {
			ps["first"][i] = 1 + math.Sin(float64(i)/5)
			ps["second"][i] = 1 + math.Cos(float64(i)/5)
		}
		return ps
	})()

	lc0 := ui.NewLineChart()
	lc0.BorderLabel = "braille-mode Line Chart"
	lc0.Data = sinps
	lc0.Width = 50
	lc0.Height = 12
	lc0.X = 0
	lc0.Y = 0
	lc0.AxesColor = ui.ColorWhite
	lc0.LineColor["first"] = ui.ColorGreen | ui.AttrBold

	lc1 := ui.NewLineChart()
	lc1.BorderLabel = "dot-mode Line Chart"
	lc1.Mode = "dot"
	lc1.Data = sinps
	lc1.Width = 26
	lc1.Height = 12
	lc1.X = 51
	lc1.DotStyle = '+'
	lc1.AxesColor = ui.ColorWhite
	lc1.LineColor["first"] = ui.ColorYellow | ui.AttrBold

	lc2 := ui.NewLineChart()
	lc2.BorderLabel = "dot-mode Line Chart"
	lc2.Mode = "dot"
	lc2.Data["first"] = sinps["first"][4:]
	lc2.Data["second"] = sinps["second"][4:]
	lc2.Width = 77
	lc2.Height = 16
	lc2.X = 0
	lc2.Y = 12
	lc2.AxesColor = ui.ColorWhite
	lc2.LineColor["first"] = ui.ColorCyan | ui.AttrBold

	ui.Render(lc0, lc1, lc2)

	ui.Handle("q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()
}
