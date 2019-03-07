// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
)

// logs all events to the termui window
// stdout can also be redirected to a file and read with `tail -f`
func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	events := ui.PollEvents()
	for {
		e := <-events
		fmt.Printf("%v", e)
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<MouseLeft>":
			return
		}
	}
}
