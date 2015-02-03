package termui

import tm "github.com/nsf/termbox-go"

type P struct {
	Div
	Text        string
	TextFgColor tm.Attribute
	TextBgColor tm.Attribute
}

func NewP(s string) P {
	return P{Div: NewDiv(), Text: s}
}

func (p P) Buffer() []Point {
	ps := p.Div.Buffer()

	(&p).sync()

	rs := str2runes(p.Text)
	i, j, k := 0, 0, 0
	for i < p.innerHeight && k < len(rs) {
		if rs[k] == '\n' || j == p.innerWidth {
			i++
			j = 0
			if rs[k] == '\n' {
				k++
			}
			continue
		}
		pi := Point{}
		pi.X = p.innerX + j
		pi.Y = p.innerY + i

		pi.Code.Ch = rs[k]
		pi.Code.Bg = p.TextBgColor
		pi.Code.Fg = p.TextFgColor

		ps = append(ps, pi)
		k++
		j++
	}
	return ps
}
