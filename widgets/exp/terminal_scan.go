// <Copyright> 2018,2019 Simon Robin Lehn. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package exp

import (
	"os"
	"strings"
)

var (
	sixelCapable, isXterm, isMlterm, isMintty, isIterm2, isUrxvt, isAlacritty, isYaft, isKitty, isMacTerm, isTmux, isScreen, isMuxed   bool
)

func scanTerminal() {
	if len(os.Getenv("XTERM_VERSION")) > 0                                          { isXterm     = true } else { isXterm     = false }
	if os.Getenv("TERM_PROGRAM") == "iTerm.app"                                     { isIterm2    = true } else { isIterm2    = false }   // https://superuser.com/a/683971
	if os.Getenv("TERM_PROGRAM") == "MacTerm"                                       { isMacTerm   = true } else { isMacTerm   = false }   // https://github.com/kmgrant/macterm/issues/3#issuecomment-458387953
	if strings.HasPrefix(os.Getenv("TERM"), "rxvt-unicode")                         { isUrxvt     = true } else { isUrxvt     = false }
	if os.Getenv("TERM") == "xterm-kitty" ||len(os.Getenv("KITTY_WINDOW_ID")) > 0   { isKitty     = true } else { isKitty     = false }
	if len(os.Getenv("MLTERM")) > 0                                                 { isMlterm    = true } else { isMlterm    = false }
	if len(os.Getenv("MINTTY_SHORTCUT")) > 0                                        { isMintty    = true } else { isMintty    = false }
	if len(os.Getenv("ALACRITTY_LOG")) > 0                                          { isAlacritty = true } else { isAlacritty = false }
	if os.Getenv("TERM") == "yaft-256color"                                         { isYaft      = true } else { isYaft      = false }   // https://github.com/uobikiemukot/yaft/blob/21b69124a2907ad6ede8f45ca96c390615e3dc0c/conf.h#L26
	if  strings.HasPrefix(os.Getenv("TERM"), "screen") && len(os.Getenv("STY")) > 0 { isScreen    = true } else { isScreen    = false }
	if (strings.HasPrefix(os.Getenv("TERM"), "screen") || strings.HasPrefix(os.Getenv("TERM"), "tmux")) &&
	   len(os.Getenv("TMUX")) > 0 || len(os.Getenv("TMUX_PANE")) > 0                { isTmux      = true } else { isTmux      = false }
	if isTmux || isScreen                                                           { isMuxed     = true } else { isMuxed     = false }

	if isYaft {
		sixelCapable = true
	} else {
		sixelCapable = false
	}
	// example query: "\033[0c"
	// possible answer from the terminal (here xterm): "\033[[?63;1;2;4;6;9;15;22c", vte(?): ...62,9;c
	// the "4" signals that the terminal is capable of sixel
	// conhost.exe knows this sequence.
	if !sixelCapable {
		termCapabilities := queryTerm(wrap("\033[0c"))
		for i, cap := range termCapabilities {
			if i == 0 || i == len(termCapabilities) - 1 {
				continue
			}
			if string(cap) == `4` {
				sixelCapable = true
			}
		}
	}
}
