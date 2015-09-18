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

	termui.Handle("/sys", func(e termui.Event) {
		k, ok := e.Data.(termui.EvtKbd)
		debug.Logf("-->%v\n", e)
		if ok && k.KeyStr == "q" {
			termui.StopLoop()
		}
	})

	termui.Handle("/timer", func(e termui.Event) {
		//debug.Logf("-->%v\n", e)
	})
	termui.Loop()
}
