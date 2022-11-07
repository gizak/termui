package termui

import (
	"strings"
	"testing"
)

func TestParseStyles(t *testing.T) {
	cells := ParseStyles("test [blue](fg:blue,bg:white,mod:bold)", NewStyle(ColorWhite))
	if len(cells) != 9 {
		t.Fatal("wrong length", len(cells))
	}
	for i := 0; i < 5; i++ {
		if cells[i].Style.Fg != ColorWhite {
			t.Fatal("wrong fg color", cells[i])
		}
		if cells[i].Style.Bg != ColorClear {
			t.Fatal("wrong bg color", cells[i])
		}
		if cells[i].Style.Modifier != ModifierClear {
			t.Fatal("wrong mod", cells[i])
		}
	}
	for i := 5; i < len(cells); i++ {
		if cells[i].Style.Fg != ColorBlue {
			t.Fatal("wrong fg color", cells[i])
		}
		if cells[i].Style.Bg != ColorWhite {
			t.Fatal("wrong bg color", cells[i])
		}
		if cells[i].Style.Modifier != ModifierBold {
			t.Fatal("wrong mod", cells[i])
		}
	}

	text := textFromCells(cells)
	if text != "test blue" {
		t.Fatal("wrong text", text)
	}

	cells = ParseStyles("[blue](fg:blue) [1]", NewStyle(ColorWhite))
	text = textFromCells(cells)
	if text != "blue [1]" {
		t.Fatal("wrong text", text)
	}

	cells = ParseStyles("[0]", NewStyle(ColorWhite))
	text = textFromCells(cells)
	if text != "[0]" {
		t.Fatal("wrong text", text)
	}

	cells = ParseStyles("[", NewStyle(ColorWhite))
	text = textFromCells(cells)
	if text != "[" {
		t.Fatal("wrong text", text)
	}

}

func textFromCells(cells []Cell) string {
	buff := []string{}
	for _, cell := range cells {
		buff = append(buff, string(cell.Rune))
	}
	return strings.Join(buff, "")
}
