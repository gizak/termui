package termui

import "testing"

func TestPar_NoBorderBackground(t *testing.T) {
	par := NewPar("a")
	par.Border = false
	par.Bg = ColorBlue
	par.TextBgColor = ColorBlue
	par.Width = 2
	par.Height = 2

	pts := par.Buffer()
	for _, p := range pts.CellMap {
		t.Log(p)
		if p.Bg != par.Bg {
			t.Errorf("expected color to be %v but got %v", par.Bg, p.Bg)
		}
	}
}
