package eltee

import (
	"github.com/eyethereal/go-config"
    "github.com/tomseago/go-eltee/api"
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
	GetF64(id string) float64
	SetF64(id string, val float64)

	ToAPI() *api.Fixture
}

type DmxFixture struct {
	name        string
	base        int
	channels    []byte
	profile     *Profile
	rootControl *FixtureControl

	controls     []*FixtureControl
	controlsById map[string]*FixtureControl

	varInts map[string]int
	varF64  map[string]float64

	overrideValues []byte
	useOverrides   bool
}

func NewDmxFixture(name string, base int, channels []byte, profile *Profile) *DmxFixture {
	f := &DmxFixture{
		name:     name,
		base:     base,
		channels: channels,
		profile:  profile,

		controls:     make([]*FixtureControl, 0),
		controlsById: make(map[string]*FixtureControl),

		varInts: make(map[string]int),
		varF64:  make(map[string]float64),

		overrideValues: make([]byte, len(channels)),
		useOverrides:   false,
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

	if f.useOverrides {
		copy(f.channels, f.overrideValues)
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
	return f.varInts[id]
}

func (f *DmxFixture) SetInt(id string, val int) {
	f.varInts[id] = val
}

func (f *DmxFixture) GetF64(id string) float64 {
	return f.varF64[id]
}

func (f *DmxFixture) SetF64(id string, val float64) {
	f.varF64[id] = val
}

func (f *DmxFixture) ToAPI()*api.Fixture {
    out := &api.Fixture{
        Name:                 f.name,
        ProfileId:            f.profile.Id,
        ControlState:         make(map[string]*api.FixtureControlState),
        VarInts:              make(map[string]int32),
        VarDoubles:           make(map[string]float64),
    }

    for k, v := range f.controlsById {
        out.ControlState[k] = v.ToAPI()
    }

    // Copy these so no ones messes with our internal stuff I guess
    for k, v := range f.varInts {
        out.VarInts[k] = int32(v)
    }

    for k, v := range f.varF64 {
        out.VarDoubles[k] = v
    }

    return out
}

func (f *DmxFixture) SetUseOverrides(use bool) {
	if f == nil {
		return
	}

	f.useOverrides = use
}

func (f *DmxFixture) GetUseOverrides() bool {
	if f == nil {
		return false
	}

	return f.useOverrides
}

func (f *DmxFixture) SetOverrides(overrides []byte) {
	if f == nil {
		return
	}

	copy(f.overrideValues, overrides)
}

func (f *DmxFixture) GetChannels() []byte {
	out := make([]byte, len(f.channels))
	copy(out, f.channels)
	return out
}

func (f *DmxFixture) LensesFrom(node *config.AclNode) {
	if node == nil {
		return
	}

	node.ForEachOrderedChild(func(cName string, lensStackNode *config.AclNode) {
		control := f.Control(cName)

		if control == nil {
			log.Warningf("Could not find a control %v to add lenses to", cName)
			return
		}

		log.Infof("Adding lens stack to %v", cName)
		control.LensStack = NewLensStackFromNode(lensStackNode)
	})
}

func (f *DmxFixture) VarsFrom(node *config.AclNode) {
	if node == nil {
		return
	}

	node.ForEachOrderedChild(func(vName string, vNode *config.AclNode) {
		log.Errorf("%v = %v", vName, vNode.AsFloat())
		f.SetF64(vName, vNode.AsFloat())
	})
}

