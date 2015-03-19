package termui

import tm "github.com/nsf/termbox-go"

/*
type container struct {
	height  int
	width   int
	BgColor Attribute
	Rows    []Row
}

type Row []Col

type Col struct {
	Blocks []ColumnBufferer
	Offset int // 0 ~ 11
	Span   int // 1 ~ 12
}

type ColumnBufferer interface {
	Bufferer
	GetHeight() int
	SetWidth(int)
	SetX(int)
	SetY(int)
}

func NewRow(cols ...Col) Row {
	return cols
}

func NewCol(span, offset int, blocks ...ColumnBufferer) Col {
	return Col{Blocks: blocks, Span: span, Offset: offset}
}

// Highest col is the height of a Row
func (r Row) GetHeight() int {
	h := 0
	for _, v := range r {
		if nh := v.GetHeight(); nh > h {
			h = nh
		}
	}
	return h
}

// Set width according to its span
func (r Row) SetWidth(w int) {
	for _, c := range r {
		c.SetWidth(int(float64(w*c.Span) / 12.0))
	}
}

// Set x y
func (r Row) SetX(x int) {
	for i := range r {
		r[i].SetX(x)
	}
}

func (r Row) SetY(y int) {
	for i := range r {
		r[i].SetY(y)
	}
}

// GetHeight recursively retrieves height of each children, then add them up.
func (c Col) GetHeight() int {
	h := 0
	for _, v := range c.Blocks {
		h += c.GetHeight()
	}
	return h
}

func (c Col) GetWidth() int {
	w := 0
	for _, v := range c.Blocks {
		if nw := v.GetWidth(); nw > w {
			w = nw
		}
	}
	return w
}

func (c Col) SetWidth(w int) {
	for i := range c.Blocks {
		c.SetWidth(w)
	}
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
			if col.GetWidth() > w {
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
*/

type LayoutBufferer interface {
	Bufferer
	GetHeight() int
	SetWidth(int)
	SetX(int)
	SetY(int)
}

// build a layout tree
type row struct {
	Cols   []*row
	Widget LayoutBufferer // only leaves hold this
	X      int
	Y      int
	Width  int
	Height int
	Span   int
	Offset int
}

func newContainer() *row {
	w, _ := tm.Size()
	r := &row{Width: w, Span: 12, X: 0, Y: 0, Cols: []*row{}}
	return r
}

func (r *row) layout() {
	r.assignWidth(r.Width)
	r.solveHeight()
	r.assignX(r.X)
	r.assignY(r.Y)
}

func (r *row) isLeaf() bool {
	return r.Cols == nil || len(r.Cols) == 0
}

func (r *row) isRenderableLeaf() bool {
	return r.isLeaf() && r.Widget != nil
}

func (r *row) assignWidth(w int) {
	cw := int(float64(w*r.Span) / 12)
	r.Width = cw

	for i, _ := range r.Cols {
		r.Cols[i].assignWidth(cw)
	}
}

// bottom up
func (r *row) solveHeight() int {
	if r.isRenderableLeaf() {
		r.Height = r.Widget.GetHeight()
		return r.Widget.GetHeight()
	}

	maxh := 0
	if !r.isLeaf() {
		for _, c := range r.Cols {
			nh := c.solveHeight()
			// when embed rows in Cols, row widgets stack up
			if r.Widget != nil {
				nh += r.Widget.GetHeight()
			}
			if nh > maxh {
				maxh = nh
			}
		}
	}

	r.Height = maxh
	return maxh
}

func (r *row) assignX(x int) {
	if r.isRenderableLeaf() {
		r.Widget.SetX(x)
	}

	if !r.isLeaf() {
		acc := 0
		for i, c := range r.Cols {
			r.Cols[i].assignX(x + acc)
			acc += c.Width
			if c.Offset != 0 {
				acc += int(float64(c.Offset*c.Width) / float64(12*c.Span))
			}
		}
	}
	r.X = x
}

func (r *row) assignY(y int) {
	r.Y = y

	if r.isRenderableLeaf() {
		r.Widget.SetY(y)
		return
	}

	for i := range r.Cols {
		acc := 0
		if r.Widget != nil {
			acc = r.Widget.GetHeight()
		}
		r.Cols[i].assignY(y + acc)
	}

}

// recursively merge all widgets buffer
func (r *row) Buffer() []Point {
	merged := []Point{}

	if r.isRenderableLeaf() {
		return r.Widget.Buffer()
	}

	if !r.isLeaf() {
		for _, c := range r.Cols {
			merged = append(merged, c.Buffer()...)
		}
	}

	return merged
}

//var Body container
