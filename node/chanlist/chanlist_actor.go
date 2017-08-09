package chanlist

import (
	"bytes"
	"s7/share/middleware/msglogger"
	"s8/node"
	"s8/util"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// actor -------------------------------------------------------
type chanlistActor struct {
	chans  map[string]*actor.PID // name pid
	revMap map[string]string     // strPID name
}

func (ca *chanlistActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
		ca.chans = map[string]*actor.PID{}
		ca.revMap = map[string]string{}
	case *node.Command:
		cmdHandler.Handle(ev.Cmd, ctx)
	case *actor.Terminated:
		name, ok := ca.revMap[ev.Who.String()]
		if !ok {
			return
		}

		delete(ca.revMap, ev.Who.String())
		delete(ca.chans, name)
	}
}

func init() {
	node.RegisterGlobal("chanlist", actor.FromProducer(func() actor.Actor {
		return &chanlistActor{}
	}).WithMiddleware(msglogger.MsgLogger))
}

// handler -------------------------------------------------------
var (
	cmdHandler = util.NewCmdHandler()
)

func init() {
	cmdHandler.Register("list", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*chanlistActor)

		var bw bytes.Buffer
		for name := range ca.chans {
			bw.WriteString(name)
			bw.WriteString(" ")
		}
		if bw.Len() > 1 {
			bw.Truncate(bw.Len() - 1)
		}
		ctx1.Respond(&node.Command{Cmd: bw.String()})
	})
	cmdHandler.Register("register", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*chanlistActor)

		chanPID := ctx1.Sender()
		name := args[0]

		ctx1.Watch(chanPID)
		ca.chans[name] = chanPID
		ca.revMap[chanPID.String()] = name
	})
}
