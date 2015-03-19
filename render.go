package termui

import tm "github.com/nsf/termbox-go"

// all renderable components should implement this
type Bufferer interface {
	Buffer() []Point
}

func Init() error {
	// Body = container{
	// 	BgColor: theme.BodyBg,
	// 	Rows:    []Row{},
	// }
	return tm.Init()
}

func Close() {
	tm.Close()
}

// render all from left to right, right could overlap on left ones
func Render(rs ...Bufferer) {
	tm.Clear(tm.ColorDefault, toTmAttr(theme.BodyBg))
	for _, r := range rs {
		buf := r.Buffer()
		for _, v := range buf {
			tm.SetCell(v.X, v.Y, v.Ch, toTmAttr(v.Fg), toTmAttr(v.Bg))
		}
	}
	tm.Flush()
}
