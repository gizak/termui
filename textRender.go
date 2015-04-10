package termui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TextRenderer adds common methods for rendering a text on screeen.
type TextRenderer interface {
	NormalizedText() string
	Render(lastColor, background Attribute) RenderedSequence
	RenderSequence(start, end int, lastColor, background Attribute) RenderedSequence
}

// TextRendererFactory is factory for creating text renderers.
type TextRendererFactory interface {
	TextRenderer(text string) TextRenderer
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
func (f MarkdownTextRendererFactory) TextRenderer(text string) TextRenderer {
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

// A PlainRenderer does not render the text at all.
type PlainRenderer struct {
	Text string
}

// NormalizedText returns the text given in
func (r PlainRenderer) NormalizedText() string {
	return r.Text
}

// RenderSequence returns a RenderedSequence that does not have any color
// sequences.
func (r PlainRenderer) RenderSequence(start, end int, lastColor, background Attribute) RenderedSequence {
	runes := []rune(r.Text)
	if end < 0 {
		end = len(runes)
	}

	runes = runes[start:end]
	var s []ColorSubsequence
	return RenderedSequence{string(runes), lastColor, background, s, nil}
}

// Render just like RenderSequence
func (r PlainRenderer) Render(lastColor, background Attribute) RenderedSequence {
	return r.RenderSequence(0, -1, lastColor, background)
}

// PlainRendererFactory is a TextRendererFactory for
// the PlainRenderer.
type PlainRendererFactory struct{}

// TextRenderer returns a PlainRenderer instance.
func (f PlainRendererFactory) TextRenderer(text string) TextRenderer {
	return PlainRenderer{text}
}

// We can't use a raw string here because \033 must not be escaped.
// I'd like to append (?<=m; i.e. lookbehind), but unfortunately,
// it is not supported. So we will need to do that manually.
var escapeRegex = "\033\\[(([0-9]{1,2}[;m])+)"
var colorEscapeCodeRegex = regexp.MustCompile(escapeRegex)
var colorEscapeCodeRegexMatchAll = regexp.MustCompile("^" + escapeRegex + "$")

// An EscapeCode is a unix ASCII Escape code.
type EscapeCode string

func (e EscapeCode) escapeNumberToColor(colorID int) (Attribute, error) {
	var color Attribute
	switch colorID {
	case 0:
		color = ColorDefault

	case 1:
		color = AttrBold

	case 4:
		color = AttrUnderline

	case 30:
		color = ColorBlack

	case 31:
		color = ColorRed

	case 32:
		color = ColorGreen

	case 33:
		color = ColorYellow

	case 34:
		color = ColorBlue

	case 35:
		color = ColorMagenta

	case 36:
		color = ColorCyan

	case 37:
		color = ColorWhite

	default:
		safeCode := e.MakeSafe()
		return 0, fmt.Errorf("Unkown/unsupported escape code: '%v'", safeCode)
	}

	return color, nil
}

// Color converts the escape code to an `Attribute` (color).
// The EscapeCode must be formatted like this:
//  - ASCII-Escape chacter (\033) + [ + Number + (;Number...) + m
// The second number is optimal. The semicolon (;) is used
// to seperate the colors.
// For example: `\033[1;31m` means: the following text is red and bold.
func (e EscapeCode) Color() (Attribute, error) {
	escapeCode := string(e)
	matches := colorEscapeCodeRegexMatchAll.FindStringSubmatch(escapeCode)
	invalidEscapeCode := func() error {
		safeCode := e.MakeSafe()
		return fmt.Errorf("%v is not a valid ASCII escape code", safeCode)
	}

	if matches == nil || escapeCode[len(escapeCode)-1] != 'm' {
		return 0, invalidEscapeCode()
	}

	color := Attribute(0)
	for _, id := range strings.Split(matches[1][:len(matches[1])-1], ";") {
		colorID, err := strconv.Atoi(id)
		if err != nil {
			return 0, invalidEscapeCode()
		}

		newColor, err := e.escapeNumberToColor(colorID)
		if err != nil {
			return 0, err
		}

		color |= newColor
	}

	return color, nil
}

// MakeSafe replace the invisible escape code chacacter (\0333)
// with \\0333 so that it will not mess up the terminal when an error
// is shown.
func (e EscapeCode) MakeSafe() string {
	return strings.Replace(string(e), "\033", "\\033", -1)
}

// Alias to `EscapeCode.MakeSafe()`
func (e EscapeCode) String() string {
	return e.MakeSafe()
}

// Raw returns the raw value of the escape code.
// Alias to string(EscapeCode)
func (e EscapeCode) Raw() string {
	return string(e)
}

// IsValid returns whether or not the syntax of the escape code is
// valid and the code is supported.
func (e EscapeCode) IsValid() bool {
	_, err := e.Color()
	return err == nil
}

// A EscapeCodeRenderer does not render the text at all.
type EscapeCodeRenderer struct {
	Text string
}

// NormalizedText strips all escape code outs (even the unkown/unsupported)
// ones.
func (r EscapeCodeRenderer) NormalizedText() string {
	matches := colorEscapeCodeRegex.FindAllStringIndex(r.Text, -1)
	text := []byte(r.Text)

	// Iterate through matches in reverse order
	for i := len(matches) - 1; i >= 0; i-- {
		start, end := matches[i][0], matches[i][1]
		if EscapeCode(text[start:end]).IsValid() {
			text = append(text[:start], text[end:]...)
		}
	}

	return string(text)
}

// RenderSequence renders the text just like Render but the start and end may
// be set. If end is -1, the end of the string will be used.
func (r EscapeCodeRenderer) RenderSequence(start, end int, lastColor, background Attribute) RenderedSequence {
	normalizedRunes := []rune(r.NormalizedText())
	if end < 0 {
		end = len(normalizedRunes)
	}

	text := []byte(r.Text)
	matches := colorEscapeCodeRegex.FindAllSubmatchIndex(text, -1)
	removed := 0
	var sequences []ColorSubsequence
	runeLength := func(length int) int {
		return len([]rune(string(text[:length])))
	}

	runes := []rune(r.Text)
	for _, theMatch := range matches {
		// Escapde code start, escape code end
		eStart := runeLength(theMatch[0]) - removed
		eEnd := runeLength(theMatch[1]) - removed
		escapeCode := EscapeCode(runes[eStart:eEnd])

		// If an error occurs (e.g. unkown escape code), we will just ignore it :)
		color, err := escapeCode.Color()
		if err != nil {
			continue
		}

		// Patch old color sequence
		if len(sequences) > 0 {
			last := &sequences[len(sequences)-1]
			last.End = eStart - start
		}

		// eEnd < 0 means the the sequence is withing the range.
		if eEnd-start >= 0 {
			// The sequence starts when the escape code ends and ends when the text
			// end. If there is another escape code, this will be patched in the
			// previous line.
			colorSeq := ColorSubsequence{color, eStart - start, end - start}
			if colorSeq.Start < 0 {
				colorSeq.Start = 0
			}

			sequences = append(sequences, colorSeq)
		}

		runes = append(runes[:eStart], runes[eEnd:]...)
		removed += eEnd - eStart
	}

	runes = runes[start:end]
	return RenderedSequence{string(runes), lastColor, background, sequences, nil}
}

// Render just like RenderSequence
func (r EscapeCodeRenderer) Render(lastColor, background Attribute) RenderedSequence {
	return r.RenderSequence(0, -1, lastColor, background)
}

// EscapeCodeRendererFactory is a TextRendererFactory for
// the EscapeCodeRenderer.
type EscapeCodeRendererFactory struct{}

// TextRenderer returns a EscapeCodeRenderer instance.
func (f EscapeCodeRendererFactory) TextRenderer(text string) TextRenderer {
	return EscapeCodeRenderer{text}
}
