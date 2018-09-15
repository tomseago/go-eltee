package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"time"
)

type DmxHarness struct {
	server *Server
	conns  map[string]DMXConn

	dmxTester *DmxTester

	frame []byte
}

func NewDmxHarness(server *Server, node *config.AclNode) *DmxHarness {
	h := &DmxHarness{
		server: server,
		conns:  make(map[string]DMXConn),

		frame: make([]byte, 512),
	}

	cConns := node.Child("connections")
	cConns.ForEachOrderedChild(func(name string, n *config.AclNode) {
		conn, err := CreateConn(name, n)
		if err == nil {
			log.Infof("Added connection %s", name)
			h.conns[name] = conn
		} else {
			log.Errorf("Error adding connection %s : %v", name, err)
		}
	})

	h.dmxTester = NewDmxTester(node.Child("tester"))

	return h
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

// Called periodically to fetch a new frame from either the tester or the main server
// and then send that out to all DMX connections
func (h *DmxHarness) SendFrame() {

	if h.dmxTester.HasTest() {
		h.dmxTester.UpdateFrame(h.frame)
	} else {
		// use the mappers to update the frame
		// fixtures, state, mappers := h.server.FrameState()
		h.server.UpdateFrame()
	}

	// Send this frame to all of our connections
	for _, c := range h.conns {
		c.SendDMX(1, h.frame)
	}
}

// The main pump function which should almost certainly be run from a go routine. It will
// repeatedly call SendFrame on the DmxHarness forever at some pre-defined frame interval.
func (h *DmxHarness) FramePump() {
	log.Infof("DmxHarness pump routine started")

	desiredDelay := 500 * time.Millisecond

	frameTimer := time.NewTimer(desiredDelay)

	for {
		<-frameTimer.C
		log.Info("tick")

		h.SendFrame()

		frameTimer.Reset(desiredDelay)
	}

}
