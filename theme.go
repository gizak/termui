// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

var StandardColors = Colors{
	ColorRed,
	ColorGreen,
	ColorYellow,
	ColorBlue,
	ColorMagenta,
	ColorCyan,
	ColorWhite,
}

var StandardStyles = Styles{
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
	Plot            PlotTheme
	List            ListTheme
	Tree            TreeTheme
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
	Bars   Colors
	Nums   Styles
	Labels Styles
}

type GaugeTheme struct {
	Bar   Color
	Label Style
}

type PlotTheme struct {
	Lines Colors
	Axes  Color
}

type ListTheme struct {
	Text Style
}

type TreeTheme struct {
	Text      Style
	Collapsed rune
	Expanded  rune
}

type ParagraphTheme struct {
	Text Style
}

type PieChartTheme struct {
	Slices Colors
}

type SparklineTheme struct {
	Title Style
	Line  Color
}

type StackedBarChartTheme struct {
	Bars   Colors
	Nums   Styles
	Labels Styles
}

type TabTheme struct {
	Active   Style
	Inactive Style
}

type TableTheme struct {
	Text Style
}

// Theme holds the default Styles and Colors for all widgets.
// You can set default widget Styles by modifying the Theme before creating the widgets.
var Theme = RootTheme{
	Default: NewStyle(ColorWhite),

	Block: BlockTheme{
		Title:  NewStyle(ColorWhite),
		Border: NewStyle(ColorWhite),
	},

	BarChart: BarChartTheme{
		Bars:   StandardColors.Clone(),
		Nums:   StandardStyles.Clone(),
		Labels: StandardStyles.Clone(),
	},

	Paragraph: ParagraphTheme{
		Text: NewStyle(ColorWhite),
	},

	PieChart: PieChartTheme{
		Slices: StandardColors.Clone(),
	},

	List: ListTheme{
		Text: NewStyle(ColorWhite),
	},

	Tree: TreeTheme{
		Text:      NewStyle(ColorWhite),
		Collapsed: COLLAPSED,
		Expanded:  EXPANDED,
	},

	StackedBarChart: StackedBarChartTheme{
		Bars:   StandardColors.Clone(),
		Nums:   StandardStyles.Clone(),
		Labels: StandardStyles.Clone(),
	},

	Gauge: GaugeTheme{
		Bar:   ColorWhite,
		Label: NewStyle(ColorWhite),
	},

	Sparkline: SparklineTheme{
		Title: NewStyle(ColorWhite),
		Line:  ColorWhite,
	},

	Plot: PlotTheme{
		Lines: StandardColors.Clone(),
		Axes:  ColorWhite,
	},

	Table: TableTheme{
		Text: NewStyle(ColorWhite),
	},

	Tab: TabTheme{
		Active:   NewStyle(ColorRed),
		Inactive: NewStyle(ColorWhite),
	},
}
