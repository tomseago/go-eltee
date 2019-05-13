package eltee

import (
	"github.com/eyethereal/go-config"
	"github.com/tomseago/go-eltee/api"
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
	name string

	controlPoints       []ControlPoint
	controlPointsByName map[string]ControlPoint

	fixturePatches       []*FixturePatch
	fixturePatchesByName map[string]*FixturePatch
	//patchesNode         *config.AclNode
}

func NewWorldState(name string) *WorldState {
	ws := &WorldState{
		name:                name,
		controlPoints:       make([]ControlPoint, 0),
		controlPointsByName: make(map[string]ControlPoint),

		fixturePatches:       make([]*FixturePatch, 0),
		fixturePatchesByName: make(map[string]*FixturePatch),
	}

	return ws
}

func NewWorldStateFromNode(name string, root *config.AclNode) *WorldState {
	ws := &WorldState{
		name: name,
	}

	ws.controlPoints, ws.controlPointsByName = CreateControlPointList(root.Child("control_points"))

	ws.fixturePatches, ws.fixturePatchesByName = CreateFixturePatchList(root.Child("fixture_patches"))

	return ws
}

func (ws *WorldState) Apply(other *WorldState) {
	if other == nil {
		return
	}

	for i := 0; i < len(other.controlPoints); i++ {
		theirs := other.controlPoints[i]
		mine := ws.controlPointsByName[theirs.Name()]
		if mine != nil {
			mine.Apply(theirs)
		}
	}

	for _, theirs := range other.fixturePatches {
		mine := ws.fixturePatchesByName[theirs.FixtureName]
		if mine != nil {
			mine.Apply(theirs)
		}
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

	fp := append(path, "control_points")

	for _, cp := range ws.controlPoints {
		cp.SetToNode(root, fp...)
	}

	pp := append(path, "fixture_patches")
	FixturePatchListToNode(ws.fixturePatches, root, pp...)
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

	for _, fp := range ws.fixturePatches {
		freshFP := fp.Copy()
		fresh.fixturePatches = append(fresh.fixturePatches, freshFP)
		fresh.fixturePatchesByName[freshFP.FixtureName] = freshFP
	}

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

func (ws *WorldState) FixturePatch(name string) *FixturePatch {
	if ws == nil {
		return nil
	}

	return ws.fixturePatchesByName[name]
}

func (ws *WorldState) FixturePatches() []*FixturePatch {
	return ws.fixturePatches
}

func (ws *WorldState) FixturePatchMap() *api.FixturePatchMap {

	fpMap := &api.FixturePatchMap{
		State:     ws.name,
		ByFixture: make(map[string]*api.FixturePatch),
	}

	for _, fp := range ws.fixturePatches {
		fpMap.ByFixture[fp.FixtureName] = fp.ToApi()
	}

	return fpMap
}

func (ws *WorldState) AddFixturePatch(fp *FixturePatch) {
	if ws == nil || fp == nil || len(fp.FixtureName) == 0 {
		return
	}

	if ws.fixturePatchesByName[fp.FixtureName] != nil {
		// exists alread
		return
	}

	ws.fixturePatchesByName[fp.FixtureName] = fp
	ws.fixturePatches = append(ws.fixturePatches, fp)
}
