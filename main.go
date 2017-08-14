package main

import (
	"flag"
	"log"
	"s7/share/net"
	"s8/actor/changate"
	"s8/actor/gate"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
)

var (
	cport = flag.Int("cport", 8000, "cluster port")
	gport = flag.Int("gport", 9000, "gate port")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	flag.Parse()

	// consul
	cp, e := consul.New()
	if e != nil {
		log.Fatal(e)
	}
	defer cp.Shutdown()

	// cluster
	addr, e := net.FindLanAddr("tcp", *cport, *cport+1000)
	if e != nil {
		log.Panic(e)
	}
	cluster.Start("mycluster", addr, cp)

	// gate
	gate.Start(*gport, *gport+1000)
	// channel gate
	changate.Start(*gport, *gport+1000)

	console.ReadLine()
}
