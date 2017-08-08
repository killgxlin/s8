package chanlist

import (
	"s8/node"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type chanlistActor struct {
}

func (g *chanlistActor) Receive(ctx actor.Context) {
	//node.LogContext(ctx)
}

func init() {
	node.RegisterGlobal("chanlist", actor.FromProducer(func() actor.Actor {
		return &chanlistActor{}
	}))
}
