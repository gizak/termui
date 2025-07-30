// <Copyright> 2019 Simon Robin Lehn. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package exp

import (
	"fmt"
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"

	. "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)


func init() {
	widgets.RegisterDrawer(
		"iterm2",
		widgets.Drawer{
			Remote:         true,
			IsEscapeString: true,
			Available:      func() bool {return isIterm2 || isMacTerm},
			Draw:           drawITerm2,
		},
	)
}

func drawITerm2(wdgt *widgets.Image, buf *Buffer) (err error) {
	wdgt.Block.Draw(buf)

	img, changed, err := resizeImage(wdgt, buf)
	if !changed || err != nil {
		return err
	}

	imageDimensions := wdgt.GetVisibleArea()

	// https://www.iterm2.com/documentation-images.html
	if isIterm2 || isMacTerm {
		buf := new(bytes.Buffer)
		if err = png.Encode(buf, img); err != nil {
			goto skipIterm2
		}
		imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
		nameBase64 := base64.StdEncoding.EncodeToString([]byte(wdgt.Block.Title))
		// 0 for stretching - 1 for no stretching
		noStretch := 0
		iterm2String := wrap(fmt.Sprintf("\033]1337;File=name=%s;inline=1;height=%d;width=%d;preserveAspectRatio=%d:%s\a", nameBase64, imageDimensions.Max.Y, nameBase64, imageDimensions.Max.X, noStretch, imgBase64))
		// for width, height:   "auto"   ||   N: N character cells   ||   Npx: N pixels   ||   N%: N percent of terminal width/height
		wdgt.Block.ANSIString = fmt.Sprintf("\033[%d;%dH%s", imageDimensions.Min.Y, imageDimensions.Min.X, iterm2String)

		return nil
	}
	skipIterm2:

	return errors.New("no method applied for ANSI drawing")
}
