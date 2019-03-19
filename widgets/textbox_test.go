package widgets

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

//TestGetRawText test simple string
func TestGetRawText(t *testing.T) {
	text := "My Sample RawText"
	tb := NewTextBox()
	tb.SetText(text)

	assert.Equal(t, text, tb.RawText())
}

//TestGetRawTextWithLBs test line breaks in the text
func TestGetRawTextWithLBs(t *testing.T) {
	text := `My Sample RawText
				with	
				line 
				breaks`

	tb := NewTextBox()
	tb.SetText(text)

	assert.Equal(t, text, tb.RawText())
}

//TestGetStyledText test styled text
func TestGetStyledText(t *testing.T) {
	text := "[red text](fg:red,mod:bold) more text [blue text](fg:blue,mod:bold) a bit more"

	tb := NewTextBox()
	tb.SetText(text)

	assert.Equal(t, text, tb.Text())
}

//TestGetStyledText2 test styled text ending in a styled string
func TestGetStyledText2(t *testing.T) {
	text := "[red text](fg:red,mod:bold) more text [blue text](fg:blue,mod:bold) a [bit more](fg:green)"

	tb := NewTextBox()
	tb.SetText(text)

	assert.Equal(t, text, tb.Text())
}

//TestGetStyledTextWithLBs test styled text with line breaks
func TestGetStyledTextWithLBs(t *testing.T) {
	text := `[red text](fg:red,mod:bold) more 
 text [blue text](fg:blue,mod:bold) a [bit more](fg:green)`

	tb := NewTextBox()
	tb.SetText(text)

	assert.Equal(t, text, tb.Text())
}