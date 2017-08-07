package term

import (
	"s8/node"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type termActor struct {
}

func (g *termActor) Receive(ctx actor.Context) {
	//node.LogContext(ctx)
}

func init() {
	node.RegisterGlobal("term", actor.FromProducer(func() actor.Actor {
		return &termActor{}
	}))
}
