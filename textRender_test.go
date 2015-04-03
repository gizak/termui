package termui

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownTextRenderer_normalizeText(t *testing.T) {
	renderer := MarkdownTextRenderer{}

	got := renderer.normalizeText("[ERROR](red,bold) Something went wrong")
	assert.Equal(t, got, "ERROR Something went wrong")

	got = renderer.normalizeText("[foo](red) hello [bar](green) world")
	assert.Equal(t, got, "foo hello bar world")

	got = renderer.normalizeText("[foo](g) hello [bar]green (world)")
	assert.Equal(t, got, "foo hello [bar]green (world)")

	// FIXME: [[ERROR]](red,bold) test should normalize to:
	// [ERROR] test
	// FIXME: Support unicode inside the error message.
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

	// Test half-rendered text (unicode)
	// FIXME: Add

	// Test inside
	renderer = MarkdownTextRenderer{"foo [foobar](red) bar"}
	got = renderer.RenderSequence(4, 10, 0, 0)
	if assertRenderSequence(t, got, 0, 0, "foobar", 1) {
		assertColorSubsequence(t, got.Sequences[0], "RED", 0, 6)
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

func TestRenderedSequence_Buffer(t *testing.T) {
	cs := []ColorSubsequence{
		{ColorRed, 3, 5},
		{ColorBlue | AttrBold, 9, 10},
	}
	sequence := RenderedSequence{"Hello world", ColorWhite, ColorBlack, cs}
	newPoint := func(char string, x, y int, colorA ...Attribute) Point {
		var color Attribute
		if colorA != nil && len(colorA) == 1 {
			color = colorA[0]
		} else {
			color = ColorWhite
		}

		return Point{[]rune(char)[0], ColorBlack, color, x, y}
	}

	expected := []Point{
		newPoint("H", 5, 7),
		newPoint("e", 6, 7),
		newPoint("l", 7, 7),
		newPoint("l", 7, 7, ColorRed),
		newPoint("o", 8, 7, ColorRed),
		newPoint(" ", 9, 7),
		newPoint("w", 10, 7),
		newPoint("o", 11, 7),
		newPoint("r", 12, 7),
		newPoint("l", 13, 7, ColorBlue|AttrBold),
		newPoint("d", 14, 7),
	}
	buffer, lastColor := sequence.Buffer(5, 7)

	assert.Equal(t, expected[:3], buffer[:3])
	assert.Equal(t, ColorWhite, lastColor)
}
