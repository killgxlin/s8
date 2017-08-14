package main

import (
	"fmt"
	"io"
	"os"
	"s7/share/util"

	kcp "github.com/xtaci/kcp-go"
)

func main() {
	addr := os.Args[1] + ":" + os.Args[2]
	fmt.Println(addr)
	c, e := kcp.Dial(addr)
	util.PanicOnErr(e)
	defer c.Close()
	go io.Copy(os.Stdout, c)
	io.Copy(c, os.Stdin)
}
