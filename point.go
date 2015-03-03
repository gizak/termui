package termui

type Point struct {
	Ch rune
	Bg Attribute
	Fg Attribute
	X  int
	Y  int
}

func newPoint(c rune, x, y int) (p Point) {
	p.Ch = c
	p.X = x
	p.Y = y
	return
}

func newPointWithAttrs(c rune, x, y int, fg, bg Attribute) Point {
	p := newPoint(c, x, y)
	p.Bg = bg
	p.Fg = fg
	return p
}
