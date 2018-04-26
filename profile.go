package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"io/ioutil"
	"math"
	"path/filepath"
	"sort"
	"strings"
)

const (
	RedIx   = 0
	GreenIx = 1
	BlueIx  = 2
	WhiteIx = 3
)

type WorldColorSettable interface {
	SetFromWorldColor(color *WorldColor)
}

type ProfileControlInstance interface {
	ProfileControl() ProfileControl
	SetDmx(channels []byte) error
}

type ProfileControl interface {
	Id() string
	Name() string

	String() string

	Instantiate() ProfileControlInstance
	SetDmx(inst ProfileControlInstance, channels []byte) error
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

// func (pcb *ProfileControlBase) String() string {
// 	return fmt.Sprintf("%v(%v)", pcb.Name, pcb.Id)
// }

// type ProfileControlRGBW struct {
//     ProfileControlBase

// }
// func NewProfileControlRGBW(id string, rootNode *config.AclNode) (*ProfileControlGroup, error) {
//     pc := &ProfileControlRGBW{
//         ProfileControlBase: MakeProfileControlBase(id, rootNode),

//     }

//     return pc, nil
// }

// /////////

type ProfileControlRGBWInstance struct {
	control *ProfileControlRGBW
	Values  [WhiteIx + 1]byte
}

func (inst *ProfileControlRGBWInstance) ProfileControl() ProfileControl {
	return inst.control
}

func (inst *ProfileControlRGBWInstance) SetDmx(channels []byte) error {
	return inst.control.SetDmx(inst, channels)
}

func (inst *ProfileControlRGBWInstance) SetFromWorldColor(color *WorldColor) {
	if inst == nil || color == nil {
		return
	}

	inst.Values[RedIx] = ByteFromFloat(color.Red)
	inst.Values[GreenIx] = ByteFromFloat(color.Green)
	inst.Values[BlueIx] = ByteFromFloat(color.Blue)
	inst.Values[WhiteIx] = ByteFromFloat(color.White)
}

type ProfileControlRGBW struct {
	ProfileControlBase

	Channels [WhiteIx + 1]int
}

func NewProfileControlRGBW(id string, rootNode *config.AclNode) (*ProfileControlRGBW, error) {
	pc := &ProfileControlRGBW{
		ProfileControlBase: MakeProfileControlBase(id, rootNode),
	}

	pc.Channels[RedIx] = rootNode.ChildAsInt("channels", "red")
	pc.Channels[GreenIx] = rootNode.ChildAsInt("channels", "green")
	pc.Channels[BlueIx] = rootNode.ChildAsInt("channels", "blue")
	pc.Channels[WhiteIx] = rootNode.ChildAsInt("channels", "white")

	return pc, nil
}

func (pc *ProfileControlRGBW) Id() string {
	return pc.id
}

func (pc *ProfileControlRGBW) Name() string {
	return pc.name
}

func (pc *ProfileControlRGBW) String() string {
	return fmt.Sprintf("RGBW %v(%v) %v", pc.name, pc.id, pc.Channels)
}

func (pc *ProfileControlRGBW) Instantiate() ProfileControlInstance {
	inst := &ProfileControlRGBWInstance{
		control: pc,
	}

	return inst
}

func (pc *ProfileControlRGBW) SetDmx(inst ProfileControlInstance, channels []byte) error {

	rgbwInst, ok := inst.(*ProfileControlRGBWInstance)
	if !ok {
		return fmt.Errorf("Tried to SetDmx on a rgbw with a bad instance type")
	}

	channels[pc.Channels[RedIx]-1] = byte(rgbwInst.Values[RedIx])
	channels[pc.Channels[GreenIx]-1] = byte(rgbwInst.Values[GreenIx])
	channels[pc.Channels[BlueIx]-1] = byte(rgbwInst.Values[BlueIx])
	channels[pc.Channels[WhiteIx]-1] = byte(rgbwInst.Values[WhiteIx])

	return nil
}

/////////

type ProfileControlLedVarInstance struct {
	control *ProfileControlLedVar
	Values  map[string]byte
}

func (inst *ProfileControlLedVarInstance) ProfileControl() ProfileControl {
	return inst.control
}

func (inst *ProfileControlLedVarInstance) SetDmx(channels []byte) error {
	return inst.control.SetDmx(inst, channels)
}

func (inst *ProfileControlLedVarInstance) SetFromWorldColor(color *WorldColor) {
	if inst == nil || color == nil {
		return
	}

	inst.Values["red"] = ByteFromFloat(color.Red)
	inst.Values["green"] = ByteFromFloat(color.Green)
	inst.Values["blue"] = ByteFromFloat(color.Blue)
	inst.Values["white"] = ByteFromFloat(color.White)
	inst.Values["amber"] = ByteFromFloat(color.Amber)
	inst.Values["pink"] = ByteFromFloat(color.Pink)
	inst.Values["uv"] = ByteFromFloat(color.UV)
}

type ProfileControlLedVar struct {
	ProfileControlBase

	Channels map[string]int
}

func NewProfileControlLedVar(id string, rootNode *config.AclNode) (*ProfileControlLedVar, error) {
	pc := &ProfileControlLedVar{
		ProfileControlBase: MakeProfileControlBase(id, rootNode),

		Channels: make(map[string]int),
	}

	rootNode.Child("leds").ForEachOrderedChild(func(name string, child *config.AclNode) {
		pc.Channels[name] = child.AsInt()
	})

	return pc, nil
}

func (pc *ProfileControlLedVar) Id() string {
	return pc.id
}

func (pc *ProfileControlLedVar) Name() string {
	return pc.name
}

func (pc *ProfileControlLedVar) String() string {
	return fmt.Sprintf("LedVar %v(%v) %v", pc.name, pc.id, pc.Channels)
}

func (pc *ProfileControlLedVar) Instantiate() ProfileControlInstance {
	inst := &ProfileControlLedVarInstance{
		control: pc,
		Values:  make(map[string]byte),
	}

	// Fill in the defaults so the outside world can know what channels there are
	for name, _ := range pc.Channels {
		inst.Values[name] = 0
	}

	return inst
}

func (pc *ProfileControlLedVar) SetDmx(inst ProfileControlInstance, channels []byte) error {

	ledVarInst, ok := inst.(*ProfileControlLedVarInstance)
	if !ok {
		return fmt.Errorf("Tried to SetDmx on a LedVar with a bad instance type")
	}

	for name, channelIx := range pc.Channels {
		if channelIx == 0 || channelIx > len(channels) {
			continue
		}
		// Adjust to network 0 index
		channelIx--

		channels[channelIx] = ledVarInst.Values[name]
	}

	return nil
}

/////////

type ProfileControlPanTiltInstance struct {
	control *ProfileControlPanTilt

	Pan   float64
	Tilt  float64
	Speed float64
}

func (inst *ProfileControlPanTiltInstance) ProfileControl() ProfileControl {
	return inst.control
}

func (inst *ProfileControlPanTiltInstance) SetDmx(channels []byte) error {
	return inst.control.SetDmx(inst, channels)
}

type ProfileControlPanTilt struct {
	ProfileControlBase

	PanCoarseCh  int
	PanFineCh    int
	TiltCoarseCh int
	TiltFineCh   int
	SpeedCh      int
}

func NewProfileControlPanTilt(id string, rootNode *config.AclNode) (*ProfileControlPanTilt, error) {
	pc := &ProfileControlPanTilt{
		ProfileControlBase: MakeProfileControlBase(id, rootNode),

		PanCoarseCh:  rootNode.ChildAsInt("pan", "coarse"),
		PanFineCh:    rootNode.ChildAsInt("pan", "fine"),
		TiltCoarseCh: rootNode.ChildAsInt("tilt", "coarse"),
		TiltFineCh:   rootNode.ChildAsInt("tilt", "fine"),

		SpeedCh: rootNode.ChildAsInt("speed"),
	}

	return pc, nil
}

func (pc *ProfileControlPanTilt) Id() string {
	return pc.id
}

func (pc *ProfileControlPanTilt) Name() string {
	return pc.name
}

func (pc *ProfileControlPanTilt) String() string {
	return fmt.Sprintf("PanTilt %v(%v) pan(%v %v) tilt(%v %v) sp(%v)", pc.name, pc.id, pc.PanCoarseCh, pc.PanFineCh, pc.TiltCoarseCh, pc.TiltFineCh, pc.SpeedCh)
}

func (pc *ProfileControlPanTilt) Instantiate() ProfileControlInstance {
	inst := &ProfileControlPanTiltInstance{
		control: pc,

		Pan:   0.5,
		Tilt:  0.5,
		Speed: 1.0,
	}

	return inst
}

func (pc *ProfileControlPanTilt) SetDmx(inst ProfileControlInstance, channels []byte) error {

	panTiltInst, ok := inst.(*ProfileControlPanTiltInstance)
	if !ok {
		return fmt.Errorf("Tried to SetDmx on a PanTilt with a bad instance type")
	}

	iPan := uint16(math.MaxUint16 * panTiltInst.Pan)
	iTilt := uint16(math.MaxUint16 * panTiltInst.Tilt)
	iSpeed := uint8(math.MaxUint8 * panTiltInst.Speed)

	if pc.PanCoarseCh > 0 && pc.PanCoarseCh <= len(channels) {
		channels[pc.PanCoarseCh-1] = byte(iPan >> 8)
		if pc.PanFineCh > 0 && pc.PanFineCh <= len(channels) {
			channels[pc.PanFineCh-1] = byte(iPan)
		}
	}

	if pc.TiltCoarseCh > 0 && pc.TiltCoarseCh <= len(channels) {
		channels[pc.TiltCoarseCh-1] = byte(iTilt >> 8)
		if pc.TiltFineCh > 0 && pc.TiltFineCh <= len(channels) {
			channels[pc.TiltFineCh-1] = byte(iTilt)
		}
	}

	if pc.SpeedCh > 0 && pc.SpeedCh <= len(channels) {
		channels[pc.SpeedCh-1] = byte(iSpeed)
	}

	return nil
}

/////////

type ProfileControlFaderInstance struct {
	control *ProfileControlFader
	Value   byte
}

func (inst *ProfileControlFaderInstance) ProfileControl() ProfileControl {
	return inst.control
}

func (inst *ProfileControlFaderInstance) SetDmx(channels []byte) error {
	return inst.control.SetDmx(inst, channels)
}

type ProfileControlFader struct {
	ProfileControlBase

	Channel int
	Range   []int
}

func NewProfileControlFader(id string, rootNode *config.AclNode) (*ProfileControlFader, error) {
	pc := &ProfileControlFader{
		ProfileControlBase: MakeProfileControlBase(id, rootNode),

		Channel: rootNode.ChildAsInt("channel"),
		Range:   rootNode.ChildAsIntList("range"),
	}

	return pc, nil
}

func (pc *ProfileControlFader) Id() string {
	return pc.id
}

func (pc *ProfileControlFader) Name() string {
	return pc.name
}

func (pc *ProfileControlFader) String() string {
	return fmt.Sprintf("Fader %v(%v) %v %v", pc.name, pc.id, pc.Channel, pc.Range)
}

func (pc *ProfileControlFader) Instantiate() ProfileControlInstance {
	inst := &ProfileControlFaderInstance{
		control: pc,
	}

	return inst
}

func (pc *ProfileControlFader) SetDmx(inst ProfileControlInstance, channels []byte) error {

	faderInst, ok := inst.(*ProfileControlFaderInstance)
	if !ok {
		return fmt.Errorf("Tried to SetDmx on a Fader with a bad instance type")
	}

	if pc.Channel != 0 && pc.Channel <= len(channels) {
		// Adjust to network 0 index
		channels[pc.Channel-1] = faderInst.Value
	}

	return nil
}

/////////

type PCEnumValue struct {
	Name    string
	Channel int
	Values  []int

	VariableName   string
	VariableOffset int
}

func (inst *PCEnumValue) SetDmx(channels []byte) error {
	//return inst.control.SetDmx(inst, channels)
	return nil
}

func NewPCEnumValue(name string, node *config.AclNode) (*PCEnumValue, error) {
	v := &PCEnumValue{
		Name:           name,
		Values:         make([]int, 0),
		VariableName:   node.ChildAsString("variable", "name"),
		VariableOffset: node.ChildAsInt("variable", "offset"),
	}

	if len(v.VariableName) > 0 {
		// That's enough
		return v, nil
	}

	valsNode := node.Child("v")
	if valsNode == nil {
		valsNode = node.Child("range")
	}

	for ix := 0; ix < valsNode.Len(); ix++ {
		v.Values = append(v.Values, valsNode.AsIntN(ix))
	}

	return v, nil
}

func (ev *PCEnumValue) String() string {
	var b strings.Builder
	b.WriteString("<'")
	b.WriteString(ev.Name)
	b.WriteString("' ")
	if len(ev.VariableName) > 0 {
		b.WriteString("${")
		b.WriteString(ev.VariableName)
		b.WriteString(" ")
		b.WriteString(fmt.Sprintf("%v", ev.VariableOffset))
		b.WriteString("}")
	} else {
		b.WriteString(fmt.Sprintf("%v", ev.Values))
	}
	b.WriteString(">")

	return b.String()
}

//

type ProfileControlEnum struct {
	ProfileControlBase

	Channel int
	Values  []*PCEnumValue
}

func NewProfileControlEnum(id string, rootNode *config.AclNode) (*ProfileControlEnum, error) {
	pc := &ProfileControlEnum{
		ProfileControlBase: MakeProfileControlBase(id, rootNode),

		Channel: rootNode.ChildAsInt("channel"),
		Values:  make([]*PCEnumValue, 0),
	}

	vNode := rootNode.Child("values")
	if vNode != nil {
		keys := vNode.OrderedChildNames
		for kix := 0; kix < len(keys); kix++ {
			name := keys[kix]

			val, err := NewPCEnumValue(name, vNode.Child(name))
			if err != nil {
				return nil, err
			}

			pc.Values = append(pc.Values, val)
		}
	}

	return pc, nil
}

func (pc *ProfileControlEnum) Id() string {
	return pc.id
}

func (pc *ProfileControlEnum) Name() string {
	return pc.name
}

func (pc *ProfileControlEnum) String() string {
	return fmt.Sprintf("Enum %v(%v) %v->%v", pc.name, pc.id, pc.Channel, pc.Values)
}

func (pc *ProfileControlEnum) Instantiate() ProfileControlInstance {
	return nil
}

func (pc *ProfileControlEnum) SetDmx(inst ProfileControlInstance, channels []byte) error {

	return nil
}

// /////////

type ProfileControlGroupInstance struct {
	control   *ProfileControlGroup
	Instances []ProfileControlInstance
}

func (inst *ProfileControlGroupInstance) ForEachControlInstance(fn func(ProfileControlInstance)) {

	for i := 0; i < len(inst.Instances); i++ {
		child := inst.Instances[i]

		fn(child)

		// Possibly recurse into it
		grpChild, ok := child.(*ProfileControlGroupInstance)
		if ok {
			grpChild.ForEachControlInstance(fn)
		}
	}
}

func (inst *ProfileControlGroupInstance) ProfileControl() ProfileControl {
	return inst.control
}

func (inst *ProfileControlGroupInstance) SetDmx(channels []byte) error {
	return inst.control.SetDmx(inst, channels)
}

type ProfileControlGroup struct {
	ProfileControlBase

	Controls     []ProfileControl
	ControlsById map[string]ProfileControl
}

func NewProfileControlGroup(id string, rootNode *config.AclNode) (*ProfileControlGroup, error) {
	pcg := &ProfileControlGroup{
		ProfileControlBase: ProfileControlBase{
			id: id,
		},

		Controls:     make([]ProfileControl, 0),
		ControlsById: make(map[string]ProfileControl),
	}

	if rootNode == nil {
		// That's it, zero value yeah...
		return pcg, nil
	}

	// Iterate children in order
	keys := rootNode.OrderedChildNames
	for kix := 0; kix < len(keys); kix++ {
		key := keys[kix]

		// We could access Children directly but it's probably not
		// super nice to do so...
		child := rootNode.Child(key)
		control, err := NewControlFromConfig(key, child)
		if err != nil {
			return nil, err
		}

		pcg.Controls = append(pcg.Controls, control)
		pcg.ControlsById[key] = control
	}

	return pcg, nil
}

func (pc *ProfileControlGroup) Id() string {
	return pc.id
}

func (pc *ProfileControlGroup) Name() string {
	return pc.name
}

func (pc *ProfileControlGroup) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Group %v(%v)", pc.name, pc.id))

	for i := 0; i < len(pc.Controls); i++ {
		child := pc.Controls[i]
		b.WriteString("\n  ")
		b.WriteString(child.String())
	}

	return b.String()
}

func (pc *ProfileControlGroup) Instantiate() ProfileControlInstance {
	inst := &ProfileControlGroupInstance{
		control:   pc,
		Instances: make([]ProfileControlInstance, len(pc.Controls)),
	}
	for i := 0; i < len(pc.Controls); i++ {
		child := pc.Controls[i]
		inst.Instances[i] = child.Instantiate()
	}
	return inst
}

func (pc *ProfileControlGroup) SetDmx(inst ProfileControlInstance, channels []byte) error {

	groupInst, ok := inst.(*ProfileControlGroupInstance)
	if !ok {
		return fmt.Errorf("Tried to SetDmx on a group with a bad instance type")
	}
	for i := 0; i < len(pc.Controls); i++ {
		child := pc.Controls[i]
		childInst := groupInst.Instances[i]
		err := child.SetDmx(childInst, channels)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pc *ProfileControlGroup) ForEachControl(fn func(ProfileControl)) {

	for i := 0; i < len(pc.Controls); i++ {
		child := pc.Controls[i]

		fn(child)
	}
}

/////////

func NewControlFromConfig(id string, node *config.AclNode) (ProfileControl, error) {
	if node == nil {
		return nil, fmt.Errorf("Can not create a control from nil node for id %v", id)
	}

	kind := node.ChildAsString("kind")
	if kind == "" {
		return nil, fmt.Errorf("kind parameter was empty for id %v", id)
	}

	log.Debugf("id=%v  %v", id, node.ColoredString())

	var control ProfileControl
	var err error

	switch kind {
	case "group":
		control, err = NewProfileControlGroup(id, node)

	case "rgbw":
		control, err = NewProfileControlRGBW(id, node)

	case "led_var":
		control, err = NewProfileControlLedVar(id, node)

	case "pan_tilt":
		control, err = NewProfileControlPanTilt(id, node)

	case "fader":
		control, err = NewProfileControlFader(id, node)

	case "enum":
		control, err = NewProfileControlEnum(id, node)

	}

	if err != nil {
		return nil, err
	}

	if control == nil {
		return nil, fmt.Errorf("Do not know how to create a control of kind '%v'", kind)
	}

	return control, nil
}

/////////

type Profile struct {
	Id string

	Name         string
	ChannelCount int

	Controls *ProfileControlGroup
}

func NewProfile(id string, rootNode *config.AclNode) (*Profile, error) {
	p := &Profile{
		Id: id,
	}

	p.Name = rootNode.DefChildAsString(id, "name")
	p.ChannelCount = rootNode.DefChildAsInt(0, "channel_count")

	var err error
	p.Controls, err = NewProfileControlGroup("", rootNode.Child("controls"))
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Profile) String() string {
	return fmt.Sprintf("%v: '%v' %v\n%v\n", p.Id, p.Name, p.ChannelCount, p.Controls)
}

/////////

type ProfileLibrary struct {
	Profiles map[string]*Profile
}

func NewProfileLibrary() *ProfileLibrary {
	lib := &ProfileLibrary{
		Profiles: make(map[string]*Profile),
	}

	return lib
}

func (lib *ProfileLibrary) LoadFile(id string, filename string) error {
	node := config.NewAclNode()

	err := node.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("While reading '%v' : %v", filename, err)
	}

	profile, err := NewProfile(id, node)
	if err != nil {
		return fmt.Errorf("While creating profile %v : %v", id, err)
	}

	lib.Profiles[id] = profile
	return nil
}

func (lib *ProfileLibrary) LoadDirectory(dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		log.Debugf("name=%v  ext=%v", file.Name(), ext)
		if ext != ".acl" {
			continue
		}

		base := filepath.Base(file.Name())
		base = base[:len(base)-4]

		full := filepath.Join(dirname, file.Name())
		log.Infof("Loading '%v' from '%v'", base, full)

		err = lib.LoadFile(base, full)
		if err != nil {
			log.Errorf("%v", err)
			// But try to load other things
		}
	}

	return nil
}

func (lib *ProfileLibrary) String() string {
	// Output in sorted order
	ids := make([]string, 0, len(lib.Profiles))
	for id, _ := range lib.Profiles {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var b strings.Builder
	b.WriteString("\n")

	for ix := 0; ix < len(ids); ix++ {
		id := ids[ix]
		profile := lib.Profiles[id]
		b.WriteString(profile.String())
		b.WriteString("\n")
	}

	return b.String()
}
