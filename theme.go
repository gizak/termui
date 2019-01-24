// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

var StandardColors = []Color{
	ColorRed,
	ColorGreen,
	ColorYellow,
	ColorBlue,
	ColorMagenta,
	ColorCyan,
	ColorWhite,
}

var StandardStyles = []Style{
	NewStyle(ColorRed),
	NewStyle(ColorGreen),
	NewStyle(ColorYellow),
	NewStyle(ColorBlue),
	NewStyle(ColorMagenta),
	NewStyle(ColorCyan),
	NewStyle(ColorWhite),
}

type RootTheme struct {
	Default Style

	Block BlockTheme

	BarChart        BarChartTheme
	Gauge           GaugeTheme
	LineChart       LineChartTheme
	List            ListTheme
	Paragraph       ParagraphTheme
	PieChart        PieChartTheme
	Sparkline       SparklineTheme
	StackedBarChart StackedBarChartTheme
	Tab             TabTheme
	Table           TableTheme
}

type BlockTheme struct {
	Title  Style
	Border Style
}

type BarChartTheme struct {
	Bars   []Color
	Nums   []Style
	Labels []Style
}

type GaugeTheme struct {
	Bar   Color
	Label Style
}

type LineChartTheme struct {
	Lines []Color
	Axes  Color
}

type ListTheme struct {
	Text Style
}

type ParagraphTheme struct {
	Text Style
}

type PieChartTheme struct {
	Slices []Color
}

type SparklineTheme struct {
	Title Style
	Line  Color
}

type StackedBarChartTheme struct {
	Bars   []Color
	Nums   []Style
	Labels []Style
}

type TabTheme struct {
	Active   Style
	Inactive Style
}

type TableTheme struct {
	Text Style
}

var Theme = RootTheme{
	Default: NewStyle(ColorWhite),

	Block: BlockTheme{
		Title:  NewStyle(ColorWhite),
		Border: NewStyle(ColorWhite),
	},

	BarChart: BarChartTheme{
		Bars:   StandardColors,
		Nums:   StandardStyles,
		Labels: StandardStyles,
	},

	Paragraph: ParagraphTheme{
		Text: NewStyle(ColorWhite),
	},

	PieChart: PieChartTheme{
		Slices: StandardColors,
	},

	List: ListTheme{
		Text: NewStyle(ColorWhite),
	},

	StackedBarChart: StackedBarChartTheme{
		Bars:   StandardColors,
		Nums:   StandardStyles,
		Labels: StandardStyles,
	},

	Gauge: GaugeTheme{
		Bar:   ColorWhite,
		Label: NewStyle(ColorWhite),
	},

	Sparkline: SparklineTheme{
		Line:  ColorBlack,
		Title: NewStyle(ColorBlue),
	},

	LineChart: LineChartTheme{
		Lines: StandardColors,
		Axes:  ColorBlue,
	},

	Table: TableTheme{
		Text: NewStyle(ColorWhite),
	},

	Tab: TabTheme{
		Active:   NewStyle(ColorRed),
		Inactive: NewStyle(ColorWhite),
	},
}
