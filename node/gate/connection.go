package gate

import (
	"log"
	"s7/share/middleware/mnet"
	cstring "s7/share/net/coder/string"
	strings "strings"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/firstrow/tcp_server"
)

func startServer(ctx actor.Context) {
	cstring.NewReadWriter()
	mnet.MakeAcceptor("0.0.0.0:0", 100, 0)
}

func handleMessage(c *tcp_server.Client, msg string) {
	log.Println(c, msg)
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	cmd := strings.Split(msg, " ")
	switch cmd[0] {
	case "csgetuser":
	case "cscreateuser":
	case "csmatch":
	}
}
