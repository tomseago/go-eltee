package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"github.com/tomseago/go-eltee/api"
	"sort"
	"strings"
)

type ControlPoint interface {
	// Control points shall be known uniquely by their name
	Name() string

	SetFromNode(node *config.AclNode)
	SetToNode(root *config.AclNode, path ...string)

	SetFromJSON(val interface{})

	ToApi() *api.ControlPoint
	SetFromApi(apiCP *api.ControlPoint)

	Apply(other ControlPoint)

	String() string
	Copy() ControlPoint
}

/*
   CreateControlPoint creates a control point from a configuration node.
   This is a switch between the string names and the NewXXXXControlPointFromNode
   functions.
*/
func CreateControlPoint(name string, node *config.AclNode) ControlPoint {
	kind := node.ChildAsString("kind")

	var cp ControlPoint
	switch kind {
	case "color":
		cp = NewColorControlPoint(name)

	case "xyz":
		cp = NewXYZControlPoint(name)

	case "enum":
		cp = NewEnumControlPoint(name)

	case "intensity":
		cp = NewIntensityControlPoint(name)
	}

	if cp == nil {
		log.Errorf("Unknown control point kind '%v'", kind)
		return nil
	}

	cp.SetFromNode(node)
	return cp
}

/*
   CreateControlPointList takes a root node which contains a map of
   names to control point definition nodes and then calls CreateControlPoint
   on each of those named definitions. The result is both an ordered list
   and an index by name of the created points.
*/
func CreateControlPointList(root *config.AclNode) ([]ControlPoint, map[string]ControlPoint) {
	list := make([]ControlPoint, 0)
	index := make(map[string]ControlPoint)

	root.ForEachOrderedChild(func(nn string, node *config.AclNode) {
		cp := CreateControlPoint(nn, node)
		if cp == nil {
			log.Errorf("Ignoring control point '%v'", nn)
			return
		}
		list = append(list, cp)
		index[nn] = cp
	})

	return list, index
}

////////////

// These types of Points represent raw data representations. We may want
// to add addtional classes that can convienently manipulate these raw
// values to do things like synthesize missing color components or translate
// positions etc.

type ColorPoint interface {
	ColorComponent(name string) float64
}
type SettableColorPoint interface {
	SetColorComponent(name string, val float64)
}

type XYZPoint interface {
	XYZ() (float64, float64, float64)
}
type SettableXYZPoint interface {
	SetXYZ(x float64, y float64, z float64)
}

type EnumPoint interface {
	// The 0 based item which is selected
	// The degree to which this item is selected. This is consulted when
	// the selected option represents a range of values
	Option() (int, float64)
}
type SettableEnumPoint interface {
	SetOption(item int, degree float64)
}

type IntensityPoint interface {
	Percent() float64
}
type SettableIntensityPoint interface {
	SetPercent(intensity float64)
}

////////////////////////////////////////

type ColorControlPoint struct {
	name       string
	Kind       string
	Components map[string]float64
}

func NewColorControlPoint(name string) *ColorControlPoint {
	cp := &ColorControlPoint{
		name: name,

		Kind:       "color",
		Components: make(map[string]float64),
	}

	return cp
}

func (cp *ColorControlPoint) Name() string {
	if cp == nil {
		return ""
	}

	return cp.name
}

func (cp *ColorControlPoint) ColorComponent(name string) float64 {
	if cp == nil {
		return 0.0
	}

	return cp.Components[name]
}

func (cp *ColorControlPoint) SetColorComponent(name string, val float64) {
	if cp == nil {
		return
	}

	cp.Components[name] = val
}

func (cp *ColorControlPoint) SetFromNode(node *config.AclNode) {
	node.ForEachOrderedChild(func(nn string, val *config.AclNode) {
		// protect reserved names, but otherwise accept everything
		// else as naming a color component
		if nn == "kind" {
			return
		}

		cp.Components[nn] = val.AsFloat()
	})
}

func (cp *ColorControlPoint) SetToNode(root *config.AclNode, path ...string) {
	fp := append(path, cp.name)

	kp := append(fp, "kind")
	root.SetValAt("color", kp...)

	for key, val := range cp.Components {
		valPath := append(fp, key)
		root.SetValAt(val, valPath...)
	}
}

func (cp *ColorControlPoint) SetFromJSON(val interface{}) {
	if cp == nil || val == nil {
		return
	}

	v, ok := val.(map[string]interface{})
	if !ok {
		return
	}

	for cName, cVal := range v {
		cp.Components[cName] = ValAsFloat(cVal)
	}
}

func (cp *ColorControlPoint) ToApi() *api.ControlPoint {
	val := &api.ControlPoint{
		Name: cp.name,
	}

	apiCP := &api.ColorPoint{
		Components: make(map[string]float64),
	}
	for k, v := range cp.Components {
		apiCP.Components[k] = v
	}
	val.Val = &api.ControlPoint_Color{apiCP}

	return val
}

func (cp *ColorControlPoint) SetFromApi(apiCP *api.ControlPoint) {
	log.Warningf("Setting %v from %v", cp, apiCP)

	if apiCP == nil || apiCP.Val == nil {
		return
	}

	apiVal := apiCP.GetColor()
	if apiVal == nil {
		return
	}

	for k, v := range apiVal.Components {
		cp.Components[k] = v
	}
}

func (cp *ColorControlPoint) Apply(other ControlPoint) {
	if cp == nil || other == nil {
		return
	}

	otherColor, ok := other.(*ColorControlPoint)
	if !ok {
		return
	}

	for name, val := range otherColor.Components {
		cp.Components[name] = val
	}
}

func (cp *ColorControlPoint) String() string {
	if cp == nil {
		return "Color(nil)"
	}

	ids := make([]string, 0, len(cp.Components))
	for id, _ := range cp.Components {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var b strings.Builder
	b.WriteString("Color(")

	for ix := 0; ix < len(ids); ix++ {
		id := ids[ix]
		val := cp.Components[id]
		b.WriteString(fmt.Sprintf("%v:%v", id, val))
		if ix < len(ids)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteByte(')')

	return b.String()
}

func (cp *ColorControlPoint) Copy() ControlPoint {
	ncp := &ColorControlPoint{
		name: cp.name,

		Kind:       "color",
		Components: make(map[string]float64),
	}

	for name, val := range cp.Components {
		ncp.Components[name] = val
	}

	return ncp
}

////////////////////////////////////////

/*
   The geometric space is defined as the XY plane being the ground plane
   with X increasing to the right and Y increasing in the forwards direction.
   Z increases vertically. This is a right handed coordinate system, but
   whereas the 'typical' isometric drawings of this space have the observer
   positioned along the left hand side, that is not a good picture of
   how the coordinates are defined for us. The relationship between the
   coordinates is the same, but the picture needs to be rotated 90 degrees
   ccw around the Z axis to match the way we envision space.

   Once everything actually works this is somewhat irrelevant, but we're going
   to keep this standard throughout the project and UI's can map onto it however
   they wish to get crazy.
*/

type XYZControlPoint struct {
	name string

	Kind string
	X    float64
	Y    float64
	Z    float64
}

func NewXYZControlPoint(name string) *XYZControlPoint {
	cp := &XYZControlPoint{
		name: name,
		Kind: "xyz",
	}

	return cp
}

func (cp *XYZControlPoint) Name() string {
	if cp == nil {
		return ""
	}

	return cp.name
}

func (cp *XYZControlPoint) XYZ() (float64, float64, float64) {
	if cp == nil {
		return 0.0, 0.0, 0.0
	}

	return cp.X, cp.Y, cp.Z
}

func (cp *XYZControlPoint) SetXYZ(x float64, y float64, z float64) {
	if cp == nil {
		return
	}

	cp.X = x
	cp.Y = y
	cp.Z = z
}

func (cp *XYZControlPoint) SetFromNode(node *config.AclNode) {
	cp.X = node.ChildAsFloat("x")
	cp.Y = node.ChildAsFloat("y")
	cp.Z = node.ChildAsFloat("z")
}

func (cp *XYZControlPoint) SetToNode(root *config.AclNode, path ...string) {
	if cp == nil {
		return
	}

	fp := append(path, cp.name)

	kp := append(fp, "kind")
	root.SetValAt("xyz", kp...)

	kp = append(fp, "x")
	root.SetValAt(cp.X, kp...)

	kp = append(fp, "y")
	root.SetValAt(cp.Y, kp...)

	kp = append(fp, "z")
	root.SetValAt(cp.Z, kp...)
}

func (cp *XYZControlPoint) SetFromJSON(val interface{}) {
	if cp == nil || val == nil {
		return
	}

	v, ok := val.(map[string]interface{})
	if !ok {
		return
	}

	cp.X = ValAsFloat(v["X"])
	cp.Y = ValAsFloat(v["Y"])
	cp.Z = ValAsFloat(v["Z"])
}

func (cp *XYZControlPoint) ToApi() *api.ControlPoint {
	val := &api.ControlPoint{
		Name: cp.name,
	}

	apiCP := &api.XYZPoint{
		X: cp.X,
		Y: cp.Y,
		Z: cp.Z,
	}
	val.Val = &api.ControlPoint_Xyz{apiCP}

	return val
}

func (cp *XYZControlPoint) SetFromApi(apiCP *api.ControlPoint) {
	if apiCP == nil || apiCP.Val == nil {
		return
	}

	apiVal := apiCP.GetXyz()
	if apiVal == nil {
		return
	}

	cp.X = apiVal.X
	cp.Y = apiVal.Y
	cp.Z = apiVal.Z
}

func (cp *XYZControlPoint) Apply(other ControlPoint) {
	if cp == nil || other == nil {
		return
	}

	otherXYZ, ok := other.(*XYZControlPoint)
	if !ok {
		return
	}

	cp.X = otherXYZ.X
	cp.Y = otherXYZ.Y
	cp.Z = otherXYZ.Z
}

func (cp *XYZControlPoint) String() string {
	if cp == nil {
		return "XYZ(nil)"
	}

	return fmt.Sprintf("XYZ(%v,%v,%v)", cp.X, cp.Y, cp.X)
}

func (cp *XYZControlPoint) Copy() ControlPoint {
	ncp := &XYZControlPoint{
		name: cp.name,

		Kind: "xyz",
		X:    cp.X,
		Y:    cp.Y,
		Z:    cp.Z,
	}

	return ncp
}

////////////////////////////////////////

type EnumControlPoint struct {
	name string

	Kind   string
	Item   int
	Degree float64
}

func NewEnumControlPoint(name string) *EnumControlPoint {
	cp := &EnumControlPoint{
		name: name,
		Kind: "enum",
	}

	return cp
}

func (cp *EnumControlPoint) Name() string {
	if cp == nil {
		return ""
	}

	return cp.name
}

func (cp *EnumControlPoint) Option() (int, float64) {
	if cp == nil {
		return 0, 0.0
	}

	return cp.Item, cp.Degree
}

func (cp *EnumControlPoint) SetOption(item int, degree float64) {
	if cp == nil {
		return
	}

	cp.Item = item
	cp.Degree = degree
}

func (cp *EnumControlPoint) SetFromNode(node *config.AclNode) {
	cp.Item = node.ChildAsInt("item")
	cp.Degree = node.ChildAsFloat("degree")
}

func (cp *EnumControlPoint) SetToNode(root *config.AclNode, path ...string) {
	if cp == nil {
		return
	}

	fp := append(path, cp.name)

	kp := append(fp, "kind")
	root.SetValAt(cp.Kind, kp...)

	kp = append(fp, "item")
	root.SetValAt(cp.Item, kp...)

	kp = append(fp, "degree")
	root.SetValAt(cp.Degree, kp...)
}

func (cp *EnumControlPoint) SetFromJSON(val interface{}) {
	if cp == nil || val == nil {
		return
	}

	v, ok := val.(map[string]interface{})
	if !ok {
		return
	}

	cp.Item = ValAsInt(v["Item"])
	cp.Degree = ValAsFloat(v["Degree"])
}

func (cp *EnumControlPoint) ToApi() *api.ControlPoint {
	val := &api.ControlPoint{
		Name: cp.name,
	}

	apiCP := &api.EnumPoint{
		Item:   int32(cp.Item),
		Degree: cp.Degree,
	}
	val.Val = &api.ControlPoint_Enm{apiCP}

	return val
}

func (cp *EnumControlPoint) SetFromApi(apiCP *api.ControlPoint) {
	if apiCP == nil || apiCP.Val == nil {
		return
	}

	apiVal := apiCP.GetEnm()
	if apiVal == nil {
		return
	}

	cp.Item = int(apiVal.Item)
	cp.Degree = apiVal.Degree
}

func (cp *EnumControlPoint) Apply(other ControlPoint) {
	if cp == nil || other == nil {
		return
	}

	otherEnum, ok := other.(*EnumControlPoint)
	if !ok {
		return
	}

	cp.Item = otherEnum.Item
	cp.Degree = otherEnum.Degree
}

func (cp *EnumControlPoint) String() string {
	if cp == nil {
		return "Enum(nil)"
	}

	return fmt.Sprintf("Enum(%v,%v)", cp.Item, cp.Degree)
}

func (cp *EnumControlPoint) Copy() ControlPoint {
	ncp := &EnumControlPoint{
		name: cp.name,

		Kind:   cp.Kind,
		Item:   cp.Item,
		Degree: cp.Degree,
	}

	return ncp
}

////////////////////////////////////////

type IntensityControlPoint struct {
	name string

	Kind      string
	Intensity float64
}

func NewIntensityControlPoint(name string) *IntensityControlPoint {
	cp := &IntensityControlPoint{
		name: name,

		Kind: "intensity",
	}

	return cp
}

func (cp *IntensityControlPoint) Name() string {
	if cp == nil {
		return ""
	}

	return cp.name
}

func (cp *IntensityControlPoint) Percent() float64 {
	if cp == nil {
		return 0.0
	}

	return cp.Intensity
}

func (cp *IntensityControlPoint) SetPercent(intensity float64) {
	if cp == nil {
		return
	}

	cp.Intensity = intensity
}

func (cp *IntensityControlPoint) SetFromNode(node *config.AclNode) {
	cp.Intensity = node.ChildAsFloat("intensity")
}

func (cp *IntensityControlPoint) SetToNode(root *config.AclNode, path ...string) {
	if cp == nil {
		return
	}

	fp := append(path, cp.name)

	kp := append(fp, "kind")
	root.SetValAt(cp.Kind, kp...)

	kp = append(fp, "intensity")
	root.SetValAt(cp.Intensity, kp...)
}

func (cp *IntensityControlPoint) SetFromJSON(val interface{}) {
	if cp == nil || val == nil {
		return
	}

	v, ok := val.(map[string]interface{})
	if !ok {
		return
	}

	cp.Intensity = ValAsFloat(v["Intensity"])
}

func (cp *IntensityControlPoint) ToApi() *api.ControlPoint {
	val := &api.ControlPoint{
		Name: cp.name,
	}

	apiCP := &api.IntensityPoint{
		Intensity: cp.Intensity,
	}
	val.Val = &api.ControlPoint_Intensity{apiCP}

	return val
}

func (cp *IntensityControlPoint) SetFromApi(apiCP *api.ControlPoint) {
	if apiCP == nil || apiCP.Val == nil {
		return
	}

	apiVal := apiCP.GetIntensity()
	if apiVal == nil {
		return
	}
	cp.Intensity = apiVal.Intensity
}

func (cp *IntensityControlPoint) Apply(other ControlPoint) {
	if cp == nil || other == nil {
		return
	}

	otherIntensity, ok := other.(*IntensityControlPoint)
	if !ok {
		return
	}

	cp.Intensity = otherIntensity.Intensity
}

func (cp *IntensityControlPoint) String() string {
	if cp == nil {
		return "Intensity(nil)"
	}

	return fmt.Sprintf("Intensity(%v)", cp.Intensity)
}

func (cp *IntensityControlPoint) Copy() ControlPoint {
	ncp := &IntensityControlPoint{
		name: cp.name,

		Kind:      cp.Kind,
		Intensity: cp.Intensity,
	}

	return ncp
}
