package main

import (
	"context"
	"fmt"
	"github.com/tomseago/go-eltee/api"
)

func init() {

	//////////////////////////////////////////////////////////////////////////////////////////
	commands["ping"] = func(lc *localContext, args []string) {

		c, err := lc.c.Client()
		if failedTo("get client", err) {
			return
		}

		smResp, err := c.Ping(context.Background(), &api.StringMsg{Val: "Heya!"})
		if failedTo("get ping response", err) {
			return
		}

		fmt.Println(smResp.Val)
	}

	help["ping"] = &helpEntry{
		short:  "Ping something",
		syntax: "server",
		man: `This will ping something else.
And that might be somewhere else.
`,
	}
}
