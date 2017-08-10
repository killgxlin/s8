package user

import (
	"fmt"
	"log"
	"path"
	"s9/grain/channel"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

//a Go struct implementing the User interface
type user struct {
	cluster.Grain
	channels map[string]bool
	connPID  *actor.PID
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
	g := channel.GetChannelGrain(r.Channel)
	res1, e1 := g.Enter(&channel.EnterRequest{User: u.ID()})
	if e1 != nil {
		e = e1
		return
	}

	u.channels[r.Channel] = true
	res.Members = res1.Members
	return
}

func (u *user) QuitChannel(r *QuitChannelRequest) (res *Unit, e error) {
	g := channel.GetChannelGrain(r.Channel)
	_, e = g.Quit(&channel.QuitRequest{User: u.ID()})
	if e == nil {
		delete(u.channels, r.Channel)
	}
	return
}

func (u *user) FireChannelEvent(r *FireChannelEventRequest) (res *Unit, e error) {
	g := channel.GetChannelGrain(r.Channel)
	_, e = g.Publish(&channel.PublishRequest{User: u.ID(), Msg: r.Msg, ToUser: r.ToUser})
	if e == nil {
		delete(u.channels, r.Channel)
	}
	return
}

func (u *user) NotifyChannelEvent(r *ChannelEventNotify) (res *Unit, e error) {
	if u.connPID != nil {
		u.connPID.Tell(r)
	}
	return
}

func (u *user) Register(r *RegisterRequest) (res *Unit, e error) {
	if u.connPID != nil {
		e = fmt.Errorf("%v already registerred", u.ID())
		return
	}
	u.connPID = r.Pid
	log.Println(*u)
	return
}

func (u *user) Unregister(r *UnregisterRequest) (res *Unit, e error) {
	if u.connPID == nil || !u.connPID.Equal(r.Pid) {
		e = fmt.Errorf("%v need not unregister", u.ID())
		return
	}

	u.connPID = nil
	log.Println(*u)
	return
}

func init() {
	channel.RegisterObserver(func(touser, channel, user, msg, enter, quit string) {
		g := GetUserGrain(touser)
		go g.NotifyChannelEvent(&ChannelEventNotify{
			Channel: channel,
			User:    user,
			Msg:     msg,
			Enter:   enter,
			Quit:    quit,
		})
	})

	UserFactory(func() User { return &user{} })
}
