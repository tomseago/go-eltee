package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
)

//////////////////////////////////////////////////////////////////////
// Package level logging
// This only needs to be in one file per-package

var log = config.Logger("eltee")

type Server struct {
	cfg *config.AclNode

	dmxHarness *DmxHarness
	library    *ProfileLibrary

	fixtures       []Fixture
	fixturesByName map[string]Fixture

	currentWS *WorldState
}

func NewServer(cfg *config.AclNode) *Server {
	s := &Server{
		cfg:     cfg,
		library: NewProfileLibrary(),

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

	// Load the fixtures instances
	fixNode := cfg.Child("fixtures")
	if fixNode == nil {
		log.Warningf("No fixtures key found")
	} else {
		base := 1
		fixNode.ForEachOrderedChild(func(name string, child *config.AclNode) {
			fixture, base := s.CreateFixture(name, child, base, s.dmxHarness.frame)
			if fixture != nil {
				s.fixtures = append(s.fixtures, fixture)
				s.fixturesByName[name] = fixture

				log.Infof("Added %v @ %v", fixture.Name(), base)
			}
		})
	}

	// Create our basic world states
	// s.currentWS = NewWorldState()
	// s.currentWS.Root = cfg.Child("world_state").Duplicate()

	// s.nextWS = s.currentWS.Duplicate()

	// Iterate through all fixtures to setup a whole bunch of mappers to go from
	// the world state to outputable values
	// s.BuildDefaultMappers()

	return s
}

func (s *Server) CreateFixture(name string, node *config.AclNode, defBase int, dmx []byte) (f Fixture, nextBase int) {

	// First lets see if we can find a profile of the right kind
	kind := node.ChildAsString("kind")

	profile := s.library.Profiles[kind]
	if profile == nil {
		log.Errorf("Unknown fixture kind '%v'", kind)
		return nil, defBase
	}

	actualBase := node.DefChildAsInt(defBase, "base")
	channels := dmx[actualBase-1 : actualBase-1+profile.ChannelCount]

	fixture := NewDmxFixture(name, actualBase, channels, profile)

	return fixture, actualBase + profile.ChannelCount
}

func CreateConn(name string, cfg *config.AclNode) (DMXConn, error) {
	if cfg == nil {
		return nil, fmt.Errorf("AclNode was nil")
	}

	kind := cfg.ChildAsString("kind")
	switch kind {
	case "olad":
		return NewOLADConn(cfg)

	case "log":
		return NewLogConn(cfg)

	case "ftdi":
		return NewFtdiConn(cfg)
	}

	return nil, fmt.Errorf("Unknown dmx kind '%v'", kind)
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
	go s.dmxHarness.FramePump()

}

// func (s *Server) FrameState() ([]Fixture, *WorldState, []StateMapper) {
// 	return s.fixtures, s.currentWS, s.defaultMappers
// }

func (s *Server) DumpFixtures() {
	log.Info("All fixtures....")
	// for _, univ := range s.universes {
	// 	all := univ.AllFixtures()
	// 	for _, fix := range all {
	// 		log.Info("Fixture [%v] @ %v", fix.Name(), fix.Base())
	// 	}
	// }
	log.Info("...done")
}
