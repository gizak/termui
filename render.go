// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import tm "github.com/nsf/termbox-go"

// all renderable components should implement this
type Bufferer interface {
	Buffer() []Point
}

func Init() error {
	Body = NewGrid()
	Body.X = 0
	Body.Y = 0
	Body.BgColor = theme.BodyBg
	defer (func() {
		w, _ := tm.Size()
		Body.Width = w
	})()
	return tm.Init()
}

func Close() {
	tm.Close()
}

func TermWidth() int {
	w, _ := tm.Size()
	return w
}

func TermHeight() int {
	_, h := tm.Size()
	return h
}

// render all from left to right, right could overlap on left ones
func Render(rs ...Bufferer) {
	tm.Clear(tm.ColorDefault, toTmAttr(theme.BodyBg))
	for _, r := range rs {
		buf := r.Buffer()
		for _, v := range buf {
			tm.SetCell(v.X, v.Y, v.Ch, toTmAttr(v.Fg), toTmAttr(v.Bg))
		}
	}
	tm.Flush()
}
