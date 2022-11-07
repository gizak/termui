package termui

import (
	"strings"
	"testing"
)

func TestPrepareStyles(t *testing.T) {
	items := PrepareStyles("test [blue](fg:blue,bg:white,mod:bold) and [red](fg:red)")
	if len(items) != 4 {
		t.Fatal("wrong length", len(items))
	}
	if items[0].Text != "test " {
		t.Fatal("wrong text", items[0].Text)
	}
	if items[0].Style != "" {
		t.Fatal("wrong text", items[0].Style)
	}
	if items[1].Text != "blue" {
		t.Fatal("wrong text", items[1].Text)
	}
	if items[1].Style != "fg:blue,bg:white,mod:bold" {
		t.Fatal("wrong text", items[1].Style)
	}
	if items[2].Text != " and " {
		t.Fatal("wrong text", items[2].Text)
	}
	if items[2].Style != "" {
		t.Fatal("wrong text", items[2].Style)
	}
	if items[3].Text != "red" {
		t.Fatal("wrong text", items[3].Text)
	}
	if items[3].Style != "fg:red" {
		t.Fatal("wrong text", items[3].Style)
	}
}

func TestParseStyles(t *testing.T) {
	cells := ParseStyles("test nothing", NewStyle(ColorWhite))
	cells = ParseStyles("test [blue](fg:blue,bg:white,mod:bold) and [red](fg:red)", NewStyle(ColorWhite))
	if len(cells) != 17 {
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
	for i := 5; i < 9; i++ {
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
	if text != "test blue and red" {
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
