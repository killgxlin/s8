package user

import (
	"s8/node"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type userActor struct {
}

func (g *userActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
	case *node.Command:
		//ctx.Respond(&node.Command{"ok"})
	}
	node.LogContext(ctx)
}

func init() {
	node.RegisterGlobal("user", actor.FromProducer(func() actor.Actor {
		return &userActor{}
	}))
}
