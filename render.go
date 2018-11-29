// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"
	"runtime/debug"
	"sync"

	tb "github.com/nsf/termbox-go"
)

// Bufferer should be implemented by all renderable components.
type Bufferer interface {
	Buffer() Buffer
}

// Init initializes termui library. This function should be called before any others.
// After initialization, the library must be finalized by 'Close' function.
func Init() error {
	if err := tb.Init(); err != nil {
		return err
	}
	tb.SetInputMode(tb.InputEsc | tb.InputMouse)
	// DefaultEvtStream = NewEvtStream()

	// sysEvtChs = make([]chan Event, 0)
	// go hookTermboxEvt()

	renderJobs = make(chan []Bufferer)
	//renderLock = new(sync.RWMutex)

	Body = NewGrid()
	Body.X = 0
	Body.Y = 0
	Body.BgColor = ThemeAttr("bg")
	Body.Width = TermWidth()

	// resizeCh := Handle("<Resize>")
	// go func() {
	// 	for e := range resizeCh {
	// 		payload := e.Payload.(Resize)
	// 		Body.Width = payload.Width
	// 	}
	// }()

	// DefaultWgtMgr = NewWgtMgr()
	// EventHook(DefaultWgtMgr.WgtHandlersHook())

	go func() {
		for bs := range renderJobs {
			render(bs...)
		}
	}()

	return nil
}

// Close finalizes termui library,
// should be called after successful initialization when termui's functionality isn't required anymore.
func Close() {
	tb.Close()
}

var renderLock sync.Mutex

func termSync() {
	renderLock.Lock()
	tb.Sync()
	termWidth, termHeight = tb.Size()
	renderLock.Unlock()
}

// TermWidth returns the current terminal's width.
func TermWidth() int {
	termSync()
	return termWidth
}

// TermHeight returns the current terminal's height.
func TermHeight() int {
	termSync()
	return termHeight
}

// Render renders all Bufferer in the given order from left to right,
// right could overlap on left ones.
func render(bs ...Bufferer) {
	defer func() {
		if e := recover(); e != nil {
			Close()
			panic(debug.Stack())
		}
	}()
	for _, b := range bs {

		buf := b.Buffer()
		// set cels in buf
		for p, c := range buf.CellMap {
			if p.In(buf.Area) {

				tb.SetCell(p.X, p.Y, c.Ch, toTmAttr(c.Fg), toTmAttr(c.Bg))

			}
		}

	}

	renderLock.Lock()
	// render
	tb.Flush()
	renderLock.Unlock()
}

func Clear() {
	tb.Clear(tb.ColorDefault, toTmAttr(ThemeAttr("bg")))
}

func clearArea(r image.Rectangle, bg Attribute) {
	for i := r.Min.X; i < r.Max.X; i++ {
		for j := r.Min.Y; j < r.Max.Y; j++ {
			tb.SetCell(i, j, ' ', tb.ColorDefault, toTmAttr(bg))
		}
	}
}

func ClearArea(r image.Rectangle, bg Attribute) {
	clearArea(r, bg)
	tb.Flush()
}

var renderJobs chan []Bufferer

func Render(bs ...Bufferer) {
	//go func() { renderJobs <- bs }()
	renderJobs <- bs
}
