// +build ignore

package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var run = true

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	rand.Seed(time.Now().UTC().UnixNano())
	randomDataAndOffset := func() (data []float64, offset float64) {
		noSlices := 1 + rand.Intn(5)
		data = make([]float64, noSlices)
		for i := range data {
			data[i] = rand.Float64()*100.0
		}
		offset = 2.0 * math.Pi * rand.Float64()
		return
	}

	pc := widgets.NewPieChart()
	pc.Title = "Pie Chart"
	pc.SetRect(5, 5, 70, 36)
	pc.Data = []float64{67.67, 11.45, 19.19, 81.00}
	pc.AngleOffset = -.5 * math.Pi
	pc.LabelFormatter = func(i int, v float64, p float64) string {
		return fmt.Sprintf("%.2f (%.0f%%)", v, p*100)
	}

	pause := func() {
		run = !run
		if run {
			pc.Title = "Pie Chart"
		} else {
			pc.Title = "Pie Chart (Stopped)"
		}
		ui.Render(pc)
	}

	ui.Render(pc)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "s":
				pause()
			}
		case <-ticker:
			if run {
				pc.Data, pc.AngleOffset = randomDataAndOffset()
				baseColors := []ui.Color{
					ui.ColorRed,
					ui.ColorGreen,
					ui.ColorYellow,
					ui.ColorBlue,
					ui.ColorMagenta,
					ui.ColorCyan,
				}

				perm := rand.Perm(len(baseColors))
				pc.Colors = make([]ui.Color, len(pc.Data))
				for i := range pc.Colors {
					pc.Colors[i] = baseColors[perm[i]]
				}
				
				ui.Render(pc)
			}
		}
	}
}
