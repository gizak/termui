// +build ignore

package main

import "github.com/gizak/termui"
import "github.com/nsf/termbox-go"

func markdownList() *termui.List {
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
	list.Height = 15
	list.Width = 26
	list.RendererFactory = termui.MarkdownTextRendererFactory{}

	return list
}

func hideList(list *termui.List) *termui.List {
	list.Border.Label = "List - Hidden"
	list.Overflow = "hidden"

	return list
}

func wrapList(list *termui.List) *termui.List {
	list.Border.Label = "List - Wrapped"
	list.Overflow = "wrap"
	list.X = 30

	return list
}

func escapeList() *termui.List {
	strs := []string{
		"[0] github.com/gizak/termui",
		"[1] 笀耔 \033[31m澉 灊灅甗 \033[0m郔镺 笀耔 澉 \033[33m灊灅甗 郔镺",
		"[2] こんにちは世界",
		"[3] keyboard.go",
		"[4] \033[31moutput\033[0m.go",
		"[5] random_out.go",
		"[6] \033[1mdashboard\033[0m.go",
		"[7] nsf/termbox-go",
		"[8] OVERFLOW!!!!!!!\033[31;1m!!!!!!!!!!!!\033[0m!!!",
	}

	list := termui.NewList()
	list.RendererFactory = termui.EscapeCodeRendererFactory{}
	list.Items = strs
	list.Height = 15
	list.Width = 26
	list.Y = 15

	return list
}

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	hiddenMarkdownList := hideList(markdownList())
	wrappedMarkdownList := wrapList(markdownList())

	hiddenEscapeList := hideList(escapeList())
	wrappedEscapeList := wrapList(escapeList())

	lists := []termui.Bufferer{
		hiddenEscapeList,
		hiddenMarkdownList,
		wrappedMarkdownList,
		wrappedEscapeList,
	}

	termui.UseTheme("helloworld")
	termui.Render(lists...)
	termbox.PollEvent()
}
