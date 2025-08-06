package main

import (
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	spinner := widgets.NewSpinner()
	spinner.Label = "loading..."
	spinner.SetRect(5, 5, 22, 8)

	spinnerRight := widgets.NewSpinner()
	spinnerRight.FormatString = "[%s] %s"
	spinnerRight.Label = "processing"
	spinnerRight.LabelOnRight = true
	spinnerRight.SetRect(5, 10, 22, 13)

	spinnerAlone := widgets.NewSpinner()
	spinnerAlone.FormatString = "%s%s"
	spinnerAlone.SetRect(5, 15, 8, 18)

	ui.Render(spinner)
	ui.Render(spinnerRight)
	ui.Render(spinnerAlone)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	uiEvents := ui.PollEvents()
	for {
		select {
		case <-ticker.C:
			spinner.Advance()
			ui.Render(spinner)
			spinnerRight.Advance()
			ui.Render(spinnerRight)
			spinnerAlone.Advance()
			ui.Render(spinnerAlone)
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		}
	}
}
