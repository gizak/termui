package termui

import tm "github.com/nsf/termbox-go"

type Renderer interface {
	Buffer() []Point
}

func Init() error {
	return tm.Init()
}

func Close() {
	tm.Close()
}

func Render(rs ...Renderer) {
	for _, r := range rs {
		buf := r.Buffer()
		for _, v := range buf {
			tm.SetCell(v.X, v.Y, v.Code.Ch, v.Code.Fg, v.Code.Bg)
		}
	}
	tm.Flush()
}
