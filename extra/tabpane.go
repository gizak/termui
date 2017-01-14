// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package extra

import (
	"unicode/utf8"

	. "github.com/gizak/termui"
)

type Tab struct {
	Label   string
	RuneLen int
	Blocks  []Bufferer
}

func NewTab(label string) *Tab {
	return &Tab{
		Label:   label,
		RuneLen: utf8.RuneCount([]byte(label))}
}

func (tab *Tab) AddBlocks(rs ...Bufferer) {
	for _, r := range rs {
		tab.Blocks = append(tab.Blocks, r)
	}
}

func (tab *Tab) Buffer() Buffer {
	buf := NewBuffer()
	for blockNum := 0; blockNum < len(tab.Blocks); blockNum++ {
		b := tab.Blocks[blockNum]
		buf.Merge(b.Buffer())
	}
	return buf
}

type Tabpane struct {
	Block
	Tabs           []Tab
	activeTabIndex int
	ActiveTabBg    Attribute
	posTabText     []int
	offTabText     int
}

func NewTabpane() *Tabpane {
	tp := Tabpane{
		Block:          *NewBlock(),
		activeTabIndex: 0,
		offTabText:     0,
		ActiveTabBg:    ThemeAttr("bg.tab.active")}
	return &tp
}

func (tp *Tabpane) SetTabs(tabs ...Tab) {
	tp.Tabs = make([]Tab, len(tabs))
	tp.posTabText = make([]int, len(tabs)+1)
	off := 0
	for i := 0; i < len(tp.Tabs); i++ {
		tp.Tabs[i] = tabs[i]
		tp.posTabText[i] = off
		off += tp.Tabs[i].RuneLen + 1 //+1 for space between tabs
	}
	tp.posTabText[len(tabs)] = off - 1 //total length of Tab's text
}

func (tp *Tabpane) SetActiveLeft() {
	if tp.activeTabIndex == 0 {
		return
	}
	tp.activeTabIndex -= 1
	if tp.posTabText[tp.activeTabIndex] < tp.offTabText {
		tp.offTabText = tp.posTabText[tp.activeTabIndex]
	}
}

func (tp *Tabpane) SetActiveRight() {
	if tp.activeTabIndex == len(tp.Tabs)-1 {
		return
	}
	tp.activeTabIndex += 1
	endOffset := tp.posTabText[tp.activeTabIndex] + tp.Tabs[tp.activeTabIndex].RuneLen
	if endOffset+tp.offTabText > tp.InnerWidth() {
		tp.offTabText = endOffset - tp.InnerWidth()
	}
}

// Checks if left and right tabs are fully visible
// if only left tabs are not visible return -1
// if only right tabs are not visible return 1
// if both return 0
// use only if fitsWidth() returns false
func (tp *Tabpane) checkAlignment() int {
	ret := 0
	if tp.offTabText > 0 {
		ret = -1
	}
	if tp.offTabText+tp.InnerWidth() < tp.posTabText[len(tp.Tabs)] {
		ret += 1
	}
	return ret
}

// Checks if all tabs fits innerWidth of Tabpane
func (tp *Tabpane) fitsWidth() bool {
	return tp.InnerWidth() >= tp.posTabText[len(tp.Tabs)]
}

func (tp *Tabpane) align() {
	if !tp.fitsWidth() && !tp.Border {
		tp.PaddingLeft += 1
		tp.PaddingRight += 1
		tp.Block.Align()
	}
}

// bridge the old Point stuct
type point struct {
	X  int
	Y  int
	Ch rune
	Fg Attribute
	Bg Attribute
}

func buf2pt(b Buffer) []point {
	ps := make([]point, 0, len(b.CellMap))
	for k, c := range b.CellMap {
		ps = append(ps, point{X: k.X, Y: k.Y, Ch: c.Ch, Fg: c.Fg, Bg: c.Bg})
	}

	return ps
}

// Adds the point only if it is visible in Tabpane.
// Point can be invisible if concatenation of Tab's texts is widther then
// innerWidth of Tabpane
func (tp *Tabpane) addPoint(ptab []point, charOffset *int, oftX *int, points ...point) []point {
	if *charOffset < tp.offTabText || tp.offTabText+tp.InnerWidth() < *charOffset {
		*charOffset++
		return ptab
	}
	for _, p := range points {
		p.X = *oftX
		ptab = append(ptab, p)
	}
	*oftX++
	*charOffset++
	return ptab
}

// Draws the point and redraws upper and lower border points (if it has one)
func (tp *Tabpane) drawPointWithBorder(p point, ch rune, chbord rune, chdown rune, chup rune) []point {
	var addp []point
	p.Ch = ch
	if tp.Border {
		p.Ch = chdown
		p.Y = tp.InnerY() - 1
		addp = append(addp, p)
		p.Ch = chup
		p.Y = tp.InnerY() + 1
		addp = append(addp, p)
		p.Ch = chbord
	}
	p.Y = tp.InnerY()
	return append(addp, p)
}

func (tp *Tabpane) Buffer() Buffer {
	if tp.Border {
		tp.Height = 3
	} else {
		tp.Height = 1
	}
	if tp.Width > tp.posTabText[len(tp.Tabs)]+2 {
		tp.Width = tp.posTabText[len(tp.Tabs)] + 2
	}
	buf := tp.Block.Buffer()
	ps := []point{}

	tp.align()
	if tp.InnerHeight() <= 0 || tp.InnerWidth() <= 0 {
		return NewBuffer()
	}
	oftX := tp.InnerX()
	charOffset := 0
	pt := point{Bg: tp.BorderBg, Fg: tp.BorderFg}
	for i, tab := range tp.Tabs {

		if i != 0 {
			pt.X = oftX
			pt.Y = tp.InnerY()
			addp := tp.drawPointWithBorder(pt, ' ', VERTICAL_LINE, HORIZONTAL_DOWN, HORIZONTAL_UP)
			ps = tp.addPoint(ps, &charOffset, &oftX, addp...)
		}

		if i == tp.activeTabIndex {
			pt.Bg = tp.ActiveTabBg
		}
		rs := []rune(tab.Label)
		for k := 0; k < len(rs); k++ {

			addp := make([]point, 0, 2)
			if i == tp.activeTabIndex && tp.Border {
				pt.Ch = ' '
				pt.Y = tp.InnerY() + 1
				pt.Bg = tp.BorderBg
				addp = append(addp, pt)
				pt.Bg = tp.ActiveTabBg
			}

			pt.Y = tp.InnerY()
			pt.Ch = rs[k]

			addp = append(addp, pt)
			ps = tp.addPoint(ps, &charOffset, &oftX, addp...)
		}
		pt.Bg = tp.BorderBg

		if !tp.fitsWidth() {
			all := tp.checkAlignment()
			pt.X = tp.InnerX() - 1

			pt.Ch = '*'
			if tp.Border {
				pt.Ch = VERTICAL_LINE
			}
			ps = append(ps, pt)

			if all <= 0 {
				addp := tp.drawPointWithBorder(pt, '<', '«', HORIZONTAL_LINE, HORIZONTAL_LINE)
				ps = append(ps, addp...)
			}

			pt.X = tp.InnerX() + tp.InnerWidth()
			pt.Ch = '*'
			if tp.Border {
				pt.Ch = VERTICAL_LINE
			}
			ps = append(ps, pt)
			if all >= 0 {
				addp := tp.drawPointWithBorder(pt, '>', '»', HORIZONTAL_LINE, HORIZONTAL_LINE)
				ps = append(ps, addp...)
			}
		}

		//draw tab content below the Tabpane
		if i == tp.activeTabIndex {
			blockPoints := buf2pt(tab.Buffer())
			for i := 0; i < len(blockPoints); i++ {
				blockPoints[i].Y += tp.Height + tp.Y
			}
			ps = append(ps, blockPoints...)
		}
	}

	for _, v := range ps {
		buf.Set(v.X, v.Y, NewCell(v.Ch, v.Fg, v.Bg))
	}
	buf.Sync()
	return buf
}
