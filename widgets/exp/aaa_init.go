// <Copyright> 2018,2019 Simon Robin Lehn. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package exp

import (
	. "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// the file name should appear at the top when alphabetically sorted (start with "aaa")
// because the init() functions are executed in alphabetic file order
func init() {
	scanTerminal()
	var drawFallback func(*widgets.Image, *Buffer) (error)
	if drbl, ok := widgets.GetDrawers()["block"]; ok {
		drawFallback = drbl.Draw
	}
	widgets.RegisterDrawer(
		"block",
		widgets.Drawer{
			Remote:         true,
			IsEscapeString: false,
			Available:      func() bool {return true},
			Draw:           func(img *widgets.Image, buf *Buffer) (err error) {
				// possible reattachments of the terminal multiplexer?
				if isMuxed {
					scanTerminal()
				}

				return drawFallback(img, buf)
			},
		},
	)
}
