package channel

import (
	"fmt"
	"log"
	"s8/grain/registry"

	"github.com/AsynkronIT/protoactor-go/cluster"
)

//a Go struct implementing the Channel interface

type channel struct {
	cluster.Grain
	members map[string]bool
}

func (c *channel) Init(id string) {
	c.Grain.Init(id)
	c.members = map[string]bool{}
}

func (c *channel) doRegistyStuff(inc bool) {
	if inc {
		if len(c.members) == 1 {
			g := registry.GetRegistryGrain("registry")
			_, e := g.Register(&registry.RegisterRequest{Channel: c.ID()})
			if e != nil {
				log.Panic(e)
			}
		}
	} else {
		if len(c.members) == 0 {
			g := registry.GetRegistryGrain("registry")
			_, e := g.Unregister(&registry.UnregisterRequest{Channel: c.ID()})
			if e != nil {
				log.Panic(e)
			}
		}
	}
}

func (c *channel) Enter(r *EnterRequest) (res *EnterResponse, e error) {
	res = &EnterResponse{}
	_, ok := c.members[r.User]

	if !ok {
		for user := range c.members {
			OnEvent(user, c.ID(), "", "", r.User, "")
		}
	}

	c.members[r.User] = true
	for user := range c.members {
		res.Members = append(res.Members, user)
	}

	c.doRegistyStuff(true)

	return
}

func (c *channel) Quit(r *QuitRequest) (res *Unit, e error) {
	_, ok := c.members[r.User]

	if ok {
		delete(c.members, r.User)
		for user := range c.members {
			OnEvent(user, c.ID(), "", "", "", r.User)
		}
	}

	c.doRegistyStuff(false)
	return
}

func (c *channel) Notify(r *NotifyRequest) (res *Unit, e error) {
	if _, ok := c.members[r.User]; !ok {
		e = fmt.Errorf("%v not member of %v", r.User, c.ID())
		return
	}

	if _, ok := c.members[r.ToUser]; ok {
		OnEvent(r.ToUser, c.ID(), r.User, r.Msg, "", "")
	} else {
		for user := range c.members {
			OnEvent(user, c.ID(), r.User, r.Msg, "", "")
		}
	}

	return
}

func init() {
	//apply DI and setup logic

	ChannelFactory(func() Channel { return &channel{} })

}
