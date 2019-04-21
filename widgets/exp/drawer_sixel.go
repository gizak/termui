// <Copyright> 2018,2019 Simon Robin Lehn. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package exp

import (
	"fmt"
	"bytes"
	"errors"

	"github.com/mattn/go-sixel"

	. "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)


func init() {
	widgets.RegisterDrawer(
		"sixel",
		widgets.Drawer{
			Remote:         true,
			IsEscapeString: true,
			Available:      func() bool {return sixelCapable},
			Draw:           drawSixel,
		},
	)
}

func drawSixel(wdgt *widgets.Image, buf *Buffer) (err error) {
	wdgt.Block.Draw(buf)

	img, changed, err := resizeImage(wdgt, buf)
	if !changed || err != nil {
		return err
	}

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
		wdgt.Block.ANSIString = fmt.Sprintf("\033[%d;%dH%s", wdgt.Inner.Min.Y + 1, wdgt.Inner.Min.X + 1, sixelString)
		// test string "HI"
		// wdgt.Block.ANSIString = fmt.Sprintf("\033[%d;%dH\033[?8452h%s", wdgt.Inner.Min.Y+1, wdgt.Inner.Min.X+1, "\033Pq#0;2;0;0;0#1;2;100;100;0#2;2;0;100;0#1~~@@vv@@~~@@~~$#2??}}GG}}??}}??-#1!14@\033\\")

		return nil
	}

	return errors.New("no method applied for ANSI drawing")
}
