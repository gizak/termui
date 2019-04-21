// <Copyright> 2019 Simon Robin Lehn. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package exp

import (
	"fmt"
	"bytes"
	"errors"
	"image/png"
	"crypto/md5"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)


var (
	tempdir string
)

func init() {
	widgets.RegisterDrawer(
		"urxvt",
		widgets.Drawer{
			Remote:         false,
			IsEscapeString: true,
			Available:      func() bool {return isUrxvt},
			Draw:           drawUrxvt,
		},
	)
}

func drawUrxvt(wdgt *widgets.Image, buf *Buffer) (err error) {
	if !isUrxvt {
		return errors.New("method not supported for this terminal type")
	}

	// wdgt.Block.Draw(buf)

	var widthPercentage, heightPercentage, CenterPosXPercentage, CenterPosYPercentage int
	termWidth, termHeight := getTermSizeInChars(true)
	if termWidth == 0 || termHeight == 0 {
		return errors.New("could not query terminal dimensions")
	}

	widthPercentage  = (100*wdgt.Inner.Dx())/termWidth
	heightPercentage = (100*wdgt.Inner.Dy())/termHeight
	maxX := wdgt.Inner.Max.X
	maxY := wdgt.Inner.Max.Y
	if termWidth < maxX {
		maxX = termWidth
	}
	if termHeight < maxY {
		maxY = termHeight
	}
	CenterPosXPercentage = 50*(wdgt.Inner.Min.X+maxX)/termWidth
	CenterPosYPercentage = 50*(wdgt.Inner.Min.Y+maxY)/termHeight

	img, changed, err := resizeImage(wdgt, buf)
	if !changed || err != nil {
		return err
	}

	bytBuf := new(bytes.Buffer)
	if err = png.Encode(bytBuf, img); err != nil {
		return errors.New("image encoding failed")
	}

	if fi, err := os.Stat(tempdir); err != nil || !fi.IsDir() {
		if tempdir, err = ioutil.TempDir("", "termui."); err != nil {
			return err
		}
	}

	// defer os.RemoveAll(dir) // clean up

	filename := filepath.Join(tempdir, fmt.Sprintf("urxvt-%x", md5.Sum(bytBuf.Bytes())) + ".png")
	if err := ioutil.WriteFile(filename, bytBuf.Bytes(), 0644); err != nil {
		return err
	}

	// "op=keep-aspect" maintains the image aspect ratio when scaling
	wdgt.Block.ANSIString = wrap(fmt.Sprintf("\033]20;%s;%dx%d+%d+%d:op=keep-aspect\a", filename, widthPercentage, heightPercentage, CenterPosXPercentage, CenterPosYPercentage))
	return nil
}
