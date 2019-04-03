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
	"strings"
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
	ASCIIonly           bool
}

var (
	sixelCapable, isIterm2, isXterm, isTmux, isScreen, isMuxed   bool
	charBoxWidthInPixels, charBoxHeightInPixels   float64
	charBoxWidthColumns,  charBoxHeightRows       int
	lastImageDimensions                           image.Rectangle
)

func init() {
	initiate()
}
func initiate() {
	if len(os.Getenv("XTERM_VERSION")) > 0                                         { isXterm  = true } else { isXterm  = false }
	if os.Getenv("TERM_PROGRAM") == "iTerm.app"                                    { isIterm2 = true } else { isIterm2 = false } // # https://superuser.com/a/683971
	// if len(os.Getenv("MINTTY_SHORTCUT")) > 0                                    { isMintty = true } else { isMintty = false } // doesn't work
	// if len(os.Getenv("MLTERM")) > 0                                             { isMlterm = true } else { isMlterm = false }
	if strings.HasPrefix(os.Getenv("TERM"), "screen") && len(os.Getenv("STY")) > 0 { isScreen = true } else { isScreen = false }
	if (strings.HasPrefix(os.Getenv("TERM"), "screen") || strings.HasPrefix(os.Getenv("TERM"), "tmux")) &&
	   len(os.Getenv("TMUX")) > 0 || len(os.Getenv("TMUX_PANE")) > 0               { isTmux   = true } else { isTmux   = false }
	if isTmux || isScreen                                                          { isMuxed  = true } else { isMuxed  = false }

	// example query: "\033[0c"
	// possible answer from the terminal (here xterm): "\033[[?63;1;2;4;6;9;15;22c", vte(?): ...62,9;c
	// the "4" signals that the terminal is capable of sixel
	// conhost.exe knows this sequence.
	sixelCapable = false
	termCapabilities := queryTerm(wrap("\033[0c"))
	for i, cap := range termCapabilities {
		if i == 0 || i == len(termCapabilities) - 1 {
			continue
		}
		if string(cap) == `4` {
			sixelCapable = true
		}
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
	// possible reattachments of the terminal multiplexer?
	if isMuxed {
		initiate()
	}

	// fall back - draw with box characters
	// possible enhancement: make use of further box characters like chafa:
	// https://hpjansson.org/chafa/
	// https://github.com/hpjansson/chafa/
	self.drawFallBack(buf)

	// overdraw with ANSI escape strings
	// sixel / iTerm2
	if !self.ASCIIonly && (sixelCapable || isIterm2) {
		////if true {
		if err := self.drawANSI(buf); err == nil {
			return
		}
	}
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
	imageHeightInRows   := self.Inner.Dy()

	// terminal size in cells and pixels and calculated terminal character box size in pixels
	var termWidthInColumns, termHeightInRows int
	var charBoxWidthInPixelsTemp, charBoxHeightInPixelsTemp float64
	termWidthInColumns, termHeightInRows, _, _, charBoxWidthInPixelsTemp, charBoxHeightInPixelsTemp, err = getTermSize()
	if err != nil {
		return err
	}
	// update if value is more precise
	if termWidthInColumns > charBoxWidthColumns {
		charBoxWidthInPixels = charBoxWidthInPixelsTemp
	}
	if termHeightInRows > charBoxHeightRows {
		charBoxHeightInPixels = charBoxHeightInPixelsTemp
	}
if isTmux {charBoxWidthInPixels, charBoxHeightInPixels = 10, 19}   // mlterm settings (temporary)

	// calculate image size in pixels
	// subtract 1 pixel for small deviations from char box size (float64)
	imageWidthInPixels  := int(float64(imageWidthInColumns) * charBoxWidthInPixels)  - 1
	imageHeightInPixels := int(float64(imageHeightInRows)   * charBoxHeightInPixels) - 1
	if imageWidthInPixels == 0 || imageHeightInPixels == 0 {
		return fmt.Errorf("could not calculate the image size in pixels")
	}

	// handle only partially displayed image
	// otherwise we get scrolling
	var needsCropX, needsCropY bool
	var imgCroppedWidth, imgCroppedHeight int
	imgCroppedWidth  = imageWidthInPixels
	imgCroppedHeight = imageHeightInPixels
	if self.Max.Y >= int(termHeightInRows) {
		var scrollExtraRows int
		// remove last 2 rows for xterm when cropped vertically to prevent scrolling
		if isXterm {
			scrollExtraRows = 2
		}
		// subtract 1 pixel for small deviations from char box size (float64)
		imgCroppedHeight = int(float64(int(termHeightInRows) - self.Inner.Min.Y - scrollExtraRows) * charBoxHeightInPixels) - 1
		needsCropY = true
	}
	if self.Max.X >= int(termWidthInColumns) {
		var scrollExtraColumns int
		imgCroppedWidth = int(float64(int(termWidthInColumns) - self.Inner.Min.X - scrollExtraColumns) * charBoxWidthInPixels) - 1
		needsCropX = true
	}

	// this is meant for comparison and for positioning in the ANSI string
	// the Min values are in cells while the Max values are in pixels
	imageDimensions := image.Rectangle{Min: image.Point{X: self.Inner.Min.X + 1, Y: self.Inner.Min.Y + 1}, Max: image.Point{X: imgCroppedWidth, Y: imgCroppedHeight}}
	// print saved ANSI string if image size and position didn't change
	if imageDimensions.Min.X == lastImageDimensions.Min.X && imageDimensions.Min.Y == lastImageDimensions.Min.Y && imageDimensions.Max.X == lastImageDimensions.Max.X && imageDimensions.Max.Y == lastImageDimensions.Max.Y {
		// reuse old ANSIString value because of unchanged image dimensions
		return nil
	}
	lastImageDimensions = imageDimensions

	// resize and crop the image //
	img := imaging.Resize(self.Image, imageWidthInPixels, imageHeightInPixels, imaging.Lanczos)
	if needsCropX || needsCropY {
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
		iterm2String := wrap(fmt.Sprintf("\033[?8452h\033]1337;File=name=%s;inline=1;height=%d;width=%d;preserveAspectRatio=%d:%s\a", nameBase64, imageDimensions.Max.Y, nameBase64, imageDimensions.Max.X, noStretch, imgBase64))
		// for width, height:   "auto"   ||   N: N character cells   ||   Npx: N pixels   ||   N%: N percent of terminal width/height
		self.Block.ANSIString = fmt.Sprintf("\033[%d;%dH%s", imageDimensions.Min.Y, imageDimensions.Min.X, iterm2String)

		return nil
	}
	skipIterm2:

	// possible enhancements:
	// kitty https://sw.kovidgoyal.net/kitty/graphics-protocol.html
	// Terminology (from Enlightenment) https://www.enlightenment.org/docs/apps/terminology.md#tycat https://github.com/billiob/terminology
	// urxvt pixbuf / ...
	//
	// Tektronix 4014, ReGis

	// sixel
	// https://vt100.net/docs/vt3xx-gp/chapter14.html
	if sixelCapable {
		byteBuf := new(bytes.Buffer)
		enc := sixel.NewEncoder(byteBuf)
		enc.Dither = true
		if err := enc.Encode(img); err != nil {
			return err
		}
		sixelString := wrap("\033[?8452h" + byteBuf.String())
		// position where the image should appear (upper left corner) + sixel
		// https://github.com/mintty/mintty/wiki/CtrlSeqs#sixel-graphics-end-position
		// "\033[?8452h" sets the cursor next right to the bottom of the image instead of below
		// this prevents vertical scrolling when the image fills the last line.
		// horizontal scrolling because of this did not happen in my test cases.
		// "\033[?80l" disables sixel scrolling if it isn't already.
		self.Block.ANSIString = fmt.Sprintf("\033[%d;%dH%s", imageDimensions.Min.Y, imageDimensions.Min.X, sixelString)
		// test string "HI"
		// self.Block.ANSIString = fmt.Sprintf("\033[%d;%dH\033[?8452h%s", self.Inner.Min.Y+1, self.Inner.Min.X+1, "\033Pq#0;2;0;0;0#1;2;100;100;0#2;2;0;100;0#1~~@@vv@@~~@@~~$#2??}}GG}}??}}??-#1!14@\033\\")

		return nil
	}

	return errors.New("no method applied for ANSI drawing")
}

func getTermSize() (termWidthInColumns, termHeightInRows, termWidthInPixels, termHeightInPixels int, charBoxWidthInPixels, charBoxHeightInPixels float64, err error) {
	// this uses a combination of TIOCGWINSZ and \033[14t , \033[18t
	// the syscall to TIOCGWINSZ only works locally

	var cx, cy, px, py int
	err = nil

	t, err := tty.Open()
	defer t.Close()

	cx, cy, px, py, err = t.SizePixel()
	// if isMuxed {
	if false {   // temporary
		// in case of split view it is better to query from the same source
		/*if cx <= 0 || cy <= 0 || px <= 0 || py <= 0 {
			cx, cy = getTermSizeInChars(true)
			px, py = getTermSizeInPixels(true)
			if cx <= 0 || cy <= 0 || px <= 0 || py <= 0 {
				cx, cy = getTermSizeInChars(false)
				px, py = getTermSizeInPixels(false)
				if cx <= 0 || cy <= 0 || px <= 0 || py <= 0 {
					return
				}	
			}	
		}*/	
	} else if err == nil {
		if cx > 0 && cy > 0 {
			if px <= 0 || py <= 0 {
				px, py = getTermSizeInPixels(true)
			}
		} else {
			if cx, cy = getTermSizeInChars(true); cx != 0 && cy != 0 {
				if px <= 0 || py <= 0 {
					px, py = getTermSizeInPixels(true)
				}
			} else {
				return
			}
		}
	}

	termWidthInColumns    = cx
	termHeightInRows      = cy
	termWidthInPixels     = px
	termHeightInPixels    = py
	charBoxWidthInPixels  = float64(px) / float64(cx)
	charBoxHeightInPixels = float64(py) / float64(cy)
	return
}

func getTermSizeInChars(needsWrap bool) (x, y int) {
	// query terminal size in character boxes
	// answer: <termHeightInRows>;<termWidthInColumns>t
	s := "\033[18t"
	if needsWrap {
		s = wrap(s)
	}
	q := queryTerm(s)
	if len(q) != 3 {
		return
	}

	if yy, err := strconv.Atoi(string(q[1])); err == nil {
		if xx, err := strconv.Atoi(string(q[2])); err == nil {
			x = xx
			y = yy
		} else {
			return
		}
	} else {
		return
	}

	return
}

func getTermSizeInPixels(needsWrap bool) (x, y int) {
	// query terminal size in pixels
	// answer: <termHeightInPixels>;<termWidthInPixels>t
	s := "\033[14t"
	if needsWrap {
		s = wrap(s)
	}
	q := queryTerm(s)

	if len(q) != 3 {
		return
	}

	if yy, err := strconv.Atoi(string(q[1])); err == nil {
		if xx, err := strconv.Atoi(string(q[2])); err == nil {
			x = xx
			y = yy
		} else {
			return
		}
	} else {
		return
	}

	return
}

func queryTerm(qs string) (ret [][]rune) {
	// temporary fix for xterm - not completely sure if still needed
	// otherwise TUI wouldn't react to any further events
	// resizing still works though
	if isXterm && qs != "\033[0c" && qs != wrap("\033[0c") {
		return
	}

	var b []rune

	t, err := tty.Open()
	if err != nil {
		return
	}

	ch := make(chan bool, 1)

	go func() {
		defer t.Close()
		// query terminal
		fmt.Printf(qs)

		for {
			r, err := t.ReadRune()
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

	var timer *time.Timer
	if isTmux {
		// tmux needs a bit more time
		timer = time.NewTimer(50000 * time.Microsecond)
	} else {
		// on my system the terminals mlterm, xterm need at least around 100 microseconds
		timer = time.NewTimer(500 * time.Microsecond)
	}
	defer timer.Stop()

	select {
	case <-ch:
		defer close(ch)
	case <-timer.C:
	}
	return
}

func wrap(s string) string {
	if !isMuxed {
		return s
	}
	if isTmux {
		return tmuxWrap(s)
	}
	return s
}

func tmuxWrap(s string) string {
	return "\033Ptmux;" + strings.Replace(s, "\033", "\033\033", -1) + "\033\\"
}

/*
// https://savannah.gnu.org/bugs/index.php?56063
func screenWrap(s string) string {}
*/
