package eltee

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"github.com/tomseago/go-eltee/api"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

/////////

// A Profile is basically a class of fixture. That is, it represents a particular model of
// light from a particular manufacturer and holds all the information about controls that
// can be manipulated for that light. It is represented by a single .acl file in the
// profiles directory. There may be multiple instances of a profile - each instance being
// known as a fixture - within a particular setup.
type Profile struct {
	Id string

	Name         string
	ChannelCount int
	DefaultData  []byte
	Controls     *GroupProfileControl
}

func NewProfile(id string, rootNode *config.AclNode) (*Profile, error) {
	p := &Profile{
		Id: id,
	}

	p.Name = rootNode.DefChildAsString(id, "name")
	p.ChannelCount = rootNode.DefChildAsInt(0, "channel_count")

	var err error
	p.Controls, err = NewGroupProfileControl("_root", rootNode.Child("controls"))
	if err != nil {
		return nil, err
	}

	defData := rootNode.ChildAsByteList("default_values")
	if len(defData) > 0 {
		p.DefaultData = defData
	}

	return p, nil
}

func (p *Profile) String() string {
	return fmt.Sprintf("%v: '%v' %v\n%v\n", p.Id, p.Name, p.ChannelCount, p.Controls)
}

/////////

// A ProfileLibrary holds all profiles known to the system. It is loaded from a
// set of .acl files.
type ProfileLibrary struct {
	Profiles map[string]*Profile
}

func NewProfileLibrary() *ProfileLibrary {
	lib := &ProfileLibrary{
		Profiles: make(map[string]*Profile),
	}

	return lib
}

// Loads a single file using the given id. When used by LoadDirectory(string) the
// id will be the base name of the .acl file
func (lib *ProfileLibrary) LoadFile(id string, filename string) error {
	node := config.NewAclNode()

	err := node.ParseFile(filename)
	if err != nil {
		return fmt.Errorf("While reading '%v' : %v", filename, err)
	}

	profile, err := NewProfile(id, node)
	if err != nil {
		return fmt.Errorf("While creating profile %v : %v", id, err)
	}

	lib.Profiles[id] = profile
	return nil
}

// Loads the library from all .acl files in a given directory
func (lib *ProfileLibrary) LoadDirectory(dirname string) error {
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

		base := filepath.Base(file.Name())
		base = base[:len(base)-4]

		full := filepath.Join(dirname, file.Name())
		log.Infof("Loading '%v' from '%v'", base, full)

		err = lib.LoadFile(base, full)
		if err != nil {
			log.Errorf("%v", err)
			// But try to load other things
		}
	}

	return nil
}

// Dumps the entire library to a reasonable string representation, mostly for debugging
func (lib *ProfileLibrary) String() string {
	// Output in sorted order
	ids := make([]string, 0, len(lib.Profiles))
	for id, _ := range lib.Profiles {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var b strings.Builder
	b.WriteString("\n")

	for ix := 0; ix < len(ids); ix++ {
		id := ids[ix]
		profile := lib.Profiles[id]
		b.WriteString(profile.String())
		b.WriteString("\n")
	}

	return b.String()
}

func (p *Profile) ToAPI() *api.Profile {
	apiP := &api.Profile{
		Id:           p.Id,
		Name:         p.Name,
		ChannelCount: int32(p.ChannelCount),
		Controls:     p.Controls.ToAPI(),
	}

	return apiP
}
