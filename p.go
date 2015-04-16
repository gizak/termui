// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

// Par displays a paragraph.
/*
  par := termui.NewPar("Simple Text")
  par.Height = 3
  par.Width = 17
  par.Border.Label = "Label"
*/
type Par struct {
	Block
	Text            string
	TextFgColor     Attribute
	TextBgColor     Attribute
	RendererFactory TextRendererFactory
}

// NewPar returns a new *Par with given text as its content.
func NewPar(s string) *Par {
	return &Par{
		Block:           *NewBlock(),
		Text:            s,
		TextFgColor:     theme.ParTextFg,
		TextBgColor:     theme.ParTextBg,
		RendererFactory: PlainRendererFactory{},
	}
}

// Buffer implements Bufferer interface.
func (p *Par) Buffer() []Point {
	ps := p.Block.Buffer()

	fg, bg := p.TextFgColor, p.TextBgColor
	sequence := p.RendererFactory.TextRenderer(p.Text).Render(fg, bg)
	runes := []rune(sequence.NormalizedText)

	y, x, n := 0, 0, 0
	for y < p.innerHeight && n < len(runes) {
		point, width := sequence.PointAt(n, x+p.innerX, y+p.innerY)

		if runes[n] == '\n' || x+width > p.innerWidth {
			y++
			x = 0 // set x = 0
			if runes[n] == '\n' {
				n++
			}

			if y >= p.innerHeight {
				ps = append(ps, newPointWithAttrs('…',
					p.innerX+p.innerWidth-1,
					p.innerY+p.innerHeight-1,
					p.TextFgColor, p.TextBgColor))
				break
			}

			continue
		}

		ps = append(ps, point)
		n++
		x += width
	}

	return p.Block.chopOverflow(ps)
}
