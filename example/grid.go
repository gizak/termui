package main

/*
import ui "github.com/gizak/termui"
import tm "github.com/nsf/termbox-go"
import "math"

import "time"

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

	ui.UseTheme("helloworld")

	spark := ui.Sparkline{}
	spark.Height = 8
	spdata := sinpsint
	spark.Data = spdata
	spark.LineColor = ui.ColorCyan
	spark.TitleColor = ui.ColorWhite

	sp := ui.NewSparklines(spark)
	sp.Height = 11
	sp.Border.Label = "Sparkline"

	lc := ui.NewLineChart()
	lc.Border.Label = "braille-mode Line Chart"
	lc.Data = sinps
	lc.Height = 11
	lc.AxesColor = ui.ColorWhite
	lc.LineColor = ui.ColorYellow | ui.AttrBold

	ui.Body.Rows = []ui.Row{
		ui.NewRow(
			ui.NewCol(sp, 6, 0, true),
			ui.NewCol(lc, 6, 0, true))}

	draw := func(t int) {
		sp.Lines[0].Data = spdata[t:]
		lc.Data = sinps[2*t:]
		ui.Render(ui.Body)
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
*/
