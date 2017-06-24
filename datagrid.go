package termui

import "strings"

/* DataGrid is like:

┌Awesome DataGrid ─────────────────────────────────────────────┐
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

	grid := termui.NewDataGrid()
	grid.Rows = data  // type [][]string
	grid.FgColor = termui.ColorWhite
	grid.BgColor = termui.ColorDefault
	grid.Height = 7
	grid.Width = 62
	grid.Y = 0
	grid.X = 0
	grid.Border = true
*/

// DataGrid tracks all the attributes of a DataGrid instance
type DataGrid struct {
	Block
	Rows        [][]string
	DataColumns []DataColumn
	FgColor     Attribute
	BgColor     Attribute
	Separator   bool
	FgColors    []Attribute
	BgColors    []Attribute
	// TextAlign   Align
}

// DataColumn ...
type DataColumn struct {
	Width     int
	Title     string
	TextAlign Align
	Formatter func(string) string // Optional function to format the column content
}

// DataRow ...
type DataRow struct {
	FgColor Attribute
	BgColor Attribute
}

// DataCell ...
type DataCell struct {
	Text       string
	Value      float64
	Message    string
	FgColor    Attribute
	BgColor    Attribute
	TextFormat string
	Height     int
}

// NewDataGrid returns a new DataGrid instance
func NewDataGrid() *DataGrid {
	grid := &DataGrid{Block: *NewBlock()}
	grid.FgColor = ColorWhite
	grid.BgColor = ColorDefault
	grid.Separator = true
	return grid
}

// CellsWidth calculates the width of a cell array and returns an int
// func cellsWidth(cells []Cell) int {
// 	width := 0
// 	for _, c := range cells {
// 		width += c.Width()
// 	}
// 	return width
// }

// Analysis generates and returns an array of []Cell that represent all columns in the DataGrid
func (grid *DataGrid) Analysis() [][]Cell {
	var rowCells [][]Cell
	length := len(grid.Rows)
	if length < 1 {
		return rowCells
	}

	// Create array of FgColors and BgColors for every row
	if len(grid.FgColors) == 0 {
		grid.FgColors = make([]Attribute, len(grid.Rows))
	}
	if len(grid.BgColors) == 0 {
		grid.BgColors = make([]Attribute, len(grid.Rows))
	}

	// cellWidths := make([]int, len(grid.Rows[0]))

	// For each row []string object with rowIndex y, set the fg and bg colors
	for y, row := range grid.Rows {
		if grid.FgColors[y] == 0 {
			grid.FgColors[y] = grid.FgColor
		}
		if grid.BgColors[y] == 0 {
			grid.BgColors[y] = grid.BgColor
		}
		// For each string in row with index x, build the string with colors
		for _, str := range row {
			cells := DefaultTxBuilder.Build(str, grid.FgColors[y], grid.BgColors[y])
			// FIXME: Datagrid is fixed width based on defined columns.
			// cw := cellsWidth(cells)
			// if cellWidths[x] < cw {
			// 	cellWidths[x] = cw
			// }
			rowCells = append(rowCells, cells)
		}
	}
	// grid.CellWidth = cellWidths
	return rowCells
}

// EvaluatePosition ...
func (grid *DataGrid) EvaluatePosition(x int, y int, coordinateX *int, coordinateY *int, cellStart *int) {
	if grid.Separator {
		*coordinateY = grid.innerArea.Min.Y + y*2
	} else {
		*coordinateY = grid.innerArea.Min.Y + y
	}
	if x == 0 {
		*cellStart = grid.innerArea.Min.X
	} else {
		*cellStart += grid.DataColumns[x-1].Width + 3
	}
	align := AlignLeft
	if x < len(grid.DataColumns) {
		align = grid.DataColumns[x].TextAlign
	}
	switch align {
	case AlignRight:
		*coordinateX = *cellStart + (grid.DataColumns[x].Width - len(grid.Rows[y][x])) + 2
	case AlignCenter:
		*coordinateX = *cellStart + (grid.DataColumns[x].Width-len(grid.Rows[y][x]))/2 + 2
	default:
		*coordinateX = *cellStart + 2
	}
}

// Buffer ...
func (grid *DataGrid) Buffer() Buffer {
	buffer := grid.Block.Buffer()
	rowCells := grid.Analysis()
	pointerX := grid.innerArea.Min.X + 2
	pointerY := grid.innerArea.Min.Y
	borderPointerX := grid.innerArea.Min.X
	// For each []string object with rowIndex y
	for y, row := range grid.Rows {
		// For each string in row array
		for x := range row {

			grid.EvaluatePosition(x, y, &pointerX, &pointerY, &borderPointerX)
			background := DefaultTxBuilder.Build(strings.Repeat(" ", grid.DataColumns[x].Width+3), grid.BgColors[y], grid.BgColors[y])
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
				dividors := DefaultTxBuilder.Build("|", grid.FgColors[y], grid.BgColors[y])
				for _, dividor := range dividors {
					buffer.Set(borderPointerX, pointerY, dividor)
				}
			}
		}

		// if grid.Separator {
		// 	border := DefaultTxBuilder.Build(strings.Repeat("─", grid.Width-2), grid.FgColor, grid.BgColor)
		// 	for i, cell := range border {
		// 		buffer.Set(i+1, pointerY+1, cell)
		// 	}
		// }
	}

	return buffer
}
