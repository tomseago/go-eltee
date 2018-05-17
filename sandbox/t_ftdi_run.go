package main

import (
	"github.com/eyethereal/go-config"
	"os"
	"time"

	"github.com/tomseago/go-eltee"
)

//////////////////////////////////////////////////////////////////////
// Package level logging
// This only needs to be in one file per-package

var log = config.Logger("eltee")

///

func writeCOBPar(data []byte, offset int) {
	data[offset] = 0
	data[offset+1] = 0
	data[offset+2] = 0

	data[offset+3] = 255
	data[offset+4] = 255
	data[offset+5] = 12
	data[offset+6] = 0
}

func writeCOBParX(data []byte, offset int) {
	data[offset] = 0
	data[offset+1] = 0
	data[offset+2] = 0

	data[offset+3] = 255
	data[offset+4] = 255
	data[offset+5] = 255
	data[offset+6] = 255
}

func main() {
	// Setup basic logging
	config.ColoredLoggingToConsole()

	ctx := eltee.NewFtdiContext()

	err := ctx.Start()
	if err != nil {
		log.Fatalf("Unable to start FTDI: %v", err)
	}

	// Now just sleep for awhile
	log.Warningf("Sleeping on main routine for a small bit")

	time.Sleep(time.Second * 1)

	log.Warningf("Writing some DMX data, then sleep for 10")

	data := make([]byte, 512)
	// COB Par 1, base 97 (dmx) 3 bytes macro, dimmer, R G B
	//writeCOBPar(data, 0)
	writeCOBPar(data, 96)

	log.Warningf("data %v", data)
	ctx.WriteDmx(data)

	time.Sleep(time.Second * 120)

	log.Warningf("Done sleeping, exiting...")
	os.Exit(0)
}
