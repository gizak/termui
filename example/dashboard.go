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

	draw := func(t int) {
		g.Percent = t % 101
		list.Items = strs[t%9:]
		ui.Render(p, list, g)
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
