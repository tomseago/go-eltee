package main

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"github.com/tomseago/go-eltee"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

//////////////////////////////////////////////////////////////////////
// Package level logging
// This only needs to be in one file per-package

var log = config.Logger("eltee")

//////

var profileLibrary = eltee.NewProfileLibrary()

func TestFile(filename string) error {
	testNode := config.NewAclNode()
	err := testNode.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("While reading '%v': %v", filename, err)
	}

	// Create the fixture
	profName := testNode.ChildAsString("profile")
	profile := profileLibrary.Profiles[profName]
	if profile == nil {
		return fmt.Errorf("Could not find profile '%v'", profName)
	}

	dmx := make([]byte, 512)
	channels := dmx[0:profile.ChannelCount]

	fixture := eltee.NewDmxFixture("test", 0, channels, profile)

	// Create all the control points
	_, cpIndex := eltee.CreateControlPointList(testNode.Child("control_points"))

	// Patch in the control points
	patches := testNode.Child("patches")
	patches.ForEachOrderedChild(func(n string, v *config.AclNode) {
		if err != nil {
			return
		}

		fixtureControl := fixture.Control(n)
		if fixtureControl == nil {
			err = fmt.Errorf("Unable to find fixture control with id '%v'", n)
			return
		}

		cpName := v.AsString()
		cp := cpIndex[cpName]
		if cp == nil {
			err = fmt.Errorf("Unable to find a control point named '%v'", cpName)
			return
		}

		fixtureControl.ControlPoint = cp
	})
	if err != nil {
		return err
	}

	// Add Lenses

	// Update the fixture using the control points
	log.Info("Updating the fixture using the control points")
	fixture.Update()

	// Check the results
	log.Info("Checking the results")
	expected := testNode.Child("expected")
	errs := make([]error, 0)
	expected.ForEachOrderedChild(func(nn string, val *config.AclNode) {
		ix, _ := strconv.Atoi(nn)

		if ix < 1 {
			err = fmt.Errorf("Channel index must begin with 1. Invalid index %v", nn)
			errs = append(errs, err)
			return
		}
		if ix > profile.ChannelCount {
			err = fmt.Errorf("Invalid expected index %v. Max is %v", nn, profile.ChannelCount)
			errs = append(errs, err)
			return
		}

		eVal := val.AsInt()
		actual := int(channels[ix-1])

		if actual != eVal {
			err = fmt.Errorf("Wrong value for channel %v. Expected %v got %v.", nn, eVal, actual)
			errs = append(errs, err)
		}
	})

	if len(errs) > 0 {
		for ix := 0; ix < len(errs); ix++ {
			log.Warning(errs[ix].Error())
		}
		log.Warningf("Full channel values: %v", channels)
		return fmt.Errorf("Found %v wrong expected channel values", len(errs))
	}

	return nil
}

func TestAllIn(dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		log.Debugf("name=%v  ext=%v", file.Name(), ext)
		if ext != ".acl" {
			continue
		}

		// base := filepath.Base(file.Name())
		// base = base[:len(base)-4]

		full := filepath.Join(dirname, file.Name())
		log.Infof("Testing '%v'", full)

		err = TestFile(full)
		if err != nil {
			log.Errorf("%v: %v", full, err)
			// But try to load other things
		} else {
			log.Infof("%v: Success", full)
		}
	}

	return nil
}

func main() {
	config.ColoredLoggingToConsole()

	err := profileLibrary.LoadDirectory("profiles")
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	log.Info("----------------------------------")
	log.Infof("Profile Library: %v", profileLibrary)
	log.Info("----------------------------------")

	TestAllIn("sandbox/control_tests")
}
