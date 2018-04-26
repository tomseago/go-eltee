package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"math"
)

func ByteFromFloat(v float64) byte {
	return byte(uint8(math.MaxUint8 * v))
}

type WorldColor struct {
	Red   float64
	Green float64
	Blue  float64

	White float64
	Amber float64
	Pink  float64

	UV float64
}

func NewWorldColor() *WorldColor {
	return &WorldColor{}
}

func (wc *WorldColor) GetFrom(root *config.AclNode, path ...string) {
	if wc == nil {
		return
	}

	fp := append(path, "red")
	wc.Red = root.ChildAsFloat(fp...)

	fp = append(path, "green")
	wc.Green = root.ChildAsFloat(fp...)

	fp = append(path, "blue")
	wc.Blue = root.ChildAsFloat(fp...)

	fp = append(path, "white")
	wc.White = root.ChildAsFloat(fp...)

	fp = append(path, "amber")
	wc.Amber = root.ChildAsFloat(fp...)

	fp = append(path, "pink")
	wc.Pink = root.ChildAsFloat(fp...)

	fp = append(path, "uv")
	wc.UV = root.ChildAsFloat(fp...)
}

func (wc *WorldColor) SetTo(root *config.AclNode, path ...string) {
	if wc == nil {
		return
	}

	fp := append(path, "red")
	root.SetValAt(wc.Red, fp...)

	fp = append(path, "green")
	root.SetValAt(wc.Green, fp...)

	fp = append(path, "blue")
	root.SetValAt(wc.Blue, fp...)

	fp = append(path, "white")
	root.SetValAt(wc.White, fp...)

	fp = append(path, "amber")
	root.SetValAt(wc.Amber, fp...)

	fp = append(path, "pink")
	root.SetValAt(wc.Pink, fp...)

	fp = append(path, "uv")
	root.SetValAt(wc.UV, fp...)
}

func (wc *WorldColor) String() string {
	return fmt.Sprintf("rgb(%v, %v, %v) w:%v a:%v p:%v u:%v",
		wc.Red, wc.Green, wc.Blue,
		wc.White, wc.Amber, wc.Pink,
		wc.UV)
}

////////

type WorldPoint struct {
	X float64
	Y float64
	Z float64
}

func NewWorldPoint() *WorldPoint {
	return &WorldPoint{}
}

func (wc *WorldPoint) GetFrom(root *config.AclNode, path ...string) {
	if wc == nil {
		return
	}

	fp := append(path, "x")
	wc.X = root.ChildAsFloat(fp...)

	fp = append(path, "y")
	wc.Y = root.ChildAsFloat(fp...)

	fp = append(path, "z")
	wc.Z = root.ChildAsFloat(fp...)
}

func (wc *WorldPoint) SetTo(root *config.AclNode, path ...string) {
	if wc == nil {
		return
	}

	fp := append(path, "x")
	root.SetValAt(wc.X, fp...)

	fp = append(path, "y")
	root.SetValAt(wc.Y, fp...)

	fp = append(path, "z")
	root.SetValAt(wc.Z, fp...)
}

func (wc *WorldPoint) String() string {
	return fmt.Sprintf("point(%v, %v, %v)",
		wc.X, wc.Y, wc.Z)
}

////////
type WorldState struct {
	Root *config.AclNode
}

func NewWorldState() *WorldState {
	ws := &WorldState{
		Root: config.NewAclNode(),
	}

	return ws
}

func (ws *WorldState) Duplicate() *WorldState {
	out := NewWorldState()
	if ws == nil {
		return out
	}
	out.Root = ws.Root.Duplicate()
	return out
}

func (ws *WorldState) GetFloat(name string) float64 {
	return ws.Root.ChildAsFloat("values", name)
}

func (ws *WorldState) SetFloat(name string, val float64) {
	ws.Root.SetValAt(val, "values", name)
}

func (ws *WorldState) GetColor(name string) *WorldColor {
	color := NewWorldColor()
	color.GetFrom(ws.Root, "colors", name)

	return color
}

func (ws *WorldState) SetColor(name string, color *WorldColor) {
	color.SetTo(ws.Root, "colors", name)
}

func (ws *WorldState) GetPoint(name string) *WorldPoint {
	point := NewWorldPoint()
	point.GetFrom(ws.Root, "points", name)

	return point
}

func (ws *WorldState) SetPoint(name string, point *WorldPoint) {
	point.SetTo(ws.Root, "points", name)
}
