// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Copyright 2018,2019 Simon R. Lehn. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png" // for encoding for iTerm2
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/mattn/go-sixel"
	"github.com/mattn/go-tty"

	. "github.com/gizak/termui/v3"
)

type Image struct {
	Block
	Image               image.Image
	Monochrome          bool
	MonochromeThreshold uint8
	MonochromeInvert    bool
}

var (
	sixelCapable, isIterm2                      bool
	charBoxWidthInPixels, charBoxHeightInPixels float64
	lastImageDimensions                         image.Rectangle
)

func init() {
	// example query: "\033[0c"
	// possible answer from the terminal (here xterm): "\033[[?63;1;2;4;6;9;15;22c"
	// the "4" signals that the terminal is capable of sixel
	termCapabilities := queryTerm("\033[0c")
	for i, cap := range termCapabilities {
		if i == 0 || i == len(termCapabilities)-1 {
			continue
		}
		if string(cap) == `4` {
			sixelCapable = true

			// terminal character box size measured in pixels
			charBoxWidthInPixels, charBoxHeightInPixels = getTermCharBoxSize()
		}
	}
	// # https://superuser.com/a/683971
	if os.Getenv("TERM_PROGRAM") == "iTerm.app" {
		isIterm2 = true
	}
}

func NewImage(img image.Image) *Image {
	return &Image{
		Block:               *NewBlock(),
		MonochromeThreshold: 128,
		Image:               img,
	}
}

func (self *Image) Draw(buf *Buffer) {
	// draw with ANSI escape strings
	// sixel / iTerm2
	if sixelCapable || isIterm2 {
		////if true {
		if err := self.drawANSI(buf); err == nil {
			return
		}
	}

	// urxvt pixbuf / ...

	// fall back - draw with box characters
	self.drawFallBack(buf)
}

func (self *Image) drawFallBack(buf *Buffer) {
	self.Block.Draw(buf)

	if self.Image == nil {
		return
	}

	bufWidth := self.Inner.Dx()
	bufHeight := self.Inner.Dy()
	imageWidth := self.Image.Bounds().Dx()
	imageHeight := self.Image.Bounds().Dy()

	if self.Monochrome {
		if bufWidth > imageWidth/2 {
			bufWidth = imageWidth / 2
		}
		if bufHeight > imageHeight/2 {
			bufHeight = imageHeight / 2
		}
		for bx := 0; bx < bufWidth; bx++ {
			for by := 0; by < bufHeight; by++ {
				ul := self.colorAverage(
					2*bx*imageWidth/bufWidth/2,
					(2*bx+1)*imageWidth/bufWidth/2,
					2*by*imageHeight/bufHeight/2,
					(2*by+1)*imageHeight/bufHeight/2,
				)
				ur := self.colorAverage(
					(2*bx+1)*imageWidth/bufWidth/2,
					(2*bx+2)*imageWidth/bufWidth/2,
					2*by*imageHeight/bufHeight/2,
					(2*by+1)*imageHeight/bufHeight/2,
				)
				ll := self.colorAverage(
					2*bx*imageWidth/bufWidth/2,
					(2*bx+1)*imageWidth/bufWidth/2,
					(2*by+1)*imageHeight/bufHeight/2,
					(2*by+2)*imageHeight/bufHeight/2,
				)
				lr := self.colorAverage(
					(2*bx+1)*imageWidth/bufWidth/2,
					(2*bx+2)*imageWidth/bufWidth/2,
					(2*by+1)*imageHeight/bufHeight/2,
					(2*by+2)*imageHeight/bufHeight/2,
				)
				buf.SetCell(
					NewCell(blocksChar(ul, ur, ll, lr, self.MonochromeThreshold, self.MonochromeInvert)),
					image.Pt(self.Inner.Min.X+bx, self.Inner.Min.Y+by),
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
				c := self.colorAverage(
					bx*imageWidth/bufWidth,
					(bx+1)*imageWidth/bufWidth,
					by*imageHeight/bufHeight,
					(by+1)*imageHeight/bufHeight,
				)
				buf.SetCell(
					NewCell(c.ch(), NewStyle(c.fgColor(), ColorBlack)),
					image.Pt(self.Inner.Min.X+bx, self.Inner.Min.Y+by),
				)
			}
		}
	}
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

func (self *Image) drawANSI(buf *Buffer) (err error) {
	self.Block.Draw(buf)

	// get dimensions //
	// terminal size measured in cells
	imageWidthInColumns := self.Inner.Dx()
	imageHeightInRows := self.Inner.Dy()

	// calculate image size in pixels
	imageWidthInPixels := int(float64(imageWidthInColumns) * charBoxWidthInPixels)
	imageHeightInPixels := int(float64(imageHeightInRows) * charBoxHeightInPixels)
	if imageWidthInPixels == 0 || imageHeightInPixels == 0 {
		return fmt.Errorf("could not calculate the image size in pixels")
	}

	termWidthInColumns, termHeightInRows := getTermSizeInChars()

	// handle only partially displayed image
	// otherwise we get scrolling
	var needsCrop bool
	imgCroppedWidth := imageWidthInPixels
	imgCroppedHeight := imageHeightInPixels
	if self.Max.X > int(termWidthInColumns)+1 {
		imgCroppedWidth = int(float64(int(termWidthInColumns)-self.Inner.Min.X-1) * charBoxWidthInPixels)
		needsCrop = true
	}
	if self.Max.Y > int(termHeightInRows)+1 {
		imgCroppedHeight = int(float64(int(termHeightInRows)-self.Inner.Min.Y-1) * charBoxHeightInPixels)
		needsCrop = true
	}

	// this is meant for comparison and for positioning in the ANSI string
	// the Min values are in cells while the Max values are in pixels
	imageDimensions := image.Rectangle{Min: image.Point{X: self.Inner.Min.X + 1, Y: self.Inner.Min.Y + 1}, Max: image.Point{X: imgCroppedWidth, Y: imgCroppedHeight}}
	// print saved ANSI string if image size and position didn't change
	if imageDimensions == lastImageDimensions {
		return nil
	}

	// resize and crop the image //
	img := imaging.Resize(self.Image, imageWidthInPixels, imageHeightInPixels, imaging.Lanczos)
	if needsCrop {
		img = imaging.Crop(img, image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: imgCroppedWidth, Y: imgCroppedHeight}})
	}

	if img.Bounds().Dx() == 0 || img.Bounds().Dy() == 0 {
		return fmt.Errorf("image size in pixels is 0")
	}

	// iTerm2
	// https://www.iterm2.com/documentation-images.html
	if isIterm2 {
		buf := new(bytes.Buffer)
		if err = png.Encode(buf, img); err != nil {
			goto skipIterm2
		}
		imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
		nameBase64 := base64.StdEncoding.EncodeToString([]byte(self.Block.Title))
		// 0 for stretching - 1 for no stretching
		noStretch := 0
		// for width, height:   "auto"   ||   N: N character cells   ||   Npx: N pixels   ||   N%: N percent of terminal width/height
		self.Block.ANSIString = fmt.Sprintf("\033]1337;File=name=%s;inline=1;height=%d;width=%d;preserveAspectRatio=%d:%s\a", nameBase64, imageDimensions.Max.Y, nameBase64, imageDimensions.Max.X, noStretch, imgBase64)

		return nil
	}
skipIterm2:

	if sixelCapable {
		byteBuf := new(bytes.Buffer)
		enc := sixel.NewEncoder(byteBuf)
		enc.Dither = true
		if err := enc.Encode(img); err != nil {
			return err
		}

		// position where the image should appear (upper left corner)
		self.Block.ANSIString = fmt.Sprintf("\033[%d;%dH%s", imageDimensions.Min.Y, imageDimensions.Min.X, byteBuf.String())
		// test string "HI"
		// self.Block.ANSIString = fmt.Sprintf("\033[%d;%dH%s", self.Inner.Min.Y+1, self.Inner.Min.X+1, "\033Pq#0;2;0;0;0#1;2;100;100;0#2;2;0;100;0#1~~@@vv@@~~@@~~$#2??}}GG}}??}}??-#1!14@\033\\")

		return nil
	}

	return errors.New("no method applied for ANSI drawing")
}

func getTermCharBoxSize() (x, y float64) {
	if cx, cy := getTermSizeInChars(); cx != 0 && cy != 0 {
		px, py := getTermSizeInPixels()
		x = float64(px) / float64(cx)
		y = float64(py) / float64(cy)
	}
	return
}

func getTermSizeInChars() (x, y uint) {
	// query terminal size in character boxes
	// answer: <termHeightInRows>;<termWidthInColumns>t
	q := queryTerm("\033[18t")

	if len(q) != 3 {
		return
	}

	if yy, err := strconv.Atoi(string(q[1])); err == nil {
		if xx, err := strconv.Atoi(string(q[2])); err == nil {
			x = uint(xx)
			y = uint(yy)
		} else {
			return
		}
	} else {
		return
	}

	return
}

func getTermSizeInPixels() (x, y uint) {
	// query terminal size in pixels
	// answer: <termHeightInPixels>;<termWidthInPixels>t
	q := queryTerm("\033[14t")

	if len(q) != 3 {
		return
	}

	if yy, err := strconv.Atoi(string(q[1])); err == nil {
		if xx, err := strconv.Atoi(string(q[2])); err == nil {
			x = uint(xx)
			y = uint(yy)
		} else {
			return
		}
	} else {
		return
	}

	return
}

func queryTerm(qs string) (ret [][]rune) {
	// temporary fix for xterm
	// otherwise TUI wouldn't react to any further events
	// resizing still works though
	if len(os.Getenv("XTERM_VERSION")) > 0 {
		return
	}

	var b []rune

	tty, err := tty.Open()
	if err != nil {
		return
	}
	defer tty.Close()

	ch := make(chan bool, 1)

	go func() {
		// query terminal
		fmt.Printf(qs)

		for {
			r, err := tty.ReadRune()
			if err != nil {
				return
			}
			// handle key event
			switch r {
			case 'c', 't':
				ret = append(ret, b)
				goto afterLoop
			case '?', ';':
				ret = append(ret, b)
				b = []rune{}
			default:
				b = append(b, r)
			}
		}
	afterLoop:
		ch <- true
	}()

	timer := time.NewTimer(50 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-ch:
		defer close(ch)
	case <-timer.C:
	}
	return
}
