// +build ignore

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
	l.RendererFactory = PlainRendererFactory{}
	return l
}

// Buffer implements Bufferer interface.
func (l *List) Buffer() []Point {
	buffer := l.Block.Buffer()

	breakLoop := func(y int) bool {
		return y+1 > l.innerArea.Dy()
	}
	y := 0

MainLoop:
	for _, item := range l.Items {
		x := 0
		bg, fg := l.ItemFgColor, l.ItemBgColor
		renderer := l.RendererFactory.TextRenderer(item)
		sequence := renderer.Render(bg, fg)

		for n := range []rune(sequence.NormalizedText) {
			point, width := sequence.PointAt(n, x+l.innerArea.Min.X, y+l.innerArea.Min.Y)

			if width+x <= l.innerArea.Dx() {
				buffer = append(buffer, point)
				x += width
			} else {
				if l.Overflow == "wrap" {
					y++
					if breakLoop(y) {
						break MainLoop
					}
					x = 0
				} else {
					dotR := []rune(dot)[0]
					dotX := l.innerArea.Dx() + l.innerArea.Min.X - charWidth(dotR)
					p := newPointWithAttrs(dotR, dotX, y+l.innerArea.Min.Y, bg, fg)
					buffer = append(buffer, p)
					break
				}
			}
		}

		y++
		if breakLoop(y) {
			break MainLoop
		}
	}

	return l.Block.chopOverflow(buffer)
}
