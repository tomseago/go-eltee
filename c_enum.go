package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"github.com/tomseago/go-eltee/api"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////
//
// Control Kind: enum
//

type EnumUpdater struct {
	// This is a reference to the slice that is held by the DmxFixture. Thus changes
	// made here go directly into the same place that the fixture references without
	// a need to copy them.

	// Since this control affects a single channel this slice is expected to be
	// only 1 element long
	channels []byte
}

func (u *EnumUpdater) Update(fc *FixtureControl) {
	if u == nil || fc == nil || fc.ControlPoint == nil {
		return
	}

	// TODO: Use the control point WasDirty() method?

	value := fc.LensStack.Observe(fc.ControlPoint)

	point, ok := value.(EnumPoint)
	if !ok {
		// not an intensity
		return
	}

	pc, ok := fc.ProfileControl.(*EnumProfileControl)
	if !ok {
		// attached to the wrong type of control. Bad
		return
	}

	item, degree := point.Option()

	// Figure out a value...
	if item >= len(pc.options) {
		item = len(pc.options) - 1
	}

	eo := pc.options[item]
	var dmxVal byte = 0
	if len(eo.VariableName) > 0 {
		iVal := fc.Fixture.GetInt(eo.VariableName) + eo.VariableOffset
		dmxVal = byte(iVal)
	} else {
		if len(eo.Values) == 1 {
			dmxVal = byte(eo.Values[0])
		} else if len(eo.Values) == 2 {
			dmxVal = byte(float64(eo.Values[1]-eo.Values[0])*degree + float64(eo.Values[0]))
		} else {
			// TODO: Map across n divisions
		}
	}

	// Set it
	u.channels[0] = dmxVal
}

////////////////

type EnumOption struct {
	Name string

	Values []int

	// An enum option can get it's value entirely from a fixture variable. This
	// is used to support the odd bit mapped values for the stinger where
	// we want to have toggle controls which control bits in a var on the fixture
	// and then an enum reads that var to figure out it's value.
	// This concept is either brilliant or insane. Not sure which. It might come
	// up with the laser also.
	VariableName   string
	VariableOffset int
}

func NewEnumOption(name string, node *config.AclNode) *EnumOption {
	eo := &EnumOption{
		Name:   name,
		Values: make([]int, 0),

		VariableName:   node.ChildAsString("variable", "name"),
		VariableOffset: node.ChildAsInt("variable", "offset"),
	}

	// If we have a variable name we ignore values
	if len(eo.VariableName) > 0 {
		return eo
	}

	// Copy the values
	for ix := 0; ix < node.Len(); ix++ {
		eo.Values = append(eo.Values, node.AsIntN(ix))
	}

	return eo
}

func (eo *EnumOption) String() string {
	var b strings.Builder
	b.WriteString("<'")
	b.WriteString(eo.Name)
	b.WriteString("' ")
	if len(eo.VariableName) > 0 {
		b.WriteString("${")
		b.WriteString(eo.VariableName)
		b.WriteString(" ")
		b.WriteString(fmt.Sprintf("%v", eo.VariableOffset))
		b.WriteString("}")
	} else {
		b.WriteString(fmt.Sprintf("%v", eo.Values))
	}
	b.WriteString(">")
	return b.String()
}

func (eo *EnumOption) ToAPI() *api.EnumProfileControlOption {
	r := &api.EnumProfileControlOption{
		Name:      eo.Name,
		Values:    make([]int32, len(eo.Values)),
		VarName:   eo.VariableName,
		VarOffset: int32(eo.VariableOffset),
	}

	for ix, val := range eo.Values {
		r.Values[ix] = int32(val)
	}

	return r
}

////////////////

type EnumProfileControl struct {
	ProfileControlBase

	channelIx int
	options   []*EnumOption
}

func NewEnumProfileControl(id string, rootNode *config.AclNode) (*EnumProfileControl, error) {
	pc := &EnumProfileControl{
		ProfileControlBase: MakeProfileControlBase(id, rootNode),

		channelIx: rootNode.ChildAsInt("channel"),
		options:   make([]*EnumOption, 0),
	}

	if pc.channelIx < 1 {
		return nil, fmt.Errorf("Channel index must be > 0")
	}

	vNode := rootNode.Child("values")
	vNode.ForEachOrderedChild(func(nn string, node *config.AclNode) {
		val := NewEnumOption(nn, node)
		pc.options = append(pc.options, val)
	})

	return pc, nil
}

func (pc *EnumProfileControl) Id() string {
	return pc.id
}

func (pc *EnumProfileControl) Name() string {
	return pc.name
}

func (pc *EnumProfileControl) Type() string {
	return "enum"
}

func (pc *EnumProfileControl) String() string {
	return fmt.Sprintf("Enum %v(%v) %v->%v", pc.name, pc.id, pc.channelIx, pc.options)
}

func (pc *EnumProfileControl) Instantiate(fixture Fixture) *FixtureControl {
	dmx, ok := fixture.(*DmxFixture)
	if !ok {
		return nil
	}

	updater := &EnumUpdater{
		channels: dmx.channels[pc.channelIx-1 : pc.channelIx],
	}

	fc := NewFixtureControl(pc, updater)

	fixture.AttachControl(pc.id, fc)

	return fc
}

func (pc *EnumProfileControl) ToAPI() *api.ProfileControl {

	aPc := &api.EnumProfileControl{
		Id:   pc.id,
		Name: pc.name,

		ChannelIx: int32(pc.channelIx),
		Options:   make([]*api.EnumProfileControlOption, 0),
	}

	for _, opt := range pc.options {
		aPc.Options = append(aPc.Options, opt.ToAPI())
	}

	aRet := &api.ProfileControl{
		Sub: &api.ProfileControl_Enm{aPc},
	}

	return aRet
}
