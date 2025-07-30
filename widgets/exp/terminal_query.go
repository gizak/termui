// <Copyright> 2018,2019 Simon Robin Lehn. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package exp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mattn/go-tty"
)

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
