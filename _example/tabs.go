// Copyright 2016 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/gizak/termui"
	"github.com/gizak/termui/extra"
)

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	termui.UseTheme("helloworld")

	header := termui.NewPar("Press q to quit, Press j or k to switch tabs")
	header.Height = 1
	header.Width = 50
	header.HasBorder = false
	header.TextBgColor = termui.ColorBlue

	tab1 := extra.NewTab("pierwszy")
	par2 := termui.NewPar("Press q to quit\nPress j or k to switch tabs\n")
	par2.Height = 5
	par2.Width = 37
	par2.Y = 0
	par2.Border.Label = "Keys"
	par2.Border.FgColor = termui.ColorYellow
	tab1.AddBlocks(par2)

	tab2 := extra.NewTab("drugi")
	bc := termui.NewBarChart()
	data := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	bclabels := []string{"S0", "S1", "S2", "S3", "S4", "S5"}
	bc.Border.Label = "Bar Chart"
	bc.Data = data
	bc.Width = 26
	bc.Height = 10
	bc.DataLabels = bclabels
	bc.TextColor = termui.ColorGreen
	bc.BarColor = termui.ColorRed
	bc.NumColor = termui.ColorYellow
	tab2.AddBlocks(bc)

	tab3 := extra.NewTab("trzeci")
	tab4 := extra.NewTab("żółw")
	tab5 := extra.NewTab("four")
	tab6 := extra.NewTab("five")

	tabpane := extra.NewTabpane()
	tabpane.Y = 1
	tabpane.Width = 30
	tabpane.HasBorder = true

	tabpane.SetTabs(*tab1, *tab2, *tab3, *tab4, *tab5, *tab6)

	evt := termui.EventCh()

	termui.Render(header, tabpane)

	for {
		select {
		case e := <-evt:
			if e.Type == termui.EventKey {
				switch e.Ch {
				case 'q':
					return
				case 'j':
					tabpane.SetActiveLeft()
					termui.Render(header, tabpane)
				case 'k':
					tabpane.SetActiveRight()
					termui.Render(header, tabpane)
				}
			}
		}
	}
}
