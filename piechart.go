package termui

import (
	"container/list"
	"image"
	"math"
)

const (
	piechartOffsetUp = -.5 * math.Pi // the northward angle
	sectorFactor     = 7.0           // circle resolution: precision vs. performance
	fullCircle       = 2.0 * math.Pi // the full circle angle
	xStretch         = 2.0           // horizontal adjustment
	solidBlock       = '░'
)

var (
	defaultColors = []Attribute{
		ColorRed,
		ColorGreen,
		ColorYellow,
		ColorBlue,
		ColorMagenta,
		ColorCyan,
		ColorWhite,
	}
)

// PieChartLabel callback, i is the current data index, v the current value
type PieChartLabel func(i int, v float64) string

type PieChart struct {
	Block
	Data        []float64     // list of data items
	Colors      []Attribute   // colors to by cycled through (see defaultColors)
	BorderColor Attribute     // color of the pie-border
	Label       PieChartLabel // callback function for labels
	Offset      float64       // which angle to start drawing at? (see piechartOffsetUp)
}

// NewPieChart Creates a new pie chart with reasonable defaults and no labels
func NewPieChart() *PieChart {
	return &PieChart{
		Block:       *NewBlock(),
		Colors:      defaultColors,
		Offset:      piechartOffsetUp,
		BorderColor: ColorDefault,
	}
}

// computes the color for a given data index
func (pc *PieChart) colorFor(i int) Attribute {
	return pc.Colors[i%len(pc.Colors)]
}

// Buffer creates the buffer for the pie chart
func (pc *PieChart) Buffer() Buffer {
	buf := pc.Block.Buffer()
	w, h := pc.innerArea.Dx(), pc.innerArea.Dy()
	center := image.Point{X: w / 2, Y: h / 2}

	// radius for the border
	r := float64(w/2/xStretch) - 1.0
	if h < w/xStretch {
		r = float64(h/2) - 1.0
	}

	// make border
	borderCircle := &circle{Point: center, radius: r}
	drawBorder := func() {
		borderCircle.draw(&buf, Cell{Ch: solidBlock, Fg: pc.BorderColor, Bg: ColorDefault})
	}
	drawBorder()

	if len(pc.Data) == 0 { // nothing to draw?
		return buf
	}

	// compute slice sizes
	sum := sum(pc.Data)
	sliceSizes := make([]float64, len(pc.Data))
	for i, v := range pc.Data {
		sliceSizes[i] = v / sum * fullCircle
	}

	// draw slice borders
	phi := pc.Offset
	for i, v := range sliceSizes {
		p := borderCircle.at(phi)
		l := line{P1: center, P2: p}
		l.draw(&buf, &Cell{Ch: solidBlock, Fg: pc.colorFor(i), Bg: ColorDefault})
		phi += v
	}

	// fill slices
	middleCircle := circle{Point: center, radius: r / 2.0}
	_, sectorSize := borderCircle.sectors()
	phi = pc.Offset
	for i, v := range sliceSizes {
		if v > sectorSize { // do not render if slice is too small
			cell := Cell{Ch: solidBlock, Fg: pc.colorFor(i), Bg: ColorDefault}
			halfSlice := phi + v/2.0
			fill(borderCircle.inner(halfSlice), &cell, &buf)
			fill(middleCircle.at(halfSlice), &cell, &buf)
			for f := phi; f < phi+v; f += sectorSize {
				line{P1: center, P2: borderCircle.inner(f)}.draw(&buf, &cell)
			}
		}
		phi += v
	}

	// labels
	if pc.Label != nil {
		drawLabel := func(p image.Point, label string, fg, bg Attribute) {
			offset := p.Add(image.Point{X: -int(round(float64(len(label)) / 2.0)), Y: 0})
			for i, v := range []rune(label) {
				buf.Set(offset.X+i, offset.Y, Cell{Ch: v, Fg: fg, Bg: bg})
			}
		}
		phi = pc.Offset
		for i, v := range sliceSizes {
			labelAt := middleCircle.at(phi + v/2.0)
			if len(pc.Data) == 1 {
				labelAt = center
			}
			drawLabel(labelAt, pc.Label(i, pc.Data[i]), pc.colorFor(i), ColorDefault)
			phi += v
		}
	}
	drawBorder()
	return buf
}

// fills empty cells from position
func fill(p image.Point, c *Cell, buf *Buffer) {
	empty := func(x, y int) bool {
		return buf.At(x, y).Ch == ' '
	}
	if !empty(p.X, p.Y) {
		return
	}
	q := list.New()
	q.PushBack(p)
	buf.Set(p.X, p.Y, *c)
	for q.Front() != nil {
		p := q.Remove(q.Front()).(image.Point)
		w, e, row := p.X, p.X, p.Y

		for empty(w-1, row) {
			w--
		}
		for empty(e+1, row) {
			e++
		}
		for x := w; x <= e; x++ {
			buf.Set(x, row, *c)
			if empty(x, row-1) {
				q.PushBack(image.Point{X: x, Y: row - 1})
			}
			if empty(x, row+1) {
				q.PushBack(image.Point{X: x, Y: row + 1})
			}
		}
	}
}

type circle struct {
	image.Point
	radius float64
}

// computes the point at a given angle phi
func (c circle) at(phi float64) image.Point {
	x := c.X + int(round(xStretch*c.radius*math.Cos(phi)))
	y := c.Y + int(round(c.radius*math.Sin(phi)))
	return image.Point{X: x, Y: y}
}

// computes the "inner" point at a given angle phi
func (c circle) inner(phi float64) image.Point {
	p := image.Point{X: 0, Y: 0}
	outer := c.at(phi)
	if c.X < outer.X {
		p.X = -1
	} else if c.X > outer.X {
		p.X = 1
	}
	if c.Y < outer.Y {
		p.Y = -1
	} else if c.Y > outer.Y {
		p.Y = 1
	}
	return outer.Add(p)
}

// computes the perimeter of a circle
func (c circle) perimeter() float64 {
	return 2.0 * math.Pi * c.radius
}

// computes the number of sectors and the size of each sector
func (c circle) sectors() (sectors float64, sectorSize float64) {
	sectors = c.perimeter() * sectorFactor
	sectorSize = fullCircle / sectors
	return
}

// draws the circle
func (c circle) draw(buf *Buffer, cell Cell) {
	sectors, sectorSize := c.sectors()
	for i := 0; i < int(round(sectors)); i++ {
		phi := float64(i) * sectorSize
		point := c.at(float64(phi))
		buf.Set(point.X, point.Y, cell)
	}
}

// a line between two points
type line struct {
	P1, P2 image.Point
}

// draws the line
func (l line) draw(buf *Buffer, cell *Cell) {
	isLeftOf := func(p1, p2 image.Point) bool {
		return p1.X <= p2.X
	}
	isTopOf := func(p1, p2 image.Point) bool {
		return p1.Y <= p2.Y
	}
	p1, p2 := l.P1, l.P2
	buf.Set(l.P2.X, l.P2.Y, Cell{Ch: '*', Fg: cell.Fg, Bg: cell.Bg})
	width, height := l.size()
	if width > height { // paint left to right
		if !isLeftOf(p1, p2) {
			p1, p2 = p2, p1
		}
		flip := 1.0
		if !isTopOf(p1, p2) {
			flip = -1.0
		}
		for x := p1.X; x <= p2.X; x++ {
			ratio := float64(height) / float64(width)
			factor := float64(x - p1.X)
			y := ratio * factor * flip
			buf.Set(x, int(round(y))+p1.Y, *cell)
		}
	} else { // paint top to bottom
		if !isTopOf(p1, p2) {
			p1, p2 = p2, p1
		}
		flip := 1.0
		if !isLeftOf(p1, p2) {
			flip = -1.0
		}
		for y := p1.Y; y <= p2.Y; y++ {
			ratio := float64(width) / float64(height)
			factor := float64(y - p1.Y)
			x := ratio * factor * flip
			buf.Set(int(round(x))+p1.X, y, *cell)
		}
	}
}

// width and height of a line
func (l line) size() (w, h int) {
	return abs(l.P2.X - l.P1.X), abs(l.P2.Y - l.P1.Y)
}

// rounds a value
func round(x float64) float64 {
	return math.Floor(x + 0.5)
}

// fold a sum
func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

// math.Abs for ints
func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}
