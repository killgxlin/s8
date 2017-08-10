package registry

import (
	"path"

	"github.com/AsynkronIT/protoactor-go/cluster"
)

//a Go struct implementing the User interface

type registry struct {
	cluster.Grain
	channels map[string]bool
}

func (re *registry) Init(id string) {
	re.Grain.Init(id)
	re.channels = map[string]bool{}
}

func (re *registry) List(r *ListRequest) (res *ListResponse, e error) {
	res = &ListResponse{}
	for channel := range re.channels {
		if ok, _ := path.Match(r.Pattern, channel); ok {
			res.Channel = append(res.Channel, channel)
		}
	}
	return
}

func (re *registry) Register(r *RegisterRequest) (res *Unit, e error) {
	if ok := re.channels[r.Channel]; ok {
		return
	}
	re.channels[r.Channel] = true
	return
}

func (re *registry) Unregister(r *UnregisterRequest) (res *Unit, e error) {
	delete(re.channels, r.Channel)
	return
}

func init() {
	//apply DI and setup logic

	RegistryFactory(func() Registry { return &registry{} })

}
