// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"strconv"
	"sync"

	tb "github.com/nsf/termbox-go"
)

/*
Here's the list of events which can be assigned handlers using Handle():
	mouse events:
		<MouseLeft> <MouseRight> <MouseMiddle>
		<MouseWheelUp> <MouseWheelDown>
	keyboard events:
		any uppercase or lowercase letter or a set of two letters like j or jj or J or JJ
		<C-d> etc
		<M-d> etc
		<Up> <Down> <Left> <Right>
		<Insert> <Delete> <Home> <End> <Previous> <Next>
		<Backspace> <Tab> <Enter> <Escape> <Space>
		<C-<Space>> etc
	terminal events:
		<Resize>
*/

type EventType int

const (
	KeyboardEvent EventType = iota
	MouseEvent
	ResizeEvent
)

type eventStream struct {
	sync.RWMutex
	handlers   map[string]func(Event)
	stopLoop   chan bool
	eventQueue chan tb.Event // list of events from termbox
	hook       func(Event)
}

var defaultES = eventStream{
	handlers:   make(map[string]func(Event)),
	stopLoop:   make(chan bool, 1),
	eventQueue: make(chan tb.Event),
	hook:       DefaultHandler,
}

// Event contains an ID used for Handle() and an optional payload.
type Event struct {
	Type    EventType
	ID      string
	Payload interface{}
}

// Mouse payload.
type Mouse struct {
	Drag bool
	X    int
	Y    int
}

// Resize payload.
type Resize struct {
	Width  int
	Height int
}

// handleEvent calls the approriate callback function if there is one.
func handleEvent(e Event) {
	if val, ok := defaultES.handlers[e.ID]; ok {
		val(e)
	}
}

// Loop gets events from termbox and passes them off to handleEvent.
// Stops when StopLoop is called.
func Loop() {
	go func() {
		for {
			defaultES.eventQueue <- tb.PollEvent()
		}
	}()

	for {
		select {
		case <-defaultES.stopLoop:
			return
		case e := <-defaultES.eventQueue:
			ne := convertTermboxEvent(e)
			defaultES.RLock()
			handleEvent(ne)
			defaultES.hook(ne)
			defaultES.RUnlock()
		}
	}
}

// StopLoop stops the event loop.
func StopLoop() {
	defaultES.stopLoop <- true
}

// Handle assigns event names to their handlers. Takes a string, strings, or a slice of strings, and a function.
func Handle(things ...interface{}) {
	function := things[len(things)-1].(func(Event))
	for _, thing := range things {
		if value, ok := thing.(string); ok {
			defaultES.Lock()
			defaultES.handlers[value] = function
			defaultES.Unlock()
		}
		if value, ok := thing.([]string); ok {
			defaultES.Lock()
			for _, name := range value {
				defaultES.handlers[name] = function
			}
			defaultES.Unlock()
		}
	}
}

func EventHook(f func(Event)) {
	defaultES.Lock()
	defaultES.hook = f
	defaultES.Unlock()
}

// convertTermboxKeyboardEvent converts a termbox keyboard event to a more friendly string format.
// Combines modifiers into the string instead of having them as additional fields in an event.
func convertTermboxKeyboardEvent(e tb.Event) Event {
	k := string(e.Ch)
	pre := ""
	mod := ""

	if e.Mod == tb.ModAlt {
		mod = "<M-"
	}
	if e.Ch == 0 {
		if e.Key > 0xFFFF-12 {
			k = "<f" + strconv.Itoa(0xFFFF-int(e.Key)+1) + ">"
		} else if e.Key > 0xFFFF-25 {
			ks := []string{"<Insert>", "<Delete>", "<Home>", "<End>", "<Previous>", "<Next>", "<Up>", "<Down>", "<Left>", "<Right>"}
			k = ks[0xFFFF-int(e.Key)-12]
		}

		if e.Key <= 0x7F {
			pre = "<C-"
			k = string('a' - 1 + int(e.Key))
			kmap := map[tb.Key][2]string{
				tb.KeyCtrlSpace:     {"C-", "<Space>"},
				tb.KeyBackspace:     {"", "<Backspace>"},
				tb.KeyTab:           {"", "<Tab>"},
				tb.KeyEnter:         {"", "<Enter>"},
				tb.KeyEsc:           {"", "<Escape>"},
				tb.KeyCtrlBackslash: {"C-", "\\"},
				tb.KeyCtrlSlash:     {"C-", "/"},
				tb.KeySpace:         {"", "<Space>"},
				tb.KeyCtrl8:         {"C-", "8"},
			}
			if sk, ok := kmap[e.Key]; ok {
				pre = sk[0]
				k = sk[1]
			}
		}
	}

	if pre != "" {
		k += ">"
	}

	id := pre + mod + k

	return Event{
		Type: KeyboardEvent,
		ID:   id,
	}
}

func convertTermboxMouseEvent(e tb.Event) Event {
	mouseButtonMap := map[tb.Key]string{
		tb.MouseLeft:      "<MouseLeft>",
		tb.MouseMiddle:    "<MouseMiddle>",
		tb.MouseRight:     "<MouseRight>",
		tb.MouseRelease:   "<MouseRelease>",
		tb.MouseWheelUp:   "<MouseWheelUp>",
		tb.MouseWheelDown: "<MouseWheelDown>",
	}

	converted, ok := mouseButtonMap[e.Key]
	if !ok {
		converted = "Unknown_Mouse_Button"
	}

	Drag := false
	if e.Mod == tb.ModMotion {
		Drag = true
	}

	return Event{
		Type: MouseEvent,
		ID:   converted,
		Payload: Mouse{
			X:    e.MouseX,
			Y:    e.MouseY,
			Drag: Drag,
		},
	}
}

// convertTermboxEvent turns a termbox event into a termui event.
func convertTermboxEvent(e tb.Event) Event {
	if e.Type == tb.EventError {
		panic(e.Err)
	}

	switch e.Type {
	case tb.EventKey:
		return convertTermboxKeyboardEvent(e)
	case tb.EventMouse:
		return convertTermboxMouseEvent(e)
	case tb.EventResize:
		return Event{
			Type: ResizeEvent,
			ID:   "<Resize>",
			Payload: Resize{
				Width:  e.Width,
				Height: e.Height,
			},
		}
	}

	return Event{}
}

var DefaultHandler = func(e Event) {
}

func ResetHandlers() {
	defaultES.Lock()
	defaultES.handlers = make(map[string]func(Event))
	defaultES.Unlock()
}

func ResetHandler(handle string) {
	defaultES.Lock()
	delete(defaultES.handlers, handle)
	defaultES.Unlock()
}
