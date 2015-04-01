// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import tm "github.com/nsf/termbox-go"
import rw "github.com/mattn/go-runewidth"

/* ---------------Port from termbox-go --------------------- */

// Attribute is printable cell's color and style.
type Attribute uint16

const (
	ColorDefault Attribute = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

const (
	AttrBold Attribute = 1 << (iota + 9)
	AttrUnderline
	AttrReverse
)

var (
	dot  = "…"
	dotw = rw.StringWidth(dot)
)

/* ----------------------- End ----------------------------- */

func toTmAttr(x Attribute) tm.Attribute {
	return tm.Attribute(x)
}

func str2runes(s string) []rune {
	return []rune(s)
}

// Here for backwards-compatibility.
func trimStr2Runes(s string, w int) []rune {
	return TrimStr2Runes(s, w)
}

// TrimStr2Runes trims string to w[-1 rune], appends …, and returns the runes
// of that string if string is grather then n. If string is small then w,
// return the runes.
func TrimStr2Runes(s string, w int) []rune {
	if w <= 0 {
		return []rune{}
	}

	sw := rw.StringWidth(s)
	if sw > w {
		return []rune(rw.Truncate(s, w, dot))
	}
	return str2runes(s)
}

// TrimStrIfAppropriate trim string to "s[:-1] + …"
// if string > width otherwise return string
func TrimStrIfAppropriate(s string, w int) string {
	if w <= 0 {
		return ""
	}

	sw := rw.StringWidth(s)
	if sw > w {
		return rw.Truncate(s, w, dot)
	}

	return s
}

func strWidth(s string) int {
	return rw.StringWidth(s)
}

func charWidth(ch rune) int {
	return rw.RuneWidth(ch)
}
