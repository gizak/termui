package termui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownTextRenderer_NormalizedText(t *testing.T) {
	renderer := MarkdownTextRenderer{}

	got := renderer.NormalizedText("[ERROR](red,bold) Something went wrong")
	assert.Equal(t, got, "ERROR Something went wrong")

	got = renderer.NormalizedText("[foo](g) hello [bar](green) world")
	assert.Equal(t, got, "foo hello bar world")

	got = renderer.NormalizedText("[foo](g) hello [bar]green (world)")
	assert.Equal(t, got, "foo hello [bar]green (world)")
}
