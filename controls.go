package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
)

// An instance of fixture control updater is attached to each fixture control and knows
// how to use the control point in the fixture control instance to update some DMX
// values (which are known to the updater).
type FixtureControlUpdater interface {
	// Causes the FixtureControlUpdater to observe the control point
	// and update it's output state
	Update(fc *FixtureControl)
}

//////

// A Fixture Control binds together a number of concepts needed to produce a DMX update.
//  * the Fixture it is attached to (which has some state data that is
//    almost certainly of interest to the Updater)
//  * the ProfileControl which tells us the shared interesting data about this control
//    like type, ranges, etc.
//  * the current ControlPoint which should be referenced to perform any updates
//  * an optional LensStack through which the ControlPoint should be observed by the
//    updater when it's trying to figure out new DMX output values
//  * the Updater which is the logic used to get from the ControlPoint to DMX output
type FixtureControl struct {
	Fixture        Fixture
	ProfileControl ProfileControl
	ControlPoint   ControlPoint
	LensStack      *LensStack

	// This is where the ProfileControl injects behavior into the FixtureControl instance
	Updater FixtureControlUpdater
}

// Creates a new fixture control which while it nominally has an updater attached it does not
// yet have a Fixture, ControlPoint, or LensStack
func NewFixtureControl(profileControl ProfileControl, updater FixtureControlUpdater) *FixtureControl {
	return &FixtureControl{
		ProfileControl: profileControl,
		Updater:        updater,
	}
}

//////

// A ProfileControl is a control on a theoretical fixture. Profile Controls can be
// used to instantiate FixtureControls which are bound to both this profile control
// and the particular Fixture given during instantiation. These resultant FixtureControls
// are the interesting things which bind together enough elements to actually do
// useful work translating from control points to dmx.
type ProfileControl interface {
	Id() string
	Name() string
	Type() string

	String() string

	// A Profile can be instantiated for a particular fixture. This instance
	// is expected to be held by the fixture
	Instantiate(fixture Fixture) *FixtureControl
}

/////////

type ProfileControlBase struct {
	id   string
	name string
}

func MakeProfileControlBase(id string, rootNode *config.AclNode) ProfileControlBase {
	pcb := ProfileControlBase{
		id:   id,
		name: rootNode.ChildAsString("name"),
	}
	return pcb
}

//////////////////

func NewControlFromConfig(id string, node *config.AclNode) (ProfileControl, error) {
	if node == nil {
		return nil, fmt.Errorf("Can not create a control from nil node for id %v", id)
	}

	kind := node.ChildAsString("kind")
	if kind == "" {
		return nil, fmt.Errorf("kind parameter was empty for id %v", id)
	}

	//log.Debugf("id=%v  %v", id, node.ColoredString())

	var control ProfileControl
	var err error

	switch kind {
	case "group":
		control, err = NewGroupProfileControl(id, node)

	case "led_var":
		control, err = NewLedVarProfileControl(id, node)

	case "intensity":
		control, err = NewIntensityProfileControl(id, node)

	case "pan_tilt":
		control, err = NewPanTiltProfileControl(id, node)

	case "enum":
		control, err = NewEnumProfileControl(id, node)

	}

	if err != nil {
		return nil, err
	}

	if control == nil {
		return nil, fmt.Errorf("Do not know how to create a control of kind '%v'", kind)
	}

	return control, nil
}
