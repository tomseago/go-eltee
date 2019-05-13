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
	Observe(fc *FixtureControl, in interface{}) interface{}
	// Kind() string
	// SetFromNode(node *config.AclNode)
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

func (l *PositionLens) Observe(fc *FixtureControl, in interface{}) interface{} {
	src, ok := in.(XYZPoint)
	if !ok {
		return in
	}

	l.source = src
	f := fc.Fixture
	l.x = f.GetF64("pos_x")
	l.y = f.GetF64("pos_y")
	l.z = f.GetF64("pos_z")
	return l
}

func (l *PositionLens) Kind() string {
	return "position"
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

func NewLensStack() *LensStack {
	return &LensStack{
		stack: make([]Lens, 0),
	}
}

func NewLensStackFromNode(base *config.AclNode) *LensStack {
	log.Infof("NewLensStackFrom %v", base)

	if base == nil {
		return nil
	}

	ls := &LensStack{
		stack: make([]Lens, 0),
	}

	for i := 0; i < base.Len(); i++ {
		lNode, ok := (base.Values[i]).(*config.AclNode)
		if !ok {
			continue
		}
		//    }

		// base.ForEachOrderedChild(func(nn string, lNode *config.AclNode) {
		lens := LensFromNode(lNode)
		if lens != nil {
			ls.stack = append(ls.stack, lens)
		}
	} //)

	log.Infof("Lens stack %v", ls)

	return ls
}

func (ls *LensStack) Observe(fc *FixtureControl, cp ControlPoint) interface{} {
	if ls == nil {
		return cp
	}

	var view interface{} = cp
	for ix := 0; ix < len(ls.stack); ix++ {
		view = ls.stack[ix].Observe(fc, view)
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

//////////////////////////////

func LensFromNode(node *config.AclNode) Lens {
	log.Debugf("LensFromNode %v", node)
	if node == nil {
		return nil
	}

	kind := node.ChildAsString("kind")
	if len(kind) == 0 {
		return nil
	}

	if kind == "position" {
		log.Infof("Returning position lens")
		return &PositionLens{}
	}

	log.Errorf("Unknown lens type %v", kind)
	return nil
}
