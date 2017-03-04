// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"
	"image/color"
)

// Image is an image widget.
type Image struct {
	Block
	Image               image.Image
	Monochrome          bool
	MonochromeThreshold uint8
	MonochromeInvert    bool
}

// NewImage returns a new image widget.
func NewImage(img image.Image) *Image {
	im := &Image{
		Block:               *NewBlock(),
		MonochromeThreshold: 128,
		Image:               img,
	}
	im.Width = 64
	im.Height = 48
	return im
}

// Image implements Bufferer.
func (im *Image) Buffer() Buffer {
	buf := im.Block.Buffer()
	bufWidth := im.innerArea.Dx()
	bufHeight := im.innerArea.Dy()
	for bx := 0; bx < bufWidth; bx++ {
		for by := 0; by < bufHeight; by++ {
			buf.Set(im.innerArea.Min.X+bx, im.innerArea.Min.Y+by, Cell{
				Ch: ' ',
				Fg: ColorDefault,
				Bg: ColorDefault,
			})
		}
	}
	if im.Image == nil {
		return buf
	}
	imageWidth := im.Image.Bounds().Dx()
	imageHeight := im.Image.Bounds().Dy()
	if im.Monochrome {
		if bufWidth > imageWidth/2 {
			bufWidth = imageWidth / 2
		}
		if bufHeight > imageHeight/2 {
			bufHeight = imageHeight / 2
		}
		for bx := 0; bx < bufWidth; bx++ {
			for by := 0; by < bufHeight; by++ {
				ul := im.colorAverage(2*bx*imageWidth/bufWidth/2, (2*bx+1)*imageWidth/bufWidth/2, 2*by*imageHeight/bufHeight/2, (2*by+1)*imageHeight/bufHeight/2)
				ur := im.colorAverage((2*bx+1)*imageWidth/bufWidth/2, (2*bx+2)*imageWidth/bufWidth/2, 2*by*imageHeight/bufHeight/2, (2*by+1)*imageHeight/bufHeight/2)
				ll := im.colorAverage(2*bx*imageWidth/bufWidth/2, (2*bx+1)*imageWidth/bufWidth/2, (2*by+1)*imageHeight/bufHeight/2, (2*by+2)*imageHeight/bufHeight/2)
				lr := im.colorAverage((2*bx+1)*imageWidth/bufWidth/2, (2*bx+2)*imageWidth/bufWidth/2, (2*by+1)*imageHeight/bufHeight/2, (2*by+2)*imageHeight/bufHeight/2)
				buf.Set(im.innerArea.Min.X+bx, im.innerArea.Min.Y+by, Cell{
					Ch: blocksChar(ul, ur, ll, lr, im.MonochromeThreshold, im.MonochromeInvert),
				})
			}
		}
	} else {
		if bufWidth > imageWidth {
			bufWidth = imageWidth
		}
		if bufHeight > imageHeight {
			bufHeight = imageHeight
		}
		for bx := 0; bx < bufWidth; bx++ {
			for by := 0; by < bufHeight; by++ {
				c := im.colorAverage(bx*imageWidth/bufWidth, (bx+1)*imageWidth/bufWidth, by*imageHeight/bufHeight, (by+1)*imageHeight/bufHeight)
				buf.Set(im.innerArea.Min.X+bx, im.innerArea.Min.Y+by, Cell{
					Ch: c.ch(),
					Fg: c.fgColor(),
					Bg: ColorBlack,
				})
			}
		}
	}
	return buf
}

func (im *Image) colorAverage(x0, x1, y0, y1 int) colorAverager {
	var c colorAverager
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			c = c.add(im.Image.At(x+im.Image.Bounds().Min.X, y+im.Image.Bounds().Min.Y))
		}
	}
	return c
}

type colorAverager struct {
	rsum, gsum, bsum, asum, count uint64
}

func (c colorAverager) add(col color.Color) colorAverager {
	r, g, b, a := col.RGBA()
	return colorAverager{
		rsum:  c.rsum + uint64(r),
		gsum:  c.gsum + uint64(g),
		bsum:  c.bsum + uint64(b),
		asum:  c.asum + uint64(a),
		count: c.count + 1,
	}
}

func (c colorAverager) RGBA() (uint32, uint32, uint32, uint32) {
	if c.count == 0 {
		return 0, 0, 0, 0
	} else {
		return uint32(c.rsum/c.count) & 0xffff,
			uint32(c.gsum/c.count) & 0xffff,
			uint32(c.bsum/c.count) & 0xffff,
			uint32(c.asum/c.count) & 0xffff
	}
}

func (c colorAverager) fgColor() Attribute {
	return palette.Convert(c).(paletteColor).attribute
}

func (c colorAverager) ch() rune {
	gray := color.GrayModel.Convert(c).(color.Gray).Y
	switch {
	case gray < 51:
		return ' '
	case gray < 102:
		return '░'
	case gray < 153:
		return '▒'
	case gray < 204:
		return '▓'
	default:
		return '█'
	}
}

func (c colorAverager) monochrome(threshold uint8, invert bool) bool {
	return c.count != 0 && (color.GrayModel.Convert(c).(color.Gray).Y < threshold != invert)
}

type paletteColor struct {
	rgba      color.RGBA
	attribute Attribute
}

func (c paletteColor) RGBA() (uint32, uint32, uint32, uint32) {
	return c.rgba.RGBA()
}

var palette = color.Palette([]color.Color{
	paletteColor{color.RGBA{0, 0, 0, 255}, ColorBlack},
	paletteColor{color.RGBA{255, 0, 0, 255}, ColorRed},
	paletteColor{color.RGBA{0, 255, 0, 255}, ColorGreen},
	paletteColor{color.RGBA{255, 255, 0, 255}, ColorYellow},
	paletteColor{color.RGBA{0, 0, 255, 255}, ColorBlue},
	paletteColor{color.RGBA{255, 0, 255, 255}, ColorMagenta},
	paletteColor{color.RGBA{0, 255, 255, 255}, ColorCyan},
	paletteColor{color.RGBA{255, 255, 255, 255}, ColorWhite},
})

var blocks = [...]rune{
	' ', '▘', '▝', '▀', '▖', '▌', '▞', '▛',
	'▗', '▚', '▐', '▜', '▄', '▙', '▟', '█',
}

func blocksChar(ul, ur, ll, lr colorAverager, threshold uint8, invert bool) rune {
	index := 0
	if ul.monochrome(threshold, invert) {
		index |= 1
	}
	if ur.monochrome(threshold, invert) {
		index |= 2
	}
	if ll.monochrome(threshold, invert) {
		index |= 4
	}
	if lr.monochrome(threshold, invert) {
		index |= 8
	}
	return blocks[index]
}
