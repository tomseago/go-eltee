package eltee

import (
	"github.com/eyethereal/go-config"
)

/*
   A Lens will be asked to observe either a raw control point or the output of
   another lens. If the lens needs the input to be of a certain type and it is
   given something that isn't of that type it should then be transparent to this
   unknown thing and pass that on up the chain.

   Usually a lens is going to save what it is asked to observe and will return
   itself as the result of the Observe call. Essentially this is creating a
   linkage of manipulations that will be applied to the source data when the
   eventual consumer asks for a particular type of viewing arrangement. That
   should be like efficient and stuff.
*/
type Lens interface {
	Observe(in interface{}) interface{}

	Kind() string
	SetFromNode(node *config.AclNode)
}

//////////////////////////////

/*
   A PositionLens is given the cartesian coordinates of a new position from which
   to observe an XYZPoint.
*/
type PositionLens struct {
	source XYZPoint

	x float64
	y float64
	z float64
}

func (l *PositionLens) Observe(in interface{}) interface{} {
	src, ok := in.(XYZPoint)
	if !ok {
		return in
	}

	l.source = src
	return l
}

func (l *PositionLens) Kind() string {
	return "position"
}

func (l *PositionLens) SetFromNode(node *config.AclNode) {
	l.x = node.ChildAsFloat("x")
	l.y = node.ChildAsFloat("y")
	l.z = node.ChildAsFloat("z")
}

func (l *PositionLens) XYZ() (float64, float64, float64) {
	if l.source == nil {
		return 0, 0, 0
	}

	sx, sy, sz := l.source.XYZ()

	return sx - l.x, sy - l.y, sz - l.z
}

//////////////////////////////

type LensStack struct {
	stack []Lens
}

func NewLenStack() *LensStack {
	return &LensStack{
		stack: make([]Lens, 0),
	}
}

func (ls *LensStack) Observe(cp ControlPoint) interface{} {
	if ls == nil {
		return cp
	}

	var view interface{} = cp
	for ix := 0; ix < len(ls.stack); ix++ {
		view = ls.stack[ix].Observe(view)
	}

	return view
}

func (ls *LensStack) AddLens(lens Lens) {
	ls.stack = append(ls.stack, lens)
}

func (ls *LensStack) Len() int {
	return len(ls.stack)
}

func (ls *LensStack) GetLens(ix int) Lens {
	return ls.stack[ix]
}
