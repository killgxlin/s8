package main

import (
	"flag"
	"log"
	"s7/share/net"
	"s8/gate"

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
	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	defer cp.Shutdown()

	// cluster
	addr, e := net.FindLanAddr(*cport, *cport+1000)
	if e != nil {
		log.Panic(e)
	}
	cluster.Start("mycluster", addr, cp)

	// gate
	gate.StartGate(*gport, *gport+1000)

	console.ReadLine()
}
