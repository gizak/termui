package termui

import "testing"

func TestPar_NoBorderBackground(t *testing.T) {
	par := NewPar("a")
	par.HasBorder = false
	par.BgColor = ColorBlue
	par.TextBgColor = ColorBlue
	par.Width = 2
	par.Height = 2

	pts := par.Buffer()
	for _, p := range pts {
		t.Log(p)
		if p.Bg != par.BgColor {
			t.Errorf("expected color to be %v but got %v", par.BgColor, p.Bg)
		}
	}
}
