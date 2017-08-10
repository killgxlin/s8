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
	msgio   = cstring.NewReadWriter()
	GatePID *actor.PID
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

func StartGate(start, end int) {
	addr, e := net.FindLanAddr(start, end)
	if e != nil {
		log.Panic(e)
	}

	prop := actor.FromProducer(func() actor.Actor {
		return &gateActor{}
	}).WithMiddleware(
		msglogger.MsgLogger,
		mnet.MakeAcceptor(addr, 100, 100),
	)
	pid, e := actor.SpawnNamed(prop, "gate")
	if e != nil {
		log.Panic(e)
	}

	GatePID = pid
}
