package main

import (
    "context"
    "errors"
    "fmt"
    "github.com/tomseago/go-eltee/api"
    "strconv"
    "strings"
)

// func printIndent(level int) {
// 	for i := 0; i < level; i++ {
// 		fmt.Printf("\t")
// 	}
// }

// func printGroup(level int, group *api.GroupProfileControl) {
// 	if level > 0 {
// 		printIndent(level)
// 		fmt.Printf("%v: %v\n", group.Id, group.Name)
// 	}

// 	for _, pc := range group.Controls {
// 		printControl(level+1, pc)
// 	}
// }

// func printEnum(level int, enum *api.EnumProfileControl) {
// 	printIndent(level)
// 	fmt.Printf("Enum %v: %v\n", enum.Id, enum.Name)
// }

// func printIntensity(level int, intensity *api.IntensityProfileControl) {
// 	printIndent(level)
// 	fmt.Printf("Intensity %v: %v\n", intensity.Id, intensity.Name)
// }

// func printPanTilt(level int, panTilt *api.PanTiltProfileControl) {
// 	printIndent(level)
// 	fmt.Printf("PanTilt %v: %v\n", panTilt.Id, panTilt.Name)
// }

// func printLedVar(level int, ledVar *api.LedVarProfileControl) {
// 	printIndent(level)
// 	fmt.Printf("LedVar %v: %v\n", ledVar.Id, ledVar.Name)
// }

// func printPoint(cp *api.ProfileControl) {
// 	group := pc.GetGroup()
// 	if group != nil {
// 		printGroup(level, group)
// 		return
// 	}

// 	enum := pc.GetEnum()
// 	if enum != nil {
// 		printEnum(level, enum)
// 		return
// 	}

// 	intensity := pc.GetIntensity()
// 	if intensity != nil {
// 		printIntensity(level, intensity)
// 		return
// 	}

// 	panTilt := pc.GetPanTilt()
// 	if panTilt != nil {
// 		printPanTilt(level, panTilt)
// 		return
// 	}

// 	ledVar := pc.GetLedVar()
// 	if ledVar != nil {
// 		printLedVar(level, ledVar)
// 		return
// 	}
// }

func parseNameVal(s string) (string, float64, error) {
	if len(s) == 0 {
		return "", 0, nil
	}

	split := strings.Split(s, ":")
	if len(split) < 2 {
		return "", 0, errors.New("String did not contain a :")
	}

	f, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return "", 0, err
	}

	return split[0], f, nil
}

func init() {

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["lscp"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) > 0 {
			stateName = args[0]
		}

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		list, err := c.ControlPoints(context.Background(), &api.StringMsg{Val: stateName})
		if failedTo("get response", err) {
			return
		}

		//fmt.Printf("%v\n", libResp)
		fmt.Printf("Got %v control points for state %v\n", len(list.Cps), list.GetState())
		for _, cp := range list.GetCps() {
			fmt.Println(cp)
		}
	}

	help["lscp"] = &helpEntry{
		short:  "Get control points for a state",
		syntax: "[state]",
		man: `Retrieves all the control points for a state. If not specified
then values for the current state are returned.`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["scolor"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) < 2 {
			fmt.Printf("Require a name and at least one color component")
			return
		}

		cp := &api.ControlPoint{
			Name: args[0],
		}

		pt := &api.ColorPoint{
			Components: make(map[string]float64),
		}
		cp.Val = &api.ControlPoint_Color{pt}

		for _, arg := range args[1:] {
			name, val, err := parseNameVal(arg)
			if err != nil {
				fmt.Printf(errColor("Unable to parse arg %v\n"), arg)
				return
			}

			pt.Components[name] = val
		}

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		cpList := &api.ControlPointList{
			Cps:    make([]*api.ControlPoint, 1),
			State:  stateName,
			Upsert: true,
		}
		cpList.Cps[0] = cp

		_, err = c.SetControlPoints(context.Background(), cpList)
		if failedTo("get response", err) {
			return
		}
	}

	help["scolor"] = &helpEntry{
		short:  "Set control point to color",
		syntax: "name (component:value)+",
		man: `Sets the named control point to a color value that is
specified by named components. Each component name is a string and each
value is a float. The common names are red, green, and blue but some
fixtures understand others like white, amber, and uv.`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["sxyz"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) < 4 {
			fmt.Printf("Require a name and 3 values")
			return
		}

		cp := &api.ControlPoint{
			Name: args[0],
		}

		pt := &api.XYZPoint{}
		cp.Val = &api.ControlPoint_Xyz{pt}

		var err error
		pt.X, err = strconv.ParseFloat(args[1], 64)
		if failedTo("parse x", err) {
			return
		}

		pt.Y, err = strconv.ParseFloat(args[2], 64)
		if failedTo("parse x", err) {
			return
		}

		pt.Z, err = strconv.ParseFloat(args[3], 64)
		if failedTo("parse x", err) {
			return
		}

		////

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		cpList := &api.ControlPointList{
			Cps:    make([]*api.ControlPoint, 1),
			State:  stateName,
			Upsert: true,
		}
		cpList.Cps[0] = cp

		_, err = c.SetControlPoints(context.Background(), cpList)
		if failedTo("get response", err) {
			return
		}
	}

	help["sxyz"] = &helpEntry{
		short:  "Set xyz control point",
		syntax: "name x y z",
		man: `Sets the named control point as an xyz point using the
three values provided.`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["senum"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) < 2 {
			fmt.Printf("Require a name and at least a value")
			return
		}

		cp := &api.ControlPoint{
			Name: args[0],
		}

		pt := &api.EnumPoint{}
		cp.Val = &api.ControlPoint_Enm{pt}

		item, err := strconv.ParseInt(args[1], 0, 32)
		if failedTo("parse item number", err) {
			return
		}
		pt.Item = int32(item)

		if len(args) > 2 {
			pt.Degree, err = strconv.ParseFloat(args[2], 64)
			if failedTo("parse degree", err) {
				return
			}
		} else {
			pt.Degree = 1.0
		}

		////

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		cpList := &api.ControlPointList{
			Cps:    make([]*api.ControlPoint, 1),
			State:  stateName,
			Upsert: true,
		}
		cpList.Cps[0] = cp

		_, err = c.SetControlPoints(context.Background(), cpList)
		if failedTo("get response", err) {
			return
		}
	}

	help["senum"] = &helpEntry{
		short:  "Set enum control point",
		syntax: "name item [degree]",
		man: `Sets the named control point as an enum with an integer item
and optionally a degree value. The degree defaults to 1.0 if not specified.`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["sintensity"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) < 2 {
			fmt.Printf("Require a name and a value")
			return
		}

		cp := &api.ControlPoint{
			Name: args[0],
		}

		pt := &api.IntensityPoint{}
		cp.Val = &api.ControlPoint_Intensity{pt}

		var err error
		pt.Intensity, err = strconv.ParseFloat(args[1], 64)
		if failedTo("parse value", err) {
			return
		}

		////

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		cpList := &api.ControlPointList{
			Cps:    make([]*api.ControlPoint, 1),
			State:  stateName,
			Upsert: true,
		}
		cpList.Cps[0] = cp

		_, err = c.SetControlPoints(context.Background(), cpList)
		if failedTo("get response", err) {
			return
		}
	}

	help["sintensity"] = &helpEntry{
		short:  "Set intensity control point",
		syntax: "name value",
		man: `Sets the named control point as an intensity point using the
values provided.`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["rmcp"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) < 1 {
			fmt.Printf("Require a name")
			return
		}

		cp := &api.ControlPoint{
			Name: args[0],
		}

		////

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		cpList := &api.ControlPointList{
			Cps:   make([]*api.ControlPoint, 1),
			State: stateName,
		}
		cpList.Cps[0] = cp

		_, err = c.RemoveControlPoints(context.Background(), cpList)
		if failedTo("get response", err) {
			return
		}
	}

	help["rmcp"] = &helpEntry{
		short:  "Remove a control point",
		syntax: "name",
		man: `Removes the named control point from the
current state.`,
	}

    //////////////////////////////////////////////////////////////////////////////////////////
    commands["cplisten"] = func(lc *localContext, args []string) {

        enable := true

        if len(args) >= 1 {
            enable, _ = strconv.ParseBool(args[0])
        }

        if enable {
            if lc.cpListener != nil {
                fmt.Printf("Already listening\n")
                return
            }

            c, err := lc.c.Client()
            if failedTo("get client", err) {
                return
            }

            ctx, cancelFunc := context.WithCancel(context.Background())

            stream, err := c.ControlPointChanges(ctx, &api.Void{})
            if failedTo("get stream", err) {
                return
            }

            lc.cpListener = NewCPListening(stream, cancelFunc)
            go lc.cpListener.Run()
        } else {
            if lc.cpListener == nil {
                fmt.Printf("Not listening\n")
                return
            }

            lc.cpListener.Stop()
            lc.cpListener = nil
        }
    }

    help["cplisten"] = &helpEntry{
        short:  "Start/Stop listening to CP changes",
        syntax: "enable",
        man: `Starts or stops listening to the changes stream for control points`,
    }
}

type cpListening struct {
    stream api.ElTee_ControlPointChangesClient
    cancelFunc context.CancelFunc
}

func NewCPListening(stream api.ElTee_ControlPointChangesClient, cancelFunc context.CancelFunc) *cpListening {
    out := &cpListening{
        stream: stream,
        cancelFunc: cancelFunc,
    }
    return out
}

func (cpl *cpListening) Run() {
    fmt.Printf("cpListening.Run starting\n\r");

    for {
        list, err := cpl.stream.Recv()
        if err != nil {
            fmt.Printf("\n\rClosing listen stream %v\n\r", err)
            break
        }

        // TODO: Add formatting
        fmt.Printf("\n\rChanges: %v\n\r", list)
    }
}

func (cpl *cpListening) Stop() {
    cpl.cancelFunc()
}
