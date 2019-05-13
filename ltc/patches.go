package main

import (
	"context"
	// "errors"
	"fmt"
	"github.com/tomseago/go-eltee/api"
	// "strconv"
	// "strings"
)

func printFCP(name string, fcp *api.FCPatch) {
	fmt.Printf("    %v: %v\n", name, fcp.GetCp())
}

func printFP(name string, fp *api.FixturePatch) {
	fmt.Printf("%v\n", name)

	for control, fcp := range fp.GetByControl() {
		printFCP(control, fcp)
	}
}

func printFPM(fpm *api.FixturePatchMap) {
	for fixture, fp := range fpm.GetByFixture() {
		printFP(fixture, fp)
	}
}

func init() {

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["lsfp"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) > 0 {
			stateName = args[0]
		}

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		fpm, err := c.FixturePatches(context.Background(), &api.StringMsg{Val: stateName})
		if failedTo("get fixture patches", err) {
			return
		}

		//fmt.Printf("%v\n", fpm)
		printFPM(fpm)
	}

	help["lsfp"] = &helpEntry{
		short:  "Get fixture patches for a state",
		syntax: "[state]",
		man:    `Returns the fixture patches for this state`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["sfp"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) < 2 {
			fmt.Printf("Require a fixture, control, and control point")
			return
		}

		fcp := &api.FCPatch{
			Cp: args[2],
		}

		fp := &api.FixturePatch{
			ByControl: make(map[string]*api.FCPatch),
		}
		fp.ByControl[args[1]] = fcp

		fpm := &api.FixturePatchMap{
			ByFixture: make(map[string]*api.FixturePatch),
			State:     stateName,
		}
		fpm.ByFixture[args[0]] = fp

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.SetFixturePatches(context.Background(), fpm)
		if failedTo("get response", err) {
			return
		}
	}

	help["sfp"] = &helpEntry{
		short:  "Set fixture patch",
		syntax: "fixture control control_point",
		man:    `Sets a fixture patch value on the current state`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["rmfp"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) < 1 {
			fmt.Printf("Requires a fixture, control")
			return
		}

		fcp := &api.FCPatch{
			Cp: "X",
		}

		fp := &api.FixturePatch{
			ByControl: make(map[string]*api.FCPatch),
		}
		fp.ByControl[args[1]] = fcp

		fpm := &api.FixturePatchMap{
			ByFixture: make(map[string]*api.FixturePatch),
			State:     stateName,
		}
		fpm.ByFixture[args[0]] = fp

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.RemoveFixturePatches(context.Background(), fpm)
		if failedTo("get response", err) {
			return
		}
	}

	help["rmfp"] = &helpEntry{
		short:  "Remove fixture patch",
		syntax: "fixture control",
		man:    `Remvoes a fixture patch value on the current state`,
	}
}
