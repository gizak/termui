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
)

type parserState uint

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

func lookLeftForBracket(s string) int {
	return strings.LastIndex(s, "[")
}

func lookRightForEndStyle(s string) int {
	return strings.Index(s, ")")
}

type StyleBlock struct {
	Start int
	End   int
}

func FindStyleBlocks(s string) []StyleBlock {
	items := []StyleBlock{}
	return items
}

func FindStylePositions(s string) []int {
	fmt.Println(s)
	index := strings.Index(s, "](")
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
		index = strings.Index(toProcess, "](")
		if index == -1 {
			break
		}
	}

	return buff

	// test [blue](fg:blue,mod:bold) and [red](fg:red) and maybe even [foo](bg:red)!

	/*
		offset := 0
		for i, token := range tokens {
			if i%2 == 0 {
				index := lookLeftForBracket(token)
				ssp := StyleStringPosition{offset, offset + index, false}
				ssps = append(ssps, ssp)
				ssp = StyleStringPosition{offset + index, offset + len(token), true}
				ssps = append(ssps, ssp)
				fmt.Println(i, "even", len(token), index)
			} else {
				index := lookRightForEndStyle(token)
				ssp := StyleStringPosition{offset, offset + index, true}
				ssps = append(ssps, ssp)
				ssp = StyleStringPosition{offset + index, offset + len(token), false}
				ssps = append(ssps, ssp)
				fmt.Println(i, "odd", len(token), index)
			}
			offset += len(token) + 2
			buff = append(buff, token)
		}
		fmt.Println(ssps)


			styleString := ""
			remainder := tokens[0]
			i := 1
			for {
				prefix, item := lookLeftForBracket(remainder)
				styleString, remainder = lookRightForEndStyle(tokens[i])
				i++
				buff = append(buff, prefix)
				buff = append(buff, item)
				buff = append(buff, styleString)
				if !strings.Contains(remainder, "*") {
					buff = append(buff, remainder)
					break
				}
				if i > len(tokens)-1 {
					break
				}
			}*/

	return buff
}

// ParseStyles parses a string for embedded Styles and returns []Cell with the correct styling.
// Uses defaultStyle for any text without an embedded style.
// Syntax is of the form [text](fg:<color>,mod:<attribute>,bg:<color>).
// Ordering does not matter. All fields are optional.
func ParseStyles(s string, defaultStyle Style) []Cell {
	cells := []Cell{}

	/*
		items := BreakByStyles(s)
		if len(items) == 1 {
			runes := []rune(s)
			for _, _rune := range runes {
				cells = append(cells, Cell{_rune, defaultStyle})
			}
			return cells
		}

		style := defaultStyle
		for i := len(items) - 1; i > -1; i-- {
			if containsColorOrMod(items[i]) {
				style = readStyle([]rune(items[i]), defaultStyle)
			} else {
				cells = append(RunesToStyledCells([]rune(items[i]), style), cells...)
				style = defaultStyle
			}
		}
	*/
	return cells
}

func ParseStyles2(s string, defaultStyle Style) []Cell {
	cells := []Cell{}
	runes := []rune(s)
	state := parserStateDefault
	styledText := []rune{}
	styleItems := []rune{}
	squareCount := 0

	reset := func() {
		styledText = []rune{}
		styleItems = []rune{}
		state = parserStateDefault
		squareCount = 0
	}

	rollback := func() {
		cells = append(cells, RunesToStyledCells(styledText, defaultStyle)...)
		cells = append(cells, RunesToStyledCells(styleItems, defaultStyle)...)
		reset()
	}

	// chop first and last runes
	chop := func(s []rune) []rune {
		return s[1 : len(s)-1]
	}

	for i, _rune := range runes {
		switch state {
		case parserStateDefault:
			if _rune == tokenBeginStyledText {
				state = parserStateStyledText
				squareCount = 1
				styledText = append(styledText, _rune)
			} else {
				cells = append(cells, Cell{_rune, defaultStyle})
			}
		case parserStateStyledText:
			switch {
			case squareCount == 0:
				switch _rune {
				case tokenBeginStyle:
					state = parserStateStyleItems
					styleItems = append(styleItems, _rune)
				default:
					rollback()
					switch _rune {
					case tokenBeginStyledText:
						state = parserStateStyledText
						squareCount = 1
						styleItems = append(styleItems, _rune)
					default:
						cells = append(cells, Cell{_rune, defaultStyle})
					}
				}
			case len(runes) == i+1:
				rollback()
				styledText = append(styledText, _rune)
			case _rune == tokenBeginStyledText:
				squareCount++
				styledText = append(styledText, _rune)
			case _rune == tokenEndStyledText:
				squareCount--
				styledText = append(styledText, _rune)
			default:
				styledText = append(styledText, _rune)
			}
		case parserStateStyleItems:
			styleItems = append(styleItems, _rune)
			if _rune == tokenEndStyle {
				style := readStyle(chop(styleItems), defaultStyle)
				cells = append(cells, RunesToStyledCells(chop(styledText), style)...)
				reset()
			} else if len(runes) == i+1 {
				rollback()
			}
		}
	}

	return cells
}
