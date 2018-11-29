# termui

[![Build Status](https://travis-ci.org/gizak/termui.svg?branch=master)](https://travis-ci.org/gizak/termui) [![Doc Status](https://godoc.org/github.com/gizak/termui?status.png)](https://godoc.org/github.com/gizak/termui)

<img src="./_examples/dashboard.gif" alt="demo cast under osx 10.10; Terminal.app; Menlo Regular 12pt.)" width="100%">

`termui` is a cross-platform, easy-to-compile, and fully-customizable terminal dashboard built on top of [termbox-go](https://github.com/nsf/termbox-go). It is inspired by [blessed-contrib](https://github.com/yaronn/blessed-contrib) and written purely in Go.

**termui is currently undergoing some API changes so make sure to check the changelog when upgrading**

## Installation

Installing from the master branch is recommended:

```bash
go get -u github.com/gizak/termui@master
```

## Usage

### Hello World

```go
package main

import ui "github.com/gizak/termui"

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewParagraph("Hello World!")
	p.Width = 25
	p.Height = 5
	ui.Render(p)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}
```

### Widgets

Click image to see the corresponding demo codes.

[<img src="./_examples/barchart.png" alt="barchart" type="image/png" width="45%">](https://github.com/gizak/termui/blob/master/_examples/barchart.go)
[<img src="./_examples/gauge.png" alt="gauge" type="image/png" width="45%">](https://github.com/gizak/termui/blob/master/_examples/gauge.go)
[<img src="./_examples/linechart.png" alt="linechart" type="image/png" width="45%">](https://github.com/gizak/termui/blob/master/_examples/linechart.go)
[<img src="./_examples/list.png" alt="list" type="image/png" width="45%">](https://github.com/gizak/termui/blob/master/_examples/list.go)
[<img src="./_examples/paragraph.png" alt="paragraph" type="image/png" width="45%">](https://github.com/gizak/termui/blob/master/_examples/par.go)
[<img src="./_examples/sparklines.png" alt="sparklines" type="image/png" width="45%">](https://github.com/gizak/termui/blob/master/_examples/sparklines.go)
[<img src="./_examples/stackedbarchart.png" alt="stackedbarchart" type="image/png" width="45%">](https://github.com/gizak/termui/blob/master/_examples/mbarchart.go)
[<img src="./_examples/table.png" alt="table" type="image/png" width="45%">](https://github.com/gizak/termui/blob/master/_examples/table.go)

### Examples

Examples can be found in [\_examples](./_examples). Run with `go run _examples/...` or run all of them consecutively with `./scripts/run_examples.py`.

## Documentation

- [godoc](https://godoc.org/github.com/gizak/termui) for code documentation
- [wiki](https://github.com/gizak/termui/wiki) for general information

## Uses

- [go-ethereum/monitorcmd](https://github.com/ethereum/go-ethereum/blob/96116758d22ddbff4dbef2050d6b63a7b74502d8/cmd/geth/monitorcmd.go)

## Related Works

- [blessed-contrib](https://github.com/yaronn/blessed-contrib)
- [tui-rs](https://github.com/fdehau/tui-rs)
- [gocui](https://github.com/jroimartin/gocui)

## License

This library is under the [MIT License](http://opensource.org/licenses/MIT)
