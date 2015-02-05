package main

import ui "github.com/gizak/termui"
import tm "github.com/nsf/termbox-go"
import "time"

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewP(":PRESS q TO QUIT DEMO\nThis is an example of termui package rendering.")
	p.Height = 4
	p.Width = 59
	p.TextFgColor = ui.ColorWhite
	p.Border.Label = "Text"
	p.Border.FgColor = ui.ColorCyan

	strs := []string{"[0] gizak/termui", "[1] editbox.go", "[2] iterrupt.go", "[3] keyboard.go", "[4] output.go", "[5] random_out.go", "[6] dashboard.go", "[7] nsf/termbox-go"}
	list := ui.NewList()
	list.Items = strs
	list.ItemFgColor = ui.ColorYellow
	list.Border.Label = "List"
	list.Height = 7
	list.Width = 25
	list.Y = 4

	g := ui.NewGauge()
	g.Percent = 50
	g.Width = 52
	g.Height = 3
	g.Y = 11
	g.Border.Label = "Gauge"
	g.BarColor = ui.ColorRed
	g.Border.FgColor = ui.ColorWhite
	g.Border.LabelFgColor = ui.ColorCyan

	spark := ui.Sparkline{}
	spark.Height = 1
	spark.Title = "srv 0:"
	spdata := []int{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}
	spark.Data = spdata
	spark.LineColor = ui.ColorCyan
	spark.TitleColor = ui.ColorWhite

	spark1 := ui.Sparkline{}
	spark1.Height = 1
	spark1.Title = "srv 1:"
	spark1.Data = spdata
	spark1.TitleColor = ui.ColorWhite
	spark1.LineColor = ui.ColorRed

	sp := ui.NewSparklines(spark, spark1)
	sp.Width = 20
	sp.Height = 6
	sp.Border.Label = "Sparkline"
	sp.Y = 14

	draw := func(t int) {
		g.Percent = t % 101
		list.Items = strs[t%9:]
		sp.Lines[0].Data = spdata[t%10:]
		ui.Render(p, list, g, sp)
	}

	evt := make(chan tm.Event)
	go func() {
		for {
			evt <- tm.PollEvent()
		}
	}()

	i := 0
	for {
		select {
		case e := <-evt:
			if e.Type == tm.EventKey && e.Ch == 'q' {
				return
			}
		default:
			draw(i)
			i++
			if i == 102 {
				return
			}
			time.Sleep(time.Second / 2)
		}
	}
}
