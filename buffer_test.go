// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"
	"testing"
)

func TestBufferUnion(t *testing.T) {
	b0 := NewBuffer()
	b1 := NewBuffer()

	b1.Area.Max.X = 100
	b1.Area.Max.Y = 100
	b0.Area.Max.X = 50
	b0.Merge(b1)
	if b0.Area.Max.X != 100 {
		t.Errorf("Buffer.Merge unions Area failed: should:%v, actual %v,%v", image.Rect(0, 0, 50, 0).Union(image.Rect(0, 0, 100, 100)), b1.Area, b0.Area)
	}
}
