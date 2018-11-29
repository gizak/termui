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

	sbc := ui.NewStackedBarChart()
	math := []int{90, 85, 90, 80}
	english := []int{70, 85, 75, 60}
	science := []int{75, 60, 80, 85}
	compsci := []int{100, 100, 100, 100}
	sbc.Data[0] = math
	sbc.Data[1] = english
	sbc.Data[2] = science
	sbc.Data[3] = compsci
	studentsName := []string{"Ken", "Rob", "Dennis", "Linus"}
	sbc.BorderLabel = "Student's Marks X-Axis=Name Y-Axis=Marks[Math,English,Science,ComputerScience] in %"
	sbc.Width = 100
	sbc.Height = 30
	sbc.Y = 0
	sbc.BarWidth = 10
	sbc.DataLabels = studentsName
	sbc.ShowScale = true //Show y_axis scale value (min and max)
	sbc.SetMax(400)

	sbc.TextColor = ui.ColorGreen    //this is color for label (x-axis)
	sbc.BarColor[3] = ui.ColorGreen  //BarColor for computerscience
	sbc.BarColor[1] = ui.ColorYellow //Bar Color for english
	sbc.NumColor[3] = ui.ColorRed    // Num color for computerscience
	sbc.NumColor[1] = ui.ColorRed    // num color for english

	//Other colors are automatically populated, btw All the students seems do well in computerscience. :p

	ui.Render(sbc)

	for {
		e := <-ui.PollEvent()
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
