package gate

import (
	"flag"
	"fmt"
	"io"
	"s8/grain/registry"
	"s8/grain/user"
	"s8/util"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// handler --------------------------------------------------
var (
	handler = util.NewCmdHandler()
)

func init() {
	handler.Register("echo", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		msg := fs.String("m", "", "message for echo")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		fmt.Fprintln(out, *msg)
	})
	handler.Register("name", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)

		name := fs.String("n", "", "name")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		if ca.name == *name {
			fmt.Fprintln(out, "same name", ca.name)
			return
		}

		if ca.name != "" {
			u := user.GetUserGrain(ca.name)
			_, e := u.Unregister(&user.UnregisterRequest{Pid: ctx1.Self()})
			if e != nil {
				fmt.Fprintln(out, e)
				return
			}
			ca.name = ""
		}

		if *name != "" {
			u := user.GetUserGrain(*name)
			_, e := u.Register(&user.RegisterRequest{Pid: ctx1.Self()})
			if e != nil {
				fmt.Fprintln(out, e)
				return
			}
			ca.name = *name
		}

		fmt.Fprintln(out, "ok")
	})
	handler.Register("allch", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		patt := fs.String("p", "*", "pattern")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		g := registry.GetRegistryGrain("registry")
		r, e := g.List(&registry.ListRequest{Pattern: *patt})
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		fmt.Fprintln(out, strings.Join(r.Channel, ","))
	})
	handler.Register("mych", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		patt := fs.String("p", "*", "pattern")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		u := user.GetUserGrain(ca.name)
		r, e := u.ListChannel(&user.ListChannelRequest{Pattern: *patt})
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		fmt.Fprintln(out, strings.Join(r.Channel, ","))
	})
	handler.Register("enter", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		ch := fs.String("c", "", "channel")
		patt := fs.String("p", "", "pattern")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		if *ch == "" && *patt == "" {
			fmt.Fprintln(out, "empty channel or pattern")
			return
		}

		chs := strings.Split(*ch, ",")
		selector := &user.ChannelSelector{Channel: chs, Pattern: *patt}

		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		u := user.GetUserGrain(ca.name)

		r, e := u.EnterChannel(&user.EnterChannelRequest{Selector: selector})
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		for i, r1 := range r.Channels {
			if i > 0 {
				fmt.Fprint(out, "|")
			}
			fmt.Fprint(out, r1.Channel, " ", r1.Err, " ", strings.Join(r1.Members, ","))
		}
		fmt.Fprintln(out)
	})
	handler.Register("quit", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		ch := fs.String("c", "", "channel")
		patt := fs.String("p", "", "pattern")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		if *ch == "" && *patt == "" {
			fmt.Fprintln(out, "empty channel or pattern")
			return
		}

		chs := strings.Split(*ch, ",")
		selector := &user.ChannelSelector{Channel: chs, Pattern: *patt}

		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		u := user.GetUserGrain(ca.name)

		_, e = u.QuitChannel(&user.QuitChannelRequest{Selector: selector})
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		fmt.Fprintln(out, *ch, "ok")
	})
	handler.Register("notify", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		ch := fs.String("c", "", "channel")
		patt := fs.String("p", "", "pattern")
		msg := fs.String("m", "", "message")
		touser := fs.String("u", "", "to which user")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		if *ch == "" && *patt == "" {
			fmt.Fprintln(out, "empty channel or pattern")
			return
		}

		if *msg == "" {
			fmt.Fprintln(out, *ch, "emtpy message")
			return
		}

		chs := strings.Split(*ch, ",")
		selector := &user.ChannelSelector{Channel: chs, Pattern: *patt}

		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		u := user.GetUserGrain(ca.name)

		_, e = u.NotifyChannelEvent(&user.NotifyChannelEventRequest{Selector: selector, Msg: *msg, ToUser: *touser})
		fmt.Fprintln(out, chs, e)
	})
	handler.Register("help", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		cmds := handler.AllCommand()
		for _, cmd := range cmds {
			fmt.Fprintln(out, cmd)
		}
	})
}
