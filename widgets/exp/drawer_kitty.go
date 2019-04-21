// <Copyright> 2019 Simon Robin Lehn. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package exp

import (
	"fmt"
	"errors"
	"bytes"
	"encoding/base64"
	"image/png"

	. "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	kittyLimit = 4096
)

var (
	// TODO: for numbering of ids
	kittyImageCount int
)

func init() {
	widgets.RegisterDrawer(
		"kitty",
		widgets.Drawer{
			Remote:         true,
			IsEscapeString: true,
			Available:      func() bool {return isKitty},
			Draw:           drawKitty,
		},
	)
}

func drawKitty(wdgt *widgets.Image, buf *Buffer) (err error) {
	if !isKitty {
		return errors.New("method not supported for this terminal type")
	}

	wdgt.Block.Draw(buf)

	//  TODO: FIX THIS
	termWidth, termHeight := getTermSizeInChars(true)
	var _ = termWidth
	/*
	if termWidth == 0 || termHeight == 0 {
		return errors.New("could not query terminal dimensions")
	}
	*/

	img, changed, err := resizeImage(wdgt, buf)
	if !changed || err != nil {
		return err
	}

	var imgHeight int
	imageDimensions := wdgt.GetVisibleArea()
	if wdgt.Inner.Max.Y < termHeight {
		imgHeight = wdgt.Inner.Dy()
	} else {
		imgHeight = termHeight-1
	}
	imgHeight = wdgt.Inner.Dy()   // TODO: REMOVE THIS CRUTCH

	// https://sw.kovidgoyal.net/kitty/graphics-protocol.html#remote-client
	// https://sw.kovidgoyal.net/kitty/graphics-protocol.html#png-data
	// https://sw.kovidgoyal.net/kitty/graphics-protocol.html#controlling-displayed-image-layout
	bytBuf := new(bytes.Buffer)
	if err = png.Encode(bytBuf, img); err != nil {
		return err
	}
	imgBase64 := base64.StdEncoding.EncodeToString(bytBuf.Bytes())
	lenImgB64 := len([]byte(imgBase64))
	// a=T           action
	// t=d           payload is (base64 encoded) data itself not a file location
	// f=100         format: 100 = PNG payload
	// o=z           data compression
	// X=...,Y=,,,   Upper left image corner in cell coordinates (starting with 1, 1)
	// c=...,r=...   image size in cell columns and rows
	// w=...,h=...   width & height (in pixels) of the image area to display   // TODO: Use this to let Kitty handle cropping!
	// z=0           z-index vertical stacking order of the image
	// m=[01]        0 last escape code chunk - 1 for all except the last
	var kittyString string
	var zIndex = 2   // draw over text
	settings := fmt.Sprintf("a=T,t=d,f=100,X=%d,Y=%d,c=%d,r=%d,z=%d,", imageDimensions.Min.X, imageDimensions.Min.Y, wdgt.Inner.Dx(), imgHeight, zIndex)
	i := 0
	for ; i < (lenImgB64-1)/kittyLimit; i++ {
		kittyString += wrap(fmt.Sprintf("\033_G%sm=1;%s\033\\", settings, imgBase64[i*kittyLimit:(i+1)*kittyLimit]))
		settings = ""
	}
	kittyString += wrap(fmt.Sprintf("\033_G%sm=0;%s\033\\", settings, imgBase64[i*kittyLimit:lenImgB64]))

	wdgt.Block.ANSIString = fmt.Sprintf("\033[%d;%dH%s", imageDimensions.Min.Y, imageDimensions.Min.X, kittyString)
	return nil
}

// TODO:
// store images with ids in Kitty
