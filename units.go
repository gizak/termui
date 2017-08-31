// writen in 2017 by cnaize

package termui

// Unit implementations of Blocks

type UBarChart struct {
	UnitInterface
	*BarChart
}

func NewUBarChart(tag string) *UBarChart {
	bch := NewBarChart()
	return &UBarChart{UnitInterface: NewUnit(tag, &bch.Block, bch), BarChart: bch}
}

type UGauge struct {
	UnitInterface
	*Gauge
}

func NewUGauge(tag string) *UGauge {
	gauge := NewGauge()
	return &UGauge{UnitInterface: NewUnit(tag, &gauge.Block, gauge), Gauge: gauge}
}

type ULineChart struct {
	UnitInterface
	*LineChart
}

func NewULineChart(tag string) *ULineChart {
	lch := NewLineChart()
	return &ULineChart{UnitInterface: NewUnit(tag, &lch.Block, lch), LineChart: lch}
}

type UMBarChart struct {
	UnitInterface
	*MBarChart
}

func NewUMBarChart(tag string) *UMBarChart {
	mbch := NewMBarChart()
	return &UMBarChart{UnitInterface: NewUnit(tag, &mbch.Block, mbch), MBarChart: mbch}
}

type UPar struct {
	UnitInterface
	*Par
}

func NewUPar(tag, text string) *UPar {
	par := NewPar(text)
	return &UPar{UnitInterface: NewUnit(tag, &par.Block, par), Par: par}
}

type UList struct {
	UnitInterface
	*List
}

func NewUList(tag string) *UList {
	list := NewList()
	return &UList{UnitInterface: NewUnit(tag, &list.Block, list), List: list}
}

type UTable struct {
	UnitInterface
	*Table
}

func NewUTable(tag string) *UTable {
	table := NewTable()
	return &UTable{UnitInterface: NewUnit(tag, &table.Block, table), Table: table}
}
