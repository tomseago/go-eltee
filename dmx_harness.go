package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"time"
)

///////////////////
// The DmxHarness holds a current DMX frame and a set of connectors that will send that
// Dmx frame somewhere. It pulls new frames from the server at a rate that it decides
// is appropriate. This isn't the same rate that the server is generating frames. It can
// do that however often it wishes (and should double buffer them).
type DmxHarness struct {
	server *Server
	conns  map[string]DMXConn

	dmxTester *DmxTester

	frame      []byte
	frameCount int

	// milliseconds between frames
	frameDelay time.Duration
}

func NewDmxHarness(server *Server, node *config.AclNode) *DmxHarness {
	h := &DmxHarness{
		server: server,
		conns:  make(map[string]DMXConn),

		frame: make([]byte, 512),
	}

	fps := time.Duration(node.DefChildAsInt(30, "frames_per_second"))
	h.frameDelay = time.Second / fps

	cConns := node.Child("connections")
	cConns.ForEachOrderedChild(func(name string, n *config.AclNode) {
		conn, err := CreateConn(name, n)
		if err == nil {
			h.AddConnection(name, conn)
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

func (h *DmxHarness) AddConnection(name string, c DMXConn) {
	log.Infof("Added connection %s", name)
	h.conns[name] = c
}

// Called periodically to fetch a new frame from either the tester or the main server
// and then send that out to all DMX connections
func (h *DmxHarness) SendFrame() {

	if h.dmxTester.HasTest() {
		h.dmxTester.UpdateFrame(h.frame)
	} else {
		h.server.UpdateFrame(h.frame)
	}

	// DEBUGGING
	h.frameCount++
	h.frame[0] = byte(h.frameCount)

	// Send this frame to all of our connections
	for _, c := range h.conns {
		c.SendDMX(1, h.frame)
	}
}

// The main pump function which should almost certainly be run from a go routine. It will
// repeatedly call SendFrame on the DmxHarness forever at some pre-defined frame interval.
func (h *DmxHarness) Start() {
	log.Infof("DmxHarness pump routine started")

	// desiredDelay := h.frameDelay * time.Millisecond

	frameTimer := time.NewTimer(h.frameDelay)

	for {
		<-frameTimer.C
		//log.Info("tick")

		h.SendFrame()

		frameTimer.Reset(h.frameDelay)
	}

}
