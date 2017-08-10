package gate

import (
	"bytes"
	fmt "fmt"
	"s7/share/middleware/mnet"
	"s8/grain/user"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// actor --------------------------------------------------
type connActor struct {
	name string
}

func (c *connActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
		mnet.SendMsg(ctx, "connected")
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
			ret, e := cmdHandler.Handle(ev, ctx)
			mnet.SendMsg(ctx, fmt.Sprintf("%v %v", ret, e))
		}
	case *user.ChannelEventNotify:
		var b bytes.Buffer
		fmt.Fprint(&b, ev.Channel)

		if ev.Enter != "" {
			fmt.Fprint(&b, " +", ev.Enter)
		}

		if ev.Quit != "" {
			fmt.Fprint(&b, " -", ev.Quit)
		}

		if ev.Msg != "" || ev.User != "" {
			fmt.Fprint(&b, " ", ev.User, ":", ev.Msg)
		}

		mnet.SendMsg(ctx, b.String())
	}
}
