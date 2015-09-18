// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
//
// Portions of this file uses [termbox-go](https://github.com/nsf/termbox-go/blob/54b74d087b7c397c402d0e3b66d2ccb6eaf5c2b4/api_common.go)
// by [authors](https://github.com/nsf/termbox-go/blob/master/AUTHORS)
// under [license](https://github.com/nsf/termbox-go/blob/master/LICENSE)

package termui

import (
	"strconv"
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

var sysEvtChs []chan Event

type EvtKbd struct {
	KeyStr string
}

func evtKbd(e termbox.Event) EvtKbd {
	ek := EvtKbd{}

	k := string(e.Ch)
	pre := ""
	mod := ""

	if e.Mod == termbox.ModAlt {
		mod = "M-"
	}
	if e.Ch == 0 {
		if e.Key > 0xFFFF-12 {
			k = "<f" + strconv.Itoa(0xFFFF-int(e.Key)+1) + ">"
		} else if e.Key > 0xFFFF-25 {
			ks := []string{"<insert>", "<delete>", "<home>", "<end>", "<previous>", "<next>", "<up>", "<down>", "<left>", "<right>"}
			k = ks[0xFFFF-int(e.Key)-12]
		}

		if e.Key <= 0x7F {
			pre = "C-"
			k = string('a' - 1 + int(e.Key))
			kmap := map[termbox.Key][2]string{
				termbox.KeyCtrlSpace:     {"C-", "<space>"},
				termbox.KeyBackspace:     {"", "<backspace>"},
				termbox.KeyTab:           {"", "<tab>"},
				termbox.KeyEnter:         {"", "<enter>"},
				termbox.KeyEsc:           {"", "<escape>"},
				termbox.KeyCtrlBackslash: {"C-", "\\"},
				termbox.KeyCtrlSlash:     {"C-", "/"},
				termbox.KeySpace:         {"", "<space>"},
				termbox.KeyCtrl8:         {"C-", "8"},
			}
			if sk, ok := kmap[e.Key]; ok {
				pre = sk[0]
				k = sk[1]
			}
		}
	}

	ek.KeyStr = pre + mod + k
	return ek
}

func crtTermboxEvt(e termbox.Event) Event {
	systypemap := map[termbox.EventType]string{
		termbox.EventKey:       "keyboard",
		termbox.EventResize:    "window",
		termbox.EventMouse:     "mouse",
		termbox.EventError:     "error",
		termbox.EventInterrupt: "interrupt",
	}
	ne := Event{From: "/sys", Time: time.Now().Unix()}
	typ := e.Type
	ne.Type = systypemap[typ]

	switch typ {
	case termbox.EventKey:
		kbd := evtKbd(e)
		ne.Path = "/sys/kbd/" + kbd.KeyStr
		ne.Data = kbd
	case termbox.EventResize:
		wnd := EvtWnd{}
		wnd.Width = e.Width
		wnd.Height = e.Height
		ne.Path = "/sys/wnd/resize"
		ne.Data = wnd
	case termbox.EventError:
		err := EvtErr(e.Err)
		ne.Path = "/sys/err"
		ne.Data = err
	case termbox.EventMouse:
		m := EvtMouse{}
		m.X = e.MouseX
		m.Y = e.MouseY
		ne.Path = "/sys/mouse"
		ne.Data = m
	}
	return ne
}

type EvtWnd struct {
	Width  int
	Height int
}

type EvtMouse struct {
	X     int
	Y     int
	Press string
}

type EvtErr error

func hookTermboxEvt() {
	for {
		e := termbox.PollEvent()

		for _, c := range sysEvtChs {
			go func(ch chan Event) {
				ch <- crtTermboxEvt(e)
			}(c)
		}
	}
}

func NewSysEvtCh() chan Event {
	ec := make(chan Event)
	sysEvtChs = append(sysEvtChs, ec)
	return ec
}

var DefaultEvtStream = NewEvtStream()

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
	srcMap      map[string]chan Event
	stream      chan Event
	wg          sync.WaitGroup
	sigStopLoop chan int
	Handlers    map[string]func(Event)
}

func NewEvtStream() *EvtStream {
	return &EvtStream{
		srcMap:      make(map[string]chan Event),
		stream:      make(chan Event),
		Handlers:    make(map[string]func(Event)),
		sigStopLoop: make(chan int),
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

	// divide by "" and rm heading ""
	sliced := func(s string) []string {
		ss := strings.Split(s, "/")

		i := 0
		for j := range ss {
			if ss[j] == "" {
				i++
			} else {
				break
			}
		}

		return ss[i:]
	}

	sa := sliced(a)
	sb := sliced(b)

	score := 0
	if len(sb) > len(sa) {
		return -1 // sb couldnt be more deeper than sa
	}

	for i, s := range sa {
		if i >= len(sb) {
			break // exhaust b
		}

		if s != sb[i] {
			return -1 // mismatch
		}
		score++
	}

	return score
}

func (es *EvtStream) Merge(name string, ec chan Event) {
	es.wg.Add(1)
	es.srcMap[name] = ec

	go func(a chan Event) {
		for n := range a {
			n.From = name
			es.stream <- n
		}
		es.wg.Done()
	}(ec)
}

func (es *EvtStream) Handle(path string, handler func(Event)) {
	es.Handlers[path] = handler
}

func (es *EvtStream) match(path string) string {
	n := 0
	pattern := ""
	for m := range es.Handlers {
		if MatchScore(path, m) < 0 {
			continue
		}
		if pattern == "" || len(m) > n {
			pattern = m
		}
	}
	return pattern
}

func (es *EvtStream) Loop() {
	for {
		select {
		case e := <-es.stream:
			if pattern := es.match(e.Path); pattern != "" {
				es.Handlers[pattern](e)
			}
		case <-es.sigStopLoop:
			return
		}
	}
}

func (es *EvtStream) StopLoop() {
	go func() { es.sigStopLoop <- 1 }()
}

func Merge(name string, ec chan Event) {
	DefaultEvtStream.Merge(name, ec)
}

func Handle(path string, handler func(Event)) {
	DefaultEvtStream.Handle(path, handler)
}

func Loop() {
	DefaultEvtStream.Loop()
}

func StopLoop() {
	DefaultEvtStream.StopLoop()
}

type EvtTimer struct {
	Duration time.Duration
	Count    uint64
}

func NewTimerCh(du time.Duration) chan Event {
	t := make(chan Event)

	go func(a chan Event) {
		n := uint64(0)
		for {
			n++
			time.Sleep(du)
			e := Event{}
			e.From = "timer"
			e.Type = "timer"
			e.Path = "/timer/" + du.String()
			e.Time = time.Now().Unix()
			e.Data = EvtTimer{
				Duration: du,
				Count:    n,
			}
			t <- e
		}
	}(t)
	return t
}

var DefualtHandler = func(e Event) {

}
