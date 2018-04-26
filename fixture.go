package eltee

import (
// "github.com/eyethereal/go-config"
)

type Fixture interface {
	Name() string

	Profile() *Profile
	ProfileControlInstance() ProfileControlInstance

	Base() int
	SetBase(base int)
}

type DmxFixture struct {
	name    string
	base    int
	profile *Profile
	inst    ProfileControlInstance
}

func NewDmxFixture(name string, profile *Profile) *DmxFixture {
	f := &DmxFixture{
		name:    name,
		profile: profile,
		inst:    profile.Controls.Instantiate(),
	}

	return f
}

func (f *DmxFixture) Name() string {
	return f.name
}

func (f *DmxFixture) Profile() *Profile {
	return f.profile
}

func (f *DmxFixture) ProfileControlInstance() ProfileControlInstance {
	return f.inst
}

func (f *DmxFixture) Base() int {
	return f.base
}

func (f *DmxFixture) SetBase(base int) {
	f.base = base
}
