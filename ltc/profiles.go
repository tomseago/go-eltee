package main

import (
	"context"
	"fmt"
	"github.com/tomseago/go-eltee/api"
)

func printIndent(level int) {
	for i := 0; i < level; i++ {
		fmt.Printf("\t")
	}
}

func printGroup(level int, group *api.GroupProfileControl) {
	if level > 0 {
		printIndent(level)
		fmt.Printf("%v: %v\n", group.Id, group.Name)
	}

	for _, pc := range group.Controls {
		printControl(level+1, pc)
	}
}

func printEnum(level int, enum *api.EnumProfileControl) {
	printIndent(level)
	fmt.Printf("Enum %v: %v\n", enum.Id, enum.Name)
}

func printIntensity(level int, intensity *api.IntensityProfileControl) {
	printIndent(level)
	fmt.Printf("Intensity %v: %v\n", intensity.Id, intensity.Name)
}

func printPanTilt(level int, panTilt *api.PanTiltProfileControl) {
	printIndent(level)
	fmt.Printf("PanTilt %v: %v\n", panTilt.Id, panTilt.Name)
}

func printLedVar(level int, ledVar *api.LedVarProfileControl) {
	printIndent(level)
	fmt.Printf("LedVar %v: %v\n", ledVar.Id, ledVar.Name)
}

func printControl(level int, pc *api.ProfileControl) {
	group := pc.GetGroup()
	if group != nil {
		printGroup(level, group)
		return
	}

	enum := pc.GetEnum()
	if enum != nil {
		printEnum(level, enum)
		return
	}

	intensity := pc.GetIntensity()
	if intensity != nil {
		printIntensity(level, intensity)
		return
	}

	panTilt := pc.GetPanTilt()
	if panTilt != nil {
		printPanTilt(level, panTilt)
		return
	}

	ledVar := pc.GetLedVar()
	if ledVar != nil {
		printLedVar(level, ledVar)
		return
	}
}

func init() {

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["lib"] = func(lc *localContext, args []string) {

		showControls := len(args) > 0

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		libResp, err := c.ProfileLibrary(context.Background(), &api.Void{})
		if failedTo("get lib response", err) {
			return
		}

		//fmt.Printf("%v\n", libResp)
		fmt.Printf("Got %v profiles\n", len(libResp.Profiles))
		for name, prof := range libResp.Profiles {
			fmt.Printf("%v\n", name)

			if showControls {
				printControl(0, prof.Controls)
			}
		}
	}

	help["lib"] = &helpEntry{
		short:  "Get the profile library",
		syntax: "",
		man: `Show all the profiles in the  profile library. Profiles define a 
class of fixture`,
	}
}
