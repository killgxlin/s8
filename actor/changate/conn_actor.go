package changate

import (
	"bufio"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	kcp "github.com/xtaci/kcp-go"
)

type connActor struct {
	c *kcp.UDPSession
}

func (ca *connActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
		self := ctx.Self()
		go func() {
			defer self.Stop()
			br := bufio.NewReader(ca.c)
			for {
				ca.c.SetReadDeadline(time.Now().Add(time.Second * 5))
				l, _, e := br.ReadLine()
				if e != nil {
					break
				}
				self.Tell(string(l))
			}
		}()
	case *actor.Stopping, *actor.Restarting:
		if ca.c != nil {
			ca.c.Close()
			ca.c = nil
		}
	case string:
		_, e := ca.c.Write([]byte(m + "\n"))
		if e != nil {
			ctx.Self().Stop()
		}
	}
}
