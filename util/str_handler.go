package util

import (
	"log"
	"runtime/debug"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Handler func(args []string, name string, pid *actor.PID, ctx interface{})

type CmdHandler struct {
	handles map[string]Handler
}

func (ch *CmdHandler) Register(typ string, h Handler) {
	if _, ok := ch.handles[typ]; ok {
		log.Panic(typ, "exist")
	}
	ch.handles[typ] = h
}

func (ch *CmdHandler) Handle(data, name string, pid *actor.PID, ctx interface{}) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Println(r)
		}
	}()
	args := strings.Fields(data)
	typ := args[0]
	args = args[1:]
	h := ch.handles[typ]
	h(args, name, pid, ctx)
}

func NewCmdHandler() *CmdHandler {
	return &CmdHandler{
		handles: map[string]Handler{},
	}
}
