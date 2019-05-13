package eltee

import (
	"github.com/eyethereal/go-config"
	"github.com/tomseago/go-eltee/api"
)

type FixturePatch struct {
	FixtureName string

	// Control point names by control
	CpsByControl map[string]string

	// TODO: Add lens stack descriptors by Control as well
}

func NewFixturePatch() *FixturePatch {
	out := &FixturePatch{
		CpsByControl: make(map[string]string),
	}

	return out
}

func CreateFixturePatchFromNode(name string, node *config.AclNode) *FixturePatch {
	out := &FixturePatch{
		FixtureName:  name,
		CpsByControl: make(map[string]string),
	}

	node.ForEachOrderedChild(func(cName string, cNode *config.AclNode) {
		cpName := cNode.ChildAsString("cp")
		if len(cpName) > 0 {
			out.CpsByControl[cName] = cpName
		}

		// TODO: Lens stacks...
	})

	return out
}

func CreateFixturePatchList(node *config.AclNode) ([]*FixturePatch, map[string]*FixturePatch) {

	list := make([]*FixturePatch, 0)
	index := make(map[string]*FixturePatch)

	node.ForEachOrderedChild(func(fixName string, fixNode *config.AclNode) {
		patch := CreateFixturePatchFromNode(fixName, fixNode)

		list = append(list, patch)
		index[fixName] = patch
	})

	return list, index
}

// Does not add to path
func FixturePatchListToNode(list []*FixturePatch, root *config.AclNode, path ...string) {
	for _, fp := range list {
		fp.SetToNode(root, path...)
	}
}

// Adds my name to path
func (fp *FixturePatch) SetToNode(root *config.AclNode, path ...string) {

	if fp == nil {
		return
	}

	myPath := append(path, fp.FixtureName)

	for cName, cpName := range fp.CpsByControl {
		cPath := append(myPath, cName)
		root.SetValAt(cpName, cPath...)
	}
}

func (fp *FixturePatch) Copy() *FixturePatch {
	if fp == nil {
		return nil
	}

	out := &FixturePatch{
		FixtureName:  fp.FixtureName,
		CpsByControl: make(map[string]string),
	}

	for cName, cpName := range fp.CpsByControl {
		out.CpsByControl[cName] = cpName
	}

	return out
}

func (fp *FixturePatch) Apply(other *FixturePatch) {
	if fp == nil || other == nil {
		return
	}

	for cName, cpName := range other.CpsByControl {
		fp.CpsByControl[cName] = cpName
	}
}

func (fp *FixturePatch) ToApi() *api.FixturePatch {
	aFP := &api.FixturePatch{
		ByControl: make(map[string]*api.FCPatch),
	}

	for cName, cpName := range fp.CpsByControl {
		aFCP := &api.FCPatch{
			Cp: cpName,
		}

		aFP.ByControl[cName] = aFCP
	}

	return aFP
}

func (fp *FixturePatch) SetFromApi(aFP *api.FixturePatch) {
	for cName, fpc := range aFP.GetByControl() {
		cpName := fpc.GetCp()
		if len(cpName) > 0 {
			fp.CpsByControl[cName] = cpName
		}
	}
}

func (fp *FixturePatch) RemoveFromApi(aFP *api.FixturePatch) {
	for cName, fpc := range aFP.GetByControl() {
		cpName := fpc.GetCp()
		if len(cpName) > 0 {
			delete(fp.CpsByControl, cName)
		}
	}
}
