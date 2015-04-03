package termui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getMDRenderer() MarkdownTextRenderer {
	return MarkdownTextRenderer{}
}

func TestMarkdownTextRenderer_NormalizedText(t *testing.T) {
	renderer := getMDRenderer()

	got := renderer.NormalizedText("[ERROR](red,bold) Something went wrong")
	assert.Equal(t, got, "ERROR Something went wrong")

	got = renderer.NormalizedText("[foo](red) hello [bar](green) world")
	assert.Equal(t, got, "foo hello bar world")

	got = renderer.NormalizedText("[foo](g) hello [bar]green (world)")
	assert.Equal(t, got, "foo hello [bar]green (world)")

	// FIXME: [[ERROR]](red,bold) test should normalize to:
	// [ERROR] test
}

func assertRenderSequence(t *testing.T, sequence RenderedSequence, last, background Attribute, text string, lenSequences int) {
	assert.Equal(t, last, sequence.LastColor)
	assert.Equal(t, background, sequence.BackgroundColor)
	assert.Equal(t, text, sequence.NormalizedText)
	assert.Equal(t, lenSequences, len(sequence.Sequences))
}

func assertColorSubsequence(t *testing.T, s ColorSubsequence, color string, start, end int) {
	assert.Equal(t, ColorSubsequence{StringToAttribute(color), start, end}, s)
}

func TestMarkdownTextRenderer_RenderSequence(t *testing.T) {
	renderer := getMDRenderer()

	got := renderer.RenderSequence("[ERROR](red,bold) something went wrong", 3, 5)
	assertRenderSequence(t, got, 3, 5, "ERROR something went wrong", 1)
	assertColorSubsequence(t, got.Sequences[0], "RED,BOLD", 0, 5)

	got = renderer.RenderSequence("[foo](red) hello [bar](green) world", 7, 2)
	assertRenderSequence(t, got, 3, 2, "foo hello bar world", 2)

	assertColorSubsequence(t, got.Sequences[0], "RED", 0, 3)
	assertColorSubsequence(t, got.Sequences[1], "GREEN", 10, 13)
}
