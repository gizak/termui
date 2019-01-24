Feel free to search/open an issue if something is missing or confusing from the changelog, since many things have been in flux.

## TODO

- moved widgets to `github.com/gizak/termui/widgets`
- rewrote widgets (check examples and code)
- rewrote grid
  - grids are instantiated locally instead of through `termui.Body`
  - grids can be nested
  - changed grid layout mechanism
    - columns and rows can be arbitrarily nested
    - column and row size is now specified as a ratio of the available space
- `Cell`s now contain a `Style` which holds a `Fg`, `Bg`, and `Modifier`
- Change `Bufferer` interface to `Drawable`
  - Add `GetRect` and `SetRect` methods to control widget sizing
  - Change `Buffer` method to `Draw`
    - `Draw` takes a `Buffer` and draws to it instead of returning a new `Buffer`
- Refactored `Theme`
  - `Theme` is now a large struct which holds the default `Styles` of everything
- Combined `TermWidth` and `TermHeight` functions into `TerminalDimensions`
- Added `Canvas` which allows for drawing braille lines to a `Buffer`
- Refactored `Block`
- Refactored `Buffer` methods
- Set `termbox-go` backend to 256 colors by default
- Decremented color numbers by 1 to match xterm colors
- Changed text parsing
  - style items changed from `fg-color` to `fg:color`
  - added mod item like `mod:reverse`

## 18/11/29

- Move Tabpane from termui/extra to termui and rename it to TabPane
- Rename PollEvent to PollEvents

## 18/11/28

- Migrated from Dep to vgo
- Overhauled the event system
  - check the wiki/examples for details
- Renamed Par widget to Paragraph
- Renamed MBarChart widget to StackedBarChart
