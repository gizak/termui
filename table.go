package termui

import (
	"fmt"
	"strings"
)

/*
	table := termui.NewTable()
	table.Rows = rows
	table.FgColor = termui.ColorWhite
	table.BgColor = termui.ColorDefault
	table.Height = 7
	table.Width = 62
	table.Y = 0
	table.X = 0
	table.Border = true
*/

type Table struct {
	Block
	Rows      [][]string
	FgColor   Attribute
	BgColor   Attribute
	FgColors  []Attribute
	BgColors  []Attribute
	Seperator bool
	TextAlign string
}

func NewTable() *Table {
	table := &Table{Block: *NewBlock()}
	table.FgColor = ColorWhite
	table.BgColor = ColorDefault
	table.TextAlign = "left"
	table.Seperator = true
	return table
}

func (table *Table) Analysis() {
	length := len(table.Rows)
	if length < 1 {
		return
	}

	if len(table.FgColors) == 0 {
		table.FgColors = make([]Attribute, len(table.Rows))
	}
	if len(table.BgColors) == 0 {
		table.BgColors = make([]Attribute, len(table.Rows))
	}

	row_width := len(table.Rows[0])
	cellWidthes := make([]int, row_width)

	for index, row := range table.Rows {
		for i, str := range row {
			if cellWidthes[i] < len(str) {
				cellWidthes[i] = len(str)
			}
		}

		if table.FgColors[index] == 0 {
			table.FgColors[index] = table.FgColor
		}

		if table.BgColors[index] == 0 {
			table.BgColors[index] = table.BgColor
		}
	}

	width_sum := 2
	for i, width := range cellWidthes {
		width_sum += (width + 2)
		for u, row := range table.Rows {
			switch table.TextAlign {
			case "right":
				row[i] = fmt.Sprintf(" %*s ", width, table.Rows[u][i])
			case "center":
				word_width := len(table.Rows[u][i])
				offset := (width - word_width) / 2
				row[i] = fmt.Sprintf(" %*s ", width, fmt.Sprintf("%-*s", offset+word_width, table.Rows[u][i]))
			default: // left
				row[i] = fmt.Sprintf(" %-*s ", width, table.Rows[u][i])
			}
		}
	}

	if table.Width == 0 {
		table.Width = width_sum
	}
}

func (table *Table) SetSize() {
	length := len(table.Rows)
	if table.Seperator {
		table.Height = length*2 + 1
	} else {
		table.Height = length + 2
	}
	table.Width = 2
	if length != 0 {
		for _, str := range table.Rows[0] {
			table.Width += len(str) + 2 + 1
		}
	}
}

func (table *Table) Buffer() Buffer {
	buffer := table.Block.Buffer()
	table.Analysis()
	for i, row := range table.Rows {
		cells := DefaultTxBuilder.Build(strings.Join(row, "|"), table.FgColors[i], table.BgColors[i])
		if table.Seperator {
			border := DefaultTxBuilder.Build(strings.Repeat("─", table.Width-2), table.FgColor, table.BgColor)
			for x, cell := range cells {
				buffer.Set(table.innerArea.Min.X+x, table.innerArea.Min.Y+i*2, cell)
			}

			for x, cell := range border {
				buffer.Set(table.innerArea.Min.X+x, table.innerArea.Min.Y+i*2+1, cell)
			}
		} else {
			for x, cell := range cells {
				buffer.Set(table.innerArea.Min.X+x, table.innerArea.Min.Y+i, cell)
			}
		}
	}
	return buffer
}
