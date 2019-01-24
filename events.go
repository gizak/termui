// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"fmt"

	tb "github.com/nsf/termbox-go"
)

/*
List of events:
	mouse events:
		<MouseLeft> <MouseRight> <MouseMiddle>
		<MouseWheelUp> <MouseWheelDown>
	keyboard events:
		any uppercase or lowercase letter like j or J
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

// PollEvents gets events from termbox, converts them, then sends them to each of its channels.
func PollEvents() <-chan Event {
	ch := make(chan Event)
	go func() {
		for {
			ch <- convertTermboxEvent(tb.PollEvent())
		}
	}()
	return ch
}

// convertTermboxKeyboardEvent converts a termbox keyboard event to a more friendly string format.
// Combines modifiers into the string instead of having them as additional fields in an event.
func convertTermboxKeyboardEvent(e tb.Event) Event {
	ID := "%s"
	if e.Mod == tb.ModAlt {
		ID = "<M-%s>"
	}

	if e.Ch != 0 {
		ID = fmt.Sprintf(ID, string(e.Ch))
	} else {
		switchExpression := func() string {
			switch e.Key {
			case tb.KeyF1:
				return "<F1>"
			case tb.KeyF2:
				return "<F2>"
			case tb.KeyF3:
				return "<F3>"
			case tb.KeyF4:
				return "<F4>"
			case tb.KeyF5:
				return "<F5>"
			case tb.KeyF6:
				return "<F6>"
			case tb.KeyF7:
				return "<F7>"
			case tb.KeyF8:
				return "<F8>"
			case tb.KeyF9:
				return "<F9>"
			case tb.KeyF10:
				return "<F10>"
			case tb.KeyF11:
				return "<F11>"
			case tb.KeyF12:
				return "<F12>"
			case tb.KeyInsert:
				return "<Insert>"
			case tb.KeyDelete:
				return "<Delete>"
			case tb.KeyHome:
				return "<Home>"
			case tb.KeyEnd:
				return "<End>"
			case tb.KeyPgup:
				return "<PageUp>"
			case tb.KeyPgdn:
				return "<PageDown>"
			case tb.KeyArrowUp:
				return "<Up>"
			case tb.KeyArrowDown:
				return "<Down>"
			case tb.KeyArrowLeft:
				return "<Left>"
			case tb.KeyArrowRight:
				return "<Right>"

			case tb.KeyCtrlTilde: // tb.KeyCtrl2 tb.KeyCtrlSpace
				// <C-~> doesn't work
				// <C-2> doesn't work
				return "<C-<Space>>"
			case tb.KeyCtrlA:
				return "<C-a>"
			case tb.KeyCtrlB:
				return "<C-b>"
			case tb.KeyCtrlC:
				return "<C-c>"
			case tb.KeyCtrlD:
				return "<C-d>"
			case tb.KeyCtrlE:
				return "<C-e>"
			case tb.KeyCtrlF:
				return "<C-f>"
			case tb.KeyCtrlG:
				return "<C-g>"
			case tb.KeyBackspace: // tb.KeyCtrlH
				// <C-h> doesn't work
				return "<C-<Backspace>>"
			case tb.KeyTab: // tb.KeyCtrlI
				// <C-i> doesn't work
				return "<Tab>"
			case tb.KeyCtrlJ:
				return "<C-j>"
			case tb.KeyCtrlK:
				return "<C-k>"
			case tb.KeyCtrlL:
				return "<C-l>"
			case tb.KeyEnter: // tb.KeyCtrlM
				// <C-m> doesn't work
				return "<Enter>"
			case tb.KeyCtrlN:
				return "<C-n>"
			case tb.KeyCtrlO:
				return "<C-o>"
			case tb.KeyCtrlP:
				return "<C-p>"
			case tb.KeyCtrlQ:
				return "<C-q>"
			case tb.KeyCtrlR:
				return "<C-r>"
			case tb.KeyCtrlS:
				return "<C-s>"
			case tb.KeyCtrlT:
				return "<C-t>"
			case tb.KeyCtrlU:
				return "<C-u>"
			case tb.KeyCtrlV:
				return "<C-v>"
			case tb.KeyCtrlW:
				return "<C-w>"
			case tb.KeyCtrlX:
				return "<C-x>"
			case tb.KeyCtrlY:
				return "<C-y>"
			case tb.KeyCtrlZ:
				return "<C-z>"
			case tb.KeyEsc: // tb.KeyCtrlLsqBracket tb.KeyCtrl3
				// <C-[> doesn't work
				// <C-3> doesn't work
				return "<Escape>"
			case tb.KeyCtrl4: // tb.KeyCtrlBackslash
				// <C-\\> doesn't work
				return "<C-4>"
			case tb.KeyCtrl5: // tb.KeyCtrlRsqBracket
				// <C-]> doesn't work
				return "<C-5>"
			case tb.KeyCtrl6:
				return "<C-6>"
			case tb.KeyCtrl7: // tb.KeyCtrlSlash tb.KeyCtrlUnderscore
				// <C-/> doesn't work
				// <C-_> doesn't work
				return "<C-7>"
			case tb.KeySpace:
				return "<Space>"
			case tb.KeyBackspace2: // tb.KeyCtrl8:
				// <C-8> doesn't work
				return "<Backspace>"
			}
			// <C--> doesn't work
			return ""
		}
		ID = fmt.Sprintf(ID, switchExpression())
	}

	return Event{
		Type: KeyboardEvent,
		ID:   ID,
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

	var event Event

	switch e.Type {
	case tb.EventKey:
		event = convertTermboxKeyboardEvent(e)
	case tb.EventMouse:
		event = convertTermboxMouseEvent(e)
	case tb.EventResize:
		event = Event{
			Type: ResizeEvent,
			ID:   "<Resize>",
			Payload: Resize{
				Width:  e.Width,
				Height: e.Height,
			},
		}
	}

	return event
}
