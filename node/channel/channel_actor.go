package channel

import (
	"log"
	"s7/share/middleware/msglogger"
	"s8/node"
	"s8/util"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

// actor --------------------------------------------------------
type channelActor struct {
	name    string
	members map[string]*actor.PID // name PID
	revMap  map[string]string     // strPID name
}

func (ca *channelActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
		listPID, e := cluster.Get("chanlist", "chanlist")
		if e != nil {
			log.Panic(e)
		}
		name, _ := util.GetStrIdFrom(ctx.Self().String(), "channel:")

		ctx.Request(listPID, &node.Command{Cmd: "register " + name})

		ca.name = name
		ca.members = map[string]*actor.PID{}
		ca.revMap = map[string]string{}

	case *node.Command:
		cmdHandler.Handle(ev.Cmd, ctx)
	}
}

func init() {
	node.RegisterGlobal("channel", actor.FromProducer(func() actor.Actor {
		return &channelActor{}
	}).WithMiddleware(msglogger.MsgLogger))
}

// handler -------------------------------------------------------
var (
	cmdHandler = util.NewCmdHandler()
)

func init() {
	cmdHandler.Register("enter", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*channelActor)

		ca.members[args[0]] = ctx1.Sender()
		ca.revMap[ctx1.Sender().String()] = args[0]

		ctx1.Respond(&node.Command{Cmd: ""})
	})

	cmdHandler.Register("exit", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*channelActor)

		name := ca.revMap[ctx1.Sender().String()]
		delete(ca.members, name)
		delete(ca.revMap, ctx1.Sender().String())

		ctx1.Respond(&node.Command{Cmd: ""})
	})

	cmdHandler.Register("notify", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*channelActor)

		log.Println(ca.revMap, ca.members)
		name := ca.revMap[ctx1.Sender().String()]

		msg := "notify " + name + ": " + args[0]
		for _, pid := range ca.members {
			log.Println(pid, msg)
			ctx1.Request(pid, &node.Command{Cmd: msg})
		}
	})
}
