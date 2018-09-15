package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"math"
)

/////////////////////////////////////////////////////////////////////////////////
//
// Control Kind: pan and tilt
//

type PanTiltUpdater struct {
	// This is a reference to the slice that is held by the DmxFixture. Thus changes
	// made here go directly into the same place that the fixture references without
	// a need to copy them.
	channels []byte
}

func (u *PanTiltUpdater) Update(fc *FixtureControl) {
	if u == nil || fc == nil || fc.ControlPoint == nil {
		return
	}

	value := fc.LensStack.Observe(fc.ControlPoint)

	xyz, ok := value.(XYZPoint)
	if !ok {
		// not an intensity
		return
	}

	pc, ok := fc.ProfileControl.(*PanTiltProfileControl)
	_ = pc
	if !ok {
		// attached to the wrong type of control. Bad
		return
	}

	x, y, z := xyz.XYZ()

	// Assume we are at 0,0,0. Find the appropriate angles to the
	// control points. If we are upside down or at a different position
	// then that's all been handled by the lens stack. It is assumed that
	// the combination of the configuration and the settings on the fixture allow
	// us to reason in a standard right hand coordinate system.

	// The XY plane is the floor, X increases to the right, Y increases to the front.
	// Z is up.
	// Pan is rotation around Z where 0 is to the right, pi/2 is forward
	// Tilt is rotation around an axis that is perpendicular to the radial line
	// (i.e. it changes as pan changes) which can also be thought of as dropping down
	// from the vertical. 0 tilt is parallel to the XY plane, pi/2 is straight up,
	// 1.5 pi is straight down.

	// theta is the angle from Z down to the radial (for tilt)
	theta := math.Atan2(math.Sqrt(x*x+y*y), z)

	// Both reverse the direction and move it by pi/2
	tiltRad := math.Pi/2.0 - theta

	// phi is the angle from X around to the radial (for pan)
	phi := math.Atan2(y, x)

	// There is no change to phi to turn it into desired standard directions. While
	// most fixtures may have 0 dmx as forward (or backward), in our right hand system
	// 0 pan is to the right - which matches the definition of phi.
	panRad := phi

	// However, the Atan2 function returns values in the range of -pi to +pi so if
	// these values are < 0 we need to add 2pi to them
	if panRad < 0 {
		panRad += 2 * math.Pi
	}

	// The conversion from radians to degrees and respecting the limits of the
	// fixture is handled by each Axis
	//log.Infof("tilt=%v pan=%v", tiltRad, panRad)
	//log.Info("Tilt")
	pc.tilt.SetRadians(u.channels, tiltRad)

	//log.Info("Pan")
	pc.pan.SetRadians(u.channels, panRad)

	// For now, always set max movement speed. Basically assume that we are producing
	// sufficient frames for smooth movement animations.
	if pc.chSpeed > 0 {
		u.channels[pc.chSpeed] = 255
	}
}

////////////////

const DEGREE_TO_RAD float64 = math.Pi / 180.0

type Axis struct {
	coarse int
	fine   int

	minRad float64
	maxRad float64
}

func NewAxis(node *config.AclNode) *Axis {
	a := &Axis{}

	log.Debugf("In NewAxis node is %v", node.ColoredString())
	if node != nil {
		a.coarse = node.ChildAsInt("coarse")
		a.fine = node.ChildAsInt("fine")

		minDegree := node.ChildAsFloat("min")
		a.minRad = minDegree * DEGREE_TO_RAD
		//log.Infof("minDegree=%v a.minRad=%v", minDegree, a.minRad)

		maxDegree := node.ChildAsFloat("max")
		a.maxRad = maxDegree * DEGREE_TO_RAD
		//log.Infof("maxDegree=%v a.maxRad=%v", maxDegree, a.maxRad)
	}

	return a
}

func (a *Axis) SetRadians(channels []byte, rad float64) {
	//log.Debugf("SetRadians(..., %v) for %v", rad, a)
	if rad < a.minRad {
		rad = a.minRad
	}

	if rad > a.maxRad {
		rad = a.maxRad
	}

	//log.Debugf("rad = %v", rad)
	p := (rad - a.minRad) / (a.maxRad - a.minRad)

	// Now convert p to a 16-bit whole number
	value := int16(math.MaxUint16 * p)
	hi := byte(value >> 8)
	lo := byte(value & 0x00FF)

	//log.Debugf("p=%v, value=%x, hi=%x lo=%x", p, value, hi, lo)

	if a.coarse > 0 {
		channels[a.coarse-1] = hi
	}
	if a.fine > 0 {
		channels[a.fine-1] = lo
	}
}

func (a *Axis) String() string {
	return fmt.Sprintf("[%v, %v](%v->%v)", a.coarse, a.fine, a.minRad, a.maxRad)
}

////////////////

type PanTiltProfileControl struct {
	ProfileControlBase

	pan  *Axis
	tilt *Axis

	chSpeed int
}

func NewPanTiltProfileControl(id string, rootNode *config.AclNode) (*PanTiltProfileControl, error) {
	pc := &PanTiltProfileControl{
		ProfileControlBase: MakeProfileControlBase(id, rootNode),

		pan:  NewAxis(rootNode.Child("pan")),
		tilt: NewAxis(rootNode.Child("tilt")),

		chSpeed: rootNode.ChildAsInt("speed"),
	}

	return pc, nil
}

func (pc *PanTiltProfileControl) Id() string {
	return pc.id
}

func (pc *PanTiltProfileControl) Name() string {
	return pc.name
}

func (pc *PanTiltProfileControl) String() string {
	return fmt.Sprintf("PanTilt %v(%v) Pan %v Tilt %v", pc.name, pc.id, pc.pan, pc.tilt)
}

func (pc *PanTiltProfileControl) Instantiate(fixture Fixture) *FixtureControl {
	dmx, ok := fixture.(*DmxFixture)
	if !ok {
		return nil
	}

	updater := &PanTiltUpdater{
		channels: dmx.channels,
	}

	fc := NewFixtureControl(pc, updater)

	fixture.AttachControl(pc.id, fc)

	return fc
}
