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

	hsbc := widgets.NewHorizontalStackedBar()
	hsbc.Title = "Student's Marks: X-Axis=Grade% (Math, English, Science, Computer Science), Y-Axis=Name"
	hsbc.Labels = []string{"Ken", "Rob", "Dennis", "Linus"}
	
	hsbc.Data = make([][]float64, 4)
	hsbc.Data[0] = []float64{90, 85, 90, 80}
	hsbc.Data[1] = []float64{70, 85, 75, 60}
	hsbc.Data[2] = []float64{75, 60, 80, 85}
	hsbc.Data[3] = []float64{100, 100, 100, 100}
	hsbc.SetRect(5, 5, 100, 30)
	hsbc.BarHeight = 2

	ui.Render(hsbc)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
