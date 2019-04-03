package main

import (
	"fmt"
	"strings"
	"os"
)

var isTmux bool

func main() {
	if len(os.Getenv("TMUX")) > 0 || len(os.Getenv("TMUX_PANE")) > 0 {
		isTmux = true
	}
	fmt.Printf("\033[%d;%dH\033[?8452h%s", 20, 20, wrap("\033Pq#0;2;0;0;0#1;2;100;100;0#2;2;0;100;0#1~~@@vv@@~~@@~~$#2??}}GG}}??}}??-#1!14@\033\\"))
}

func wrap(s string) string {
	if isTmux {
		return tmuxWrap(s)
	}
	return s
}

func tmuxWrap(s string) string {
	return "\x1bPtmux;" + strings.Replace(s, "\x1b", "\x1b\x1b", -1) + "\x1b\\"
}
