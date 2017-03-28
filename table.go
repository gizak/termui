// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"strings"
	"sync"
)

/* Table is like:

┌Awesome Table ────────────────────────────────────────────────┐
│  Col0          | Col1 | Col2 | Col3  | Col4  | Col5  | Col6  |
│──────────────────────────────────────────────────────────────│
│  Some Item #1  | AAA  | 123  | CCCCC | EEEEE | GGGGG | IIIII |
│──────────────────────────────────────────────────────────────│
│  Some Item #2  | BBB  | 456  | DDDDD | FFFFF | HHHHH | JJJJJ |
└──────────────────────────────────────────────────────────────┘

Datapoints are a two dimensional array of strings: [][]string

Example:
	data := [][]string{
		{"Col0", "Col1", "Col3", "Col4", "Col5", "Col6"},
		{"Some Item #1", "AAA", "123", "CCCCC", "EEEEE", "GGGGG", "IIIII"},
		{"Some Item #2", "BBB", "456", "DDDDD", "FFFFF", "HHHHH", "JJJJJ"},
	}

	table := termui.NewTable()
	table.Rows = data  // type [][]string
	table.FgColor = termui.ColorWhite
	table.BgColor = termui.ColorDefault
	table.Height = 7
	table.Width = 62
	table.Y = 0
	table.X = 0
	table.Border = true
*/

// Table tracks all the attributes of a Table instance
type Table struct {
	Block
	Rows      [][]string
	CellWidth []int
	FgColor   Attribute
	BgColor   Attribute
	FgColors  []Attribute
	BgColors  []Attribute
	Separator bool
	TextAlign Align
	sync.Mutex
}

// NewTable returns a new Table instance
func NewTable() *Table {
	table := &Table{Block: *NewBlock()}
	table.FgColor = ColorWhite
	table.BgColor = ColorDefault
	table.Separator = true
	return table
}

// CellsWidth calculates the width of a cell array and returns an int
func cellsWidth(cells []Cell) int {
	width := 0
	for _, c := range cells {
		width += c.Width()
	}
	return width
}

func (table *Table) analysis() [][]Cell {
	var rowCells [][]Cell
	length := len(table.Rows)
	if length < 1 {
		return rowCells
	}

	if len(table.FgColors) == 0 {
		table.FgColors = make([]Attribute, len(table.Rows))
	}
	if len(table.BgColors) == 0 {
		table.BgColors = make([]Attribute, len(table.Rows))
	}

	cellWidths := make([]int, len(table.Rows[0]))

	for y, row := range table.Rows {
		if table.FgColors[y] == 0 {
			table.FgColors[y] = table.FgColor
		}
		if table.BgColors[y] == 0 {
			table.BgColors[y] = table.BgColor
		}
		for x, str := range row {
			cells := DefaultTxBuilder.Build(str, table.FgColors[y], table.BgColors[y])
			cw := cellsWidth(cells)
			if cellWidths[x] < cw {
				cellWidths[x] = cw
			}
			rowCells = append(rowCells, cells)
		}
	}
	table.CellWidth = cellWidths
	return rowCells
}

// Analysis generates and returns an array of []Cell that represent all columns in the Table
func (table *Table) Analysis() [][]Cell {
	table.Lock()
	defer table.Unlock()
	return table.analysis()
}

// SetSize calculates the table size and sets the internal value
func (table *Table) SetSize() {
	table.Lock()
	defer table.Unlock()
	length := len(table.Rows)
	if table.Separator {
		table.Height = length*2 + 1
	} else {
		table.Height = length + 2
	}
	table.Width = 2
	if length != 0 {
		for _, cellWidth := range table.CellWidth {
			table.Width += cellWidth + 3
		}
	}
}

func (table *Table) calculatePosition(x int, y int, coordinateX *int, coordinateY *int, cellStart *int) {
	if table.Separator {
		*coordinateY = table.innerArea.Min.Y + y*2
	} else {
		*coordinateY = table.innerArea.Min.Y + y
	}
	if x == 0 {
		*cellStart = table.innerArea.Min.X
	} else {
		*cellStart += table.CellWidth[x-1] + 3
	}

	switch table.TextAlign {
	case AlignRight:
		*coordinateX = *cellStart + (table.CellWidth[x] - len(table.Rows[y][x])) + 2
	case AlignCenter:
		*coordinateX = *cellStart + (table.CellWidth[x]-len(table.Rows[y][x]))/2 + 2
	default:
		*coordinateX = *cellStart + 2
	}
}

// CalculatePosition ...
func (table *Table) CalculatePosition(x int, y int, coordinateX *int, coordinateY *int, cellStart *int) {
	table.Lock()
	defer table.Unlock()
	table.calculatePosition(x, y, coordinateX, coordinateY, cellStart)
}

// Buffer ...
func (table *Table) Buffer() Buffer {
	table.Lock()
	defer table.Unlock()
	buffer := table.Block.Buffer()
	rowCells := table.analysis()
	pointerX := table.innerArea.Min.X + 2
	pointerY := table.innerArea.Min.Y
	borderPointerX := table.innerArea.Min.X
	for y, row := range table.Rows {
		for x := range row {
			table.calculatePosition(x, y, &pointerX, &pointerY, &borderPointerX)
			background := DefaultTxBuilder.Build(strings.Repeat(" ", table.CellWidth[x]+3), table.BgColors[y], table.BgColors[y])
			cells := rowCells[y*len(row)+x]
			for i, back := range background {
				buffer.Set(borderPointerX+i, pointerY, back)
			}

			coordinateX := pointerX
			for _, printer := range cells {
				buffer.Set(coordinateX, pointerY, printer)
				coordinateX += printer.Width()
			}

			if x != 0 {
				dividors := DefaultTxBuilder.Build("|", table.FgColors[y], table.BgColors[y])
				for _, dividor := range dividors {
					buffer.Set(borderPointerX, pointerY, dividor)
				}
			}
		}

		if table.Separator {
			border := DefaultTxBuilder.Build(strings.Repeat("─", table.Width-2), table.FgColor, table.BgColor)
			for i, cell := range border {
				buffer.Set(i+1, pointerY+1, cell)
			}
		}
	}

	return buffer
}

func (table *Table) SetRows(rows [][]string) {
	table.Lock()
	defer table.Unlock()

	oldlen := len(table.Rows)
	table.Rows = rows
	nrNewRows := len(table.Rows) - oldlen
	nrNewColors := len(table.Rows) - len(table.FgColors) /* FgColors and BgColors are in sync */

	/* if there is a positive delta between the current number of colors and then number we expect allocate them */
	/* we intentionally do not deallocate unnecessary colors. They are not used and we keep them "chached" */
	/* caching avoids reallocation and it is relatively unlikely that a table starts very big, decreades a lot and stays like that */
	if nrNewColors > 0 {
		newfgs := make([]Attribute, nrNewColors)
		newbgs := make([]Attribute, nrNewColors)

		table.FgColors = append(table.FgColors, newfgs...)
		table.BgColors = append(table.BgColors, newbgs...)
	}

	/* always reset the colors of additional rows */
	if nrNewRows > 0 {
		for i := 0; i < nrNewRows; i++ {
			table.FgColors[oldlen+i] = table.FgColor
			table.BgColors[oldlen+i] = table.BgColor
		}
	}
}
