package widgets_test

import (
	"image"
	"testing"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func TestEmptyListPageChange(t *testing.T) {
	l := widgets.NewList()
	l.SetRect(0, 0, 10, 10)
	buff := termui.NewBuffer(image.Rect(0, 0, 10, 10))
	l.ScrollDown()
	l.Draw(buff)
	l.ScrollBottom()
	l.Draw(buff)
}
