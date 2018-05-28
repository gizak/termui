// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import "testing"

func TestReadAttr(t *testing.T) {
	m := MarkdownTxBuilder{}
	m.baseFg = ColorCyan | AttrUnderline
	m.baseBg = ColorBlue | AttrBold
	fg, bg := m.readAttr("fg-red,bg-reverse")
	if fg != ColorRed|AttrUnderline || bg != ColorBlue|AttrBold|AttrReverse {
		t.Error("readAttr failed")
	}
}

func TestMTBParse(t *testing.T) {
	/*
		str := func(cs []Cell) string {
			rs := make([]rune, len(cs))
			for i := range cs {
				rs[i] = cs[i].Ch
			}
			return string(rs)
		}
	*/

	tbls := [][]string{
		{"hello world", "hello world"},
		{"[hello](fg-red) world", "hello world"},
		{"[[hello]](bg-red) world", "[hello] world"},
		{"[1] hello world", "[1] hello world"},
		{"[[1]](bg-white) [hello] world", "[1] [hello] world"},
		{"[hello world]", "[hello world]"},
		{"", ""},
		{"[hello world)", "[hello world)"},
		{"[0] [hello](bg-red)[ world](fg-blue)!", "[0] hello world!"},
	}

	m := MarkdownTxBuilder{}
	m.baseFg = ColorWhite
	m.baseBg = ColorDefault
	for _, s := range tbls {
		m.reset()
		m.parse(s[0])
		res := string(m.plainTx)
		if s[1] != res {
			t.Errorf("\ninput :%s\nshould:%s\noutput:%s", s[0], s[1], res)
		}
	}

	m.reset()
	m.parse("[0] [hello](bg-red)[ world](fg-blue)")
	if len(m.markers) != 2 &&
		m.markers[0].st == 4 &&
		m.markers[0].ed == 11 &&
		m.markers[0].fg == ColorWhite &&
		m.markers[0].bg == ColorRed {
		t.Error("markers dismatch")
	}

	m2 := NewMarkdownTxBuilder()
	cs := m2.Build("[0] [hellob-e) wrd]fgblue)!", ColorWhite, ColorBlack)
	cs = m2.Build("[0] [hello](bg-red) [world](fg-blue)!", ColorWhite, ColorBlack)
	if cs[4].Ch != 'h' && cs[4].Bg != ColorRed && cs[4].Fg != ColorWhite {
		t.Error("dismatch in Build")
	}
}
