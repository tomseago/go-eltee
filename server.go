package eltee


import (
    "github.com/eyethereal/go-config"
    "github.com/nickysemenza/gola"
    "strconv"
)

//////////////////////////////////////////////////////////////////////
// Package level logging
// This only needs to be in one file per-package

var log = config.Logger("eltee")



type Server struct {
    cfg *config.AclNode
    client *gola.OlaClient
    universes []Universe
}


func NewServer(cfg *config.AclNode) (*Server) {
    s := &Server {
        cfg: cfg,
    }

    // Get our concept of the world set up
    s.universes = make([]Universe, 0)
    cUnis := cfg.Child("universes")
    if cUnis != nil {
        for name, cUni := range cUnis.Children {
            num, _ := strconv.Atoi(name)
            uni := NewUniverse(num)
            s.universes = append(s.universes, uni)

            cFixtures := cUni.Child("fixtures")
            if cFixtures != nil {
                for name, cFixture := range cFixtures.Children {
                    fix := CreateFixture(name, cFixture)
                    uni.addFixture(fix)
                }
            }
        }
    }

    // Open the network stuff
    hostaddr := cfg.DefChildAsString("localhost:9010", "hostaddr")
    log.Infof("Opening connection to %v", hostaddr)
    s.client = gola.New(hostaddr)

    return s
}

func CreateFixture(name string, cfg *config.AclNode) *Fixture {
    if cfg == nil {
        return nil
    }

    kind := cfg.AsString()

    _ = kind
    
    return nil
}

func (s *Server) Start() {

}