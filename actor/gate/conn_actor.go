package gate

import (
	"bytes"
	fmt "fmt"
	"myutil/actor/plugin/net"
	"s8/grain/user"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// actor --------------------------------------------------
type connActor struct {
	name string
}

func (c *connActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
		net.SendMsg(ctx, "connected")
	case *actor.Stopping, *actor.Restarting:
		if c.name != "" {
			u := user.GetUserGrain(c.name)
			_, e := u.Unregister(&user.UnregisterRequest{Pid: ctx.Self()})
			if e != nil {
				return
			}
			c.name = ""
		}

	case *net.ConnectionEvent:
		ctx.Self().Stop()
	case string:
		if m != "" {
			ret, e := handler.Handle(m, ctx)
			net.SendMsg(ctx, fmt.Sprintf("%v %v", ret, e))
		}
	case *user.ChannelEventRequest:
		var b bytes.Buffer
		fmt.Fprint(&b, m.Channel)

		if m.Enter != "" {
			fmt.Fprint(&b, " +", m.Enter)
		}

		if m.Quit != "" {
			fmt.Fprint(&b, " -", m.Quit)
		}

		if m.Msg != "" || m.User != "" {
			fmt.Fprint(&b, " ", m.User, ":", m.Msg)
		}

		net.SendMsg(ctx, b.String())
	}
}
