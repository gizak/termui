package termui

import "unicode/utf8"
import "strings"

func str2runes(s string) []rune {
	n := utf8.RuneCountInString(s)
	ss := strings.Split(s, "")

	rs := make([]rune, n)
	for i := 0; i < n; i++ {
		r, _ := utf8.DecodeRuneInString(ss[i])
		rs[i] = r
	}
	return rs
}

func trimStr2Runes(s string, w int) []rune {
	rs := str2runes(s)
	if w <= 0 {
		return []rune{}
	}
	if len(rs) > w {
		rs = rs[:w]
		rs[w-1] = '…'
	}
	return rs
}
