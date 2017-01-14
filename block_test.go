// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"testing"
)

func TestBlockFloat(t *testing.T) {
	Init()
	defer Close()

	b := NewBlock()
	b.X = 10
	b.Y = 20

	b.Float = AlignCenter
	b.Align()
}

func TestBlockInnerBounds(t *testing.T) {
	Init()
	defer Close()

	b := NewBlock()
	b.X = 10
	b.Y = 11
	b.Width = 12
	b.Height = 13

	assert := func(name string, x, y, w, h int) {
		t.Log(name)
		area := b.InnerBounds()
		cx := area.Min.X
		cy := area.Min.Y
		cw := area.Dx()
		ch := area.Dy()

		if cx != x {
			t.Errorf("expected x to be %d but got %d", x, cx)
		}
		if cy != y {
			t.Errorf("expected y to be %d but got %d\n%+v", y, cy, area)
		}
		if cw != w {
			t.Errorf("expected width to be %d but got %d", w, cw)
		}
		if ch != h {
			t.Errorf("expected height to be %d but got %d", h, ch)
		}
	}

	b.Border = false
	assert("no border, no padding", 10, 11, 12, 13)

	b.Border = true
	assert("border, no padding", 11, 12, 10, 11)

	b.PaddingBottom = 2
	assert("border, 2b padding", 11, 12, 10, 9)

	b.PaddingTop = 3
	assert("border, 2b 3t padding", 11, 15, 10, 6)

	b.PaddingLeft = 4
	assert("border, 2b 3t 4l padding", 15, 15, 6, 6)

	b.PaddingRight = 5
	assert("border, 2b 3t 4l 5r padding", 15, 15, 1, 6)
}
