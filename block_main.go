// +build ignore

package main

import "github.com/gizak/termui"

func main() {
	termui.Init()

	termui.UseTheme("helloworld")
	b := termui.NewBlock()
	b.Width = 20
	b.Height = 30
	b.BorderLabel = "[HELLO](fg-red,bg-white) [WORLD](fg-blue,bg-green)"

	termui.Render(b)
	<-termui.EventCh()
	termui.Close()
}
