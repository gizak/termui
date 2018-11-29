// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"math"
	"time"

	ui "github.com/gizak/termui"
)

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	sinps := (func() []float64 {
		n := 400
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5)
		}
		return ps
	})()
	sinpsint := (func() []int {
		ps := make([]int, len(sinps))
		for i, v := range sinps {
			ps[i] = int(100*v + 10)
		}
		return ps
	})()

	spark := ui.Sparkline{}
	spark.Height = 8
	spdata := sinpsint
	spark.Data = spdata[:100]
	spark.LineColor = ui.ColorCyan
	spark.TitleColor = ui.ColorWhite

	sp := ui.NewSparklines(spark)
	sp.Height = 11
	sp.BorderLabel = "Sparkline"

	lc := ui.NewLineChart()
	lc.BorderLabel = "braille-mode Line Chart"
	lc.Data["default"] = sinps
	lc.Height = 11
	lc.AxesColor = ui.ColorWhite
	lc.LineColor["default"] = ui.ColorYellow | ui.AttrBold

	gs := make([]*ui.Gauge, 3)
	for i := range gs {
		gs[i] = ui.NewGauge()
		//gs[i].LabelAlign = ui.AlignCenter
		gs[i].Height = 2
		gs[i].Border = false
		gs[i].Percent = i * 10
		gs[i].PaddingBottom = 1
		gs[i].BarColor = ui.ColorRed
	}

	ls := ui.NewList()
	ls.Border = false
	ls.Items = []string{
		"[1] Downloading File 1",
		"", // == \newline
		"[2] Downloading File 2",
		"",
		"[3] Uploading File 3",
	}
	ls.Height = 5

	par := ui.NewPar("<> This row has 3 columns\n<- Widgets can be stacked up like left side\n<- Stacked widgets are treated as a single widget")
	par.Height = 5
	par.BorderLabel = "Demonstration"

	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, sp),
			ui.NewCol(6, 0, lc)),
		ui.NewRow(
			ui.NewCol(3, 0, ls),
			ui.NewCol(3, 0, gs[0], gs[1], gs[2]),
			ui.NewCol(6, 0, par)))

	// calculate layout
	ui.Body.Align()

	ui.Render(ui.Body)

	tickerCount := 1
	for {
		select {
		case e := <-ui.PollEvent():
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				ui.Body.Width = payload.Width
				ui.Body.Align()
				ui.Clear()
				ui.Render(ui.Body)
			}
		case <-time.NewTicker(time.Second).C:
			if tickerCount > 103 {
				return
			}
			for _, g := range gs {
				g.Percent = (g.Percent + 3) % 100
			}
			sp.Lines[0].Data = spdata[:100+tickerCount]
			lc.Data["default"] = sinps[2*tickerCount:]
			ui.Render(ui.Body)
			tickerCount++
		}
	}
}
