package termui

import (
	"regexp"
	"strings"
)

// TextRender adds common methods for rendering a text on screeen.
type TextRender interface {
	NormalizedText(text string) string
	RenderSequence(start, end int, lastColor, background Attribute) RenderedSequence
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
type MarkdownTextRenderer struct {
	Text string
}

// NormalizedText returns the text the user will see (without colors).
// It strips out all formatting option and only preserves plain text.
func (r MarkdownTextRenderer) NormalizedText() string {
	return r.normalizeText(r.Text)
}

func (r MarkdownTextRenderer) normalizeText(text string) string {
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

/*
RenderSequence renders the sequence `text` using a markdown-like syntax:
`[hello](red) world` will become: `hello world` where hello is red.

You may also specify other attributes such as bold text:
`[foo](YELLOW, BOLD)` will become `foo` in yellow, bold text.


For all available combinations, colors, and attribute, see: `StringToAttribute`.

This method returns a RenderedSequence
*/
func (r MarkdownTextRenderer) RenderSequence(start, end int, lastColor, background Attribute) RenderedSequence {
	text := r.Text
	if end == -1 {
		end = len(r.NormalizedText())
	}

	getMatch := func(s string) []int {
		return markdownPattern.FindStringSubmatchIndex(strings.ToLower(s))
	}

	var sequences []ColorSubsequence
	for match := getMatch(text); match != nil; match = getMatch(text) {
		// Check if match is in the start/end range.

		matchStart, matchEnd := match[0], match[1]
		colorStart, colorEnd := match[4], match[5]
		contentStart, contentEnd := match[2], match[3]

		color := StringToAttribute(text[colorStart:colorEnd])
		content := text[contentStart:contentEnd]
		theSequence := ColorSubsequence{color, contentStart - 1, contentEnd - 1}

		if start < theSequence.End && end > theSequence.Start {
			// Make the sequence relative and append.
			theSequence.Start -= start
			if theSequence.Start < 0 {
				theSequence.Start = 0
			}

			theSequence.End -= start
			if theSequence.End < 0 {
				theSequence.End = 0
			} else if theSequence.End > end-start {
				theSequence.End = end - start
			}

			sequences = append(sequences, theSequence)
		}

		text = text[:matchStart] + content + text[matchEnd:]
	}

	if end == -1 {
		end = len(text)
	}
	return RenderedSequence{text[start:end], lastColor, background, sequences}
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

// ColorSubsequencesToMap creates a map with all colors that from the
// subsequences.
func ColorSubsequencesToMap(sequences []ColorSubsequence) map[int]Attribute {
	result := make(map[int]Attribute)
	for _, theSequence := range sequences {
		for i := theSequence.Start; i < theSequence.End; i++ {
			result[i] = theSequence.Color
		}
	}

	return result
}

// Buffer returns the colorful formatted buffer and the last color that was
// used.
func (s *RenderedSequence) Buffer(x, y int) ([]Point, Attribute) {
	buffer := make([]Point, 0, len(s.NormalizedText)) // This is just an assumtion

	colors := ColorSubsequencesToMap(s.Sequences)
	for i, r := range []rune(s.NormalizedText) {
		color, ok := colors[i]
		if !ok {
			color = s.LastColor
		}

		p := Point{r, s.BackgroundColor, color, x, y}
		buffer = append(buffer, p)
		x += charWidth(r)
	}

	return buffer, s.LastColor
}
