// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var r *Row

func TestRowWidth(t *testing.T) {
	p0 := NewBlock()
	p0.Height = 1
	p1 := NewBlock()
	p1.Height = 1
	p2 := NewBlock()
	p2.Height = 1
	p3 := NewBlock()
	p3.Height = 1

	/* test against tree:

	       r
	     /  \
	   0:w   1
	        / \
	     10:w 11
	          /
	        110:w
	         /
	       1100:w
	*/

	r = NewRow(
		NewCol(6, 0, p0),
		NewCol(6, 0,
			NewRow(
				NewCol(6, 0, p1),
				NewCol(6, 0, p2, p3))))

	r.assignWidth(100)
	if r.Width != 100 ||
		(r.Cols[0].Width) != 50 ||
		(r.Cols[1].Width) != 50 ||
		(r.Cols[1].Cols[0].Width) != 25 ||
		(r.Cols[1].Cols[1].Width) != 25 ||
		(r.Cols[1].Cols[1].Cols[0].Width) != 25 ||
		(r.Cols[1].Cols[1].Cols[0].Cols[0].Width) != 25 {
		t.Error("assignWidth fails")
	}
}

func TestRowHeight(t *testing.T) {
	spew.Dump()

	if (r.solveHeight()) != 2 ||
		(r.Cols[1].Cols[1].Height) != 2 ||
		(r.Cols[1].Cols[1].Cols[0].Height) != 2 ||
		(r.Cols[1].Cols[0].Height) != 1 {
		t.Error("solveHeight fails")
	}
}

func TestAssignXY(t *testing.T) {
	r.assignX(0)
	r.assignY(0)
	if (r.Cols[0].X) != 0 ||
		(r.Cols[1].Cols[0].X) != 50 ||
		(r.Cols[1].Cols[1].X) != 75 ||
		(r.Cols[1].Cols[1].Cols[0].X) != 75 ||
		(r.Cols[1].Cols[0].Y) != 0 ||
		(r.Cols[1].Cols[1].Cols[0].Y) != 0 ||
		(r.Cols[1].Cols[1].Cols[0].Cols[0].Y) != 1 {
		t.Error("assignXY fails")
	}
}
