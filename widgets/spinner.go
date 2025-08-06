package widgets

import (
	"fmt"
	"image"

	ui "github.com/gizak/termui/v3"
)

// Spinner it's simple animated widget
type Spinner struct {
	ui.Block
	Frames       []string
	Index        int
	Label        string
	LabelOnRight bool
	FormatString string
	TextStyle    ui.Style
}

func NewSpinner() *Spinner {
	return &Spinner{
		Block:        *ui.NewBlock(),
		Frames:       []string{"|", "/", "-", "\\", "*"},
		Index:        0,
		FormatString: "%s [%s]",
		TextStyle:    ui.NewStyle(ui.ColorWhite),
	}
}

func (s *Spinner) Advance() {
	s.Index = (s.Index + 1) % len(s.Frames)
}

func (s *Spinner) Draw(buf *ui.Buffer) {
	s.Block.Draw(buf)

	if len(s.Frames) == 0 {
		return
	}

	symbol := fmt.Sprintf(s.FormatString, s.Frames[s.Index], "")
	if len(s.Label) > 0 {
		if s.LabelOnRight {
			symbol = fmt.Sprintf(s.FormatString, s.Frames[s.Index], s.Label)
		} else {
			symbol = fmt.Sprintf(s.FormatString, s.Label, s.Frames[s.Index])
		}
	}
	x := s.Inner.Min.X
	y := s.Inner.Min.Y

	buf.SetString(symbol, s.TextStyle, image.Pt(x, y))
}
