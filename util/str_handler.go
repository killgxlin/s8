package util

import (
	"log"
	"strings"
)

type Handler func(args []string, ctx interface{})

type CmdHandler struct {
	handles map[string]Handler
}

func (ch *CmdHandler) Register(typ string, h Handler) {
	if _, ok := ch.handles[typ]; ok {
		log.Panic(typ, "exist")
	}
	ch.handles[typ] = h
}

func (ch *CmdHandler) Handle(cmd string, ctx interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	args := strings.Fields(cmd)
	typ := args[0]
	args = args[1:]
	h := ch.handles[typ]
	h(args, ctx)
}

func NewCmdHandler() *CmdHandler {
	return &CmdHandler{
		handles: map[string]Handler{},
	}
}
