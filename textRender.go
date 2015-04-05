package termui

import (
	"regexp"
	"strings"
)

// TextRender adds common methods for rendering a text on screeen.
type TextRender interface {
	NormalizedText() string
	Render(lastColor, background Attribute) RenderedSequence
	RenderSequence(start, end int, lastColor, background Attribute) RenderedSequence
}

// TextRendererFactory is factory for creating text renderers.
type TextRendererFactory interface {
	TextRenderer(text string) TextRender
}

// MarkdownRegex is used by MarkdownTextRenderer to determine how to format the
// text.
const MarkdownRegex = `(?:\[([^]]+)\])\(([a-z\s,]+)\)`

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

// Returns the position considering unicode characters.
func posUnicode(text string, pos int) int {
	return len([]rune(text[:pos]))
}

/*
Render renders the sequence `text` using a markdown-like syntax:
`[hello](red) world` will become: `hello world` where hello is red.

You may also specify other attributes such as bold text:
`[foo](YELLOW, BOLD)` will become `foo` in yellow, bold text.


For all available combinations, colors, and attribute, see: `StringToAttribute`.

This method returns a RenderedSequence
*/
func (r MarkdownTextRenderer) Render(lastColor, background Attribute) RenderedSequence {
	return r.RenderSequence(0, -1, lastColor, background)
}

// RenderSequence renders the text just like Render but the start and end can
// be specified.
func (r MarkdownTextRenderer) RenderSequence(start, end int, lastColor, background Attribute) RenderedSequence {
	text := r.Text
	if end == -1 {
		end = len([]rune(r.NormalizedText()))
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
		theSequence.Start = posUnicode(text, contentStart) - 1
		theSequence.End = posUnicode(text, contentEnd) - 1

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

	runes := []rune(text)[start:end]
	return RenderedSequence{string(runes), lastColor, background, sequences, nil}
}

// MarkdownTextRendererFactory is a TextRendererFactory for
// the MarkdownTextRenderer.
type MarkdownTextRendererFactory struct{}

// TextRenderer returns a MarkdownTextRenderer instance.
func (f MarkdownTextRendererFactory) TextRenderer(text string) TextRender {
	return MarkdownTextRenderer{text}
}

// RenderedSequence is a string sequence that is capable of returning the
// Buffer used by termui for displaying the colorful string.
type RenderedSequence struct {
	NormalizedText  string
	LastColor       Attribute
	BackgroundColor Attribute
	Sequences       []ColorSubsequence

	// Use the color() method for getting the correct value.
	mapSequences map[int]Attribute
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

func (s *RenderedSequence) colors() map[int]Attribute {
	if s.mapSequences == nil {
		s.mapSequences = ColorSubsequencesToMap(s.Sequences)
	}

	return s.mapSequences
}

// Buffer returns the colorful formatted buffer and the last color that was
// used.
func (s *RenderedSequence) Buffer(x, y int) ([]Point, Attribute) {
	buffer := make([]Point, 0, len(s.NormalizedText)) // This is just an assumtion

	for i := range []rune(s.NormalizedText) {
		p, width := s.PointAt(i, x, y)
		buffer = append(buffer, p)
		x += width
	}

	return buffer, s.LastColor
}

// PointAt returns the point at the position of n. The x, and y coordinates
// are used for placing the point at the right position.
// Since some charaters are wider (like `一`), this method also returns the
// width the point actually took.
// This is important for increasing the x property properly.
func (s *RenderedSequence) PointAt(n, x, y int) (Point, int) {
	color, ok := s.colors()[n]
	if !ok {
		color = s.LastColor
	}

	char := []rune(s.NormalizedText)[n]
	return Point{char, s.BackgroundColor, color, x, y}, charWidth(char)
}

// A NoopRenderer does not render the text at all.
type NoopRenderer struct {
	Text string
}

// NormalizedText returns the text given in
func (r NoopRenderer) NormalizedText() string {
	return r.Text
}

// RenderSequence returns a RenderedSequence that does not have any color
// sequences.
func (r NoopRenderer) RenderSequence(start, end int, lastColor, background Attribute) RenderedSequence {
	runes := []rune(r.Text)
	if end < 0 {
		end = len(runes)
	}

	runes = runes[start:end]
	var s []ColorSubsequence
	return RenderedSequence{string(runes), lastColor, background, s, nil}
}

// Render just like RenderSequence
func (r NoopRenderer) Render(lastColor, background Attribute) RenderedSequence {
	return r.RenderSequence(0, -1, lastColor, background)
}

// NoopRendererFactory is a TextRendererFactory for
// the NoopRenderer.
type NoopRendererFactory struct{}

// TextRenderer returns a NoopRenderer instance.
func (f NoopRendererFactory) TextRenderer(text string) TextRender {
	return NoopRenderer{text}
}
