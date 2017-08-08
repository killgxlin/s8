package channel

import (
	"log"
	"s8/node"
	"s8/node/term"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

type channelActor struct {
}

func (g *channelActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		termPID, e := cluster.Get("term", "term")
		if e != nil {
			log.Println(e)
			return
		}

		termPID.Request(&term.RegLocal{}, ctx.Self())
	}
	//node.LogContext(ctx)
}

func init() {
	node.RegisterGlobal("channel", actor.FromProducer(func() actor.Actor {
		return &channelActor{}
	}))
}