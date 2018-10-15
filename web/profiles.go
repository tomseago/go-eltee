package web

import (
	"github.com/tomseago/go-eltee"
)

type ReqProfilesControl struct {
	Id   string
	Name string
	Type string

	Controls []*ReqProfilesControl `json:",omitempty"`
}

type ReqProfilesProfile struct {
	Id           string
	Name         string
	ChannelCount int

	Controls *ReqProfilesControl
}

type ReqProfilesResp struct {
	Msg  string
	Body map[string]*ReqProfilesProfile
}

func recurseIntoGroup(grp *eltee.GroupProfileControl) *ReqProfilesControl {
	out := &ReqProfilesControl{
		Id:   grp.Id(),
		Name: grp.Name(),
		Type: grp.Type(),

		Controls: make([]*ReqProfilesControl, 0),
	}

	for _, control := range grp.Controls {
		var child *ReqProfilesControl

		subGroup, ok := control.(*eltee.GroupProfileControl)
		if ok {
			// It is a sub group
			child = recurseIntoGroup(subGroup)
		} else {
			// Not a sub group, just a regular control
			child = &ReqProfilesControl{
				Id:   control.Id(),
				Name: control.Name(),
				Type: control.Type(),
			}

			// todo: Express details about what exact channel does what
			// for this control so that this information can be annotated
			// into the UI in a useful way. This probably gets complicated
			// though so skipping it for now. Probably want to pass a
			// map[int]string down into the profile control so that each
			// type can annotate itself appropriately
		}

		out.Controls = append(out.Controls, child)
	}

	return out
}

func (client *WebServerSocketClient) HandleReqProfiles(msg *ClientMessage) {
	profiles := client.web.s.GetProfiles()

	resp := &ReqProfilesResp{
		Msg:  "profileList",
		Body: make(map[string]*ReqProfilesProfile),
	}

	for id, prof := range profiles {
		resp.Body[id] = &ReqProfilesProfile{
			Id:           prof.Id,
			Name:         prof.Name,
			ChannelCount: prof.ChannelCount,

			Controls: recurseIntoGroup(prof.Controls),
		}
	}

	client.WriteJSON(resp)
}
