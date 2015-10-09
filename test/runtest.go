package main

import (
	"fmt"
	"os"

	"github.com/gizak/termui"
	"github.com/gizak/termui/debug"
)

func main() {
	// run as client
	if len(os.Args) > 1 {
		fmt.Print(debug.ConnectAndListen())
		return
	}

	// run as server
	go func() { panic(debug.ListenAndServe()) }()

	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	//termui.UseTheme("helloworld")
	b := termui.NewBlock()
	b.Width = 20
	b.Height = 30
	b.BorderLabel = "[HELLO](fg-red,bg-white) [WORLD](fg-blue,bg-green)"

	termui.SendBufferToRender(b)

	termui.Handle("/sys", func(e termui.Event) {
		k, ok := e.Data.(termui.EvtKbd)
		debug.Logf("->%v\n", e)
		if ok && k.KeyStr == "q" {
			termui.StopLoop()
		}
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		t := e.Data.(termui.EvtTimer)

		if t.Count%2 == 0 {
			b.BorderLabel = "[HELLO](fg-red,bg-green) [WORLD](fg-blue,bg-white)"
		} else {
			b.BorderLabel = "[HELLO](fg-blue,bg-white) [WORLD](fg-red,bg-green)"
		}

		termui.SendBufferToRender(b)

	})
	termui.Loop()
}
