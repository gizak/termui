// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

// List displays []string as its items,
// it has a Overflow option (default is "hidden"), when set to "hidden",
// the item exceeding List's width is truncated, but when set to "wrap",
// the overflowed text breaks into next line.
/*
  strs := []string{
		"[0] github.com/gizak/termui",
		"[1] editbox.go",
		"[2] iterrupt.go",
		"[3] keyboard.go",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] nsf/termbox-go"}

  ls := termui.NewList()
  ls.Items = strs
  ls.ItemFgColor = termui.ColorYellow
  ls.Border.Label = "List"
  ls.Height = 7
  ls.Width = 25
  ls.Y = 0
*/
type List struct {
	Block
	Items           []string
	Overflow        string
	ItemFgColor     Attribute
	ItemBgColor     Attribute
	RendererFactory TextRendererFactory
}

// NewList returns a new *List with current theme.
func NewList() *List {
	l := &List{Block: *NewBlock()}
	l.Overflow = "hidden"
	l.ItemFgColor = theme.ListItemFg
	l.ItemBgColor = theme.ListItemBg
	l.RendererFactory = NoopRendererFactory{}
	return l
}

// Buffer implements Bufferer interface.
func (l *List) Buffer() []Point {
	ps := l.Block.Buffer()
	switch l.Overflow {
	case "wrap":
		y := 0
		for _, item := range l.Items {
			x := 0

			renderer := l.RendererFactory.TextRenderer(item)
			sequence := renderer.Render(l.ItemFgColor, l.ItemBgColor)
			for n := range []rune(sequence.NormalizedText) {
				point, width := sequence.PointAt(n, x+l.innerX, y+l.innerY)

				if width+x <= l.innerWidth {
					ps = append(ps, point)
					x += width
				} else {
					y++
					x = 0
				}
			}
			y++
		}

	case "hidden":
		trimItems := l.Items
		if len(trimItems) > l.innerHeight {
			trimItems = trimItems[:l.innerHeight]
		}

		for y, item := range trimItems {
			text := TrimStrIfAppropriate(item, l.innerWidth)
			render := l.RendererFactory.TextRenderer(text)
			sequence := render.RenderSequence(0, -1, l.ItemFgColor, l.ItemBgColor)
			t, _ := sequence.Buffer(l.innerX, y+l.innerY)
			ps = append(ps, t...)
		}
	}

	return l.Block.chopOverflow(ps)
}
