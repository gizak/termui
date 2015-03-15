package termui

import tm "github.com/nsf/termbox-go"

type container struct {
	//Height  int
	//Width   int
	BgColor Attribute
	Rows    []Row
}

type Row []Col

type Col struct {
	ColumnBufferer
	Offset int // 0 ~ 11
	Span   int // 1 ~ 12
	Sticky bool
}

type ColumnBufferer interface {
	Bufferer
	GetHeight() int
	GetWidth() int
	SetWidth(int)
	SetX(int)
	SetY(int)
}

func NewRow(cols ...Col) Row {
	return cols
}

func NewCol(block ColumnBufferer, span, offset int, sticky bool) Col {
	return Col{ColumnBufferer: block, Span: span, Sticky: sticky, Offset: offset}
}

func (c container) Buffer() []Point {
	ps := []Point{}
	maxw, _ := tm.Size()

	y := 0
	for _, row := range c.Rows {
		x := 0
		maxHeight := 0

		for _, col := range row {
			if h := col.GetHeight(); h > maxHeight {
				maxHeight = h
			}

			w := int(float64(maxw*(col.Span+col.Offset)) / 12.0)
			if col.Sticky || col.GetWidth() > w {
				col.SetWidth(w)
			}

			col.SetY(y)
			col.SetX(x)
			ps = append(ps, col.Buffer()...)
			x += w + int(float64(maxw*col.Offset)/12)
		}
		y += maxHeight
	}
	return ps
}

var Body container
