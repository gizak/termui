package termui

import (
	"testing"
	"strings"
)

// TODO: More tests!

const TESTING_LINE = "testing"

func TestInput_SingleLine_NoNewLines(t *testing.T) {
	input := NewInput("", false)
	input.addString(TESTING_LINE)
	input.addString(NEW_LINE)

	if input.Text() != TESTING_LINE {
		t.Errorf("Expected test to only contains %s (%d) but found %s (%d)",
			TESTING_LINE, len(TESTING_LINE), input.Text(), len(input.Text()))
	}

	if strings.HasSuffix(input.Text(), NEW_LINE) {
		t.Error("Unexpected newline at the end of TEXT")
	}
}

func TestInput_MultiLine_SingleEntry(t *testing.T) {
	input := NewInput("", true)
	input.addString(TESTING_LINE)

	if len(input.Lines()) != 1 {
		t.Errorf("Invalid number of lines in input, expected 1 but found %d", len(input.Lines()))
	}
	if input.Text() != TESTING_LINE {
		t.Errorf("Expected test to only contains %s (%d) but found %s (%d)",
			TESTING_LINE, len(TESTING_LINE), input.Text(), len(input.Text()))
	}
}

func TestInput_MultiLine_ShiftLineDown(t *testing.T) {
	input := NewInput("", true)
	input.addString(TESTING_LINE)
	input.addString(NEW_LINE)
	input.addString(TESTING_LINE)

	if len(input.Lines()) != 2 {
		t.Errorf("Expected 2 lines in input but found %d", len(input.Lines()))
	}
	if input.cursorLineIndex != 1 {
		t.Errorf("Expected line cursor to be at index 1, found it at %d", input.cursorLineIndex)
	}

	input.moveUp()
	input.addString(NEW_LINE)

	if len(input.Lines()) != 3 {
		t.Errorf("Expected 3 lines in input but found %d", len(input.Lines()))
	}
	if input.Lines()[0] != "" {
		t.Errorf("Expected first line to be blank but found %s", input.Lines()[0])
	}
}

func TestInput_MultiLine_MoveLeftToPreviousLine(t *testing.T) {
	input := NewInput("", true)
	input.addString(TESTING_LINE)
	input.addString(NEW_LINE)
	input.addString(TESTING_LINE)

	if input.cursorLinePos != len(TESTING_LINE) {
		t.Errorf("Expected cursor to be at position %d, found it at %d", len(TESTING_LINE), input.cursorLinePos)
	}

	// reset to 0 and move left
	input.cursorLinePos = 0
	input.moveLeft()

	if input.cursorLineIndex != 0 {
		t.Errorf("Expcted line cursor to be at index 0, found it at %d", input.cursorLineIndex)
	}
	if input.cursorLinePos != len(TESTING_LINE) {
		t.Errorf("Expected cursor to be at position %d, found it at %d", len(TESTING_LINE), input.cursorLinePos)
	}
}