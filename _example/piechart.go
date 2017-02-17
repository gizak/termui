// +build ignore

package main

import (
	"fmt"
	"github.com/gizak/termui"
	"math"
	"math/rand"
	"time"
)

func main() {
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()
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

	pc := termui.NewPieChart()
	pc.BorderLabel = "Pie Chart"
	pc.Width = 70
	pc.Height = 36
	pc.Data = []float64{.25, .25, .25, .25}
	pc.Offset = -.5 * math.Pi
	pc.Label = func(i int, v float64) string {
		return fmt.Sprintf("%.02f", v)
	}

	termui.Handle("/timer/1s", func(e termui.Event) {
		if run {
			pc.Data, pc.Offset = randomDataAndOffset()
			termui.Render(pc)
		}
	})

	termui.Handle("/sys/kbd/s", func(termui.Event) {
		run = !run
		if run {
			pc.BorderLabel = "Pie Chart"
		} else {
			pc.BorderLabel = "Pie Chart (Stopped)"
		}
		termui.Render(pc)
	})

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Render(pc)
	termui.Loop()
}
