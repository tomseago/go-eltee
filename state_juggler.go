package eltee

import (
	"errors"
	"github.com/eyethereal/go-config"
	"io/ioutil"
	"path/filepath"
)

type StateJuggler struct {
	// This one contains the default values and is loaded from the
	// "control_points" directory
	base *WorldState

	current *WorldState

	statesByName   map[string]*WorldState
	fixturesByName map[string]Fixture
}

func NewStateJuggler(fixturesByName map[string]Fixture) *StateJuggler {
	sj := &StateJuggler{
		statesByName:   make(map[string]*WorldState),
		fixturesByName: fixturesByName,
	}

	return sj
}

func (sj *StateJuggler) StateNames() []string {
	out := make([]string, 0, len(sj.statesByName))
	for key := range sj.statesByName {
		out = append(out, key)
	}
	return out
}

// func (sj *StateJuggler) BaseFrom(root *config.AclNode) {
// 	sj.base = NewWorldStateFromNode("BASE", root)

// 	sj.current = sj.base.Copy()
// }

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
	if len(name) == 0 {
		return sj.current
	}

	return sj.statesByName[name]
}

/**
This loads all states from the given directoy. After loading the state named
"base" is copied as the current state. The fixtures are not yet mapped to
the control points in this state though, so that has to be done afterwards.
*/
func (sj *StateJuggler) LoadDirectory(dirname string) error {
	if sj == nil {
		return errors.New("nil StateJuggler")
	}

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
		if ext != ".acl" {
			continue
		}

		base := filepath.Base(file.Name())
		base = base[:len(base)-4]

		full := filepath.Join(dirname, file.Name())
		err = sj.LoadFile(base, full)
		if err != nil {
			log.Errorf("%v", err)
		}
	}

	if sj.statesByName["base"] == nil {
		return errors.New("No base state was defined")
	}

	sj.current = sj.statesByName["base"].Copy()
	sj.current.name = "CURRENT"
	sj.patchFixtures(sj.current)

	return nil
}

func (sj *StateJuggler) LoadFile(stateName string, file string) error {
	stateNode := config.NewAclNode()
	err := stateNode.ParseFile(file)
	if err != nil {
		log.Warningf("Unable to load state file '%v': %v", file, err)
		return err
	}

	state := NewWorldStateFromNode(stateName, stateNode)

	sj.statesByName[stateName] = state

	return nil
}

func (sj *StateJuggler) patchFixtures(fromState *WorldState) {
	if sj == nil || fromState == nil {
		return
	}

	fromState.patchesNode.ForEachOrderedChild(func(fixName string, fixPatches *config.AclNode) {
		fixture := sj.fixturesByName[fixName]
		if fixture == nil {
			log.Warningf("Patching: Unable to find fixture named '%v'", fixName)
			return
		}

		fixPatches.ForEachOrderedChild(func(fcId string, fcNode *config.AclNode) {
			fixtureControl := fixture.Control(fcId)
			if fixtureControl == nil {
				log.Warningf("Patching: Fixture '%v' does not have a control with id '%v'", fixName, fcId)
				return
			}

			// Get the control point, if any
			cpName := fcNode.ChildAsString("cp")
			if len(cpName) > 0 {
				cp := sj.current.ControlPoint(cpName)
				if cp == nil {
					log.Warningf("Patching: Could not find control point '%v' to patch to fixture '%v' control '%v'", cpName, fixName, fcId)
				} else {
					fixtureControl.ControlPoint = cp
				}
			}

			// TODO: Add the lens stack
		})
	})
}
