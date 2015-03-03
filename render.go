package termui

import tm "github.com/nsf/termbox-go"

type Bufferer interface {
	Buffer() []Point
}

func Init() error {
	return tm.Init()
}

func Close() {
	tm.Close()
}

func Render(rs ...Bufferer) {
	for _, r := range rs {
		buf := r.Buffer()
		for _, v := range buf {
			tm.SetCell(v.X, v.Y, v.Ch, toTmAttr(v.Fg), toTmAttr(v.Bg))
		}
	}
	tm.Flush()
}
