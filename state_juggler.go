package eltee

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/eyethereal/go-config"
	"io/ioutil"
	"os"
	"path/filepath"
)

// The StateJuggler does exactly as its name implies - it juggles a list of
// current and potential WorldStates. A WorldState has two parts, a list of
// named control points and a set of FixturePatches which are connections between
// FixtureControls and ControlPoints.
//
// A base WorldState defines defaults which hold until any other state is
// applied. A current WorldState contains the singular active state for the
// system. It's control points can be modified directly or another, typically
// partial, WorldState can be "applied" onto the current state. Applying one
// WorldState onto another overwrites shared values but doesn't delete anything
// (although that should maybe get added). Applying one state onto another
// may also cause FixturePatches to change.
type StateJuggler struct {
	// This one contains the default values and is loaded from the
	// "control_points" directory
	base *WorldState

	current *WorldState

	statesByName   map[string]*WorldState
	fixturesByName map[string]Fixture

	lastStateApplied string
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
	if len(name) == 0 || name == "CURRENT" {
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

func (sj *StateJuggler) SaveFile(stateName string, file string) error {
	if sj == nil {
		return errors.New("nil StateJuggler")
	}

	state := sj.statesByName[stateName]
	if state == nil {
		return fmt.Errorf("No state named '%v'", stateName)
	}

	stateNode := config.NewAclNode()
	state.SetToNode(stateNode)

	fp, err := os.Create(file)
	if err != nil {
		return err
	}

	bio := bufio.NewWriter(fp)
	stateNode.StringTo(bio, "  ", 0, false)
	bio.Flush()
	fp.Close()

	return nil
}

func (sj *StateJuggler) patchFixtures(fromState *WorldState) {
	if sj == nil || fromState == nil {
		return
	}

	for _, fp := range fromState.fixturePatches {
		fixture := sj.fixturesByName[fp.FixtureName]
		if fixture == nil {
			log.Warningf("Patching: Unable to find fixture named '%v'", fp.FixtureName)
			return
		}

		for cName, cpName := range fp.CpsByControl {
			fixtureControl := fixture.Control(cName)
			if fixtureControl == nil {
				log.Warningf("Patching: Fixture '%v' does not have a control with id '%v'", fp.FixtureName, cName)
				return
			}

			if len(cpName) == 0 || cpName == "_" {
				// Unpatch it
				fixtureControl.ControlPoint = nil
			} else {
				cp := sj.current.ControlPoint(cpName)
				if cp == nil {
					log.Warningf("Patching: Could not find control point '%v' to patch to fixture '%v' control '%v'", cpName, fp.FixtureName, cName)
				} else {
					fixtureControl.ControlPoint = cp
				}
			}
		}

		// TODO: Lens stacks...
	}

	// fromState.patchesNode.ForEachOrderedChild(func(fixName string, fixPatches *config.AclNode) {
	// 	fixture := sj.fixturesByName[fixName]
	// 	if fixture == nil {
	// 		log.Warningf("Patching: Unable to find fixture named '%v'", fixName)
	// 		return
	// 	}

	// 	fixPatches.ForEachOrderedChild(func(fcId string, fcNode *config.AclNode) {
	// 		fixtureControl := fixture.Control(fcId)
	// 		if fixtureControl == nil {
	// 			log.Warningf("Patching: Fixture '%v' does not have a control with id '%v'", fixName, fcId)
	// 			return
	// 		}

	// 		// Get the control point, if any
	// 		cpName := fcNode.ChildAsString("cp")
	// 		if len(cpName) > 0 {
	// 			cp := sj.current.ControlPoint(cpName)
	// 			if cp == nil {
	// 				log.Warningf("Patching: Could not find control point '%v' to patch to fixture '%v' control '%v'", cpName, fixName, fcId)
	// 			} else {
	// 				fixtureControl.ControlPoint = cp
	// 			}
	// 		}

	// 		// TODO: Add the lens stack
	// 	})
	// })
}

func (sj *StateJuggler) LoadableStateNames(dirname string) ([]string, error) {
	if sj == nil {
		return nil, errors.New("nil StateJuggler")
	}

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	out := make([]string, 0)
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

		out = append(out, base)
	}

	return out, nil
}

func (sj *StateJuggler) LoadLoadableState(dirname string, stateName string) error {
	if sj == nil {
		return errors.New("nil StateJuggler")
	}

	filename := filepath.Join(dirname, stateName+".acl")
	return sj.LoadFile(stateName, filename)
}

func (sj *StateJuggler) SaveLoadableState(dirname string, stateName string) error {
	if sj == nil {
		return errors.New("nil StateJuggler")
	}

	filename := filepath.Join(dirname, stateName+".acl")
	return sj.SaveFile(stateName, filename)
}

func (sj *StateJuggler) SaveAll(dirname string) {
	if sj == nil {
		return
	}

	for name, _ := range sj.statesByName {
		if name == "CURRENT" || name == "base" {
			log.Infof("Not saving state named %v", name)
			continue
		}

		e := sj.SaveLoadableState(dirname, name)
		if e != nil {
			log.Warningf("While saving state %v: %v", name, e)
		} else {
			log.Debugf("Saved state %v", name)
		}
	}
}

func (sj *StateJuggler) ApplyState(stateName string) error {
	if sj == nil {
		return errors.New("nil StateJuggler")
	}

	if stateName == "CURRENT" {
		// Silly, but whatever
		return nil
	}

	state := sj.statesByName[stateName]
	if state == nil {
		return fmt.Errorf("Could not find state named '%v'", stateName)
	}

	sj.current.Apply(state)

	// Always repatch the fixture object references
	sj.patchFixtures(sj.current)

	sj.lastStateApplied = stateName

	return nil
}

func (sj *StateJuggler) AddState(stateName string) error {
	if sj == nil {
		return errors.New("nil StateJuggler")
	}

	state := sj.State(stateName)
	if state != nil {
		return fmt.Errorf("State '%v' already exists", stateName)
	}

	state = NewWorldState(stateName)
	sj.statesByName[stateName] = state

	return nil
}

func (sj *StateJuggler) RemoveState(stateName string) error {
	if sj == nil {
		return errors.New("nil StateJuggler")
	}

	if len(stateName) == 0 {
		return errors.New("Must specify a state name")
	}

	if stateName == "CURRENT" || stateName == "base" {
		return fmt.Errorf("Can not remove state %v", stateName)
	}

	delete(sj.statesByName, stateName)

	return nil
}

func (sj *StateJuggler) CopyStateTo(src string, dest string) error {
	if sj == nil {
		return errors.New("nil StateJuggler")
	}

	sState := sj.State(src)
	if sState == nil {
		return fmt.Errorf("Could not find source state '%v'", src)
	}
	if len(dest) == 0 || dest == "CURRENT" || dest == "base" || dest == src {
		return fmt.Errorf("Can not copy to state '%v'", dest)
	}

	// In case it exists, simply remove it
	sj.RemoveState(dest)

	sj.statesByName[dest] = sState.Copy()

	return nil
}

func (sj *StateJuggler) MoveStateTo(src string, dest string) error {
	// A move is just a copy and a remove of the source
	err := sj.CopyStateTo(src, dest)
	if err != nil {
		return err
	}

	// If it's an invalid name, that's okay
	sj.RemoveState(src)

	return nil
}

func (sj *StateJuggler) ApplyStateTo(src string, dest string) error {
	if sj == nil {
		return errors.New("nil StateJuggler")
	}

	if dest == src {
		// Nothing to do
		return nil
	}

	sState := sj.State(src)
	if sState == nil {
		return fmt.Errorf("Could not find source state '%v'", src)
	}
	dState := sj.State(dest)
	if dState == nil {
		return fmt.Errorf("Could not find destination state '%v'", dest)
	}

	dState.Apply(sState)

	return nil
}
