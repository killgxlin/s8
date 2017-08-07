package gate

import (
	"s8/node"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type gateActor struct {
}

func (g *gateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		//startServer(ctx)
	}
	node.LogContext(ctx)
}

func init() {
	node.RegisterLocal("gate", actor.FromProducer(func() actor.Actor {
		return &gateActor{}
	}))
}
