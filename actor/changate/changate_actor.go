package changate

import (
	"gamelib/actor/plugin/logger"
	"gamelib/base/net/util"
	butil "gamelib/base/util"
	"log"
	stdnet "net"

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

				p := actor.FromInstance(&connActor{c: c}).WithMiddleware(logger.MsgLogger)
				ctx.SpawnPrefix(p, "conn")
			}
		}()
	case *actor.Stopping, *actor.Restarting:
		cg.l.Close()
	}
}

func Start(start, end int) {
	addr, e := util.FindLanAddr("udp", start, end)
	butil.PanicOnErr(e)
	l, e := kcp.Listen(addr)
	prop := actor.FromProducer(func() actor.Actor {
		return &changateActor{l: l.(*kcp.Listener)}
	}).WithMiddleware(
		logger.MsgLogger,
	)
	ChangatePID, e = actor.SpawnNamed(prop, "changate")
	if e != nil {
		log.Panic(e)
	}
}
