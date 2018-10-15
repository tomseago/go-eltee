package web

import (
	"github.com/tomseago/go-eltee"
)

type ReqFixturesFixture struct {
	Profile      string
	UseOverrides bool
}

type ReqFixturesResp struct {
	Msg  string
	Body map[string]*ReqFixturesFixture
}

func (client *WebServerSocketClient) HandleReqFixtures(msg *ClientMessage) {
	fixtures := client.web.s.GetFixtures()

	resp := &ReqFixturesResp{
		Msg:  "fixtureList",
		Body: make(map[string]*ReqFixturesFixture),
	}

	for name, fix := range fixtures {
		resp.Body[name] = &ReqFixturesFixture{
			Profile: fix.Profile().Id,
		}

		dmx, ok := fix.(*eltee.DmxFixture)
		if ok {
			resp.Body[name].UseOverrides = dmx.GetUseOverrides()
		}
	}

	client.WriteJSON(resp)
}
