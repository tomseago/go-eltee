package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"sort"
	"strings"
)

type ControlPoint interface {
	// Control points shall be known uniquely by their name
	Name() string
	// ObserveAsColor() ColorValues
	// ObserveAsXYZ() XYZValues

	SetFromNode(node *config.AclNode)
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
	Intensity() float64
}
type SettableIntensityPoint interface {
	SetIntensity(intensity float64)
}

////////////////////////////////////////

type ColorControlPoint struct {
	name string

	components map[string]float64
}

func NewColorControlPoint(name string) *ColorControlPoint {
	cp := &ColorControlPoint{
		name: name,

		components: make(map[string]float64),
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

	return cp.components[name]
}

func (cp *ColorControlPoint) SetColorComponent(name string, val float64) {
	if cp == nil {
		return
	}

	cp.components[name] = val
}

func (cp *ColorControlPoint) SetFromNode(node *config.AclNode) {
	node.ForEachOrderedChild(func(nn string, val *config.AclNode) {
		// protect reserved names, but otherwise accept everything
		// else as naming a color component
		if nn == "kind" {
			return
		}

		cp.components[nn] = val.AsFloat()
	})
}

func (cp *ColorControlPoint) String() string {
	if cp == nil {
		return "ColorCP(nil)"
	}

	ids := make([]string, 0, len(cp.components))
	for id, _ := range cp.components {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var b strings.Builder
	b.WriteString("ColorCP(")

	for ix := 0; ix < len(ids); ix++ {
		id := ids[ix]
		val := cp.components[id]
		b.WriteString(fmt.Sprintf("%v:%v", id, val))
		if ix < len(ids)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteByte(')')

	return b.String()
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

	x float64
	y float64
	z float64
}

func NewXYZControlPoint(name string) *XYZControlPoint {
	cp := &XYZControlPoint{
		name: name,
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

	return cp.x, cp.y, cp.z
}

func (cp *XYZControlPoint) SetXYZ(x float64, y float64, z float64) {
	if cp == nil {
		return
	}

	cp.x = x
	cp.y = y
	cp.z = z
}

func (cp *XYZControlPoint) SetFromNode(node *config.AclNode) {
	cp.x = node.ChildAsFloat("x")
	cp.y = node.ChildAsFloat("y")
	cp.z = node.ChildAsFloat("z")
}

////////////////////////////////////////

type EnumControlPoint struct {
	name string

	item   int
	degree float64
}

func NewEnumControlPoint(name string) *EnumControlPoint {
	cp := &EnumControlPoint{
		name: name,
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

	return cp.item, cp.degree
}

func (cp *EnumControlPoint) SetOption(item int, degree float64) {
	if cp == nil {
		return
	}

	cp.item = item
	cp.degree = degree
}

func (cp *EnumControlPoint) SetFromNode(node *config.AclNode) {
	cp.item = node.ChildAsInt("item")
	cp.degree = node.ChildAsFloat("degree")
}

////////////////////////////////////////

type IntensityControlPoint struct {
	name string

	intensity float64
}

func NewIntensityControlPoint(name string) *IntensityControlPoint {
	cp := &IntensityControlPoint{
		name: name,
	}

	return cp
}

func (cp *IntensityControlPoint) Name() string {
	if cp == nil {
		return ""
	}

	return cp.name
}

func (cp *IntensityControlPoint) Intensity() float64 {
	if cp == nil {
		return 0.0
	}

	return cp.intensity
}

func (cp *IntensityControlPoint) SetIntensity(intensity float64) {
	if cp == nil {
		return
	}

	cp.intensity = intensity
}

func (cp *IntensityControlPoint) SetFromNode(node *config.AclNode) {
	cp.intensity = node.ChildAsFloat("intensity")
}
