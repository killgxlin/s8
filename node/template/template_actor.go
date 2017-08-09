package TEMPLATE

import (
	"s7/share/middleware/msglogger"
	"s8/node"
	"s8/util"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type TEMPLATEActor struct {
	name string
}

func (ua *TEMPLATEActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
	case *node.Command:
		cmdHandler.Handle(ev.Data, ev.Name, ev.Pid, ctx)
	}
}

func init() {
	node.RegisterGlobal("TEMPLATE", actor.FromProducer(func() actor.Actor {
		return &TEMPLATEActor{}
	}).WithMiddleware(msglogger.MsgLogger))

}

var (
	cmdHandler = util.NewCmdHandler()
)

func init() {
	cmdHandler.Register("EXAMPLE", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ua := ctx1.Actor().(*TEMPLATEActor)

		ctx1.Respond(&node.Command{Data: ua.name})
	})
}
