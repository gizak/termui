package termui

import "testing"

func TestStr2Rune(t *testing.T) {
	s := "你好,世界."
	rs := str2runes(s)
	if len(rs) != 6 {
		t.Error()
	}
}
