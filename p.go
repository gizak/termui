package termui

type Par struct {
	Block
	Text        string
	TextFgColor Attribute
	TextBgColor Attribute
}

func NewPar(s string) *Par {
	return &Par{
		Block:       *NewBlock(),
		Text:        s,
		TextFgColor: theme.ParTextFg,
		TextBgColor: theme.ParTextBg}
}

func (p *Par) Buffer() []Point {
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

			if i >= p.innerHeight {
				ps = append(ps, newPointWithAttrs('…',
					p.innerX+p.innerWidth-1,
					p.innerY+p.innerHeight-1,
					p.TextFgColor, p.TextBgColor))
				break
			}

			continue
		}
		pi := Point{}
		pi.X = p.innerX + j
		pi.Y = p.innerY + i

		pi.Ch = rs[k]
		pi.Bg = p.TextBgColor
		pi.Fg = p.TextFgColor

		ps = append(ps, pi)

		k++
		j++
	}
	return ps
}
