package termui

import "strings"

/* DataGrid is like the output of the unix top command.

Example:

	columns := []termui.DataColumn{
		{12, `Market`, ui.AlignLeft, nil},
		{10, `High`, ui.AlignRight, nil},
		{10, `Low`, ui.AlignRight, nil},
		{10, `Ask`, ui.AlignRight, nil},
		{10, `Bid`, ui.AlignRight, nil},
		{7, `Buys`, ui.AlignRight, nil},
		{7, `Sells`, ui.AlignRight, nil},
		{10, `Volume`, ui.AlignRight, nil},
		{10, `Last`, ui.AlignRight, nil},
	}

	grid := termui.NewDataGrid()
	grid.FgColor = termui.ColorYellow
	grid.BgColor = termui.ColorDefault
	grid.Height = 50
	grid.Width = 120
	grid.Y = 0
	grid.X = 0
	grid.Border = false

	var results [][]string
	// Set value of results with your data
	grid.Rows = results
	termui.Render(grid)
*/

// DataGrid tracks all the attributes of a DataGrid instance
type DataGrid struct {
	Block
	Rows        [][]string
	DataColumns []DataColumn
	FgColor     Attribute
	BgColor     Attribute
	ShowBorder  bool
	ShowHeader  bool
}

// DataColumn ...
type DataColumn struct {
	Width     int
	Title     string
	TextAlign Align
	Formatter func(string) string // Optional function to format the column content
}

// DataRow ... TBD
type DataRow struct {
	FgColor Attribute
	BgColor Attribute
}

// DataCell ... TBD
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
	grid.ShowBorder = false
	grid.ShowHeader = true
	return grid
}

// Analysis generates and returns an array of []Cell that represent all cells in the DataGrid
func (grid *DataGrid) Analysis() [][]Cell {
	var rowCells [][]Cell

	// For each row []string object with rowIndex y, set the fg and bg colors
	for _, row := range grid.Rows {
		// For each string in row with index x, build the string with colors
		for _, str := range row {
			// FIXME: use column def
			cells := DefaultTxBuilder.Build(str, grid.FgColor, grid.BgColor)
			rowCells = append(rowCells, cells)
		}
	}
	return rowCells
}

// Buffer is a required func that renders the text content into termui
func (grid *DataGrid) Buffer() Buffer {
	buffer := grid.Block.Buffer()
	rowCells := grid.Analysis()
	pointerX := grid.innerArea.Min.X + 1
	pointerY := grid.innerArea.Min.Y
	startPointerX := grid.innerArea.Min.X

	if grid.ShowHeader {
		for x, column := range grid.DataColumns {
			grid.positionText(column.Title, x, &pointerX, &startPointerX)
			cells := DefaultTxBuilder.Build(column.Title, grid.FgColor, grid.BgColor)
			for i, cell := range cells {
				buffer.Set(pointerX+i, pointerY, cell)
			}
		}
	}

	// For each []string object with rowIndex y
	for y, row := range grid.Rows {
		// For each string in row array
		for x := range row {
			if grid.ShowHeader {
				pointerY = y + grid.Y + 1
			}
			grid.positionText(grid.Rows[y][x], x, &pointerX, &startPointerX)
			bgWidth := grid.DataColumns[x].Width
			if grid.ShowBorder {
				bgWidth += 3
			}
			bgCells := DefaultTxBuilder.Build(strings.Repeat(" ", bgWidth), grid.BgColor, grid.BgColor)

			cells := rowCells[y*len(row)+x]
			for i, bgCell := range bgCells {
				buffer.Set(startPointerX+i, pointerY, bgCell)
			}

			coordinateX := pointerX
			for _, printer := range cells {
				buffer.Set(coordinateX, pointerY, printer)
				coordinateX += printer.Width()
			}

			if x != 0 && grid.ShowBorder {
				dividors := DefaultTxBuilder.Build("|", grid.FgColor, grid.BgColor)
				for _, dividor := range dividors {
					buffer.Set(startPointerX, pointerY, dividor)
				}
			}
		}
	}
	return buffer
}

func (grid *DataGrid) positionText(text string, x int, coordinateX *int, cellStart *int) {
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
		*coordinateX = *cellStart + (grid.DataColumns[x].Width - len(text)) + 2
	case AlignCenter:
		*coordinateX = *cellStart + (grid.DataColumns[x].Width-len(text))/2 + 2
	default:
		*coordinateX = *cellStart + 1
	}
}
