package gate

import (
	"fmt"
	"log"
	"s7/share/middleware/mnet"
	"s8/node"
	"s8/util"
	"strings"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

// actor --------------------------------------------------
type connActor struct {
	name string
}

func (c *connActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
		mnet.SendMsg(ctx, "hello")
	case *mnet.ConnectionEvent:
		ctx.Self().Stop()
	case string:
		cmdHandler.Handle(ev, "", nil, ctx)
	}
}

// handler --------------------------------------------------
var (
	cmdHandler = util.NewCmdHandler()
)

func init() {
	cmdHandler.Register("eco", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		mnet.SendMsg(ctx.(actor.Context), strings.Join(args, " "))
	})
	cmdHandler.Register("nam", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)

		name = args[0]

		userPID, e := cluster.Get("user:"+name, "user")
		if e != nil {
			log.Panic(e)
		}
		ctx1.RequestFuture(
			userPID,
			&node.Command{
				Data: "register",
				Name: name,
				Pid:  ctx1.Self(),
			},
			time.Second,
		)

		mnet.SendMsg(ctx1, fmt.Sprintf("nam %v %v", "ok", 0))

		ca.name = name
	})
	cmdHandler.Register("cls", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)

		listPID, e := cluster.Get("chanlist", "chanlist")
		if e != nil {
			log.Panic(e)
		}

		rep, e := ctx1.RequestFuture(
			listPID,
			&node.Command{Data: "list"},
			time.Second,
		).Result()
		if e != nil {
			log.Panic(e)
		}

		err := "ok"
		ret := ""

		cmd, ok := rep.(*node.Command)
		if !ok || cmd == nil {
			err = "unknown"
		} else {
			ret = cmd.Data
		}

		mnet.SendMsg(ctx1, fmt.Sprintf("cls %v %v", err, ret))
	})
	cmdHandler.Register("cen", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)

		chanPID, e := cluster.Get("channel:"+args[0], "channel")
		if e != nil {
			log.Panic(e)
		}

		rep, e := ctx1.RequestFuture(
			chanPID,
			&node.Command{Data: "enter " + ca.name},
			time.Second,
		).Result()
		if e != nil {
			log.Panic(e)
		}

		err := "ok"
		ret := ""

		cmd, ok := rep.(*node.Command)
		if !ok || cmd == nil {
			err = "unknown"
		} else {
			ret = cmd.Data
		}

		mnet.SendMsg(ctx1, fmt.Sprintf("cen %v %v", err, ret))
	})
	cmdHandler.Register("cex", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)

		chanPID, e := cluster.Get("channel:"+args[0], "channel")
		if e != nil {
			log.Panic(e)
		}

		rep, e := ctx1.RequestFuture(
			chanPID,
			&node.Command{Data: "exit"},
			time.Second,
		).Result()
		if e != nil {
			log.Panic(e)
		}

		err := "ok"
		ret := ""

		cmd, ok := rep.(*node.Command)
		if !ok || cmd == nil {
			err = "unknown"
		} else {
			ret = cmd.Data
		}

		mnet.SendMsg(ctx1, fmt.Sprintf("cex %v %v", err, ret))
	})
	cmdHandler.Register("cnt", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)

		chanPID, e := cluster.Get("channel:"+args[0], "channel")
		if e != nil {
			log.Panic(e)
		}

		ctx1.Request(chanPID, &node.Command{Data: "notify " + args[1]})
	})
	cmdHandler.Register("notify", func(args []string, name string, pid *actor.PID, ctx interface{}) {
		ctx1 := ctx.(actor.Context)

		mnet.SendMsg(ctx1, strings.Join(args, " "))
	})
}
