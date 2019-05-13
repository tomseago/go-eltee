package main

import (
	"context"
	"fmt"
	"github.com/tomseago/go-eltee/api"
)

func init() {

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["cs"] = func(lc *localContext, args []string) {

		if len(args) < 1 {
			lc.stateName = ""
			return
		}

		lc.stateName = args[0]
	}

	help["cs"] = &helpEntry{
		short:  "Change state name",
		syntax: "[stateName]",
		man: `Change the default state name used for commands
that are state specific.
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["lsst"] = func(lc *localContext, args []string) {

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		resp, err := c.StateNames(context.Background(), &api.Void{})
		if failedTo("get response", err) {
			return
		}

		fmt.Println(resp.List)
	}

	help["lsst"] = &helpEntry{
		short:  "List state names",
		syntax: "",
		man: `List the names of loaded states
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["apply"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) > 0 {
			stateName = args[0]
		}

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.ApplyState(context.Background(), &api.StringMsg{Val: stateName})
		if failedTo("apply state", err) {
			return
		}

		fmt.Printf("Applied %v\n", stateName)
	}

	help["apply"] = &helpEntry{
		short:  "Apply a named state",
		syntax: "[stateName]",
		man: `Applies the named state to the current values
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["load"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) > 0 {
			stateName = args[0]
		}

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.LoadLoadableState(context.Background(), &api.StringMsg{Val: stateName})
		if failedTo("load state", err) {
			return
		}

		fmt.Printf("Loaded state %v\n", stateName)
	}

	help["load"] = &helpEntry{
		short:  "Load a loadable state",
		syntax: "[stateName]",
		man: `Loads the state with the given name from the
loadable states directory.
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["save"] = func(lc *localContext, args []string) {

		stateName := lc.stateName

		if len(args) > 0 {
			stateName = args[0]
		}

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.SaveLoadableState(context.Background(), &api.StringMsg{Val: stateName})
		if failedTo("save state", err) {
			return
		}

		fmt.Printf("Saved state %v\n", stateName)
	}

	help["save"] = &helpEntry{
		short:  "Save a loadable state",
		syntax: "[stateName]",
		man: `Saves the state with the given name to the
loadable states directory.
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["saveall"] = func(lc *localContext, args []string) {

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.SaveAllStates(context.Background(), &api.Void{})
		if failedTo("save all states", err) {
			return
		}

		fmt.Printf("Saved all states\n")
	}

	help["saveall"] = &helpEntry{
		short:  "Save all loadable states",
		syntax: "",
		man: `Saves all states to the
loadable states directory overwriting anything previous.
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["lslst"] = func(lc *localContext, args []string) {

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		resp, err := c.LoadableStateNames(context.Background(), &api.Void{})
		if failedTo("get response", err) {
			return
		}

		fmt.Println(resp.List)
	}

	help["lslst"] = &helpEntry{
		short:  "List loadable state names",
		syntax: "",
		man: `List the names of loadable states in the
loadable states directory.
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["addst"] = func(lc *localContext, args []string) {

		if len(args) < 1 {
			fmt.Printf("A state name is required\n")
			return
		}
		stateName := args[0]

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.AddState(context.Background(), &api.StringMsg{Val: stateName})
		if failedTo("add state", err) {
			return
		}
	}

	help["addst"] = &helpEntry{
		short:  "Add a new empty state",
		syntax: "stateName",
		man: `Creates a new empty state that control points can then
be set into. The new state is only in memory and not saved.
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["rmst"] = func(lc *localContext, args []string) {

		if len(args) < 1 {
			fmt.Printf("A state name is required\n")
			return
		}
		stateName := args[0]

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.RemoveState(context.Background(), &api.StringMsg{Val: stateName})
		if failedTo("remove state", err) {
			return
		}
	}

	help["rmst"] = &helpEntry{
		short:  "Remove a state",
		syntax: "stateName",
		man: `Removes the state from memory. If it was previously
saved then it is still on disk.
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["cpst"] = func(lc *localContext, args []string) {

		if len(args) == 0 {
			fmt.Printf("A destination name is required\n")
			return
		}

		src := lc.stateName
		var dest string

		if len(args) == 1 {
			dest = args[0]
		} else {
			src = args[0]
			dest = args[1]
		}

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.CopyStateTo(context.Background(), &api.SrcDest{Src: src, Dest: dest})
		if failedTo("copy state", err) {
			return
		}
	}

	help["cpst"] = &helpEntry{
		short:  "Copy a state to another name",
		syntax: "[src] dest",
		man: `Copies the state from one name to another.
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["mvst"] = func(lc *localContext, args []string) {

		if len(args) == 0 {
			fmt.Printf("A destination name is required\n")
			return
		}

		src := lc.stateName
		var dest string

		if len(args) == 1 {
			dest = args[0]
		} else {
			src = args[0]
			dest = args[1]
		}

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.MoveStateTo(context.Background(), &api.SrcDest{Src: src, Dest: dest})
		if failedTo("move state", err) {
			return
		}
	}

	help["mvst"] = &helpEntry{
		short:  "Move a state to another name",
		syntax: "[src] dest",
		man: `Like a copy followed by a remove.
`,
	}

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["applyto"] = func(lc *localContext, args []string) {

		if len(args) == 0 {
			fmt.Printf("A destination name is required\n")
			return
		}

		src := lc.stateName
		var dest string

		if len(args) == 1 {
			dest = args[0]
		} else {
			src = args[0]
			dest = args[1]
		}

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		_, err = c.ApplyStateTo(context.Background(), &api.SrcDest{Src: src, Dest: dest})
		if failedTo("apply state", err) {
			return
		}
	}

	help["applyto"] = &helpEntry{
		short:  "Apply one state to another",
		syntax: "[src] dest",
		man: `Like the apply to current command but between states.
`,
	}

}
