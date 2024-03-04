package widgets

import (
	"context"
	"fmt"

	. "github.com/proullon/termui/v3"
)

type Field struct {
	Name   string
	Text   string
	cursor int
}

type Form struct {
	Paragraph

	TextFgColor Style
	TextBgColor Style
	Fields      []*Field
	Done        bool

	ctx            context.Context
	onExitCallback func(ctx context.Context, validated bool, fields []*Field)
	fieldIdx       int
}

func NewField(name string, value string) *Field {
	f := &Field{
		Name: name,
		Text: value,
	}

	return f
}

func NewForm(ctx context.Context, name string, callback func(ctx context.Context, validated bool, fields []*Field), fields ...*Field) *Form {

	f := &Form{
		Paragraph:      *NewParagraph(),
		TextFgColor:    Theme.Paragraph.Text,
		ctx:            ctx,
		onExitCallback: callback,
	}

	for _, field := range fields {
		f.Fields = append(f.Fields, field)
	}

	f.Border = true
	f.Title = name
	f._draw()
	return f
}

// Handle input event
// <Up> or <Down> to switch field
// <Enter> will switch to next field or exit scan if last field
func (f *Form) Handle(e Event) {
	switch e.ID {
	case "<Down>":
		if f.fieldIdx < len(f.Fields)-1 {
			f.fieldIdx++
		}
	case "<Up>":
		if f.fieldIdx > 0 {
			f.fieldIdx--
		}
	case "<Enter>":
		if f.fieldIdx < len(f.Fields)-1 {
			f.fieldIdx++
		} else {
			f.Done = true
			f.onExitCallback(f.ctx, true, f.Fields)
		}
	case "<Esc>":
		f.Done = true
		f.onExitCallback(f.ctx, false, f.Fields)
	default:
		f.Fields[f.fieldIdx].Handle(e)
	}
	f._draw()
	Render(f)
}

// Handle field input
func (f *Field) Handle(e Event) {
	switch e.ID {
	case "<Backspace>":
		if f.cursor > 0 {
			f.Text = f.Text[:f.cursor-1] + f.Text[f.cursor:]
		} else if len(f.Text) > 1 {
			f.Text = f.Text[1:]
		} else {
			f.Text = ""
		}
		if f.cursor > 0 {
			f.cursor--
		}
		if f.cursor > len(f.Text) {
			f.cursor = len(f.Text)
		}
	case "<Right>":
		if f.cursor < len(f.Text) {
			f.cursor++
		}
	case "<Left>":
		if f.cursor > 0 {
			f.cursor--
		}
	case "<Space>":
		f.Text = f.Text[:f.cursor] + " " + f.Text[f.cursor:]
		f.cursor++
	default:
		f.Text = f.Text[:f.cursor] + e.ID + f.Text[f.cursor:]
		f.cursor++
	}
}

// IsDone return wether user has entered <Esc> or <Enter>
func (f *Form) IsDone() bool {
	return f.Done
}

func (f *Form) _draw() {
	cursor := "[|](bg:white)"
	f.Text = ""
	for i, field := range f.Fields {
		txt := field.Text
		if i == f.fieldIdx && field.Text != "" {
			begin := field.Text[:field.cursor]
			end := field.Text[field.cursor:]
			txt = begin + cursor + end
		} else if i == f.fieldIdx {
			txt = cursor
		}
		f.Text += fmt.Sprintf("%s: %s\n", field.Name, txt)
	}
}
