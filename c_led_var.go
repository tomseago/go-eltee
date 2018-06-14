package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
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

	value := fc.LensStack.Observe(fc.ControlPoint)

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

		u.channels[channelIx] = ByteFromFloat(color.ColorComponent(name))
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
