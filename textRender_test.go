package termui

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextRender_TestInterface(t *testing.T) {
	var inter *TextRender

	assert.Implements(t, inter, new(MarkdownTextRenderer))
	assert.Implements(t, inter, new(NoopRenderer))
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

	got = renderer.normalizeText("[[foo]](red,white) bar")
	assert.Equal(t, renderer.normalizeText(got), "[foo] bar")
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

func getTestNoopRenderer() NoopRenderer {
	return NoopRenderer{"[Hello](red) \x1b[31mworld"}
}

func TestNoopRenderer_NormalizedText(t *testing.T) {
	r := getTestNoopRenderer()
	assert.Equal(t, "[Hello](red) \x1b[31mworld", r.NormalizedText())
	assert.Equal(t, "[Hello](red) \x1b[31mworld", r.Text)
}

func TestNoopRenderer_Render(t *testing.T) {
	renderer := getTestNoopRenderer()
	got := renderer.Render(5, 7)
	assertRenderSequence(t, got, 5, 7, "[Hello](red) \x1b[31mworld", 0)
}

func TestNoopRenderer_RenderSequence(t *testing.T) {
	renderer := getTestNoopRenderer()
	got := renderer.RenderSequence(3, 5, 9, 1)
	assertRenderSequence(t, got, 9, 1, "ll", 0)
}

func TestPosUnicode(t *testing.T) {
	// Every characters takes 3 bytes
	text := "你好世界"
	require.Equal(t, "你好", text[:6])
	assert.Equal(t, 2, posUnicode(text, 6))
}
