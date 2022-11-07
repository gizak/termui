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

func processToken(token, previous string) (string, string) {
	fmt.Println("1", token)
	index := strings.Index(token, ")")
	if index == -1 {
		return "", ""
	}
	styleString := token[0:index]
	restOfString := token[index:]
	return styleString, restOfString
}

func lookLeftForBracket(s string) (string, string) {
	index := strings.LastIndex(s, "[")
	return s[0:index], s[index+1:]
}

func lookRightForEndStyle(s string) (string, string) {
	index := strings.Index(s, ")")
	return s[0:index], s[index+1:]
}

func BreakByStyles(s string) []string {
	tokens := strings.Split(s, "](")
	if len(tokens) == 1 {
		return tokens
	}

	buff := []string{}
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
		if !strings.Contains(remainder, "[") {
			buff = append(buff, remainder)
			break
		}
	}

	return buff
}

func containsColorOrMod(s string) bool {
	if strings.Contains(s, "fg:") {
		return true
	}
	if strings.Contains(s, "bg:") {
		return true
	}
	if strings.Contains(s, "mod:") {
		return true
	}

	return false
}

// ParseStyles parses a string for embedded Styles and returns []Cell with the correct styling.
// Uses defaultStyle for any text without an embedded style.
// Syntax is of the form [text](fg:<color>,mod:<attribute>,bg:<color>).
// Ordering does not matter. All fields are optional.
func ParseStyles(s string, defaultStyle Style) []Cell {
	cells := []Cell{}

	fmt.Println("11")
	items := BreakByStyles(s)
	fmt.Println("11", len(items))
	if len(items) == 1 {
		runes := []rune(s)
		for _, _rune := range runes {
			cells = append(cells, Cell{_rune, defaultStyle})
		}
		return cells
	}

	//test  blue fg:blue,bg:white,mod:bold  and  red fg:red  and maybe even  foo bg:red !
	for i := len(items) - 1; i > -1; i-- {
		if containsColorOrMod(items[i]) {
		} else {
			fmt.Println(items[i])
		}
	}

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
