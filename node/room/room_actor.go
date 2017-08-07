package room

import (
	"s8/node"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type roomActor struct {
}

func (g *roomActor) Receive(ctx actor.Context) {
	//node.LogContext(ctx)
}

func init() {
	node.RegisterGlobal("room", actor.FromProducer(func() actor.Actor {
		return &roomActor{}
	}))
}
