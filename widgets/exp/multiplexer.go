// <Copyright> 2019 Simon Robin Lehn. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package exp

import (
	"strings"
)

func wrap(s string) string {
	if !isMuxed {
		return s
	}
	if isTmux {
		return tmuxWrap(s)
	}
	return s
}

func tmuxWrap(s string) string {
	return "\033Ptmux;" + strings.Replace(s, "\033", "\033\033", -1) + "\033\\"
}

/*
// https://savannah.gnu.org/bugs/index.php?56063
func screenWrap(s string) string {}
*/
