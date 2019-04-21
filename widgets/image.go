// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"image"
	"image/color"
	"sync"
	"sync/atomic"

	. "github.com/gizak/termui/v3"
)

type Drawer struct {
	Remote         bool   // if the drawer can be used for remote terminals
	IsEscapeString bool
	Available      func() (bool)
	Draw           func(*Image, *Buffer) (error)
}

var (
	drawersMu     sync.Mutex
	atomicDrawers atomic.Value
)

// from https://golang.org/src/image/format.go?s=1069:1193#L31
func RegisterDrawer(nameNew string, drawerNew Drawer) {
	drawersMu.Lock()
	// drawers, _ := atomicDrawers.Load().([]Drawer)
	// atomicDrawers.Store(append(drawers, dr))
	drawers, _ := atomicDrawers.Load().(map[string]Drawer)
	drawersNew := make(map[string]Drawer)
	for name, drawer := range drawers {
		drawersNew[name] = drawer
	}
	drawersNew[nameNew] = drawerNew
	atomicDrawers.Store(drawersNew)
	drawersMu.Unlock()
}

func GetDrawers() map[string]Drawer {
	if drawers, ok := atomicDrawers.Load().(map[string]Drawer); ok {
		return drawers
	}
	return map[string]Drawer{}
}

func init() {
	RegisterDrawer(
		"block",
		Drawer{
			Remote:         true,
			IsEscapeString: false,
			Available:      func() bool {return true},
			Draw:           drawBlocks,
		},
	)
}

type Image struct {
	Block
	Image                   image.Image
	Monochrome              bool
	MonochromeThreshold     uint8
	MonochromeInvert        bool
	visibleSubImagePixels   image.Rectangle
}

func NewImage(img image.Image) *Image {
	return &Image{
		Block:                 *NewBlock(),
		MonochromeThreshold:   128,
		Image:                 img,
		visibleSubImagePixels: image.Rectangle{},
	}
}

func (self *Image) Draw(buf *Buffer) {
	drawers := GetDrawers()

	// fall back - draw with box characters atomicDrawers.Load().(map[string]Drawer)]["blocks"]
	// possible enhancement: make use of further box characters like chafa:
	// https://hpjansson.org/chafa/
	// https://github.com/hpjansson/chafa/
	if drbl, ok := drawers["block"]; ok {
		drbl.Draw(self, buf)
	}

	for name, dr := range drawers {
		if name != "block" && dr.Available() {
			dr.Draw(self, buf)
		}
	}
}

func drawBlocks(img *Image, buf *Buffer) (err error) {
	img.Block.Draw(buf)

	if img.Image == nil {
		return
	}

	bufWidth := img.Inner.Dx()
	bufHeight := img.Inner.Dy()
	imageWidth := img.Image.Bounds().Dx()
	imageHeight := img.Image.Bounds().Dy()

	if img.Monochrome {
		if bufWidth > imageWidth/2 {
			bufWidth = imageWidth / 2
		}
		if bufHeight > imageHeight/2 {
			bufHeight = imageHeight / 2
		}
		for bx := 0; bx < bufWidth; bx++ {
			for by := 0; by < bufHeight; by++ {
				ul := img.colorAverage(
					2*bx*imageWidth/bufWidth/2,
					(2*bx+1)*imageWidth/bufWidth/2,
					2*by*imageHeight/bufHeight/2,
					(2*by+1)*imageHeight/bufHeight/2,
				)
				ur := img.colorAverage(
					(2*bx+1)*imageWidth/bufWidth/2,
					(2*bx+2)*imageWidth/bufWidth/2,
					2*by*imageHeight/bufHeight/2,
					(2*by+1)*imageHeight/bufHeight/2,
				)
				ll := img.colorAverage(
					2*bx*imageWidth/bufWidth/2,
					(2*bx+1)*imageWidth/bufWidth/2,
					(2*by+1)*imageHeight/bufHeight/2,
					(2*by+2)*imageHeight/bufHeight/2,
				)
				lr := img.colorAverage(
					(2*bx+1)*imageWidth/bufWidth/2,
					(2*bx+2)*imageWidth/bufWidth/2,
					(2*by+1)*imageHeight/bufHeight/2,
					(2*by+2)*imageHeight/bufHeight/2,
				)
				buf.SetCell(
					NewCell(blocksChar(ul, ur, ll, lr, img.MonochromeThreshold, img.MonochromeInvert)),
					image.Pt(img.Inner.Min.X+bx, img.Inner.Min.Y+by),
				)
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
				c := img.colorAverage(
					bx*imageWidth/bufWidth,
					(bx+1)*imageWidth/bufWidth,
					by*imageHeight/bufHeight,
					(by+1)*imageHeight/bufHeight,
				)
				buf.SetCell(
					NewCell(c.ch(), NewStyle(c.fgColor(), ColorBlack)),
					image.Pt(img.Inner.Min.X+bx, img.Inner.Min.Y+by),
				)
			}
		}
	}
	return
}

// measured in pixels
func (self *Image) SetVisibleArea(area image.Rectangle) {
	self.visibleSubImagePixels = area
}

// measured in pixels
func (self *Image) GetVisibleArea() image.Rectangle {
	return self.visibleSubImagePixels
}

func (self *Image) colorAverage(x0, x1, y0, y1 int) colorAverager {
	var c colorAverager
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			c = c.add(
				self.Image.At(
					x+self.Image.Bounds().Min.X,
					y+self.Image.Bounds().Min.Y,
				),
			)
		}
	}
	return c
}

type colorAverager struct {
	rsum, gsum, bsum, asum, count uint64
}

func (self colorAverager) add(col color.Color) colorAverager {
	r, g, b, a := col.RGBA()
	return colorAverager{
		rsum:  self.rsum + uint64(r),
		gsum:  self.gsum + uint64(g),
		bsum:  self.bsum + uint64(b),
		asum:  self.asum + uint64(a),
		count: self.count + 1,
	}
}

func (self colorAverager) RGBA() (uint32, uint32, uint32, uint32) {
	if self.count == 0 {
		return 0, 0, 0, 0
	}
	return uint32(self.rsum/self.count) & 0xffff,
		uint32(self.gsum/self.count) & 0xffff,
		uint32(self.bsum/self.count) & 0xffff,
		uint32(self.asum/self.count) & 0xffff
}

func (self colorAverager) fgColor() Color {
	return palette.Convert(self).(paletteColor).attribute
}

func (self colorAverager) ch() rune {
	gray := color.GrayModel.Convert(self).(color.Gray).Y
	switch {
	case gray < 51:
		return SHADED_BLOCKS[0]
	case gray < 102:
		return SHADED_BLOCKS[1]
	case gray < 153:
		return SHADED_BLOCKS[2]
	case gray < 204:
		return SHADED_BLOCKS[3]
	default:
		return SHADED_BLOCKS[4]
	}
}

func (self colorAverager) monochrome(threshold uint8, invert bool) bool {
	return self.count != 0 && (color.GrayModel.Convert(self).(color.Gray).Y < threshold != invert)
}

type paletteColor struct {
	rgba      color.RGBA
	attribute Color
}

func (self paletteColor) RGBA() (uint32, uint32, uint32, uint32) {
	return self.rgba.RGBA()
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
	return IRREGULAR_BLOCKS[index]
}
