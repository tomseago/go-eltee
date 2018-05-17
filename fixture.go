package eltee

import (
// "github.com/eyethereal/go-config"
)

type Fixture interface {
	// Say my name... It is unique
	Name() string

	// What type of a fixture are we?
	Profile() *Profile

	// Controls() returns the root group of the FixtureControl hierarchy, which
	// is itself a FixtureControl
	Controls() *FixtureControl

	// // The base is only relevant for DMX things, but we keep it in the generic
	// // interface because the relationship between DmxFixture instances is
	// // relevant
	// Base() int
	// SetBase(base int)
}

type DmxFixture struct {
	name     string
	base     int
	channels []byte
	profile  *Profile
	controls *FixtureControl
}

func NewDmxFixture(name string, base int, channels []byte, profile *Profile) *DmxFixture {
	f := &DmxFixture{
		name:     name,
		base:     base,
		channels: channels,
		profile:  profile,
	}
	f.controls = profile.Controls.Instantiate(f)

	return f
}

func (f *DmxFixture) Name() string {
	return f.name
}

func (f *DmxFixture) Profile() *Profile {
	return f.profile
}

func (f *DmxFixture) Controls() *FixtureControl {
	return f.controls
}

// func (f *DmxFixture) Base() int {
// 	return f.base
// }

// func (f *DmxFixture) SetBase(base int) {
// 	f.base = base
// }
