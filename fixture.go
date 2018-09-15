package eltee

import (
	"sort"
)

// "github.com/eyethereal/go-config"

type Fixture interface {
	// Say my name... It is unique
	Name() string

	// What type of a fixture are we?
	Profile() *Profile

	// Controls() returns the root group of the FixtureControl hierarchy, which
	// is itself a FixtureControl
	RootControl() *FixtureControl

	// // The base is only relevant for DMX things, but we keep it in the generic
	// // interface because the relationship between DmxFixture instances is
	// // relevant
	// Base() int
	// SetBase(base int)

	AttachControl(id string, fc *FixtureControl)
	ForEachFixtureControl(func(string, *FixtureControl))
	Control(id string) *FixtureControl

	Update()

	GetInt(id string) int
	SetInt(id string, val int)
}

type DmxFixture struct {
	name        string
	base        int
	channels    []byte
	profile     *Profile
	rootControl *FixtureControl

	controls     []*FixtureControl
	controlsById map[string]*FixtureControl

	variables map[string]int
}

func NewDmxFixture(name string, base int, channels []byte, profile *Profile) *DmxFixture {
	f := &DmxFixture{
		name:     name,
		base:     base,
		channels: channels,
		profile:  profile,

		controls:     make([]*FixtureControl, 0),
		controlsById: make(map[string]*FixtureControl),

		variables: make(map[string]int),
	}
	f.rootControl = profile.Controls.Instantiate(f)

	toCopy := profile.DefaultData
	if toCopy != nil {
		if len(toCopy) > len(channels) {
			toCopy = channels[0:len(channels)]
		}
		copy(channels, toCopy)
	}

	return f
}

func (f *DmxFixture) Name() string {
	return f.name
}

func (f *DmxFixture) Profile() *Profile {
	return f.profile
}

func (f *DmxFixture) RootControl() *FixtureControl {
	return f.rootControl
}

func (f *DmxFixture) AttachControl(id string, fc *FixtureControl) {
	if f == nil || fc == nil {
		return
	}

	f.controls = append(f.controls, fc)
	f.controlsById[id] = fc

	// And a reference loop - yay!!!
	fc.Fixture = f
}

func (f *DmxFixture) ForEachFixtureControl(fn func(string, *FixtureControl)) {
	if f == nil {
		return
	}

	var ids []string
	for id, _ := range f.controlsById {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		fn(id, f.controlsById[id])
	}
}

func (f *DmxFixture) Control(id string) *FixtureControl {
	if f == nil {
		return nil
	}

	return f.controlsById[id]
}

func (f *DmxFixture) Update() {
	if f == nil {
		return
	}

	// Instead of doing this update, we just tell the group to update and
	// that cascades down. This makes sense if we want to enable and disable
	// groups. Or we could ignore groups. Or not have groups do recursive updates.
	// For now, it seems like, let the root group do it is the way to go.
	//
	// f.ForEachFixtureControl(func(id string, fc *FixtureControl) {
	// 	fc.Updater.Update(fc)
	// })

	f.rootControl.Updater.Update(f.rootControl)
}

func (f *DmxFixture) GetInt(id string) int {
	return f.variables[id]
}

func (f *DmxFixture) SetInt(id string, val int) {
	f.variables[id] = val
}
