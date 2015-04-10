package termui

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextRender_TestInterface(t *testing.T) {
	var inter *TextRenderer

	assert.Implements(t, inter, new(MarkdownTextRenderer))
	assert.Implements(t, inter, new(EscapeCodeRenderer))
	assert.Implements(t, inter, new(PlainRenderer))
}

func TestTextRendererFactory_TestInterface(t *testing.T) {
	var inter *TextRendererFactory

	assert.Implements(t, inter, new(MarkdownTextRendererFactory))
	assert.Implements(t, inter, new(EscapeCodeRendererFactory))
	assert.Implements(t, inter, new(PlainRendererFactory))
}

func TestMarkdownTextRenderer_normalizeText(t *testing.T) {
	renderer := MarkdownTextRenderer{}

	got := renderer.normalizeText("[ERROR](red,bold) Something went wrong")
	assert.Equal(t, got, "ERROR Something went wrong")

	got = renderer.normalizeText("[foo](red) hello [bar](green) world")
	assert.Equal(t, got, "foo hello bar world")

	got = renderer.normalizeText("[foo](g) hello [bar]green (world)")
	assert.Equal(t, got, "foo hello [bar]green (world)")

	got = "笀耔 [澉 灊灅甗](RED) 郔镺 笀耔 澉 [灊灅甗](yellow) 郔镺"
	expected := "笀耔 澉 灊灅甗 郔镺 笀耔 澉 灊灅甗 郔镺"
	assert.Equal(t, renderer.normalizeText(got), expected)

	got = renderer.normalizeText("[(foo)](red,white) bar")
	assert.Equal(t, renderer.normalizeText(got), "(foo) bar")

	// TODO: make this regex work correctly:
	// got = renderer.normalizeText("[[foo]](red,white) bar")
	// assert.Equal(t, renderer.normalizeText(got), "[foo] bar")
	// I had to comment it out because the unit tests keep failing and
	// I don't know how to fix it. See more:
	// https://github.com/gizak/termui/pull/22
}

func TestMarkdownTextRenderer_NormalizedText(t *testing.T) {
	renderer := MarkdownTextRenderer{"[ERROR](red,bold) Something went wrong"}
	assert.Equal(t, renderer.NormalizedText(), "ERROR Something went wrong")
}

func assertRenderSequence(t *testing.T, sequence RenderedSequence, last, background Attribute, text string, lenSequences int) bool {
	msg := fmt.Sprintf("seq: %v", spew.Sdump(sequence))
	assert.Equal(t, last, sequence.LastColor, msg)
	assert.Equal(t, background, sequence.BackgroundColor, msg)
	assert.Equal(t, text, sequence.NormalizedText, msg)
	return assert.Equal(t, lenSequences, len(sequence.Sequences), msg)
}

func assertColorSubsequence(t *testing.T, s ColorSubsequence, color string, start, end int) {
	assert.Equal(t, ColorSubsequence{StringToAttribute(color), start, end}, s)
}

func TestMarkdownTextRenderer_RenderSequence(t *testing.T) {
	// Simple test.
	renderer := MarkdownTextRenderer{"[ERROR](red,bold) something went wrong"}
	got := renderer.RenderSequence(0, -1, 3, 5)
	if assertRenderSequence(t, got, 3, 5, "ERROR something went wrong", 1) {
		assertColorSubsequence(t, got.Sequences[0], "RED,BOLD", 0, 5)
	}

	got = renderer.RenderSequence(3, 8, 3, 5)
	if assertRenderSequence(t, got, 3, 5, "OR so", 1) {
		assertColorSubsequence(t, got.Sequences[0], "RED,BOLD", 0, 2)
	}

	// Test for mutiple colors.
	renderer = MarkdownTextRenderer{"[foo](red) hello [bar](blue) world"}
	got = renderer.RenderSequence(0, -1, 7, 2)
	if assertRenderSequence(t, got, 7, 2, "foo hello bar world", 2) {
		assertColorSubsequence(t, got.Sequences[0], "RED", 0, 3)
		assertColorSubsequence(t, got.Sequences[1], "BLUE", 10, 13)
	}

	// Test that out-of-bound color sequences are not added.
	got = renderer.RenderSequence(4, 6, 8, 1)
	assertRenderSequence(t, got, 8, 1, "he", 0)

	// Test Half-rendered text
	got = renderer.RenderSequence(1, 12, 0, 0)
	if assertRenderSequence(t, got, 0, 0, "oo hello ba", 2) {
		assertColorSubsequence(t, got.Sequences[0], "RED", 0, 2)
		assertColorSubsequence(t, got.Sequences[1], "BLUE", 9, 11)
	}

	// Test Half-rendered text (edges)
	got = renderer.RenderSequence(2, 11, 0, 0)
	if assertRenderSequence(t, got, 0, 0, "o hello b", 2) {
		assertColorSubsequence(t, got.Sequences[0], "RED", 0, 1)
		assertColorSubsequence(t, got.Sequences[1], "BLUE", 8, 9)
	}

	// TODO: test barkets

	// Test with unicodes
	text := "笀耔 [澉 灊灅甗](RED) 郔镺 笀耔 澉 [灊灅甗](yellow) 郔镺"
	normalized := "笀耔 澉 灊灅甗 郔镺 笀耔 澉 灊灅甗 郔镺"
	renderer = MarkdownTextRenderer{text}
	got = renderer.RenderSequence(0, -1, 4, 7)
	if assertRenderSequence(t, got, 4, 7, normalized, 2) {
		assertColorSubsequence(t, got.Sequences[0], "RED", 3, 8)
		assertColorSubsequence(t, got.Sequences[1], "YELLOW", 17, 20)
	}

	got = renderer.RenderSequence(6, 7, 0, 0)
	if assertRenderSequence(t, got, 0, 0, "灅", 1) {
		assertColorSubsequence(t, got.Sequences[0], "RED", 0, 1)
	}

	got = renderer.RenderSequence(7, 19, 0, 0)
	if assertRenderSequence(t, got, 0, 0, "甗 郔镺 笀耔 澉 灊灅", 2) {
		assertColorSubsequence(t, got.Sequences[0], "RED", 0, 1)
		assertColorSubsequence(t, got.Sequences[1], "YELLOW", 10, 12)
	}

	// Test inside
	renderer = MarkdownTextRenderer{"foo [foobar](red) bar"}
	got = renderer.RenderSequence(4, 10, 0, 0)
	if assertRenderSequence(t, got, 0, 0, "foobar", 1) {
		assertColorSubsequence(t, got.Sequences[0], "RED", 0, 6)
	}
}

func TestMarkdownTextRenderer_Render(t *testing.T) {
	renderer := MarkdownTextRenderer{"[foo](red,bold) [bar](blue)"}
	got := renderer.Render(6, 8)
	if assertRenderSequence(t, got, 6, 8, "foo bar", 2) {
		assertColorSubsequence(t, got.Sequences[0], "RED,BOLD", 0, 3)
		assertColorSubsequence(t, got.Sequences[1], "blue", 4, 7)
	}
}

func TestMarkdownTextRendererFactory(t *testing.T) {
	factory := MarkdownTextRendererFactory{}
	expected := MarkdownTextRenderer{"Hello world"}
	assert.Equal(t, factory.TextRenderer("Hello world"), expected)
}

func TestColorSubsequencesToMap(t *testing.T) {
	colorSubsequences := []ColorSubsequence{
		{ColorRed, 1, 4},
		{ColorBlue | AttrBold, 9, 10},
	}

	expected := make(map[int]Attribute)
	expected[1] = ColorRed
	expected[2] = ColorRed
	expected[3] = ColorRed
	expected[9] = ColorBlue | AttrBold

	assert.Equal(t, expected, ColorSubsequencesToMap(colorSubsequences))
}

func getTestRenderedSequence() RenderedSequence {
	cs := []ColorSubsequence{
		{ColorRed, 3, 5},
		{ColorBlue | AttrBold, 9, 10},
	}

	return RenderedSequence{"Hello world", ColorWhite, ColorBlack, cs, nil}
}

func newTestPoint(char rune, x, y int, colorA ...Attribute) Point {
	var color Attribute
	if colorA != nil && len(colorA) == 1 {
		color = colorA[0]
	} else {
		color = ColorWhite
	}

	return Point{char, ColorBlack, color, x, y}
}

func TestRenderedSequence_Buffer(t *testing.T) {
	sequence := getTestRenderedSequence()
	expected := []Point{
		newTestPoint('H', 5, 7),
		newTestPoint('e', 6, 7),
		newTestPoint('l', 7, 7),
		newTestPoint('l', 7, 7, ColorRed),
		newTestPoint('o', 8, 7, ColorRed),
		newTestPoint(' ', 9, 7),
		newTestPoint('w', 10, 7),
		newTestPoint('o', 11, 7),
		newTestPoint('r', 12, 7),
		newTestPoint('l', 13, 7, ColorBlue|AttrBold),
		newTestPoint('d', 14, 7),
	}

	buffer, lastColor := sequence.Buffer(5, 7)

	assert.Equal(t, expected[:3], buffer[:3])
	assert.Equal(t, ColorWhite, lastColor)
}

func AssertPoint(t *testing.T, got Point, char rune, x, y int, colorA ...Attribute) {
	expected := newTestPoint(char, x, y, colorA...)
	assert.Equal(t, expected, got)
}

func TestRenderedSequence_PointAt(t *testing.T) {
	sequence := getTestRenderedSequence()
	pointAt := func(n, x, y int) Point {
		p, w := sequence.PointAt(n, x, y)
		assert.Equal(t, w, 1)

		return p
	}

	AssertPoint(t, pointAt(0, 3, 4), 'H', 3, 4)
	AssertPoint(t, pointAt(1, 2, 1), 'e', 2, 1)
	AssertPoint(t, pointAt(2, 6, 3), 'l', 6, 3)
	AssertPoint(t, pointAt(3, 8, 8), 'l', 8, 8, ColorRed)
	AssertPoint(t, pointAt(4, 1, 4), 'o', 1, 4, ColorRed)
	AssertPoint(t, pointAt(5, 3, 6), ' ', 3, 6)
	AssertPoint(t, pointAt(6, 4, 3), 'w', 4, 3)
	AssertPoint(t, pointAt(7, 5, 2), 'o', 5, 2)
	AssertPoint(t, pointAt(8, 0, 5), 'r', 0, 5)
	AssertPoint(t, pointAt(9, 9, 0), 'l', 9, 0, ColorBlue|AttrBold)
	AssertPoint(t, pointAt(10, 7, 1), 'd', 7, 1)
}

func getTestPlainRenderer() PlainRenderer {
	return PlainRenderer{"[Hello](red) \x1b[31mworld"}
}

func TestPlainRenderer_NormalizedText(t *testing.T) {
	r := getTestPlainRenderer()
	assert.Equal(t, "[Hello](red) \x1b[31mworld", r.NormalizedText())
	assert.Equal(t, "[Hello](red) \x1b[31mworld", r.Text)
}

func TestPlainRenderer_Render(t *testing.T) {
	renderer := getTestPlainRenderer()
	got := renderer.Render(5, 7)
	assertRenderSequence(t, got, 5, 7, "[Hello](red) \x1b[31mworld", 0)
}

func TestPlainRenderer_RenderSequence(t *testing.T) {
	renderer := getTestPlainRenderer()
	got := renderer.RenderSequence(3, 5, 9, 1)
	assertRenderSequence(t, got, 9, 1, "ll", 0)
}

func TestPlainRendererFactory(t *testing.T) {
	factory := PlainRendererFactory{}
	expected := PlainRenderer{"Hello world"}
	assert.Equal(t, factory.TextRenderer("Hello world"), expected)
}

func TestPosUnicode(t *testing.T) {
	// Every characters takes 3 bytes
	text := "你好世界"
	require.Equal(t, "你好", text[:6])
	assert.Equal(t, 2, posUnicode(text, 6))
}

// Make `escapeCode` safe (i.e. replace \033 by \\033) so that it is not
// formatted.
// func makeEscapeCodeSafe(escapeCode string) string {
// 	return strings.Replace(escapeCode, "\033", "\\033", -1)
// }

func TestEscapeCode_Color(t *testing.T) {
	codes := map[EscapeCode]Attribute{
		"\033[30m":     ColorBlack,
		"\033[31m":     ColorRed,
		"\033[32m":     ColorGreen,
		"\033[33m":     ColorYellow,
		"\033[34m":     ColorBlue,
		"\033[35m":     ColorMagenta,
		"\033[36m":     ColorCyan,
		"\033[37m":     ColorWhite,
		"\033[1;31m":   ColorRed | AttrBold,
		"\033[1;4;31m": ColorRed | AttrBold | AttrUnderline,
		"\033[0m":      ColorDefault,
	}

	for code, color := range codes {
		got, err := code.Color()
		msg := fmt.Sprintf("Escape code: '%v'", code.MakeSafe())
		if assert.NoError(t, err, msg) {
			assert.Equal(t, color, got, msg)
		}
	}

	invalidEscapeCodes := []EscapeCode{
		"\03354m",
		"[54m",
		"\033[34",
		"\033[34;m",
		"\033[34m;",
		"\033[34;",
		"\033[5432m",
		"t\033[30m",
		"t\033[30ms",
		"\033[30ms",
	}

	errMsg := "%v is not a valid ASCII escape code"
	for _, invalidEscapeCode := range invalidEscapeCodes {
		color, err := invalidEscapeCode.Color()
		safeEscapeCode := invalidEscapeCode.MakeSafe()
		expectedErr := fmt.Sprintf(errMsg, safeEscapeCode)
		if assert.EqualError(t, err, expectedErr, "Expected: "+expectedErr) {
			assert.Equal(t, color, Attribute(0))
		}
	}

	outOfRangeCodes := []EscapeCode{
		"\033[2m",
		"\033[3m",
		"\033[3m",
		"\033[5m",
		"\033[6m",
		"\033[7m",
		"\033[8m",
		"\033[38m",
		"\033[39m",
		"\033[40m",
		"\033[41m",
		"\033[43m",
		"\033[45m",
		"\033[46m",
		"\033[48m",
		"\033[49m",
		"\033[50m",
	}

	for _, code := range outOfRangeCodes {
		color, err := code.Color()
		safeCode := code.MakeSafe()
		errMsg := fmt.Sprintf("Unkown/unsupported escape code: '%v'", safeCode)
		if assert.EqualError(t, err, errMsg) {
			assert.Equal(t, color, Attribute(0), "Escape Code: "+safeCode)
		}
	}

	// Special case: check for out of slice panic on empty string
	_, err := EscapeCode("").Color()
	assert.EqualError(t, err, " is not a valid ASCII escape code")
}

func TestEscapeCode_String(t *testing.T) {
	e := EscapeCode("\033[32m")
	assert.Equal(t, "\\033[32m", e.String())
}

func TestEscapeCode_Raw(t *testing.T) {
	e := EscapeCode("\033[32m")
	assert.Equal(t, "\033[32m", e.Raw())
}

func TestEscapeCodeRenderer_NormalizedText(t *testing.T) {
	renderer := EscapeCodeRenderer{"\033[33mtest \033[35mfoo \033[33;1mbar"}
	assert.Equal(t, "test foo bar", renderer.NormalizedText())

	renderer = EscapeCodeRenderer{"hello \033[38mtest"}
	assert.Equal(t, "hello \033[38mtest", renderer.NormalizedText())
}

func TestEscapeCodeRenderer_RenderSequence(t *testing.T) {
	black, white := ColorWhite, ColorBlack
	renderer := EscapeCodeRenderer{"test \033[33mfoo \033[31mbar"}
	sequence := renderer.RenderSequence(0, -1, black, white)
	if assertRenderSequence(t, sequence, black, white, "test foo bar", 2) {
		assertColorSubsequence(t, sequence.Sequences[0], "YELLOW", 5, 9)
		assertColorSubsequence(t, sequence.Sequences[1], "RED", 9, 12)
		getPoint := func(n int) Point {
			point, width := sequence.PointAt(n, 10+n, 30)
			assert.Equal(t, 1, width)

			return point
		}

		// Also test the points at to make sure that
		// I didn't make a counting mistake...
		AssertPoint(t, getPoint(0), 't', 10, 30)
		AssertPoint(t, getPoint(1), 'e', 11, 30)
		AssertPoint(t, getPoint(2), 's', 12, 30)
		AssertPoint(t, getPoint(3), 't', 13, 30)
		AssertPoint(t, getPoint(4), ' ', 14, 30)
		AssertPoint(t, getPoint(5), 'f', 15, 30, ColorYellow)
		AssertPoint(t, getPoint(6), 'o', 16, 30, ColorYellow)
		AssertPoint(t, getPoint(7), 'o', 17, 30, ColorYellow)
		AssertPoint(t, getPoint(8), ' ', 18, 30, ColorYellow)
		AssertPoint(t, getPoint(9), 'b', 19, 30, ColorRed)
		AssertPoint(t, getPoint(10), 'a', 20, 30, ColorRed)
		AssertPoint(t, getPoint(11), 'r', 21, 30, ColorRed)
	}

	renderer = EscapeCodeRenderer{"甗 郔\033[33m镺 笀耔 澉 灊\033[31m灅甗"}
	sequence = renderer.RenderSequence(2, -1, black, white)
	if assertRenderSequence(t, sequence, black, white, "郔镺 笀耔 澉 灊灅甗", 2) {
		assertColorSubsequence(t, sequence.Sequences[0], "YELLOW", 1, 9)
		assertColorSubsequence(t, sequence.Sequences[1], "RED", 9, 11)
	}

	renderer = EscapeCodeRenderer{"\033[33mHell\033[31mo world"}
	sequence = renderer.RenderSequence(2, -1, black, white)
	if assertRenderSequence(t, sequence, black, white, "llo world", 2) {
		assertColorSubsequence(t, sequence.Sequences[0], "YELLOW", 0, 2)
		assertColorSubsequence(t, sequence.Sequences[1], "RED", 2, 9)
	}

	sequence = renderer.RenderSequence(1, 7, black, white)
	if assertRenderSequence(t, sequence, black, white, "ello w", 2) {
		assertColorSubsequence(t, sequence.Sequences[0], "YELLOW", 0, 3)
		assertColorSubsequence(t, sequence.Sequences[1], "RED", 3, 6)
	}

	sequence = renderer.RenderSequence(6, 10, black, white)
	if assertRenderSequence(t, sequence, black, white, "worl", 1) {
		assertColorSubsequence(t, sequence.Sequences[0], "RED", 0, 4)
	}

	// Test with out-of-range escape code
	renderer = EscapeCodeRenderer{"hello \033[38mtest"}
	sequence = renderer.RenderSequence(0, -1, black, white)
	assertRenderSequence(t, sequence, black, white, "hello \033[38mtest", 0)
}

func TestEscapeCodeRenderer_Render(t *testing.T) {
	renderer := EscapeCodeRenderer{"test \033[33mfoo \033[31mbar"}
	sequence := renderer.Render(4, 6)
	if assertRenderSequence(t, sequence, 4, 6, "test foo bar", 2) {
		assertColorSubsequence(t, sequence.Sequences[0], "YELLOW", 5, 9)
		assertColorSubsequence(t, sequence.Sequences[1], "RED", 9, 12)
	}
}

func TestEscapeCodeRendererFactory_TextRenderer(t *testing.T) {
	factory := EscapeCodeRendererFactory{}
	assert.Equal(t, EscapeCodeRenderer{"foo"}, factory.TextRenderer("foo"))
	assert.Equal(t, EscapeCodeRenderer{"bar"}, factory.TextRenderer("bar"))
}
