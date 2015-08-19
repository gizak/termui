// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
//
// Portions of this file uses [termbox-go](https://github.com/nsf/termbox-go/blob/54b74d087b7c397c402d0e3b66d2ccb6eaf5c2b4/api_common.go)
// by [authors](https://github.com/nsf/termbox-go/blob/master/AUTHORS)
// under [license](https://github.com/nsf/termbox-go/blob/master/LICENSE)

package termui

import (
	"strings"

	"github.com/nsf/termbox-go"
)

//import "github.com/nsf/termbox-go"

var evtChs = make([]chan Event, 0)

// EventCh returns an output-only event channel.
// This function can be called many times (multiplexer).
func EventCh() <-chan Event {
	out := make(chan Event)
	evtChs = append(evtChs, out)
	return out
}

// turn on event listener
func evtListen() {
	go func() {
		for {
			e := termbox.PollEvent()
			// dispatch
			for _, c := range evtChs {
				go func(ch chan Event) {
					ch <- uiEvt(e)
				}(c)
			}
		}
	}()
}

type Event struct {
	Type string
	Uri  string
	From string
	To   string
	Data interface{}
	Time int
}

type evtCtl struct {
	in      chan Event
	out     chan Event
	suspend chan int
	recover chan int
	close   chan int
}

//
type EvtStream struct {
	srcMap   map[string]Event
	stream   chan Event
	cache    map[string][]func(Event)
	Handlers map[string]func(Event)
}

func newEvtCtl() evtCtl {
	ec := evtCtl{}
	ec.in = make(chan Event)
	ec.suspend = make(chan int)
	ec.recover = make(chan int)
	ec.close = make(chan int)
	ec.out = make(chan Event)
	return ec
}

func NewEvtStream() EvtStream {
	return EvtStream{
		srcMap: make(map[string]Event),
		stream: make(chan Event),
	}
}

// a: /sys/bell
// b: /sys
// score: 1
//
// a: /sys
// b: /usr
// score: -1
//
// a: /sys
// b: /
// score: 0
func MatchScore(a, b string) int {
	sa := strings.Split(a, "/")
	sb := strings.Split(b, "/")

	score := -1
	for i, s := range sa {
		if i >= len(sb) {
			break
		}

		if s != sb[i] {
			return -1
		}
		score++
	}

	return score
}

/*
func (es *EvtStream) hookup() {

}

func (es EvtStream) Subscribe(uri string) chan Event {

}
*/
