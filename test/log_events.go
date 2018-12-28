// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package main

import (
	"fmt"

	ui "github.com/gizak/termui"
)

// logs all events to the termui window
// stdout can also be redirected to a file and read with `tail -f`
func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	for {
		e := <-ui.PollEvents()
		fmt.Printf("%v", e)
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
