package eltee

import (
	"github.com/eyethereal/go-config"
	"path"
)

//////////////////////////////////////////////////////////////////////
// Package level logging
// This only needs to be in one file per-package

var log = config.Logger("eltee")

///////////////////

// The Server is the main object which gets instantiated and then creates everything
// else. Generally there is only one, but as it is an object a larger server might
// embed this server along with other servers, or multiple servers for some crazy
// reason.
type Server struct {
	cfg *config.AclNode

	dmxHarness *DmxHarness
	library    *ProfileLibrary

	backingFrame []byte

	fixtures       []Fixture
	fixturesByName map[string]Fixture

	stateJuggler *StateJuggler
	// controlPoints       []ControlPoint
	// controlPointsByName map[string]ControlPoint

	inputAdapters []*InputAdapterRegistration

	apiServer *apiServer
}

func NewServer(cfg *config.AclNode) *Server {
	s := &Server{
		cfg:          cfg,
		library:      NewProfileLibrary(),
		backingFrame: make([]byte, 512),

		fixtures:       make([]Fixture, 0),
		fixturesByName: make(map[string]Fixture),
	}

	var err error

	// DMX Output connections
	s.dmxHarness = NewDmxHarness(s, cfg.Child("dmx"))

	// Load our library of profiles
	//err := s.library.LoadFile("donner_pinspot", "profiles/donner_pinspot.acl")
	err = s.library.LoadDirectory("profiles")
	if err != nil {
		log.Errorf("%v", err)
	}

	log.Debugf("***** Profile Library *******")
	log.Debugf("%v", s.library)
	log.Debugf("*****************************")

	defaultsNode := cfg.Child("defaults")
	if defaultsNode == nil {
		defaultsNode = config.NewAclNode()
	}

	// Load the fixtures instances. Fixtures are instances of profiles
	defFixturesFilename := defaultsNode.DefChildAsString("default", "fixtures")
	defFixturesFilename = path.Join("fixtures", defFixturesFilename) + ".acl"
	fixNode := config.NewAclNode()
	err = fixNode.ParseFile(defFixturesFilename)
	if err != nil {
		log.Warningf("Unable to load default fixtures file '%v': %v", defFixturesFilename, err)
	} else {
		fixNode = fixNode.Child("fixtures")
		if fixNode == nil {
			log.Warningf("File '%v' did not contain a fixtures node", defFixturesFilename)
		} else {
			base := 1
			fixNode.ForEachOrderedChild(func(name string, child *config.AclNode) {
				fixture, base := s.CreateFixture(name, child, base)
				if fixture != nil {
					s.fixtures = append(s.fixtures, fixture)
					s.fixturesByName[name] = fixture

					log.Infof("Added %v @ %v", fixture.Name(), base)
				}
			})
		}
	}

	// Load the initial set of control points
	s.stateJuggler = NewStateJuggler()

	defCPFilename := defaultsNode.DefChildAsString("default", "control_points")
	defCPFilename = path.Join("control_points", defCPFilename) + ".acl"
	cpNode := config.NewAclNode()
	err = cpNode.ParseFile(defCPFilename)
	if err != nil {
		log.Warningf("Unable to load default control points file '%v': %v", defCPFilename, err)
	} else {
		cpNode = cpNode.Child("control_points")
		if cpNode == nil {
			log.Warningf("File '%v' did not contain a control_points node", defCPFilename)
		} else {
			// s.controlPoints, s.controlPointsByName = CreateControlPointList(cpNode)
			s.stateJuggler.BaseFrom(cpNode)
		}
	}

	// Load an initial mapping between fixtures and control points
	defPatchFilename := defaultsNode.DefChildAsString("default", "patches")
	defPatchFilename = path.Join("patches", defPatchFilename) + ".acl"
	patchesNode := config.NewAclNode()
	err = patchesNode.ParseFile(defPatchFilename)
	if err != nil {
		log.Warningf("Unable to load default patches file '%v': %v", defPatchFilename, err)
	} else {
		patchesNode = patchesNode.Child("patches")
		if patchesNode == nil {
			log.Warningf("File '%v' did not contain a patches node", defPatchFilename)
		} else {
			// The order of names is fixture -> fixture_control -> control_point & lens_stack

			patchesNode.ForEachOrderedChild(func(fixName string, fixPatches *config.AclNode) {
				fixture := s.fixturesByName[fixName]
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
						cp := s.stateJuggler.CurrentCP(cpName)
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
	}

	s.apiServer = NewApiServer(s)

	return s
}

// Create a new fixture with the given node from the config file
func (s *Server) CreateFixture(name string, node *config.AclNode, defBase int) (f Fixture, nextBase int) {

	// First lets see if we can find a profile of the right kind
	kind := node.ChildAsString("kind")

	profile := s.library.Profiles[kind]
	if profile == nil {
		log.Errorf("Unknown fixture kind '%v'", kind)
		return nil, defBase
	}

	actualBase := node.DefChildAsInt(defBase, "base")
	channels := s.backingFrame[actualBase-1 : actualBase-1+profile.ChannelCount]

	fixture := NewDmxFixture(name, actualBase, channels, profile)

	return fixture, actualBase + profile.ChannelCount
}

// func (s *Server) BuildDefaultMappers() {
// 	for fIx := 0; fIx < len(s.fixtures); fIx++ {
// 		fixture := s.fixtures[fIx]

// 		inst := fixture.ProfileControlInstance()
// 		groupInst, ok := inst.(*ProfileControlGroupInstance)
// 		if !ok {
// 			continue
// 		}
// 		groupInst.ForEachControlInstance(func(inst ProfileControlInstance) {
// 			if inst == nil {
// 				log.Warningf("%v has nil PCI", fixture.Name())
// 			} else {
// 				pc := inst.ProfileControl()
// 				log.Infof("%v %v", fixture.Name(), pc)

// 				// Can it do color?
// 				// colorSettable, ok := inst.(WorldColorSettable)
// 				// if ok {
// 				// 	m := NewColorMapper("default")
// 				// 	s.defaultMappers = append(s.defaultMappers, m)
// 				// }
// 			}
// 		})
// 	}
// }

func (s *Server) Start() {
	go s.dmxHarness.Start()
	go s.apiServer.Start()
}

// A function to be called by the DmxHarness to update the frame slice that it
// contains with new values. Essentially this is a pull of new data into the
// DmxHarness at the time that it will then immediately be wanting to send the
// data out. If the data was pre-calculated then lovely. In a simple first
// implementation this is where we will actually do the updates to the frame
// values based on the current time.
//
// The job to be done in the function is the first part of what is described
// as the main looper
//
// 1. Polling animators for control point updates
// 2. Polling input adapters for control point changes
// 3. Commit all the control point changes (which collects the dirty ones)
// 4. Publishing changed control point values to subscribed input adaptors
// 5. Publishing changed control point values to observing fixture controls
// 6. Poking all Fixtures to have them generate their DMX output
//

func (s *Server) UpdateBackingFrame() {
	// 1. Poll animators
	// todo: poll animators

	// 2. Poll all input adapters
	for _, reg := range s.inputAdapters {
		reg.ia.UpdateControlPoints()
	}

	// // 3. Commit all the control point updates
	// dirtyNames := make([]string, 0)
	// for _, cp := range s.controlPoints {
	// 	if cp.Commit() {
	// 		dirtyNames = append(dirtyNames, cp.Name())
	// 	}
	// }

	// TODO: Get rid of dirtyNames???

	// 4. Tell Input adapters that control points might have changes
	for _, reg := range s.inputAdapters {
		reg.ia.ObserveControlPoints()
	}

	// 5 & 6. Tell fixtures to update check for updates and do new DMX
	for _, fix := range s.fixtures {
		fix.Update()
	}
}

// Updates the passed in frame from the backing frame. In the future we should
// maybe decouple the rates???
func (s *Server) UpdateFrame(frame []byte) {
	s.UpdateBackingFrame()

	copy(frame, s.backingFrame)
}

// func (s *Server) FrameState() ([]Fixture, *WorldState, []StateMapper) {
// 	return s.fixtures, s.currentWS, s.defaultMappers
// }

func (s *Server) RegisterInputAdapter(name string, adapter InputAdapter) {
	registration := &InputAdapterRegistration{
		ia:   adapter,
		name: name,
	}

	s.inputAdapters = append(s.inputAdapters, registration)

	adapterAsDmx, ok := adapter.(DMXConn)
	if ok {
		s.dmxHarness.AddConnection(name, adapterAsDmx)
	}
}

func (s *Server) DumpFixtures() {
	log.Info("--------- All fixtures....")
	for _, fix := range s.fixtures {
		log.Infof("Fixture [%v] † %v", fix.Name(), fix.Profile().Name)
		// Controls and mappings
		fix.ForEachFixtureControl(func(id string, fc *FixtureControl) {
			if fc.ControlPoint != nil {
				log.Debugf("  %v <- %v", id, fc.ControlPoint.Name())
			} else {
				// Don't bother showing the root group
				if id != "_root" {
					log.Debugf("  %v ... unattached", id)
				}
			}
		})
	}
	log.Info("--------- fixutres done")
}

func (s *Server) DumpControlPoints() {
	s.stateJuggler.DumpControlPoints()
	// log.Info("--------- All control points....")
	// for _, cp := range s.controlPoints {
	// 	log.Infof("CP [%v] %v", cp.Name(), cp)
	// }
	// log.Info("--------- control points done")
}

func (s *Server) GetFixtures() map[string]Fixture {
	// TODO: Full encapsulation so callers can't accidentally fuck up the fixtures slice???
	return s.fixturesByName
}

func (s *Server) GetProfiles() map[string]*Profile {
	// TODO: Full encapsulation so callers can't accidentally fuck up the  slice???
	return s.library.Profiles
}

func (s *Server) Juggler() *StateJuggler {
	return s.stateJuggler
}

// func (s *Server) GetControlPoints() map[string]ControlPoint {
// 	// TODO: Full encapsulation so callers can't accidentally fuck up the fixtures slice???
// 	return s.controlPointsByName
// }
