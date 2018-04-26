package eltee

type Universe interface {
	SetVal(index uint16, val byte)
	Val(index uint16) byte

	Commit() (num int, data *[512]byte)

	AddFixture(fixture Fixture)
	AllFixtures() []Fixture
}

type UniverseData struct {
	num   int
	data  [512]byte
	dirty bool

	fixtures []Fixture
}

func NewUniverse(num int) Universe {
	r := &UniverseData{
		num:      num,
		fixtures: make([]Fixture, 0),
	}
	return r
}

func (u *UniverseData) SetVal(index uint16, val byte) {
	u.data[index] = val
	u.dirty = true
}

func (u *UniverseData) Val(index uint16) byte {
	return u.data[index]
}

func (u *UniverseData) Commit() (num int, data *[512]byte) {
	u.dirty = false
	return u.num, &u.data
}

func (u *UniverseData) AddFixture(fixture Fixture) {
	if fixture == nil {
		return
	}

	u.fixtures = append(u.fixtures, fixture)
}

func (u *UniverseData) AllFixtures() []Fixture {
	if u == nil {
		return make([]Fixture, 0)
	}

	out := make([]Fixture, 0, len(u.fixtures))
	copy(out, u.fixtures)
	return out
}
