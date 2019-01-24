# termui

[<img src="./_assets/demo.gif" alt="demo cast under osx 10.10; Terminal.app; Menlo Regular 12pt.)" width="100%">](./_examples/demo.go)

termui is a cross-platform and fully-customizable terminal dashboard and widget library built on top of [termbox-go](https://github.com/nsf/termbox-go). It is inspired by [blessed-contrib](https://github.com/yaronn/blessed-contrib) and written purely in Go.

## Installation

Installing from the master branch is recommended:

```bash
go get -u github.com/gizak/termui@master
```

**Note**: termui is currently undergoing API changes so make sure to check the changelog when upgrading.
If you upgrade and notice something is missing or don't like a change, revert the upgrade and open an issue.

## Widgets

- [BarChart](./_examples/barchart.go)
- [Canvas](./_examples/canvas.go)
- [Gauge](./_examples/gauge.go)
- [LineChart](./_examples/linechart.go)
- [List](./_examples/list.go)
- [Paragraph](./_examples/paragraph.go)
- [PieChart](./_examples/piechart.go)
- [Sparkline](./_examples/sparkline.go)
- [StackedBarChart](./_examples/stacked_barchart.go)
- [Table](./_examples/table.go)
- [Tabs](./_examples/tabs.go)

Run an example with `go run _examples/{example}.go` or run all of them consecutively with `make run-examples`.

## Documentation

- [wiki](https://github.com/gizak/termui/wiki)

## Uses

- [cjbassi/gotop](https://github.com/cjbassi/gotop)
- [go-ethereum/monitorcmd](https://github.com/ethereum/go-ethereum/blob/96116758d22ddbff4dbef2050d6b63a7b74502d8/cmd/geth/monitorcmd.go)

## Related Works

- [blessed-contrib](https://github.com/yaronn/blessed-contrib)
- [tui-rs](https://github.com/fdehau/tui-rs)
- [gocui](https://github.com/jroimartin/gocui)

## License

[MIT](http://opensource.org/licenses/MIT)
