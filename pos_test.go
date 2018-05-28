// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"
	"testing"
)

func TestAlignArea(t *testing.T) {
	p := image.Rect(0, 0, 100, 100)
	c := image.Rect(10, 10, 20, 20)

	nc := AlignArea(p, c, AlignLeft)
	if nc.Min.X != 0 || nc.Max.Y != 20 {
		t.Errorf("AlignLeft failed:\n%+v", nc)
	}

	nc = AlignArea(p, c, AlignCenter)
	if nc.Min.X != 45 || nc.Max.Y != 55 {
		t.Error("AlignCenter failed")
	}

	nc = AlignArea(p, c, AlignBottom|AlignRight)
	if nc.Min.X != 90 || nc.Max.Y != 100 {
		t.Errorf("AlignBottom|AlignRight failed\n%+v", nc)
	}
}

func TestMoveArea(t *testing.T) {
	a := image.Rect(10, 10, 20, 20)
	a = MoveArea(a, 5, 10)
	if a.Min.X != 15 || a.Min.Y != 20 || a.Max.X != 25 || a.Max.Y != 30 {
		t.Error("MoveArea failed")
	}
}
