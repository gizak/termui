package termui

type colorScheme struct {
	BodyBg            Attribute
	BlockBg           Attribute
	HasBorder         bool
	BorderFg          Attribute
	BorderBg          Attribute
	BorderLabelTextFg Attribute
	BorderLabelTextBg Attribute
	ParTextFg         Attribute
	ParTextBg         Attribute
	SparklineLine     Attribute
	SparklineTitle    Attribute
	GaugeBar          Attribute
	GaugePercent      Attribute
	LineChartLine     Attribute
	LineChartAxes     Attribute
	ListItemFg        Attribute
	ListItemBg        Attribute
	BarChartBar       Attribute
	BarChartText      Attribute
	BarChartNum       Attribute
}

// default color scheme depends on the user's terminal setting.
var themeDefault = colorScheme{HasBorder: true}

var themeHelloWorld = colorScheme{
	BodyBg:            ColorBlack,
	BlockBg:           ColorBlack,
	HasBorder:         true,
	BorderFg:          ColorWhite,
	BorderBg:          ColorBlack,
	BorderLabelTextBg: ColorBlack,
	BorderLabelTextFg: ColorGreen,
	ParTextBg:         ColorBlack,
	ParTextFg:         ColorWhite,
	SparklineLine:     ColorMagenta,
	SparklineTitle:    ColorWhite,
	GaugeBar:          ColorRed,
	GaugePercent:      ColorWhite,
	LineChartLine:     ColorYellow | AttrBold,
	LineChartAxes:     ColorWhite,
	ListItemBg:        ColorBlack,
	ListItemFg:        ColorYellow,
	BarChartBar:       ColorRed,
	BarChartNum:       ColorWhite,
	BarChartText:      ColorCyan,
}

var theme = themeDefault // global dep

func UseTheme(th string) {
	switch th {
	case "helloworld":
		theme = themeHelloWorld
	default:
		theme = themeDefault
	}
}
