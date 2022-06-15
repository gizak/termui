package main

import (
	"context"
	"fmt"
	"log"
	"os"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	f0 := widgets.NewForm(
		context.Background(),
		"Add user",
		func(ctx context.Context, validated bool, fields []*widgets.Field) {
			ui.Close()
			fmt.Printf("Form validated: %t. Values:\n", validated)
			for _, f := range fields {
				fmt.Printf("%s: %s\n", f.Name, f.Text)
			}
			os.Exit(0)
		},
		widgets.NewField("First name", ""),
		widgets.NewField("Surname", ""),
		widgets.NewField("Test", "Default value"),
	)
	f0.SetRect(0, 0, 50, 20)

	p2 := widgets.NewParagraph()
	p2.Title = "Multiline"
	p2.Text = "Simple colored text\nwith label. It [can be](fg:red) multilined with \\n or [break automatically](fg:red,fg:bold)"
	p2.SetRect(0, 50, 35, 55)
	p2.BorderStyle.Fg = ui.ColorYellow

	ui.Render(f0, p2)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			return
		default:
			if !f0.IsDone() {
				f0.Handle(e)
			}
		}
	}
}
