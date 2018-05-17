package eltee

type ControlPoint interface {
	// Control points shall be known uniquely by their name
	Name() string
    Observe() interface{}
}

type FixturePatch struct {
	patch map[FixtureControl]ControlPoint
}

type Switchboard struct {
}

type LensStack interface {
    Observe(cp ControlPoint) interface{}
}

type ColorValues interface {
    Component(name string) float64
}