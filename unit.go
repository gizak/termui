// writen in 2017 by cnaize

package termui

import (
	"fmt"
	"sort"
)

// Unit is a layered Block with ability to store children
// with relative to the Unit coordinates
type UnitInterface interface {
	Buffers() []Bufferer
	Align()

	// only callbacks, to manage units use AddUnitTo() and RemoveUnitFrom()
	Add(parent UnitInterface)
	Remove(parent UnitInterface)
	AddChild(child UnitInterface) error
	RemoveChild(child UnitInterface) error

	Tag() string
	Root() UnitInterface
	Parent() UnitInterface
	Children() []UnitInterface

	// only in the unit (not recurcive)
	ChildByTag(tag string) (UnitInterface, error)

	LocalX() int
	LocalY() int
	WorldX() int
	WorldY() int
	WorldCenterX() int
	WorldCenterY() int

	// set coordinates respect to parent or world instead
	SetLocalX(x int)
	SetLocalY(y int)
	SetLocalCenterX(x int)
	SetLocalCenterY(y int)

	ZOrder() int
	SetZOrder(z int)
}

func AddUnitTo(unit UnitInterface, to UnitInterface) error {
	if to == nil {
		if unit != nil {
			return fmt.Errorf("adding to nil unit")
		}
		return fmt.Errorf("both units are nil")
	}

	if unit.Parent() != nil {
		unit.Remove(unit.Parent())
	}
	if err := to.AddChild(unit); err != nil {
		return err
	}
	unit.Add(to)
	return nil
}

func RemoveUnitFrom(unit UnitInterface, from UnitInterface) error {
	if from == nil {
		if unit != nil {
			return fmt.Errorf("removing from nil unit")
		}
	}

	if err := from.RemoveChild(unit); err != nil {
		return err
	}
	unit.Remove(from)
	return nil
}

// Warning!
// - don't use termui.Align for rendring
// - use SetLocalX(Y) instead of direct X(Y)
type Unit struct {
	*Block
	tag      string
	parent   UnitInterface
	children []UnitInterface
	localx   int
	localy   int
	zorder   int
	selfbuf  Bufferer
}

func NewUnit(tag string, block *Block, selfbuf Bufferer) *Unit {
	block.Float = AlignNone
	return &Unit{Block: block, tag: tag, selfbuf: selfbuf}
}

func (u *Unit) Buffers() []Bufferer {
	if !u.Display {
		return []Bufferer{}
	}

	u.Align()
	bs := []Bufferer{u.selfbuf}

	sort.Slice(u.children, func(i int, j int) bool {
		return u.children[i].ZOrder() < u.children[j].ZOrder()
	})
	for _, c := range u.Children() {
		if !u.Display {
			continue
		}
		bs = append(bs, c.Buffers()...)
	}
	return bs
}

func (u *Unit) Align() {
	u.Block.Align()
	for _, c := range u.Children() {
		c.Align()
	}
}

func (u *Unit) Add(parent UnitInterface) {
	u.parent = parent
	u.SetLocalX(0)
	u.SetLocalY(0)
}

func (u *Unit) Remove(parent UnitInterface) {
	u.parent = nil
	u.SetLocalX(u.Block.X)
	u.SetLocalY(u.Block.Y)
}

func (u *Unit) AddChild(child UnitInterface) error {
	if child == nil {
		return fmt.Errorf("adding nil child")
	}

	for _, c := range u.Children() {
		if child.Tag() == c.Tag() {
			return fmt.Errorf("already have same child")
		}
	}
	u.children = append(u.children, child)
	return nil
}

func (u *Unit) RemoveChild(child UnitInterface) error {
	if child == nil {
		return fmt.Errorf("removing nil child")
	}

	for i, c := range u.Children() {
		if c == child {
			u.children = append(u.children[:i], u.children[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("the child not exists")
}

func (u *Unit) Tag() string {
	return u.tag
}

func (u *Unit) Root() UnitInterface {
	root := u.Parent()
	if root.Parent() != nil {
		root = root.Root()
	}
	return root
}

func (u *Unit) Parent() UnitInterface {
	return u.parent
}

func (u *Unit) Children() []UnitInterface {
	return u.children
}

func (u *Unit) ChildByTag(tag string) (UnitInterface, error) {
	for _, c := range u.Children() {
		if c.Tag() == tag {
			return c, nil
		}
	}
	return nil, fmt.Errorf("child not found")
}

func (u *Unit) LocalX() int {
	return u.localx
}

func (u *Unit) LocalY() int {
	return u.localy
}

func (u *Unit) WorldCenterX() int {
	return u.WorldX() + (u.Width / 2)
}

func (u *Unit) WorldCenterY() int {
	return u.WorldY() + (u.Height / 2)
}

func (u *Unit) WorldX() int {
	return u.Block.X
}

func (u *Unit) WorldY() int {
	return u.Block.Y
}

func (u *Unit) SetLocalX(x int) {
	u.localx = x

	worldX := u.LocalX()
	parent := u.Parent()
	for parent != nil {
		worldX += parent.LocalX()
		parent = parent.Parent()
	}
	u.Block.X = worldX
	for _, c := range u.Children() {
		c.SetLocalX(c.LocalX())
	}
}

func (u *Unit) SetLocalY(y int) {
	u.localy = y

	worldY := u.LocalY()
	parent := u.Parent()
	for parent != nil {
		worldY += parent.LocalY()
		parent = parent.Parent()
	}
	u.Block.Y = worldY
	for _, c := range u.Children() {
		c.SetLocalY(c.LocalY())
	}
}

func (u *Unit) SetLocalCenterX(x int) {
	u.SetLocalX(x - u.Width/2)
}

func (u *Unit) SetLocalCenterY(y int) {
	u.SetLocalY(y - u.Height/2)
}

func (u *Unit) ZOrder() int {
	return u.zorder
}

func (u *Unit) SetZOrder(z int) {
	u.zorder = z
}
