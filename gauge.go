package termui

import "strconv"

type Gauge struct {
	Block
	Percent      int
	BarColor     Attribute
	PercentColor Attribute
}

func NewGauge() *Gauge {
	g := &Gauge{Block: *NewBlock(), PercentColor: ColorWhite, BarColor: ColorGreen}
	g.Width = 12
	g.Height = 5
	return g
}

func (g *Gauge) Buffer() []Point {
	ps := g.Block.Buffer()

	w := g.Percent * g.innerWidth / 100
	s := strconv.Itoa(g.Percent) + "%"
	rs := str2runes(s)

	prx := g.innerX + g.innerWidth/2 - 1
	pry := g.innerY + g.innerHeight/2

	// plot bar
	for i := 0; i < g.innerHeight; i++ {
		for j := 0; j < w; j++ {
			p := Point{}
			p.X = g.innerX + j
			p.Y = g.innerY + i
			p.Code.Ch = ' '
			p.Code.Bg = toTmAttr(g.BarColor)
			ps = append(ps, p)
		}
	}

	// plot percentage
	for i, v := range rs {
		p := Point{}
		p.X = prx + i
		p.Y = pry
		p.Code.Ch = v
		p.Code.Fg = toTmAttr(g.PercentColor)
		if w > g.innerWidth/2-1+i {
			p.Code.Bg = toTmAttr(g.BarColor)
		} else {
			p.Code.Bg = toTmAttr(g.Block.BgColor)
		}
		ps = append(ps, p)
	}
	return ps
}
