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

	p0 := ui.NewParagraph("Borderless Text")
	p0.Height = 1
	p0.Width = 20
	p0.Y = 1
	p0.Border = false

	p1 := ui.NewParagraph("你好，世界。")
	p1.Height = 3
	p1.Width = 17
	p1.X = 20
	p1.BorderLabel = "标签"

	p2 := ui.NewParagraph("Simple colored text\nwith label. It [can be](fg-red) multilined with \\n or [break automatically](fg-red,fg-bold)")
	p2.Height = 5
	p2.Width = 37
	p2.Y = 4
	p2.BorderLabel = "Multiline"
	p2.BorderFg = ui.ColorYellow

	p3 := ui.NewParagraph("Long text with label and it is auto trimmed.")
	p3.Height = 3
	p3.Width = 37
	p3.Y = 9
	p3.BorderLabel = "Auto Trim"

	ui.Render(p0, p1, p2, p3)

	for {
		e := <-ui.PollEvent()
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
