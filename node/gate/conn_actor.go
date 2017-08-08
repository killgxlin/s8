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
	token   string
	account string
	userid  uint64
}

func (c *connActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
		mnet.SendMsg(ctx, "hello")
	case *mnet.ConnectionEvent:
		ctx.Self().Stop()
	case string:
		cmdHandler.Handle(ev, ctx)
	}
}

// handler --------------------------------------------------
var (
	cmdHandler = util.NewCmdHandler()
)

func init() {
	cmdHandler.Register("eco", func(args []string, ctx interface{}) {
		mnet.SendMsg(ctx.(actor.Context), strings.Join(args, " "))
	})
	cmdHandler.Register("au", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)

		token := args[0]
		log.Printf("token %v\n", token)
		mnet.SendMsg(ctx1, fmt.Sprintf("au %v %v", "ok", 0))

		ca.token = args[0]
		ca.account = ca.token
		ca.userid = 1
	})
	cmdHandler.Register("upa", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)

		userPID, e := cluster.Get(fmt.Sprintf("user:%v", ca.userid), "user")
		if e != nil {
			log.Panic(e)
		}

		ud, e := userPID.RequestFuture(
			&node.Command{Userid: ca.userid, Cmd: "parent"},
			time.Second,
		).Result()
		if e != nil {
			log.Panic(e)
		}

		err := "ok"
		ret := ""

		cmd, ok := ud.(*node.Command)
		if !ok || cmd == nil {
			err = "unknown"
		} else if cmd.Cmd == "" {
			err = "noexist"
		} else {
			ret = cmd.Cmd
		}

		mnet.SendMsg(ctx1, fmt.Sprintf("upa %v %v", err, ret))
	})
	cmdHandler.Register("ulo", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)

		userPID, e := cluster.Get(fmt.Sprintf("user:%v", ca.userid), "user")
		if e != nil {
			log.Panic(e)
		}

		ud, e := userPID.RequestFuture(
			&node.Command{Userid: ca.userid, Cmd: "load"},
			time.Second,
		).Result()
		if e != nil {
			log.Panic(e)
		}

		err := "ok"
		ret := ""

		cmd, ok := ud.(*node.Command)
		if !ok || cmd == nil {
			err = "unknown"
		} else if cmd.Cmd == "" {
			err = "noexist"
		} else {
			ret = cmd.Cmd
		}

		mnet.SendMsg(ctx1, fmt.Sprintf("ulo %v %v", err, ret))
	})
	cmdHandler.Register("ucr", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)

		userPID, e := cluster.Get(fmt.Sprintf("user:%v", ca.userid), "user")
		if e != nil {
			log.Panic(e)
		}

		rep, e := userPID.RequestFuture(
			&node.Command{
				Userid: ca.userid,
				Cmd:    strings.Join(append([]string{"create"}, args...), " "),
			},
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
		} else if cmd.Cmd == "" {
			err = "noexist"
		} else {
			ret = cmd.Cmd
		}

		mnet.SendMsg(ctx1, fmt.Sprintf("ucr %v %v", err, ret))
	})
	cmdHandler.Register("cls", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)

		listPID, e := cluster.Get("chanlist", "chanlist")
		if e != nil {
			log.Panic(e)
		}

		rep, e := listPID.RequestFuture(&node.Command{Cmd: "list"}, time.Second).Result()
		if e != nil {
			log.Panic(e)
		}

		err := "ok"
		ret := ""

		cmd, ok := rep.(*node.Command)
		if !ok || cmd == nil {
			err = "unknown"
		} else {
			ret = cmd.Cmd
		}

		mnet.SendMsg(ctx1, fmt.Sprintf("cls %v %v", err, ret))
	})
	cmdHandler.Register("ccr", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)

		listPID, e := cluster.Get("chanlist", "chanlist")
		if e != nil {
			log.Panic(e)
		}

		rep, e := listPID.RequestFuture(&node.Command{Cmd: "create"}, time.Second).Result()
		if e != nil {
			log.Panic(e)
		}

		err := "ok"
		ret := ""

		cmd, ok := rep.(*node.Command)
		if !ok || cmd == nil {
			err = "unknown"
		} else {
			ret = cmd.Cmd
		}

		mnet.SendMsg(ctx1, fmt.Sprintf("ccr %v %v", err, ret))
	})
	cmdHandler.Register("cde", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)

		listPID, e := cluster.Get("chanlist", "chanlist")
		if e != nil {
			log.Panic(e)
		}

		rep, e := listPID.RequestFuture(&node.Command{Cmd: "delete " + args[0]}, time.Second).Result()
		if e != nil {
			log.Panic(e)
		}

		err := "ok"
		ret := ""

		cmd, ok := rep.(*node.Command)
		if !ok || cmd == nil {
			err = "unknown"
		} else {
			ret = cmd.Cmd
		}

		mnet.SendMsg(ctx1, fmt.Sprintf("cde %v %v", err, ret))
	})
	cmdHandler.Register("cen", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)

		chanPID, e := cluster.Get(args[0], "channel")
		if e != nil {
			log.Panic(e)
		}

		rep, e := chanPID.RequestFuture(&node.Command{Cmd: "enter"}, time.Second).Result()
		if e != nil {
			log.Panic(e)
		}

		err := "ok"
		ret := ""

		cmd, ok := rep.(*node.Command)
		if !ok || cmd == nil {
			err = "unknown"
		} else {
			ret = cmd.Cmd
		}

		mnet.SendMsg(ctx1, fmt.Sprintf("cen %v %v", err, ret))
	})
	cmdHandler.Register("cex", func(args []string, ctx interface{}) {
	})
}
