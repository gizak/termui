package termui

import (
	"github.com/mitchellh/go-wordwrap" // LEFT UP TO PKG MAINTAINER TO DECIDE HOW TO VENDOR; is MIT LICENSED
	"regexp"
	"strings"
)

// TextBuilder is a minimal interface to produce text []Cell using specific syntax (markdown).
type TextBuilder interface {
	Build(s string, fg, bg Attribute) []Cell
	BuildWrap(s string, fg, bg Attribute, wl uint) []Cell
}

// DefaultTxBuilder is set to be MarkdownTxBuilder.
var DefaultTxBuilder = NewMarkdownTxBuilder()

// MarkdownTxBuilder implements TextBuilder interface, using markdown syntax.
type MarkdownTxBuilder struct {
	baseFg  Attribute
	baseBg  Attribute
	plainTx []rune
	markers []marker
}

type marker struct {
	st int
	ed int
	fg Attribute
	bg Attribute
}

var colorMap = map[string]Attribute{
	"red":     ColorRed,
	"blue":    ColorBlue,
	"black":   ColorBlack,
	"cyan":    ColorCyan,
	"white":   ColorWhite,
	"default": ColorDefault,
	"green":   ColorGreen,
	"magenta": ColorMagenta,
}

var attrMap = map[string]Attribute{
	"bold":      AttrBold,
	"underline": AttrUnderline,
	"reverse":   AttrReverse,
}

func rmSpc(s string) string {
	reg := regexp.MustCompile(`\s+`)
	return reg.ReplaceAllString(s, "")
}

// readAttr translates strings like `fg-red,fg-bold,bg-white` to fg and bg Attribute
func (mtb MarkdownTxBuilder) readAttr(s string) (Attribute, Attribute) {
	fg := mtb.baseFg
	bg := mtb.baseBg

	updateAttr := func(a Attribute, attrs []string) Attribute {
		for _, s := range attrs {
			// replace the color
			if c, ok := colorMap[s]; ok {
				a &= 0xFF00 // erase clr 0 ~ 8 bits
				a |= c      // set clr
			}
			// add attrs
			if c, ok := attrMap[s]; ok {
				a |= c
			}
		}
		return a
	}

	ss := strings.Split(s, ",")
	fgs := []string{}
	bgs := []string{}
	for _, v := range ss {
		subs := strings.Split(v, "-")
		if len(subs) > 1 {
			if subs[0] == "fg" {
				fgs = append(fgs, subs[1])
			}
			if subs[0] == "bg" {
				bgs = append(bgs, subs[1])
			}
		}
	}

	fg = updateAttr(fg, fgs)
	bg = updateAttr(bg, bgs)
	return fg, bg
}

func (mtb *MarkdownTxBuilder) reset() {
	mtb.plainTx = []rune{}
	mtb.markers = []marker{}
}

// parse streams and parses text into normalized text and render sequence.
func (mtb *MarkdownTxBuilder) parse(str string) {
	rs := str2runes(str)
	normTx := []rune{}
	square := []rune{}
	brackt := []rune{}
	accSquare := false
	accBrackt := false
	cntSquare := 0

	reset := func() {
		square = []rune{}
		brackt = []rune{}
		accSquare = false
		accBrackt = false
		cntSquare = 0
	}
	// pipe stacks into normTx and clear
	rollback := func() {
		normTx = append(normTx, square...)
		normTx = append(normTx, brackt...)
		reset()
	}
	// chop first and last
	chop := func(s []rune) []rune {
		return s[1 : len(s)-1]
	}

	for i, r := range rs {
		switch {
		// stacking brackt
		case accBrackt:
			brackt = append(brackt, r)
			if ')' == r {
				fg, bg := mtb.readAttr(string(chop(brackt)))
				st := len(normTx)
				ed := len(normTx) + len(square) - 2
				mtb.markers = append(mtb.markers, marker{st, ed, fg, bg})
				normTx = append(normTx, chop(square)...)
				reset()
			} else if i+1 == len(rs) {
				rollback()
			}
		// stacking square
		case accSquare:
			switch {
			// squares closed and followed by a '('
			case cntSquare == 0 && '(' == r:
				accBrackt = true
				brackt = append(brackt, '(')
			// squares closed but not followed by a '('
			case cntSquare == 0:
				rollback()
				if '[' == r {
					accSquare = true
					cntSquare = 1
					brackt = append(brackt, '[')
				} else {
					normTx = append(normTx, r)
				}
			// hit the end
			case i+1 == len(rs):
				square = append(square, r)
				rollback()
			case '[' == r:
				cntSquare++
				square = append(square, '[')
			case ']' == r:
				cntSquare--
				square = append(square, ']')
			// normal char
			default:
				square = append(square, r)
			}
		// stacking normTx
		default:
			if '[' == r {
				accSquare = true
				cntSquare = 1
				square = append(square, '[')
			} else {
				normTx = append(normTx, r)
			}
		}
	}

	mtb.plainTx = normTx
}

// BuildWrap implements TextBuilder interface and will naively wrap the plain
// text string to length specified by wl while preserving the fg and bg
// attributes
func (mtb MarkdownTxBuilder) BuildWrap(s string, fg, bg Attribute, wl uint) []Cell {

	// get the []Cell from plain Build
	tmpCell := mtb.Build(s, fg, bg)

	// get the plaintext
	mtb.parse(s)
	plain := string(mtb.plainTx)

	// wrap
	plainWrapped := wordwrap.WrapString(plain, wl)

	// find differences and insert
	finalCell := tmpCell // finalcell will get the inserts and is what is returned

	plainRune := []rune(plain)
	plainWrappedRune := []rune(plainWrapped)
	trigger := "go"
	plainRuneNew := plainRune

	for trigger != "stop" {
		plainRune = plainRuneNew
		for i, _ := range plainRune {
			if plainRune[i] == plainWrappedRune[i] {
				trigger = "stop"
			} else if plainRune[i] != plainWrappedRune[i] && plainWrappedRune[i] == 10 {
				trigger = "go"
				cell := Cell{10, 0, 0}

				// insert a cell into the []Cell in correct position
				tmpCell = append(tmpCell, Cell{0, 0, 0})
				copy(tmpCell[i+1:], tmpCell[i:])
				tmpCell[i] = cell

				// insert the newline into plain so we avoid indexing errors
				plainRuneNew = append(plainRune, 10)
				copy(plainRuneNew[i+1:], plainRuneNew[i:])
				plainRuneNew[i] = plainWrappedRune[i]

				// restart the inner for loop until plain and plain wrapped are
				// the same; yeah, it's inefficient, but the text amounts
				// should be small
				break

			} else if plainRune[i] != plainWrappedRune[i] &&
				plainWrappedRune[i-1] == 10 && // if the prior rune is a newline
				plainRune[i] == 32 { // and this rune is a space
				trigger = "go"
				// need to delete plainRune[i] because it gets rid of an extra
				// space
				plainRuneNew = append(plainRune[:i], plainRune[i+1:]...)
				break

			} else {
				trigger = "stop" // stops the outer for loop
			}
		}
	}

	finalCell = tmpCell

	return finalCell
}

// Build implements TextBuilder interface.
func (mtb MarkdownTxBuilder) Build(s string, fg, bg Attribute) []Cell {
	mtb.baseFg = fg
	mtb.baseBg = bg
	mtb.reset()
	mtb.parse(s)
	cs := make([]Cell, len(mtb.plainTx))
	for i := range cs {
		cs[i] = Cell{Ch: mtb.plainTx[i], Fg: fg, Bg: bg}
	}
	for _, mrk := range mtb.markers {
		for i := mrk.st; i < mrk.ed; i++ {
			cs[i].Fg = mrk.fg
			cs[i].Bg = mrk.bg
		}
	}

	return cs
}

// NewMarkdownTxBuilder returns a TextBuilder employing markdown syntax.
func NewMarkdownTxBuilder() TextBuilder {
	return MarkdownTxBuilder{}
}
