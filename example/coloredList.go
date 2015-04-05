// +build ignore

package main

import "github.com/gizak/termui"
import "github.com/nsf/termbox-go"

func commonList() *termui.List {
	strs := []string{
		"[0] github.com/gizak/termui",
		"[1] 笀耔 [澉 灊灅甗](RED) 郔镺 笀耔 澉 [灊灅甗](yellow) 郔镺",
		"[2] こんにちは世界",
		"[3] keyboard.go",
		"[4] [output](RED).go",
		"[5] random_out.go",
		"[6] [dashboard](BOLD).go",
		"[7] nsf/termbox-go",
		"[8] OVERFLOW!!!!!!![!!!!!!!!!!!!](red,bold)!!!"}

	list := termui.NewList()
	list.Items = strs
	list.Height = 20
	list.Width = 25
	list.RendererFactory = termui.MarkdownTextRendererFactory{}

	return list
}

func listHidden() *termui.List {
	list := commonList()
	list.Border.Label = "List - Hidden"
	list.Overflow = "hidden"

	return list
}

func listWrap() *termui.List {
	list := commonList()
	list.Border.Label = "List - Wrapped"
	list.Overflow = "wrap"

	return list
}

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	hiddenList := listHidden()
	wrappedList := listWrap()
	wrappedList.X = 30

	termui.UseTheme("helloworld")
	termui.Render(hiddenList, wrappedList)
	termbox.PollEvent()
}
