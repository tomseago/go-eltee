package main

import (
	"github.com/eyethereal/go-config"

	"github.com/tomseago/go-eltee"
)

//////////////////////////////////////////////////////////////////////
// Package level logging
// This only needs to be in one file per-package

var log = config.Logger("eltee")

//////////

func main() {
	// Setup basic logging
	config.ColoredLoggingToConsole()

	x := eltee.FtdiLibVersion()
	log.Infof("X is %v", x)

	fc := eltee.NewFtdiContext()
	fc.FindAll()
}
