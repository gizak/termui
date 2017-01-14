// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package main

import "github.com/gizak/termui"

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()
	rows_1 := [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"孙嘉你好吗", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "10", "11"},
	}

	table_1 := termui.NewTable()
	table_1.Rows = rows_1
	table_1.FgColor = termui.ColorWhite
	table_1.BgColor = termui.ColorDefault
	table_1.Y = 0
	table_1.X = 0
	table_1.Width = 62
	table_1.Height = 7

	termui.Render(table_1)

	rows := [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"Foundations", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "11", "11"},
	}

	table := termui.NewTable()
	table.Rows = rows
	table.FgColor = termui.ColorWhite
	table.BgColor = termui.ColorDefault
	table.TextAlign = termui.AlignCenter
	table.Separator = false
	table.Analysis()
	table.SetSize()
	table.BgColors[2] = termui.ColorRed
	table.Y = 20
	table.X = 0
	table.Border = true

	termui.Render(table)
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()
}
