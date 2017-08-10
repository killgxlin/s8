package gate

import (
	"bytes"
	"fmt"
	"s7/share/middleware/mnet"
	"s8/grain/registry"
	"s8/grain/user"
	"s8/util"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// actor --------------------------------------------------
type connActor struct {
	name string
}

func (c *connActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
		mnet.SendMsg(ctx, "hello")
	case *actor.Stopping, *actor.Restarting:
		if c.name != "" {
			u := user.GetUserGrain(c.name)
			_, e := u.Unregister(&user.UnregisterRequest{Pid: ctx.Self()})
			if e != nil {
				return
			}
			c.name = ""
		}

	case *mnet.ConnectionEvent:
		ctx.Self().Stop()
	case string:
		if ev != "" {
			cmdHandler.Handle(ev, "", nil, ctx)
		}
	case *user.ChannelEventNotify:
		response(ctx, "notify", ev.Channel, ev.User, ev.Msg, ev.Enter, ev.Quit)
	}
}

func response(ctx interface{}, args ...interface{}) {
	var b bytes.Buffer
	for i, a := range args {
		b.WriteString(fmt.Sprint(a))
		if i != len(args)-1 {
			b.WriteString(" ")
		}
	}
	mnet.SendMsg(ctx.(actor.Context), b.String())
}

// handler --------------------------------------------------
var (
	cmdHandler = util.NewCmdHandler()
)

func init() {
	cmdHandler.Register("eco", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		_, _ = ctx1, ca
		response(ctx, args)
	})
	cmdHandler.Register("nam", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)

		name = args[0]
		if name == "_" {
			name = ""
		}

		if ca.name == name {
			response(ctx, "nam", name, "same name")
			return
		}

		if ca.name != "" {
			u := user.GetUserGrain(ca.name)
			_, e := u.Unregister(&user.UnregisterRequest{Pid: ctx1.Self()})
			if e != nil {
				response(ctx, "nam", name, e)
				return
			}
			ca.name = ""
		}

		if name != "" {
			u := user.GetUserGrain(name)
			_, e := u.Register(&user.RegisterRequest{Pid: ctx1.Self()})
			if e != nil {
				response(ctx, "nam", name, e)
				return
			}
			ca.name = name
		}

		response(ctx, "nam", name, "ok")
	})
	cmdHandler.Register("rlsc", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		g := registry.GetRegistryGrain("registry")
		r, e := g.List(&registry.ListRequest{Pattern: args[0]})

		response(ctx, "rlsc", args[0], strings.Join(r.Channel, ","), e)
	})
	cmdHandler.Register("ulsc", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		u := user.GetUserGrain(ca.name)
		r, e := u.ListChannel(&user.ListChannelRequest{Pattern: args[0]})
		response(ctx, "ulsc", args[0], strings.Join(r.Channel, ","), e)
	})
	cmdHandler.Register("cen", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		u := user.GetUserGrain(ca.name)
		r, e := u.EnterChannel(&user.EnterChannelRequest{Channel: args[0]})
		response(ctx, "cen", args[0], strings.Join(r.Members, ","), e)
	})
	cmdHandler.Register("cqu", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		u := user.GetUserGrain(ca.name)
		_, e := u.QuitChannel(&user.QuitChannelRequest{Channel: args[0]})
		response(ctx, "cqu", args[0], e)
	})
	cmdHandler.Register("cev", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		u := user.GetUserGrain(ca.name)
		_, e := u.FireChannelEvent(&user.FireChannelEventRequest{Channel: args[0], Msg: args[1], ToUser: args[2]})
		response(ctx, "cev", args[0], args[1], args[2], e)
	})
}
