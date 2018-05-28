// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package main

import ui "github.com/gizak/termui"

func main() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewPar("Press q to QUIT THE DEMO. [There](fg-blue) are other things [that](fg-red) are going to fit in here I think. What do you think? Now is the time for all good [men to](bg-blue) come to the aid of their country. [This is going to be one really really really long line](fg-green) that is going to go together and stuffs and things. Let's see how this thing renders out.\n    Here is a new paragraph and stuffs and things. There should be a tab indent at the beginning of the paragraph. Let's see if that worked as well.")
	p.WrapLength = 48 // this should be at least p.Width - 2
	p.Height = 20
	p.Width = 50
	p.Y = 2
	p.X = 20
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = "Text Box with Wrapping"
	p.BorderFg = ui.ColorCyan
	//p.Border = false

	ui.Render(p)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()
}
