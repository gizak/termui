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

	bc := ui.NewMBarChart()
	math := []int{90, 85, 90, 80}
	english := []int{70, 85, 75, 60}
	science := []int{75, 60, 80, 85}
	compsci := []int{100, 100, 100, 100}
	bc.Data[0] = math
	bc.Data[1] = english
	bc.Data[2] = science
	bc.Data[3] = compsci
	studentsName := []string{"Ken", "Rob", "Dennis", "Linus"}
	bc.BorderLabel = "Student's Marks X-Axis=Name Y-Axis=Marks[Math,English,Science,ComputerScience] in %"
	bc.Width = 100
	bc.Height = 30
	bc.Y = 0
	bc.BarWidth = 10
	bc.DataLabels = studentsName
	bc.ShowScale = true //Show y_axis scale value (min and max)
	bc.SetMax(400)

	bc.TextColor = ui.ColorGreen    //this is color for label (x-axis)
	bc.BarColor[3] = ui.ColorGreen  //BarColor for computerscience
	bc.BarColor[1] = ui.ColorYellow //Bar Color for english
	bc.NumColor[3] = ui.ColorRed    // Num color for computerscience
	bc.NumColor[1] = ui.ColorRed    // num color for english

	//Other colors are automatically populated, btw All the students seems do well in computerscience. :p

	ui.Render(bc)

	ui.Handle("q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()

}
