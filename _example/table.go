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
	rows1 := [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"你好吗", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "10", "11"},
	}

	table1 := termui.NewTable()
	table1.Rows = rows1
	table1.FgColor = termui.ColorWhite
	table1.BgColor = termui.ColorDefault
	table1.Y = 0
	table1.X = 0
	table1.Width = 62
	table1.Height = 7

	termui.Render(table1)

	rows2 := [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"Foundations", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "11", "11"},
	}

	table2 := termui.NewTable()
	table2.Rows = rows2
	table2.FgColor = termui.ColorWhite
	table2.BgColor = termui.ColorDefault
	table2.TextAlign = termui.AlignCenter
	table2.Separator = false
	table2.Analysis()
	table2.SetSize()
	table2.BgColors[2] = termui.ColorRed
	table2.Y = 10
	table2.X = 0
	table2.Border = true

	termui.Render(table2)
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()
}
