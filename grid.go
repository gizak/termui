// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

// Bufferers that can be manipulated by Grid
type LayoutBufferer interface {
	Bufferer
	GetHeight() int
	SetWidth(int)
	SetX(int)
	SetY(int)
}

// build a layout tree
type row struct {
	Cols   []*row         //children
	Widget LayoutBufferer // root
	X      int
	Y      int
	Width  int
	Height int
	Span   int
	Offset int
}

func (r *row) calcLayout() {
	r.assignWidth(r.Width)
	r.Height = r.solveHeight()
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
	r.SetWidth(cw)

	for i, _ := range r.Cols {
		r.Cols[i].assignWidth(cw)
	}
}

// bottom up, return r's total height
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
	r.SetX(x)

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
}

func (r *row) assignY(y int) {
	r.SetY(y)

	if r.isLeaf() {
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

func (r row) GetHeight() int {
	return r.Height
}

func (r *row) SetX(x int) {
	r.X = x
	if r.Widget != nil {
		r.Widget.SetX(x)
	}
}

func (r *row) SetY(y int) {
	r.Y = y
	if r.Widget != nil {
		r.Widget.SetY(y)
	}
}

func (r *row) SetWidth(w int) {
	r.Width = w
	if r.Widget != nil {
		r.Widget.SetWidth(w)
	}
}

// recursively merge all widgets buffer
func (r *row) Buffer() []Point {
	merged := []Point{}

	if r.isRenderableLeaf() {
		return r.Widget.Buffer()
	}

	// for those are not leaves but have a renderable widget
	if r.Widget != nil {
		merged = append(merged, r.Widget.Buffer()...)
	}

	// collect buffer from children
	if !r.isLeaf() {
		for _, c := range r.Cols {
			merged = append(merged, c.Buffer()...)
		}
	}

	return merged
}

type Grid struct {
	Rows    []*row
	Width   int
	X       int
	Y       int
	BgColor Attribute
}

func NewGrid(rows ...*row) *Grid {
	return &Grid{Rows: rows}
}

func (g *Grid) AddRows(rs ...*row) {
	g.Rows = append(g.Rows, rs...)
}

func NewRow(cols ...*row) *row {
	rs := &row{Span: 12, Cols: cols}
	return rs
}

// NewCol accepts: widgets are LayoutBufferer or
//                 widgets is A NewRow
// Note that if multiple widgets are provided, they will stack up in the col
func NewCol(span, offset int, widgets ...LayoutBufferer) *row {
	r := &row{Span: span, Offset: offset}

	if widgets != nil && len(widgets) == 1 {
		wgt := widgets[0]
		nw, isRow := wgt.(*row)
		if isRow {
			r.Cols = nw.Cols
		} else {
			r.Widget = wgt
		}
		return r
	}

	r.Cols = []*row{}
	ir := r
	for _, w := range widgets {
		nr := &row{Span: 12, Widget: w}
		ir.Cols = []*row{nr}
		ir = nr
	}

	return r
}

// Calculate each rows' layout
func (g *Grid) Align() {
	h := 0
	for _, r := range g.Rows {
		r.SetWidth(g.Width)
		r.SetX(g.X)
		r.SetY(g.Y + h)
		r.calcLayout()
		h += r.GetHeight()
	}
}

func (g Grid) Buffer() []Point {
	ps := []Point{}
	for _, r := range g.Rows {
		ps = append(ps, r.Buffer()...)
	}
	return ps
}

var Body *Grid
