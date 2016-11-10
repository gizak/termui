package main

import "github.com/gizak/termui"

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	rows := [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"I love Veronica", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "11", "11"},
	}

	table := termui.NewTable()
	table.Rows = rows
	table.FgColor = termui.ColorWhite
	table.BgColor = termui.ColorDefault
	table.TextAlign = "center"
	table.Analysis()
	table.SetSize()
	table.Y = 0
	table.X = 0
	table.Border = true

	termui.Render(table)
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()
}
