package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	// "github.com/nickysemenza/gola"
	"net/http"
	"strconv"
	"strings"
)

type DMXConn interface {
	SendDMX(universe int, data []byte)
}

type OLADConn struct {
	hostaddr string
	// client   *gola.OlaClient

	//connected bool

	url string
}

func NewOLADConn(cfg *config.AclNode) (DMXConn, error) {
	hostaddr := cfg.DefChildAsString("localhost:9090", "hostaddr")
	conn := &OLADConn{
		hostaddr: hostaddr,
	}

	conn.url = fmt.Sprintf("http://%v/set_dmx", hostaddr)
	return conn, nil
	// return conn, conn.AttemptOpen()
}

// func (c *OLADConn) AttemptOpen() error {
// 	if c == nil {
// 		return fmt.Errorf("Can not open a nil OLADConn")
// 	}

// 	log.Infof("Opening connection to %v", c.hostaddr)
// 	c.client = gola.New(c.hostaddr)
// 	c.connected = true
// 	return nil
// }

// func (c *OLADConn) SendDMX(universe int, data []byte) {
// 	if c == nil || c.connected == false {
// 		return
// 	}

// 	status, err := c.client.SendDmx(universe, data)

// 	_ = status

// 	if err != nil {
// 		log.Errorf("Error sending DMX. Considering this dead: %v", err)
// 		c.connected = false
// 		return
// 	}
// }

func (c *OLADConn) SendDMX(universe int, data []byte) {
	var b strings.Builder

	b.WriteString("u=")
	b.WriteString(strconv.Itoa(universe))
	b.WriteString("&d=")
	for i := 0; i < len(data); i++ {
		b.WriteString(strconv.Itoa(int(data[i])))

		if i < len(data)-1 {
			b.WriteString(",")
		}
	}

	dataStr := b.String()
	// fmt.Printf("data = %v\n", dataStr)
	reader := strings.NewReader(dataStr)

	// log.Infof("url = %v", c.url)
	resp, err := http.Post(c.url, "application/x-www-form-urlencoded", reader)
	_ = resp
	_ = err
	// fmt.Printf("resp=%v\nerr=%v", resp, err)
}

///////////

type LogConn struct {
}

func NewLogConn(cfg *config.AclNode) (DMXConn, error) {
	conn := &LogConn{}

	return conn, nil
}

func (c *LogConn) SendDMX(Universe int, data []byte) {
	if c == nil {
		return
	}

	toLog := len(data)
	if toLog > 30 {
		toLog = 30
	}

	var b strings.Builder
	b.WriteString("SendDMX [")
	for i := 0; i < toLog; i++ {
		b.WriteString(fmt.Sprintf("%d ", data[i]))
	}
	b.WriteString("]")

	log.Debugf(b.String())
}

///////////

type FtdiConn struct {
	ctx *FtdiContext
}

func NewFtdiConn(cfg *config.AclNode) (DMXConn, error) {
	conn := &FtdiConn{
		ctx: NewFtdiContext(),
	}

	// No reason to start it straight away
	err := conn.ctx.Start()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *FtdiConn) SendDMX(Universe int, data []byte) {
	if c == nil {
		return
	}

	// TODO: Honor universe mappings here???

	c.ctx.WriteDmx(data)
}
