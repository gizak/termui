// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
//
// Portions of this file uses [termbox-go](https://github.com/nsf/termbox-go/blob/54b74d087b7c397c402d0e3b66d2ccb6eaf5c2b4/api_common.go)
// by [authors](https://github.com/nsf/termbox-go/blob/master/AUTHORS)
// under [license](https://github.com/nsf/termbox-go/blob/master/LICENSE)

package termui

import "testing"

var ps = []string{
	"",
	"/",
	"/a",
	"/b",
	"/a/c",
	"/a/b",
	"/a/b/c",
	"/a/b/c/d",
	"/a/b/c/d/"}

func TestMatchScore(t *testing.T) {
	chk := func(a, b string, s bool) {
		if c := isPathMatch(a, b); c != s {
			t.Errorf("\na:%s\nb:%s\nshould:%t\nactual:%t", a, b, s, c)
		}
	}

	chk(ps[1], ps[1], true)
	chk(ps[1], ps[2], true)
	chk(ps[2], ps[1], false)
	chk(ps[4], ps[1], false)
	chk(ps[6], ps[2], false)
	chk(ps[4], ps[5], false)
}

func TestCrtEvt(t *testing.T) {

}
