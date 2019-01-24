// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"log"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table1 := widgets.NewTable()
	table1.Rows = [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"你好吗", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "10", "11"},
	}
	table1.TextStyle = ui.NewStyle(ui.ColorWhite)
	table1.SetRect(0, 0, 60, 10)

	ui.Render(table1)

	table2 := widgets.NewTable()
	table2.Rows = [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"Foundations", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "11", "11"},
	}
	table2.TextStyle = ui.NewStyle(ui.ColorWhite)
	table2.TextAlign = ui.AlignCenter
	table2.RowSeparator = false
	table2.SetRect(0, 10, 20, 20)

	ui.Render(table2)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
