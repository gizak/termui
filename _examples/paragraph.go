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

	p0 := widgets.NewParagraph()
	p0.Text = "Borderless Text"
	p0.SetRect(0, 0, 20, 5)
	p0.Border = false

	p1 := widgets.NewParagraph()
	p1.Title = "标签"
	p1.Text = "你好，世界。"
	p1.SetRect(20, 0, 35, 5)

	p2 := widgets.NewParagraph()
	p2.Title = "Multiline"
	p2.Text = "Simple colored text\nwith label. It [can be](fg:red) multilined with \\n or [break automatically](fg:red,fg:bold)"
	p2.SetRect(0, 5, 35, 10)
	p2.BorderStyle.Fg = ui.ColorYellow

	p3 := widgets.NewParagraph()
	p3.Title = "Auto Trim"
	p3.Text = "Long text with label and it is auto trimmed."
	p3.SetRect(0, 10, 40, 15)

	p4 := widgets.NewParagraph()
	p4.Title = "Text Box with Wrapping"
	p4.Text = "Press q to QUIT THE DEMO. [There](fg:blue,mod:bold) are other things [that](fg:red) are going to fit in here I think. What do you think? Now is the time for all good [men to](bg:blue) come to the aid of their country. [This is going to be one really really really long line](fg:green) that is going to go together and stuffs and things. Let's see how this thing renders out.\n    Here is a new paragraph and stuffs and things. There should be a tab indent at the beginning of the paragraph. Let's see if that worked as well."
	p4.SetRect(40, 0, 70, 20)
	p4.BorderStyle.Fg = ui.ColorBlue

	ui.Render(p0, p1, p2, p3, p4)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
