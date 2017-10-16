package gate

import (
	"log"
	"s7/share/middleware/mnet"
	"s7/share/middleware/msglogger"
	"s7/share/net"
	cstring "s7/share/net/coder/string"

	"github.com/AsynkronIT/protoactor-go/actor"
)

var (
	msgio = cstring.NewReadWriter()
)

type gateActor struct {
}

func (g *gateActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
	case *mnet.AcceptorEvent:
		if m.E != nil {
			ctx.Self().Stop()
			return
		}
		log.Printf("gate port:%v", m.GetPort())
		if m.C != nil {
			p := actor.FromInstance(&connActor{}).WithMiddleware(
				msglogger.MsgLogger,
				mnet.MakeConnection(m.C, msgio, true, true, 60*60*24),
			)
			ctx.SpawnPrefix(p, "conn")
		}
	}
}

func StartGate(start, end int) {
	addr, e := net.FindLanAddr("tcp", start, end)
	if e != nil {
		log.Panic(e)
	}

	prop := actor.FromProducer(func() actor.Actor {
		return &gateActor{}
	}).WithMiddleware(
		msglogger.MsgLogger,
		mnet.MakeAcceptor(addr, 100, 100),
	)
	_, e = actor.SpawnNamed(prop, "gate")
	if e != nil {
		log.Panic(e)
	}
}
