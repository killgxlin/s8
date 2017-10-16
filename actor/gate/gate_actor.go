package gate

import (
	"log"
	"myutil/actor/plugin/logger"
	anet "myutil/actor/plugin/net"
	"myutil/base/net/coder/string"
	"myutil/base/net/util"

	"github.com/AsynkronIT/protoactor-go/actor"
)

var (
	msgio = string.NewReadWriter()
)

type gateActor struct {
}

func (g *gateActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
	case *anet.AcceptorEvent:
		if m.E != nil {
			ctx.Self().Stop()
			return
		}
		log.Printf("gate port:%v", m.GetPort())
		if m.C != nil {
			p := actor.FromInstance(&connActor{}).WithMiddleware(
				logger.MsgLogger,
				anet.MakeConnection(m.C, msgio, true, true, 60*60*24),
			)
			ctx.SpawnPrefix(p, "conn")
		}
	}
}

func StartGate(start, end int) {
	addr, e := util.FindLanAddr("tcp", start, end)
	if e != nil {
		log.Panic(e)
	}

	prop := actor.FromProducer(func() actor.Actor {
		return &gateActor{}
	}).WithMiddleware(
		logger.MsgLogger,
		anet.MakeAcceptor(addr, 100, 100),
	)
	_, e = actor.SpawnNamed(prop, "gate")
	if e != nil {
		log.Panic(e)
	}
}
