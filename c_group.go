package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"github.com/tomseago/go-eltee/api"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////
//
// Control Kind: group
//

type GroupUpdater struct {
	children []*FixtureControl
}

func (u *GroupUpdater) Update(fc *FixtureControl) {
	// At the moment we don't look at a control point for groups. We might want
	// to make them enableable or some such, but the group thing is currently
	// not all that exciting outside of having a hierarchy of controls (for no
	// real reason...)
	if u == nil || fc == nil {
		return
	}

	u.ForEachFixtureControl(func(child *FixtureControl) {
		child.Updater.Update(child)
	})
}

func (u *GroupUpdater) ForEachFixtureControl(fn func(*FixtureControl)) {

	for i := 0; i < len(u.children); i++ {
		child := u.children[i]

		fn(child)

		// Possibly recurse into it
		grpChild, ok := child.Updater.(*GroupUpdater)
		if ok {
			grpChild.ForEachFixtureControl(fn)
		}
	}
}

////////////////

type GroupProfileControl struct {
	ProfileControlBase

	Controls     []ProfileControl
	ControlsById map[string]ProfileControl
}

func NewGroupProfileControl(id string, rootNode *config.AclNode) (*GroupProfileControl, error) {
	pcg := &GroupProfileControl{
		ProfileControlBase: ProfileControlBase{
			id: id,
		},

		Controls:     make([]ProfileControl, 0),
		ControlsById: make(map[string]ProfileControl),
	}

	if rootNode == nil {
		// That's it, zero value yeah...
		return pcg, nil
	}

	// Iterate children in order
	keys := rootNode.OrderedChildNames
	for kix := 0; kix < len(keys); kix++ {
		key := keys[kix]

		// We could access Children directly but it's probably not
		// super nice to do so...
		child := rootNode.Child(key)
		control, err := NewControlFromConfig(key, child)
		if err != nil {
			return nil, err
		}

		pcg.Controls = append(pcg.Controls, control)
		pcg.ControlsById[key] = control
	}

	return pcg, nil
}

func (pc *GroupProfileControl) Id() string {
	return pc.id
}

func (pc *GroupProfileControl) Name() string {
	return pc.name
}

func (pc *GroupProfileControl) Type() string {
	return "group"
}

func (pc *GroupProfileControl) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Group %v(%v)", pc.name, pc.id))

	for i := 0; i < len(pc.Controls); i++ {
		child := pc.Controls[i]
		b.WriteString("\n  ")
		b.WriteString(child.String())
	}

	return b.String()
}

func (pc *GroupProfileControl) Instantiate(fixture Fixture) *FixtureControl {
	updater := &GroupUpdater{
		children: make([]*FixtureControl, len(pc.Controls)),
	}

	for i := 0; i < len(pc.Controls); i++ {
		childPC := pc.Controls[i]
		updater.children[i] = childPC.Instantiate(fixture)
	}

	fc := NewFixtureControl(pc, updater)

	fixture.AttachControl(pc.id, fc)

	return fc
}

func (pc *GroupProfileControl) ForEachControl(fn func(ProfileControl)) {

	for i := 0; i < len(pc.Controls); i++ {
		child := pc.Controls[i]

		fn(child)
	}
}

func (pc *GroupProfileControl) ToAPI() *api.ProfileControl {

	aPc := &api.GroupProfileControl{
		Id:   pc.id,
		Name: pc.name,

		Controls: make([]*api.ProfileControl, len(pc.Controls)),
	}

	for i := 0; i < len(pc.Controls); i++ {
		aPc.Controls[i] = pc.Controls[i].ToAPI()
	}

	aRet := &api.ProfileControl{
		Sub: &api.ProfileControl_Group{aPc},
	}

	return aRet
}
