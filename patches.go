package eltee

import (
// "github.com/eyethereal/go-config"
)

type FixturePatch struct {
	FixtureId      string
	ControlId      string
	ControlPointId string
}

type FixturePatchSet struct {
	Name           string
	FixturePatches []FixturePatch
}

type InputPatch struct {
	AdapterId        string
	AdapterControlId string
	ControlPointId   string
}
