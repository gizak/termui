// +build ignore

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	ui "github.com/gizak/termui"
)

func main() {
	ui.Init()
	defer ui.Close()

	rand.Seed(time.Now().UTC().UnixNano())
	randomDataAndOffset := func() (data []float64, offset float64) {
		noSlices := 1 + rand.Intn(5)
		data = make([]float64, noSlices)
		for i := range data {
			data[i] = rand.Float64()
		}
		offset = 2.0 * math.Pi * rand.Float64()
		return
	}
	run := true

	pc := ui.NewPieChart()
	pc.BorderLabel = "Pie Chart"
	pc.Width = 70
	pc.Height = 36
	pc.Data = []float64{.25, .25, .25, .25}
	pc.Offset = -.5 * math.Pi
	pc.Label = func(i int, v float64) string {
		return fmt.Sprintf("%.02f", v)
	}

	ui.Handle("/timer/1s", func(e ui.Event) {
		if run {
			pc.Data, pc.Offset = randomDataAndOffset()
			ui.Render(pc)
		}
	})

	ui.Handle("/sys/kbd/s", func(ui.Event) {
		run = !run
		if run {
			pc.BorderLabel = "Pie Chart"
		} else {
			pc.BorderLabel = "Pie Chart (Stopped)"
		}
		ui.Render(pc)
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Render(pc)

	ui.Loop()
}
