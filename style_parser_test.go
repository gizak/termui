package termui

import (
	"fmt"
	"strings"
	"testing"
)

func TestFindStylePositions(t *testing.T) {
	items := FindStylePositions("test [blue](fg:blue,mod:bold) and [red](fg:red) and maybe even [foo](bg:red)!")
	if len(items) != 3 {
		t.Fatal("wrong length", len(items))
	}
	if items[0] != 10 {
		t.Fatal("wrong index", items[0])
	}
	if items[1] != 38 {
		t.Fatal("wrong index", items[1])
	}
	if items[2] != 67 {
		t.Fatal("wrong index", items[2])
	}
}

func TestBreakByStylesComplex(t *testing.T) {
	items := BreakByStyles("test [blue](fg:blue,mod:bold) and [red](fg:red) and maybe even [foo](bg:red)!")
	if len(items) != 10 {
		t.Fatal("wrong length", len(items))
	}
	text := strings.Join(items, " ")
	if text != "test  blue fg:blue,bg:white,mod:bold  and  red fg:red  and maybe even  foo bg:red !" {
		t.Fatal("wrong text", text)
	}
}
func TestBreakByStylesSimpler(t *testing.T) {
	items := BreakByStyles("[blue](fg:blue) [1]")
	// [blue](fg:blue) [1]
	// 012345678901234
	// 0-14, rest
	fmt.Println(items, len(items))
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
