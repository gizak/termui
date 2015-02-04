package termui

type P struct {
	Block
	Text        string
	TextFgColor Attribute
	TextBgColor Attribute
}

func NewP(s string) *P {
	return &P{Block: *NewBlock(), Text: s}
}

func (p *P) Buffer() []Point {
	ps := p.Block.Buffer()

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
		pi.Code.Bg = toTmAttr(p.TextBgColor)
		pi.Code.Fg = toTmAttr(p.TextFgColor)

		ps = append(ps, pi)
		k++
		j++
	}
	return ps
}
