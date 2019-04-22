package eltee

import (
	"github.com/eyethereal/go-config"
	"math"
	"strings"
)

func ByteFromFloat(v float64) byte {
	return byte(uint8(math.MaxUint8 * v))
}

/*
	A WorldState is a collection of named control points. Only one of these
	is the actual reality but they exist in limbo, can be modified, and can
	be copied into each other.
*/

type WorldState struct {
	name                string
	controlPoints       []ControlPoint
	controlPointsByName map[string]ControlPoint
	patchesNode         *config.AclNode
}

func NewWorldState(name string) *WorldState {
	ws := &WorldState{
		name:                name,
		controlPoints:       make([]ControlPoint, 0),
		controlPointsByName: make(map[string]ControlPoint),
	}

	return ws
}

func NewWorldStateFromNode(name string, root *config.AclNode) *WorldState {
	ws := &WorldState{
		name: name,
	}

	ws.controlPoints, ws.controlPointsByName = CreateControlPointList(root.Child("control_points"))

	ws.patchesNode = root.Child("patches")

	return ws
}

func (ws *WorldState) Apply(other *WorldState) {
	if other == nil {
		return
	}

	for i := 0; i < len(other.controlPoints); i++ {
		theirs := other.controlPoints[i]
		mine := ws.controlPointsByName[theirs.Name()]
		mine.Apply(theirs)
	}
}

func (ws *WorldState) String() string {
	if ws == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(ws.name)
	b.WriteString(" {\n")
	for i := 0; i < len(ws.controlPoints); i++ {
		b.WriteString(ws.controlPoints[i].String())
		b.WriteByte('\n')
	}
	b.WriteByte('}')

	return b.String()
}

func (ws *WorldState) AddControlPoint(name string, cp ControlPoint) {
	if ws == nil {
		return
	}

	ws.RemoveControlPoint(name)
	ws.controlPoints = append(ws.controlPoints, cp)
	ws.controlPointsByName[name] = cp
}

func (ws *WorldState) RemoveControlPoint(name string) {
	if ws == nil {
		return
	}

	cp := ws.controlPointsByName[name]
	if cp == nil {
		return
	}

	ws.controlPointsByName[name] = nil

	for i, v := range ws.controlPoints {
		if v == cp {
			ws.controlPoints[i] = ws.controlPoints[len(ws.controlPoints)-1]
			ws.controlPoints = ws.controlPoints[:len(ws.controlPoints)-1]
			return
		}
	}
}

// Calls SetTo on all ControlPoints but does not add it's own
// name to the path (the control points are expected to do that though)
func (ws *WorldState) SetToNode(root *config.AclNode, path ...string) {
	if ws == nil {
		return
	}

	for _, cp := range ws.controlPoints {
		cp.SetToNode(root, path...)
	}
}

func (ws *WorldState) Copy() *WorldState {
	if ws == nil {
		return nil
	}

	fresh := NewWorldState(ws.name)

	for _, cp := range ws.controlPoints {
		freshCP := cp.Copy()
		fresh.controlPoints = append(fresh.controlPoints, freshCP)
		fresh.controlPointsByName[freshCP.Name()] = freshCP
	}

	// Inefficient, so don't want to copy a whole lot...
	fresh.patchesNode = ws.patchesNode.Duplicate()

	return fresh
}

func (ws *WorldState) ControlPoint(name string) ControlPoint {
	if ws == nil {
		return nil
	}

	return ws.controlPointsByName[name]
}

func (ws *WorldState) ControlPoints() []ControlPoint {
	return ws.controlPoints
}

// func (ws *WorldState) Duplicate() *WorldState {
// 	out := NewWorldState()
// 	if ws == nil {
// 		return out
// 	}
// 	out.Root = ws.Root.Duplicate()
// 	return out
// }

// func (ws *WorldState) GetFloat(name string) float64 {
// 	return ws.Root.ChildAsFloat("values", name)
// }

// func (ws *WorldState) SetFloat(name string, val float64) {
// 	ws.Root.SetValAt(val, "values", name)
// }

// func (ws *WorldState) GetColor(name string) *WorldColor {
// 	color := NewWorldColor()
// 	color.GetFrom(ws.Root, "colors", name)

// 	return color
// }

// func (ws *WorldState) SetColor(name string, color *WorldColor) {
// 	color.SetTo(ws.Root, "colors", name)
// }

// func (ws *WorldState) GetPoint(name string) *WorldPoint {
// 	point := NewWorldPoint()
// 	point.GetFrom(ws.Root, "points", name)

// 	return point
// }

// func (ws *WorldState) SetPoint(name string, point *WorldPoint) {
// 	point.SetTo(ws.Root, "points", name)
// }

// type WorldColor struct {
// 	Red   float64
// 	Green float64
// 	Blue  float64

// 	White float64
// 	Amber float64
// 	Pink  float64

// 	UV float64
// }

// func NewWorldColor() *WorldColor {
// 	return &WorldColor{}
// }

// func (wc *WorldColor) GetFrom(root *config.AclNode, path ...string) {
// 	if wc == nil {
// 		return
// 	}

// 	fp := append(path, "red")
// 	wc.Red = root.ChildAsFloat(fp...)

// 	fp = append(path, "green")
// 	wc.Green = root.ChildAsFloat(fp...)

// 	fp = append(path, "blue")
// 	wc.Blue = root.ChildAsFloat(fp...)

// 	fp = append(path, "white")
// 	wc.White = root.ChildAsFloat(fp...)

// 	fp = append(path, "amber")
// 	wc.Amber = root.ChildAsFloat(fp...)

// 	fp = append(path, "pink")
// 	wc.Pink = root.ChildAsFloat(fp...)

// 	fp = append(path, "uv")
// 	wc.UV = root.ChildAsFloat(fp...)
// }

// func (wc *WorldColor) SetTo(root *config.AclNode, path ...string) {
// 	if wc == nil {
// 		return
// 	}

// 	fp := append(path, "red")
// 	root.SetValAt(wc.Red, fp...)

// 	fp = append(path, "green")
// 	root.SetValAt(wc.Green, fp...)

// 	fp = append(path, "blue")
// 	root.SetValAt(wc.Blue, fp...)

// 	fp = append(path, "white")
// 	root.SetValAt(wc.White, fp...)

// 	fp = append(path, "amber")
// 	root.SetValAt(wc.Amber, fp...)

// 	fp = append(path, "pink")
// 	root.SetValAt(wc.Pink, fp...)

// 	fp = append(path, "uv")
// 	root.SetValAt(wc.UV, fp...)
// }

// func (wc *WorldColor) String() string {
// 	return fmt.Sprintf("rgb(%v, %v, %v) w:%v a:%v p:%v u:%v",
// 		wc.Red, wc.Green, wc.Blue,
// 		wc.White, wc.Amber, wc.Pink,
// 		wc.UV)
// }

// ////////

// type WorldPoint struct {
// 	X float64
// 	Y float64
// 	Z float64
// }

// func NewWorldPoint() *WorldPoint {
// 	return &WorldPoint{}
// }

// func (wc *WorldPoint) GetFrom(root *config.AclNode, path ...string) {
// 	if wc == nil {
// 		return
// 	}

// 	fp := append(path, "x")
// 	wc.X = root.ChildAsFloat(fp...)

// 	fp = append(path, "y")
// 	wc.Y = root.ChildAsFloat(fp...)

// 	fp = append(path, "z")
// 	wc.Z = root.ChildAsFloat(fp...)
// }

// func (wc *WorldPoint) SetTo(root *config.AclNode, path ...string) {
// 	if wc == nil {
// 		return
// 	}

// 	fp := append(path, "x")
// 	root.SetValAt(wc.X, fp...)

// 	fp = append(path, "y")
// 	root.SetValAt(wc.Y, fp...)

// 	fp = append(path, "z")
// 	root.SetValAt(wc.Z, fp...)
// }

// func (wc *WorldPoint) String() string {
// 	return fmt.Sprintf("point(%v, %v, %v)",
// 		wc.X, wc.Y, wc.Z)
// }

// ////////
// type WorldState struct {
// 	Root *config.AclNode
// }

// func NewWorldState() *WorldState {
// 	ws := &WorldState{
// 		Root: config.NewAclNode(),
// 	}

// 	return ws
// }

// func (ws *WorldState) Duplicate() *WorldState {
// 	out := NewWorldState()
// 	if ws == nil {
// 		return out
// 	}
// 	out.Root = ws.Root.Duplicate()
// 	return out
// }

// func (ws *WorldState) GetFloat(name string) float64 {
// 	return ws.Root.ChildAsFloat("values", name)
// }

// func (ws *WorldState) SetFloat(name string, val float64) {
// 	ws.Root.SetValAt(val, "values", name)
// }

// func (ws *WorldState) GetColor(name string) *WorldColor {
// 	color := NewWorldColor()
// 	color.GetFrom(ws.Root, "colors", name)

// 	return color
// }

// func (ws *WorldState) SetColor(name string, color *WorldColor) {
// 	color.SetTo(ws.Root, "colors", name)
// }

// func (ws *WorldState) GetPoint(name string) *WorldPoint {
// 	point := NewWorldPoint()
// 	point.GetFrom(ws.Root, "points", name)

// 	return point
// }

// func (ws *WorldState) SetPoint(name string, point *WorldPoint) {
// 	point.SetTo(ws.Root, "points", name)
// }
