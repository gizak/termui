package termui

const TOP_RIGHT = '┐'
const VERTICAL_LINE = '│'
const HORIZONTAL_LINE = '─'
const TOP_LEFT = '┌'
const BOTTOM_RIGHT = '┘'
const BOTTOM_LEFT = '└'

type border struct {
	X       int
	Y       int
	Width   int
	Height  int
	FgColor Attribute
	BgColor Attribute
}

type hline struct {
	X       int
	Y       int
	Length  int
	FgColor Attribute
	BgColor Attribute
}

type vline struct {
	X       int
	Y       int
	Length  int
	FgColor Attribute
	BgColor Attribute
}

func (l hline) Buffer() []Point {
	pts := make([]Point, l.Length)
	for i := 0; i < l.Length; i++ {
		pts[i].X = l.X + i
		pts[i].Y = l.Y
		pts[i].Code.Ch = HORIZONTAL_LINE
		pts[i].Code.Bg = toTmAttr(l.BgColor)
		pts[i].Code.Fg = toTmAttr(l.FgColor)
	}
	return pts
}

func (l vline) Buffer() []Point {
	pts := make([]Point, l.Length)
	for i := 0; i < l.Length; i++ {
		pts[i].X = l.X
		pts[i].Y = l.Y + i
		pts[i].Code.Ch = VERTICAL_LINE
		pts[i].Code.Bg = toTmAttr(l.BgColor)
		pts[i].Code.Fg = toTmAttr(l.FgColor)
	}
	return pts
}

func (b border) Buffer() []Point {
	if b.Width < 2 || b.Height < 2 {
		return nil
	}
	pts := make([]Point, 2*b.Width+2*b.Height-4)

	pts[0].X = b.X
	pts[0].Y = b.Y
	pts[0].Code.Fg = toTmAttr(b.FgColor)
	pts[0].Code.Bg = toTmAttr(b.BgColor)
	pts[0].Code.Ch = TOP_LEFT

	pts[1].X = b.X + b.Width - 1
	pts[1].Y = b.Y
	pts[1].Code.Fg = toTmAttr(b.FgColor)
	pts[1].Code.Bg = toTmAttr(b.BgColor)
	pts[1].Code.Ch = TOP_RIGHT

	pts[2].X = b.X
	pts[2].Y = b.Y + b.Height - 1
	pts[2].Code.Fg = toTmAttr(b.FgColor)
	pts[2].Code.Bg = toTmAttr(b.BgColor)
	pts[2].Code.Ch = BOTTOM_LEFT

	pts[3].X = b.X + b.Width - 1
	pts[3].Y = b.Y + b.Height - 1
	pts[3].Code.Fg = toTmAttr(b.FgColor)
	pts[3].Code.Bg = toTmAttr(b.BgColor)
	pts[3].Code.Ch = BOTTOM_RIGHT

	copy(pts[4:], (hline{b.X + 1, b.Y, b.Width - 2, b.FgColor, b.BgColor}).Buffer())
	copy(pts[4+b.Width-2:], (hline{b.X + 1, b.Y + b.Height - 1, b.Width - 2, b.FgColor, b.BgColor}).Buffer())
	copy(pts[4+2*b.Width-4:], (vline{b.X, b.Y + 1, b.Height - 2, b.FgColor, b.BgColor}).Buffer())
	copy(pts[4+2*b.Width-4+b.Height-2:], (vline{b.X + b.Width - 1, b.Y + 1, b.Height - 2, b.FgColor, b.BgColor}).Buffer())

	return pts
}

type labeledBorder struct {
	border
	Label        string
	LabelFgColor Attribute
	LabelBgColor Attribute
}

func (lb labeledBorder) Buffer() []Point {
	ps := lb.border.Buffer()
	maxTxtW := lb.Width - 2
	rs := trimStr2Runes(lb.Label, maxTxtW)

	for i := 0; i < len(rs); i++ {
		p := Point{}
		p.X = lb.X + 1 + i
		p.Y = lb.Y
		p.Code.Ch = rs[i]
		p.Code.Fg = toTmAttr(lb.LabelFgColor)
		p.Code.Bg = toTmAttr(lb.LabelBgColor)
		ps = append(ps, p)
	}

	return ps
}
