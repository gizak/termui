package termui

import "strings"

type List struct {
	Block
	Items       []string
	Overflow    string
	ItemFgColor Attribute
	ItemBgColor Attribute
}

func NewList() *List {
	l := &List{Block: *NewBlock()}
	l.Overflow = "hidden"
	return l
}

func (l *List) Buffer() []Point {
	ps := l.Block.Buffer()
	switch l.Overflow {
	case "wrap":
		rs := str2runes(strings.Join(l.Items, "\n"))
		i, j, k := 0, 0, 0
		for i < l.innerHeight && k < len(rs) {
			if rs[k] == '\n' || j == l.innerWidth {
				i++
				j = 0
				if rs[k] == '\n' {
					k++
				}
				continue
			}
			pi := Point{}
			pi.X = l.innerX + j
			pi.Y = l.innerY + i

			pi.Code.Ch = rs[k]
			pi.Code.Bg = toTmAttr(l.ItemBgColor)
			pi.Code.Fg = toTmAttr(l.ItemFgColor)

			ps = append(ps, pi)
			k++
			j++
		}

	case "hidden":
		trimItems := l.Items
		if len(trimItems) > l.innerHeight {
			trimItems = trimItems[:l.innerHeight]
		}
		for i, v := range trimItems {
			rs := trimStr2Runes(v, l.innerWidth)

			j := 0
			for _, vv := range rs {
				p := Point{}
				p.X = l.innerX + j
				p.Y = l.innerY + i

				p.Code.Ch = vv
				p.Code.Bg = toTmAttr(l.ItemBgColor)
				p.Code.Fg = toTmAttr(l.ItemFgColor)

				ps = append(ps, p)
				j++
			}
		}
	}
	return ps
}
