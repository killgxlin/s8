package channel

import (
	"s7/share/middleware/msglogger"
	"s8/node"
	"s8/util"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// actor --------------------------------------------------------
type channelActor struct {
	members map[string]*actor.PID // strPID pid
}

func (ca *channelActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
		ca.members = map[string]*actor.PID{}
	case *node.Command:
		cmdHandler.Handle(ev.Cmd, ctx)
	}
}

func init() {
	node.RegisterGlobal("channel", actor.FromProducer(func() actor.Actor {
		return &channelActor{}
	}).WithMiddleware(msglogger.MsgLogger))
}

// handler -------------------------------------------------------
var (
	cmdHandler = util.NewCmdHandler()
)

func init() {
	cmdHandler.Register("enter", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*channelActor)

		_ = ca

		ctx1.Respond(&node.Command{Cmd: ""})
	})

	cmdHandler.Register("exit", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*channelActor)

		_ = ca

		ctx1.Respond(&node.Command{Cmd: ""})
	})

}
