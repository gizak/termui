// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import "testing"

var cmap = map[string]Attribute{
	"fg":           ColorWhite,
	"bg":           ColorDefault,
	"border.fg":    ColorWhite,
	"label.fg":     ColorGreen,
	"par.fg":       ColorYellow,
	"par.label.bg": ColorWhite,
}

func TestLoopUpAttr(t *testing.T) {
	tbl := []struct {
		name   string
		should Attribute
	}{
		{"par.label.bg", ColorWhite},
		{"par.label.fg", ColorGreen},
		{"par.bg", ColorDefault},
		{"bar.border.fg", ColorWhite},
		{"bar.label.bg", ColorDefault},
	}

	for _, v := range tbl {
		if lookUpAttr(cmap, v.name) != v.should {
			t.Error(v.name)
		}
	}
}
