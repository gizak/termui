// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"fmt"
	"strings"
)

const (
	tokenFg       = "fg"
	tokenBg       = "bg"
	tokenModifier = "mod"

	tokenItemSeparator  = ","
	tokenValueSeparator = ":"

	tokenBeginStyledText = '['
	tokenEndStyledText   = ']'

	tokenBeginStyle = '('
	tokenEndStyle   = ')'

	tokenStyleKey = "]("
)

type parserState uint

type StyleBlock struct {
	Start int
	End   int
}

const (
	parserStateDefault parserState = iota
	parserStateStyleItems
	parserStateStyledText
)

// StyleParserColorMap can be modified to add custom color parsing to text
var StyleParserColorMap = map[string]Color{
	"red":     ColorRed,
	"blue":    ColorBlue,
	"black":   ColorBlack,
	"cyan":    ColorCyan,
	"yellow":  ColorYellow,
	"white":   ColorWhite,
	"clear":   ColorClear,
	"green":   ColorGreen,
	"magenta": ColorMagenta,
}

var modifierMap = map[string]Modifier{
	"bold":      ModifierBold,
	"underline": ModifierUnderline,
	"reverse":   ModifierReverse,
}

// readStyle translates an []rune like `fg:red,mod:bold,bg:white` to a style
func readStyle(runes []rune, defaultStyle Style) Style {
	style := defaultStyle
	split := strings.Split(string(runes), tokenItemSeparator)
	for _, item := range split {
		pair := strings.Split(item, tokenValueSeparator)
		if len(pair) == 2 {
			switch pair[0] {
			case tokenFg:
				style.Fg = StyleParserColorMap[pair[1]]
			case tokenBg:
				style.Bg = StyleParserColorMap[pair[1]]
			case tokenModifier:
				style.Modifier = modifierMap[pair[1]]
			}
		}
	}
	return style
}

// this will start at ]( and look backwards to find the [ and forward
// to find the ) and record these Start and End indexes in a StyleBlock
func findStartEndOfStyle(pos int, runes []rune) StyleBlock {
	current := pos
	sb := StyleBlock{0, 0}
	for {
		current--
		if runes[current] == tokenBeginStyledText {
			sb.Start = current
			break
		}
	}
	current = pos
	for {
		current++
		if runes[current] == tokenEndStyle {
			sb.End = current
			break
		}
	}
	return sb
}

// if are string is "foo [thing](style) foo [more](style)"
// this will return "foo ", "[thing](style)", " foo ", "[more](style)"
func breakBlocksIntoStrings(s string) []string {
	buff := []string{}
	blocks := findStyleBlocks(s)
	if len(blocks) == 0 {
		return buff
	}
	startEnd := len(s)
	for i := len(blocks) - 1; i >= 0; i-- {
		b := blocks[i]
		item := s[b.End+1 : startEnd]
		if item != "" {
			buff = append([]string{item}, buff...)
		}
		item = s[b.Start : b.End+1]
		buff = append([]string{item}, buff...)
		startEnd = b.Start
	}
	item := s[0:startEnd]
	if item != "" {
		buff = append([]string{item}, buff...)
	}
	return buff
}

// loop through positions and make [] of StyleBlocks
func findStyleBlocks(s string) []StyleBlock {
	items := []StyleBlock{}
	runes := []rune(s)

	positions := findStylePositions(s)
	for _, pos := range positions {
		sb := findStartEndOfStyle(pos, runes)
		items = append(items, sb)
	}
	return items
}

// uses tokenStyleKey ]( which tells us we have both a [text] and a (style)
// if are string is "foo [thing](style) foo [more](style)"
// this func will return a list of two ints: the index of the first ]( and
// the index of the next one
func findStylePositions(s string) []int {
	index := strings.Index(s, tokenStyleKey)
	if index == -1 {
		return []int{}
	}

	buff := []int{}

	toProcess := s
	offset := 0
	for {
		buff = append(buff, index+offset)
		toProcess = toProcess[index+1:]
		offset += index + 1
		index = strings.Index(toProcess, tokenStyleKey)
		if index == -1 {
			break
		}
	}

	return buff
}

func containsStyle(s string) bool {
	if strings.HasPrefix(s, string(tokenBeginStyledText)) &&
		strings.HasSuffix(s, string(tokenEndStyle)) &&
		strings.Contains(s, string(tokenEndStyledText)) &&
		strings.Contains(s, string(tokenBeginStyle)) {
		return true
	}
	return false
}

// [text](style) will return text
func extractTextFromBlock(item string) string {
	index := strings.Index(item, string(tokenEndStyledText))
	return item[1:index]
}

// [text](style) will return style
func extractStyleFromBlock(item string) string {
	index := strings.Index(item, string(tokenBeginStyle))
	return item[index+1 : len(item)-1]
}

// ParseStyles parses a string for embedded Styles and returns []Cell with the correct styling.
// Uses defaultStyle for any text without an embedded style.
// Syntax is of the form [text](fg:<color>,mod:<attribute>,bg:<color>).
// Ordering does not matter. All fields are optional.
func ParseStyles(s string, defaultStyle Style) []Cell {
	cells := []Cell{}

	items := breakBlocksIntoStrings(s)
	if len(items) == 0 {
		return RunesToStyledCells([]rune(s), defaultStyle)
	}

	for _, item := range items {
		if containsStyle(item) {
			text := extractTextFromBlock(item)
			styleText := extractStyleFromBlock(item)
			fmt.Println("|" + text + "|" + styleText + "|")
			style := readStyle([]rune(styleText), defaultStyle)
			cells = append(cells, RunesToStyledCells([]rune(text), style)...)
		} else {
			cells = append(cells, RunesToStyledCells([]rune(item), defaultStyle)...)
		}
	}
	return cells
}
