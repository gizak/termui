// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "List"
	l.Rows = []string{
		"[0] github.com/gizak/termui/v3",
		"[1] [你好，世界](fg:blue)",
		"[2] [こんにちは世界](fg:red)",
		"[3] [color](fg:white,bg:green) output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
	}
	l.ActiveBorderStyle = ui.NewStyle(ui.ColorRed)
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 8)
	l.Active = true

	l2 := widgets.NewList()
	l2.Title = "List 2"
	l2.Rows = []string{
		"foo",
		"bar",
		"baz",
	}
	l2.ActiveBorderStyle = ui.NewStyle(ui.ColorRed)
	l2.TextStyle = ui.NewStyle(ui.ColorYellow)
	l2.WrapText = false
	l2.SetRect(28, 0, 53, 8)

	uiWidgets := []*widgets.List{l, l2}

	for _, w := range(uiWidgets) {
		ui.Render(w)
	}

	var activeWidget *widgets.List

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		for _, w := range(uiWidgets) {
			if w.Active {
				activeWidget = w
				break
			}
		}

		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			activeWidget.ScrollDown()
		case "k", "<Up>":
			activeWidget.ScrollUp()
		case "<C-d>":
			activeWidget.ScrollHalfPageDown()
		case "<C-u>":
			activeWidget.ScrollHalfPageUp()
		case "<C-f>":
			activeWidget.ScrollPageDown()
		case "<C-b>":
			activeWidget.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				activeWidget.ScrollTop()
			}
		case "<Home>":
			activeWidget.ScrollTop()
		case "G", "<End>":
			activeWidget.ScrollBottom()
		case "<Tab>":
			for _, w := range(uiWidgets) {
				w.Active = !w.Active
			}
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		for _, w := range(uiWidgets) {
			ui.Render(w)
		}
	}
}
