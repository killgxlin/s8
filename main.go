package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"s7/share/middleware/msglogger"
	"s7/share/profile"
	"s8/node"
	"strings"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"

	_ "s8/node/gate"
	_ "s8/node/match"
	_ "s8/node/room"
	_ "s8/node/term"
	_ "s8/node/user"
)

var (
	port   = flag.Int("port", 2000, "port to listen")
	launch = flag.String("launch", "", "serivice to launch")
)

func main() {
	profile.StartWebTrace()

	msglogger.Enable(true)
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	types := strings.Split(*launch, ",")
	if len(types) == 0 || len(types) == 1 && types[0] == "" {
		types = node.AllLocalNames()
		types = append(types, node.AllGlobalNames()...)
	}

	node.RegToCluster(types)

	p, e := consul.New()
	if e != nil {
		log.Fatal(e)
	}
	cluster.Start("s8", getLanAddr(*port), p)
	node.StartLocals(types)

	for {
		l, e := console.ReadLine()
		if e == io.EOF {
			break
		}
		if l == "" {
			continue
		}

		parts := strings.Split(l, ":")
		if len(parts) == 0 {
			continue
		}

		pid, e := cluster.Get(l, parts[0])
		if e == actor.ErrTimeout {
			continue
		}
		if e != nil || pid == nil {
			log.Fatal(e, pid)
		}
		pid.Tell(&node.Command{Cmd: l})
	}
}

func getLanAddr(port int) string {
	as, e := net.InterfaceAddrs()
	if e != nil {
		log.Fatal(e)
	}
	for _, a := range as {
		ta := a.(*net.IPNet)
		if ta == nil || len(ta.IP.To4()) != net.IPv4len || ta.IP.IsLoopback() {
			continue
		}
		return fmt.Sprintf("%v:%v", ta.IP.String(), port)
	}
	return fmt.Sprintf("%v:%v", "0000", port)
}

/*
cluster:user,room,term,match
local:gate
*/
