package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"github.com/tomseago/go-eltee/api"
)

/////////////////////////////////////////////////////////////////////////////////
//
// Control Kind: intensity
//

type IntensityUpdater struct {
	// This is a reference to the slice that is held by the DmxFixture. Thus changes
	// made here go directly into the same place that the fixture references without
	// a need to copy them.

	// Since this control affects a single channel this slice is expected to be
	// only 1 element long
	channels []byte
}

func (u *IntensityUpdater) Update(fc *FixtureControl) {
	if u == nil || fc == nil || fc.ControlPoint == nil {
		return
	}

	// TODO: Use the control point WasDirty() method?

	value := fc.LensStack.Observe(fc.ControlPoint)

	intensity, ok := value.(IntensityPoint)
	if !ok {
		// not an intensity
		return
	}

	pc, ok := fc.ProfileControl.(*IntensityProfileControl)
	_ = pc
	if !ok {
		// attached to the wrong type of control. Bad
		return
	}

	val := intensity.Percent()

	// TODO: Scaling. Min and max could be stored in the profile control. For now assume 0 to 255.

	u.channels[0] = byte(255.0 * val)
}

////////////////

type IntensityProfileControl struct {
	ProfileControlBase

	// In the future we may store min and max type values
	channelIx int
}

func NewIntensityProfileControl(id string, rootNode *config.AclNode) (*IntensityProfileControl, error) {
	pc := &IntensityProfileControl{
		ProfileControlBase: MakeProfileControlBase(id, rootNode),

		channelIx: rootNode.ChildAsInt("channel"),
	}

	if pc.channelIx < 1 {
		return nil, fmt.Errorf("Channel index must be > 0")
	}

	return pc, nil
}

func (pc *IntensityProfileControl) Id() string {
	return pc.id
}

func (pc *IntensityProfileControl) Name() string {
	return pc.name
}

func (pc *IntensityProfileControl) Type() string {
	return "intensity"
}

func (pc *IntensityProfileControl) String() string {
	return fmt.Sprintf("Intensity %v(%v) %v", pc.name, pc.id, pc.channelIx)
}

func (pc *IntensityProfileControl) Instantiate(fixture Fixture) *FixtureControl {
	dmx, ok := fixture.(*DmxFixture)
	if !ok {
		return nil
	}

	updater := &IntensityUpdater{
		channels: dmx.channels[pc.channelIx-1 : pc.channelIx],
	}

	fc := NewFixtureControl(pc, updater)

	fixture.AttachControl(pc.id, fc)

	return fc
}

func (pc *IntensityProfileControl) ToAPI() *api.ProfileControl {
	aPc := &api.IntensityProfileControl{
		Id:   pc.id,
		Name: pc.name,

		ChannelIx: int32(pc.channelIx),
	}

	aRet := &api.ProfileControl{
		Sub: &api.ProfileControl_Intensity{aPc},
	}

	return aRet
}
