package eltee

import (
	"github.com/eyethereal/go-config"
)

type StateJuggler struct {
	// This one contains the default values and is loaded from the
	// "control_points" directory
	base *WorldState

	current *WorldState

	statesByName map[string]*WorldState
}

func NewStateJuggler() *StateJuggler {
	sj := &StateJuggler{
		statesByName: make(map[string]*WorldState),
	}

	return sj
}

func (sj *StateJuggler) BaseFrom(root *config.AclNode) {
	sj.base = NewWorldStateFromNode("BASE", root)

	sj.current = sj.base.Copy()
}

func (sj *StateJuggler) CurrentCP(name string) ControlPoint {
	if sj.current == nil {
		return nil
	}

	return sj.current.ControlPoint(name)
}

func (sj *StateJuggler) DumpControlPoints() {
	log.Info("--------- All control points....")
	for _, cp := range sj.current.controlPoints {
		log.Infof("CP [%v] %v", cp.Name(), cp)
	}
	log.Info("--------- control points done")
}

func (sj *StateJuggler) Current() *WorldState {
	return sj.current
}

func (sj *StateJuggler) State(name string) *WorldState {
	return sj.statesByName[name]
}
