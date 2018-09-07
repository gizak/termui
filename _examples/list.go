// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import ui "github.com/gizak/termui"

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	strs := []string{
		"[0] github.com/gizak/termui",
		"[1] [你好，世界](fg-blue)",
		"[2] [こんにちは世界](fg-red)",
		"[3] [color output](fg-white,bg-green)",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] nsf/termbox-go"}

	ls := ui.NewList()
	ls.Items = strs
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = 7
	ls.Width = 25
	ls.Y = 0

	ui.Render(ls)

	ui.Handle("q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()
}
