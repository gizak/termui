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

	par0 := ui.NewPar("Borderless Text")
	par0.Height = 1
	par0.Width = 20
	par0.Y = 1
	par0.Border = false

	par1 := ui.NewPar("你好，世界。")
	par1.Height = 3
	par1.Width = 17
	par1.X = 20
	par1.BorderLabel = "标签"

	par2 := ui.NewPar("Simple colored text\nwith label. It [can be](fg-red) multilined with \\n or [break automatically](fg-red,fg-bold)")
	par2.Height = 5
	par2.Width = 37
	par2.Y = 4
	par2.BorderLabel = "Multiline"
	par2.BorderFg = ui.ColorYellow

	par3 := ui.NewPar("Long text with label and it is auto trimmed.")
	par3.Height = 3
	par3.Width = 37
	par3.Y = 9
	par3.BorderLabel = "Auto Trim"

	ui.Render(par0, par1, par2, par3)

	ui.Handle("q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()

}
