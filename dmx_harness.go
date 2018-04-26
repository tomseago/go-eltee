package eltee

import (
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

func (h *DmxHarness) SendFrame() {

	if h.dmxTester.HasTest() {
		h.dmxTester.UpdateFrame(h.frame)
	} else {
		// use the mappers to update the frame
		// fixtures, state, mappers := h.server.FrameState()

	}

	// Send this frame to all of our connections
	for _, c := range h.conns {
		c.SendDMX(1, h.frame)
	}
}

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
