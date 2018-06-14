package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
)

type FixtureControlUpdater interface {
	// Causes the FixtureControlUpdater to observe the control point
	// and update it's output state
	Update(fc *FixtureControl)
}

//////

type FixtureControl struct {
	Fixture        Fixture
	ProfileControl ProfileControl
	ControlPoint   ControlPoint
	LensStack      *LensStack

	// This is where the ProfileControl injects behavior into the FixtureControl instance
	Updater FixtureControlUpdater
}

func NewFixtureControl(profileControl ProfileControl, updater FixtureControlUpdater) *FixtureControl {
	return &FixtureControl{
		ProfileControl: profileControl,
		Updater:        updater,
	}
}

//////

// A ProfileControl is a control on a theoretical fixture.
type ProfileControl interface {
	Id() string
	Name() string

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
