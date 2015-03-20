// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import "testing"

func TestStr2Rune(t *testing.T) {
	s := "你好,世界."
	rs := str2runes(s)
	if len(rs) != 6 {
		t.Error()
	}
}
