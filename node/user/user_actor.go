package user

import (
	"log"
	"s7/share/middleware/msglogger"
	"s8/node"
	"s8/util"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type userActor struct {
	name string
}

func (ua *userActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
		var e error
		ua.name, e = util.GetStrIdFrom(ctx.Self().Id, "user:")
		if e != nil {
			log.Panic(e)
		}
	case *node.Command:
		cmdHandler.Handle(ev.Data, ctx)
	}
}

func init() {
	node.RegisterGlobal("user", actor.FromProducer(func() actor.Actor {
		return &userActor{}
	}).WithMiddleware(msglogger.MsgLogger))

}

var (
	cmdHandler = util.NewCmdHandler()
)

func init() {
	cmdHandler.Register("chanenter", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ua := ctx1.Actor().(*userActor)

		ctx1.Respond(&node.Command{Cmd: ua.data})
	})
	cmdHandler.Register("chanexit", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ua := ctx1.Actor().(*userActor)

		ua.data = strings.Join(args, " ")
		ctx1.Respond(&node.Command{Cmd: ua.data})
	})
	cmdHandler.Register("chanevent", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)

		ctx1.Respond(&node.Command{Cmd: ctx1.Parent().String()})
	})
}
