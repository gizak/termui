# termui
Go terminal dashboard. Inspired by [blessed-contrib](https://github.com/yaronn/blessed-contrib), but purely in Go.

Cross-platform, easy to compile, and fully-customizable.

__Demo:__

<img src="./example/screencast.gif" alt="demo" width="600">

## Installation

	go get github.com/gizak/termui

## Usage

Each component's layout is a bit like HTML block, which has border and padding. 

The `Border` property can be chosen to hide or display (with its border label), when it comes to display, in this case the space it takes is counted as padding space (i.e. `PaddingTop=PaddingBottom=PaddingLeft=PaddingRight=1`).

`````go
	import ui "github.com/gizak/termui" // <- ui shortcut, optional

	func main() {
		err := ui.Init()
		if err != nil {
			panic(err)
		}
		defer ui.Close()

		p := ui.NewP(":PRESS q TO QUIT DEMO")
		p.Height = 3
		p.Width = 50
		p.TextFgColor = ui.ColorWhite
		p.Border.Label = "Text Box"
		p.Border.FgColor = ui.ColorCyan

		g := ui.NewGauge()
		g.Percent = 50
		g.Width = 50
		g.Height = 3
		g.Y = 11
		g.Border.Label = "Gauge"
		g.BarColor = ui.ColorRed
		g.Border.FgColor = ui.ColorWhite
		g.Border.LabelFgColor = ui.ColorCyan

		ui.Render(p, g)

		// event handler...
	}
`````

Note that components can be overlapped (I'd rather call this as a feature...), `Render(rs ...Renderer)` renders its args from left to right (i.e. each component's weight is arising from left to right).

## Widgets

_APIs are subject to change, docs will be added after 2 or 3 commits_

## GoDoc

[godoc](https://godoc.org/github.com/gizak/termui).

## License
This library is under the [MIT License](http://opensource.org/licenses/MIT)
