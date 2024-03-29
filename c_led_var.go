package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"github.com/tomseago/go-eltee/api"
)

/////////////////////////////////////////////////////////////////////////////////
//
// Control Kind: led_var
//

type LedVarUpdater struct {
	// This is a reference to the slice that is held by the DmxFixture. Thus changes
	// made here go directly into the same place that the fixture references without
	// a need to copy them.
	channels []byte
}

func (u *LedVarUpdater) Update(fc *FixtureControl) {
	if u == nil || fc == nil || fc.ControlPoint == nil {
		return
	}

	// TODO: Use the control point WasDirty() method?

	value := fc.LensStack.Observe(fc, fc.ControlPoint)

	color, ok := value.(ColorPoint)
	if !ok {
		// not a color
		return
	}

	pc, ok := fc.ProfileControl.(*LedVarProfileControl)
	if !ok {
		// attached to the wrong type of control. Bad
		return
	}

	for name, channelIx := range pc.ColorMap {
		if channelIx == 0 || channelIx > len(u.channels) {
			// Invalid channel number
			continue
		}
		// Adjust to network 0 index
		channelIx--

		cc := color.ColorComponent(name)
		u.channels[channelIx] = ByteFromFloat(cc)
		// log.Infof("cc %v %v = %v", name, cc, u.channels[channelIx])
	}

}

////////////////

type LedVarProfileControl struct {
	ProfileControlBase

	ColorMap map[string]int
}

func NewLedVarProfileControl(id string, rootNode *config.AclNode) (*LedVarProfileControl, error) {
	pc := &LedVarProfileControl{
		ProfileControlBase: MakeProfileControlBase(id, rootNode),

		ColorMap: make(map[string]int),
	}

	rootNode.Child("leds").ForEachOrderedChild(func(name string, child *config.AclNode) {
		pc.ColorMap[name] = child.AsInt()
	})

	return pc, nil
}

func (pc *LedVarProfileControl) Id() string {
	return pc.id
}

func (pc *LedVarProfileControl) Name() string {
	return pc.name
}

func (pc *LedVarProfileControl) Type() string {
	return "led_var"
}

func (pc *LedVarProfileControl) String() string {
	return fmt.Sprintf("LedVar %v(%v) %v", pc.name, pc.id, pc.ColorMap)
}

func (pc *LedVarProfileControl) Instantiate(fixture Fixture) *FixtureControl {
	dmx, ok := fixture.(*DmxFixture)
	if !ok {
		return nil
	}

	updater := &LedVarUpdater{
		channels: dmx.channels,
	}

	fc := NewFixtureControl(pc, updater)

	fixture.AttachControl(pc.id, fc)

	return fc
}

func (pc *LedVarProfileControl) ToAPI() *api.ProfileControl {
	aPc := &api.LedVarProfileControl{
		Id:   pc.id,
		Name: pc.name,

		ColorMap: make(map[string]int32),
	}

	for k, v := range pc.ColorMap {
		aPc.ColorMap[k] = int32(v)
	}

	aRet := &api.ProfileControl{
		Sub: &api.ProfileControl_LedVar{aPc},
	}

	return aRet
}
