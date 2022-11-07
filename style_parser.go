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

type PreparedStyle struct {
	Text  string
	Style string
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
	// "test [blue](fg:blue,bg:white,mod:bold) and [red](fg:red)"
	tokens := strings.Split(s, "](")
	if len(tokens) == 1 {
		return tokens
	}

	styleString := ""
	remainder := tokens[0]
	i := 1
	for {
		prefix, item := lookLeftForBracket(remainder)
		styleString, remainder = lookRightForEndStyle(tokens[i])
		i++
		fmt.Println(i, prefix)
		fmt.Println(i, item)
		fmt.Println(i, styleString)
		fmt.Println(i, remainder)
		if !strings.Contains(remainder, "[") {
			break
		}
	}

	buffer := []string{}

	return buffer
}

func PrepareStyles(s string) []PreparedStyle {
	items := []PreparedStyle{}
	tokens := strings.Split(s, "](")
	if len(tokens) == 1 {
		// easy case, not styled string
		ps := PreparedStyle{s, ""}
		return []PreparedStyle{ps}
	}

	fmt.Println(strings.Join(tokens, "|"))
	return items
}

// ParseStyles parses a string for embedded Styles and returns []Cell with the correct styling.
// Uses defaultStyle for any text without an embedded style.
// Syntax is of the form [text](fg:<color>,mod:<attribute>,bg:<color>).
// Ordering does not matter. All fields are optional.
func ParseStyles(s string, defaultStyle Style) []Cell {
	//test [blue](fg:blue,bg:white,mod:bold)
	cells := []Cell{}

	tokens := strings.Split(s, "](")
	if len(tokens) == 1 {
		// easy case, not styled string
		return cells
	}

	styleString, rest := processToken(tokens[len(tokens)-1], "")
	fmt.Println("2", styleString, rest)
	for i := len(tokens) - 2; i >= 0; i-- {
		styleString, rest = processToken(tokens[i], styleString)
		fmt.Println("3", styleString, rest)
	}

	//1 fg:red)
	//1 fg:blue,bg:white,mod:bold) and [red
	//1 test [blue

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
