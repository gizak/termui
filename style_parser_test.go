package termui

import (
	"testing"
)

func TestParseStyles(t *testing.T) {
	cells := ParseStyles("test [blue](fg:blue)", NewStyle(ColorWhite))
	if cells[0].Style.Fg != 7 {
		t.Fatal("normal text fg is wrong")
	}
	if cells[5].Style.Fg != 4 {
		t.Fatal("blue text fg is wrong")
	}

}
