package main

import (
	"github.com/eyethereal/go-config"
	"os"
	"os/signal"

	"github.com/tomseago/go-eltee"
	"github.com/tomseago/go-eltee/web"
)

//////////////////////////////////////////////////////////////////////
// Package level logging
// This only needs to be in one file per-package

var log = config.Logger("eltee")

//////////////////////////////////////////////////////////////////////

func main() {

	// Load the configuration
	info, err := os.Stat("config")
	if err == nil && info.IsDir() {
		os.Chdir("config")
	}
	cfg := config.LoadACLConfig("eltee", "ELTEE")

	log.Warningf("***********************************************************************************")
	log.Warningf("ElTee build %v - Starting Up", cfg.PrettyVersion())
	log.Warningf("***********************************************************************************")

	server := eltee.NewServer(cfg)

	server.DumpFixtures()
	server.DumpControlPoints()

	// Start up the web server
	webServer := web.NewWebServer(cfg.Child("web"), server)
	server.RegisterInputAdapter("web", webServer)

	///////
	server.Start()

	// Wait for sigKill???
	c := make(chan os.Signal, 3)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until we get something
	s := <-c
	log.Warningf("Received signal %v. Exiting", s)
}
