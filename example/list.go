package main

import "github.com/gizak/termui"
import "github.com/nsf/termbox-go"

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	termui.UseTheme("helloworld")

	strs := []string{
		"[0] github.com/gizak/termui",
		"[1] editbox.go",
		"[2] iterrupt.go",
		"[3] keyboard.go",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] nsf/termbox-go"}

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.Border.Label = "List"
	ls.Height = 7
	ls.Width = 25
	ls.Y = 0

	termui.Render(ls)

	termbox.PollEvent()
}
