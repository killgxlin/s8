package changate

import (
	"log"
	stdnet "net"
	"s7/share/middleware/msglogger"
	"s7/share/net"
	"s7/share/util"

	"github.com/AsynkronIT/protoactor-go/actor"
	kcp "github.com/xtaci/kcp-go"
)

var (
	ChangatePID *actor.PID
)

type changateActor struct {
	l *kcp.Listener
}

func (cg *changateActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		go func() {
			for {
				c, e := cg.l.AcceptKCP()
				if e != nil {
					if e, ok := e.(stdnet.Error); ok && !e.Temporary() {
						return
					}
					continue
				}

				p := actor.FromInstance(&connActor{c: c}).WithMiddleware(msglogger.MsgLogger)
				ctx.SpawnPrefix(p, "conn")
			}
		}()
	case *actor.Stopping, *actor.Restarting:
		cg.l.Close()
	}
}

func Start(start, end int) {
	addr, e := net.FindLanAddr("udp", start, end)
	util.PanicOnErr(e)
	l, e := kcp.Listen(addr)
	prop := actor.FromProducer(func() actor.Actor {
		return &changateActor{l: l.(*kcp.Listener)}
	}).WithMiddleware(
		msglogger.MsgLogger,
	)
	ChangatePID, e = actor.SpawnNamed(prop, "changate")
	if e != nil {
		log.Panic(e)
	}
}
