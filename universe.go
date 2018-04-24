package eltee

type Universe interface {
    setVal(index uint16, val byte)
    val(index uint16) byte

    commit() (num int, data *[512]byte)

    addFixture(fixture Fixture)
}

type UniverseData struct {
    num int
    data [512]byte
    dirty bool

    fixtures []Fixture
}

func NewUniverse(num int) Universe {
    r := &UniverseData{
        num: num,
        fixtures: make([]Fixture, 0),
    }
    return r
}

func (u *UniverseData) setVal(index uint16, val byte) {
    u.data[index] = val
    u.dirty = true
}

func (u *UniverseData) val(index uint16) byte {
    return u.data[index]
}

func (u *UniverseData) commit() (num int, data *[512]byte) {
    u.dirty = false
    return u.num, &u.data
}

func (u *UniverseData) addFixture(fixture Fixture) {
    if fixture == nil {
        return
    }

    u.fixtures = append(u.fixtures, fixture)
}
