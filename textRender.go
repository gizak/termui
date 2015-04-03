package termui

import (
	"regexp"
	"strings"
)

// TextRender adds common methods for rendering a text on screeen.
type TextRender interface {
	NormalizedText(text string) string
	RenderSequence(text string, lastColor, background Attribute) RenderedSequence
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
	return r.RenderSequence(text, 0, 0).NormalizedText
}

/*
RenderSequence renders the sequence `text` using a markdown-like syntax:
`[hello](red) world` will become: `hello world` where hello is red.

You may also specify other attributes such as bold text:
`[foo](YELLOW, BOLD)` will become `foo` in yellow, bold text.


For all available combinations, colors, and attribute, see: `StringToAttribute`.

This method returns a RenderedSequence
*/
func (r MarkdownTextRenderer) RenderSequence(text string, lastColor, background Attribute) RenderedSequence {
	getMatch := func(s string) []int {
		return markdownPattern.FindStringSubmatchIndex(strings.ToLower(s))
	}

	var sequences []ColorSubsequence
	for match := getMatch(text); match != nil; match = getMatch(text) {
		start, end := match[0], match[1]
		colorStart, colorEnd := match[4], match[5]
		contentStart, contentEnd := match[2], match[3]

		color := StringToAttribute(text[colorStart:colorEnd])
		content := text[contentStart:contentEnd]
		theSequence := ColorSubsequence{color, contentStart - 1, contentEnd - 1}

		sequences = append(sequences, theSequence)
		text = text[:start] + content + text[end:]
	}

	return RenderedSequence{text, lastColor, background, sequences}
}

// RenderedSequence is a string sequence that is capable of returning the
// Buffer used by termui for displaying the colorful string.
type RenderedSequence struct {
	NormalizedText  string
	LastColor       Attribute
	BackgroundColor Attribute
	Sequences       []ColorSubsequence
}

// A ColorSubsequence represents a color for the given text span.
type ColorSubsequence struct {
	Color Attribute
	Start int
	End   int
}

// Buffer returns the colorful formatted buffer and the last color that was
// used.
func (s *RenderedSequence) Buffer(x, y int) ([]Point, Attribute) {
	return nil, s.LastColor
}
