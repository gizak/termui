// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
//
// Portions copyright 2017 Patrick Devine <patrick@immense.ly>

// +build ignore

package main

import ui "github.com/gizak/termui"

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	strs := []ui.Item{
		ui.Item{"github.com/gizak/termui", "github.com/gizak/termui"},
		ui.Item{"你好，世界", "[你好，世界](fg-blue)"},
		ui.Item{"こんにちは世界", "[こんにちは世界](fg-red)"},
		ui.Item{"color output", "[color output](fg-white,bg-green)"},
		ui.Item{"output.go", "output.go"},
		ui.Item{"random_out.go", "random_out.go"},
		ui.Item{"dashboard.go", "dashboard.go"},
		ui.Item{"nsf/termbox-go", "nsf/termbox-go"},
	}

	ls := ui.NewListBox()
	ls.Items = strs
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = 7
	ls.Width = 20
	ls.Y = 0

	ui.Render(ls)
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		ls.Up()
		ui.Render(ls)
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		ls.Down()
		ui.Render(ls)
	})
	ui.Loop()

}
