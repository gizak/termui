package termui

import (
	"regexp"
	"strings"
)

// TextRender adds common methods for rendering a text on screeen.
type TextRender interface {
	NormalizedText(text string) string
	RenderSequence(text string, lastColor, background Attribute) RenderedSubsequence
}

type subSecequence struct {
	start int
	end   int
	color Attribute
}

// MarkdownRegex is used by MarkdownTextRenderer to determine how to format the
// text.
const MarkdownRegex = `(?:\[([[a-z]+)\])\(([a-z\s,]+)\)`

// unexported because a pattern can't be a constant and we don't want anyone
// messing with the regex.
var markdownPattern = regexp.MustCompile(MarkdownRegex)

// MarkdownTextRenderer is used for rendering the text with colors using
// markdown-like syntax.
// See: https://github.com/gizak/termui/issues/4#issuecomment-87270635
type MarkdownTextRenderer struct{}

// NormalizedText returns the text the user will see (without colors).
// It strips out all formatting option and only preserves plain text.
func (r MarkdownTextRenderer) NormalizedText(text string) string {
	lText := strings.ToLower(text)
	indexes := markdownPattern.FindAllStringSubmatchIndex(lText, -1)

	// Interate through indexes in reverse order.
	for i := len(indexes) - 1; i >= 0; i-- {
		theIndex := indexes[i]
		start, end := theIndex[0], theIndex[1]
		contentStart, contentEnd := theIndex[2], theIndex[3]

		text = text[:start] + text[contentStart:contentEnd] + text[end:]
	}

	return text
}

// RenderedSubsequence is a string sequence that is capable of returning the
// Buffer used by termui for displaying the colorful string.
type RenderedSubsequence struct {
	RawText         string
	NormalizedText  string
	LastColor       Attribute
	BackgroundColor Attribute

	sequences subSecequence
}

// Buffer returns the colorful formatted buffer and the last color that was
// used.
func (s *RenderedSubsequence) Buffer(x, y int) ([]Point, Attribute) {
	// var buffer []Point
	// dx := 0
	// for _, r := range []rune(s.NormalizedText) {
	// 	p := Point{
	// 		Ch: r,
	// 		X:  x + dx,
	// 		Y:  y,
	// 		Fg: Attribute(rand.Intn(8)),
	// 		Bg: background,
	// 	}
	//
	// 	buffer = append(buffer, p)
	// 	dx += charWidth(r)
	// }
	//
	// return buffer
	return nil, s.LastColor
}
