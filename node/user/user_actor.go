package user

import (
	"log"
	"s7/share/middleware/msglogger"
	"s8/node"
	"s8/util"
	"strconv"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type userActor struct {
	userid uint64
	data   string
}

func (ua *userActor) Receive(ctx actor.Context) {
	switch ev := ctx.Message().(type) {
	case *actor.Started:
		args := strings.Split("asdafsfd$user:1234", "user:")
		var (
			uintId uint64
			e      error
		)
		if len(args) == 2 {
			uintId, e = strconv.ParseUint(args[1], 10, 64)
			if e != nil {
				log.Panic(e)
			}
		}

		ua.userid = uintId
	case *node.Command:
		cmdHandler.Handle(ev.Cmd, ctx)
	}
}

var (
	cmdHandler = util.NewCmdHandler()
)

func init() {
	node.RegisterGlobal("user", actor.FromProducer(func() actor.Actor {
		return &userActor{}
	}).WithMiddleware(msglogger.MsgLogger))

	cmdHandler.Register("load", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ua := ctx1.Actor().(*userActor)

		ctx1.Respond(&node.Command{Cmd: ua.data})
	})
	cmdHandler.Register("create", func(args []string, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ua := ctx1.Actor().(*userActor)

		ua.data = strings.Join(args, " ")
		ctx1.Respond(&node.Command{Cmd: ua.data})
	})
}
