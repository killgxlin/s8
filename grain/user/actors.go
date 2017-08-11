package user

import (
	"fmt"
	"log"
	"path"
	"s8/grain/channel"
	"s8/util"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

//a Go struct implementing the User interface
type user struct {
	cluster.Grain
	channels map[string]bool
	connpid  *actor.PID
}

func (u *user) Init(id string) {
	u.Grain.Init(id)
	u.channels = map[string]bool{}
}

func (u *user) ListChannel(r *ListChannelRequest) (res *ListChannelResponse, e error) {
	res = &ListChannelResponse{}
	for channel := range u.channels {
		if ok, _ := path.Match(r.Pattern, channel); ok {
			res.Channel = append(res.Channel, channel)
		}
	}
	return
}

func (u *user) EnterChannel(r *EnterChannelRequest) (res *EnterChannelResponse, e error) {
	res = &EnterChannelResponse{}

	for _, chnn := range r.Selector.Channel {
		resChannel := &EnterChannelResponse_Channel{Channel: chnn}
		res.Channels = append(res.Channels, resChannel)

		g := channel.GetChannelGrain(chnn)
		res1, e1 := g.Enter(&channel.EnterRequest{User: u.ID()})
		if e1 != nil {
			resChannel.Err = e.Error()
		}
		if res1 != nil {
			resChannel.Members = res1.Members
		}

		if e1 == nil {
			u.channels[chnn] = true
		}
	}

	return
}

func (u *user) QuitChannel(r *QuitChannelRequest) (res *Unit, e error) {
	quit := func(chnn string) {
		g := channel.GetChannelGrain(chnn)
		_, e = g.Quit(&channel.QuitRequest{User: u.ID()})
		if e == nil {
			delete(u.channels, chnn)
		}
	}
	toquit := util.SelectMap(u.channels, r.Selector.Pattern)
	for _, chnn := range r.Selector.Channel {
		if _, ok := u.channels[chnn]; ok {
			toquit[chnn] = true
		}
	}
	for chnn := range toquit {
		quit(chnn)
	}

	return
}

func (u *user) NotifyChannelEvent(r *NotifyChannelEventRequest) (res *Unit, e error) {
	notify := func(chnn string) {
		g := channel.GetChannelGrain(chnn)
		_, e = g.Notify(&channel.NotifyRequest{User: u.ID(), Msg: r.Msg, ToUser: r.ToUser})
		if e != nil {
			delete(u.channels, chnn)
		}
	}
	tonotify := util.SelectMap(u.channels, r.Selector.Pattern)
	for _, chnn := range r.Selector.Channel {
		if _, ok := u.channels[chnn]; ok {
			tonotify[chnn] = true
		}
	}
	for chnn := range tonotify {
		notify(chnn)
	}

	return
}

func (u *user) OnChannelEvent(r *ChannelEventRequest) (res *Unit, e error) {
	if u.connpid != nil {
		u.connpid.Tell(r)
	}
	return
}

func (u *user) Register(r *RegisterRequest) (res *Unit, e error) {
	if u.connpid != nil {
		e = fmt.Errorf("%v already registerred", u.ID())
		return
	}
	u.connpid = r.Pid
	log.Println(*u)
	return
}

func (u *user) Unregister(r *UnregisterRequest) (res *Unit, e error) {
	if u.connpid == nil || !u.connpid.Equal(r.Pid) {
		e = fmt.Errorf("%v need not unregister", u.ID())
		return
	}

	u.connpid = nil
	log.Println(*u)
	return
}

func init() {
	channel.RegisterEventHandler(func(touser, channel, user, msg, enter, quit string) {
		g := GetUserGrain(touser)
		go g.OnChannelEvent(&ChannelEventRequest{
			Channel: channel,
			User:    user,
			Msg:     msg,
			Enter:   enter,
			Quit:    quit,
		})
	})

	UserFactory(func() User { return &user{} })
}
