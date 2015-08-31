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
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

type Event struct {
	Type string
	Path string
	From string
	To   string
	Data interface{}
	Time int64
}

var sysevt struct {
	chs []chan Event
}

func newSysEvtFromTb(e termbox.Event) Event {
	ne := Event{From: "/sys", Time: time.Now().Unix()}
	return ne
}

func hookSysEvt() {
	sysevt.chs = make([]chan Event, 0)
	for {
		e := termbox.PollEvent()
		for _, c := range sysevt.chs {
			// shorten?
			go func(ch chan Event, ev Event) { ch <- ev }(c, newSysEvtFromTb(e))
		}
	}
}

func NewSysEvtCh() chan Event {
	ec := make(chan Event)
	sysevt.chs = append(sysevt.chs, ec)
	return ec
}

/*
type evtCtl struct {
	in      chan Event
	out     chan Event
	suspend chan int
	recover chan int
	close   chan int
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

*/
//
type EvtStream struct {
	srcMap   map[string]chan Event
	stream   chan Event
	cache    map[string][]func(Event)
	wg       sync.WaitGroup
	Handlers map[string]func(Event)
}

func NewEvtStream() *EvtStream {
	return &EvtStream{
		srcMap: make(map[string]chan Event),
		stream: make(chan Event),
	}
}

func (es *EvtStream) Init() {

	go func() {
		es.wg.Wait()
		close(es.stream)
	}()
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

func (es *EvtStream) Merge(ec chan Event) {
	es.wg.Add(1)

	go func(a chan Event) {
		for n := range ec {
			es.stream <- n
		}
		wg.Done()
	}(ec)
}

/*
func (es *EvtStream) hookup() {

}

func (es EvtStream) Subscribe(uri string) chan Event {

}
*/
