package main

import (
	"github.com/gizak/termui"
	//"fmt"
	//"os"
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

    tab1 := termui.NewTab("pierwszy")
	par2 := termui.NewPar("Press q to quit\nPress j or k to switch tabs\n")
	par2.Height = 5
	par2.Width = 37
	par2.Y = 0
	par2.Border.Label = "Keys"
	par2.Border.FgColor = termui.ColorYellow
	tab1.AddBlocks(par2)

    tab2 := termui.NewTab("drugi")
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

    tab3 := termui.NewTab("trzeci")
    tab4 := termui.NewTab("żółw")
    tab5 := termui.NewTab("four")
    tab6 := termui.NewTab("five")

    tabpane := termui.NewTabpane()
	tabpane.Y = 1
    tabpane.Width = 30
    tabpane.HasBorder = true

    tabpane.SetTabs(*tab1, *tab2, *tab3, *tab4, *tab5, *tab6)

    termui.Render(header, tabpane)

	evt := termui.EventCh()
	for {
		select {
		case e := <- evt:
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
