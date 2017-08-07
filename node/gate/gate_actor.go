package gate

import (
	"log"
	"s7/share/middleware/mnet"
	"s7/share/middleware/msglogger"
	cstring "s7/share/net/coder/string"
	"s8/node"

	"github.com/AsynkronIT/protoactor-go/actor"
)

var (
	msgio = cstring.NewReadWriter()
)

type gateActor struct {
}

func (g *gateActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
	case *mnet.AcceptorEvent:
		if ev.E != nil {
			ctx.Self().Stop()
			return
		}
		log.Printf("gate port:%v", ev.GetPort())
		if ev.C != nil {
			prop := actor.FromInstance(&connActor{}).WithMiddleware(
				msglogger.MsgLogger,
				mnet.MakeConnection(ev.C, msgio, true, true),
			)
			ctx.SpawnPrefix(prop, "conn")
		}
	}
}

func init() {
	prop := actor.FromProducer(func() actor.Actor {
		return &gateActor{}
	}).WithMiddleware(
		msglogger.MsgLogger,
		mnet.MakeAcceptor("0.0.0.0:1234", 100, 100),
	)
	node.RegisterLocal("gate", prop)
}
