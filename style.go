package termui

import "strings"

// Color is an integer from -1 to 255
// -1 = ColorClear
// 0-255 = Xterm colors
type Color int

// ColorClear clears the Fg or Bg color of a Style
const ColorClear Color = -1

// Basic terminal colors
const (
	ColorBlack   Color = 0
	ColorRed     Color = 1
	ColorGreen   Color = 2
	ColorYellow  Color = 3
	ColorBlue    Color = 4
	ColorMagenta Color = 5
	ColorCyan    Color = 6
	ColorWhite   Color = 7
)

type Modifier uint

const (
	// ModifierClear clears any modifiers
	ModifierClear     Modifier = 0
	ModifierBold      Modifier = 1 << 9
	ModifierUnderline Modifier = 1 << 10
	ModifierReverse   Modifier = 1 << 11
)

// Style represents the style of one terminal cell
type Style struct {
	Fg       Color
	Bg       Color
	Modifier Modifier
}

// StyleClear represents a default Style, with no colors or modifiers
var StyleClear = Style{
	Fg:       ColorClear,
	Bg:       ColorClear,
	Modifier: ModifierClear,
}

// NewStyle takes 1 to 3 arguments
// 1st argument = Fg
// 2nd argument = optional Bg
// 3rd argument = optional Modifier
func NewStyle(fg Color, args ...interface{}) Style {
	bg := ColorClear
	modifier := ModifierClear
	if len(args) >= 1 {
		bg = args[0].(Color)
	}
	if len(args) == 2 {
		modifier = args[1].(Modifier)
	}
	return Style{
		fg,
		bg,
		modifier,
	}
}

//String returns a string representation of a Style
func (self Style) String() string {
	styles := make([]string, 0)

	if color, ok := textColorMap[self.Fg]; ok  && self.Fg !=  StyleClear.Fg {
		styles = append(styles, tokenFg + tokenValueSeparator + color)
	}

	if color, ok := textColorMap[self.Bg]; ok && self.Bg !=  StyleClear.Bg {
		styles = append(styles, tokenBg + tokenValueSeparator + color)
	}

	if mod, ok := textModifierMap[self.Modifier]; ok  && self.Modifier !=  StyleClear.Modifier {
		styles = append(styles, tokenModifier + tokenValueSeparator + mod)
	}

	return strings.Join(styles, tokenItemSeparator)
}
