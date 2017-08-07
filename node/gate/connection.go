package gate

import (
	"log"
	strings "strings"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/firstrow/tcp_server"
)

func startServer(ctx actor.Context) {
	server := tcp_server.New("0.0.0.0:0")
	server.OnNewClient(func(c *tcp_server.Client) {
	})
	server.OnNewMessage(handleMessage)
	server.OnClientConnectionClosed(func(c *tcp_server.Client, e error) {
	})
	go server.Listen()
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
