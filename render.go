// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"fmt"
	"image"
	"sync"

	tb "github.com/nsf/termbox-go"
)

type Drawable interface {
	GetRect() image.Rectangle
	SetRect(int, int, int, int)
	Draw(*Buffer)
	sync.Locker
	GetANSIString() string
}

func Render(items ...Drawable) {
	// draw background, etc for items with ANSI escape strings
	for _, item := range items {
		if len(item.GetANSIString()) > 0 {
			continue
		}
		buf := NewBuffer(item.GetRect())
		item.Lock()
		item.Draw(buf)
		item.Unlock()
		for point, cell := range buf.CellMap {
			if point.In(buf.Rectangle) {
				tb.SetCell(
					point.X, point.Y,
					cell.Rune,
					tb.Attribute(cell.Style.Fg+1)|tb.Attribute(cell.Style.Modifier), tb.Attribute(cell.Style.Bg+1),
				)
			}
		}
	}
	tb.Flush()

	// draw images, etc over the already filled cells with ANSI escape strings (sixel, ...)
	for _, item := range items {
		if ansiString := item.GetANSIString(); len(ansiString) > 0 {
			fmt.Printf("%s", ansiString)
			continue
		}
	}

	// draw items without ANSI strings last in case the ANSI escape strings ended messed up
	for _, item := range items {
		if len(item.GetANSIString()) == 0 {
			continue
		}
		buf := NewBuffer(item.GetRect())
		item.Lock()
		item.Draw(buf)
		item.Unlock()
		for point, cell := range buf.CellMap {
			if point.In(buf.Rectangle) {
				tb.SetCell(
					point.X, point.Y,
					cell.Rune,
					tb.Attribute(cell.Style.Fg+1)|tb.Attribute(cell.Style.Modifier), tb.Attribute(cell.Style.Bg+1),
				)
			}
		}
	}
	tb.Flush()
}
