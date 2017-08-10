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
	cmdHandler = util.NewCmdHandler()
)

func init() {
	cmdHandler.Register("echo", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		msg := fs.String("m", "", "message for echo")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		fmt.Fprintln(out, *msg)
	})
	cmdHandler.Register("name", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
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
	cmdHandler.Register("allch", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		pattern := fs.String("p", "*", "pattern")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		g := registry.GetRegistryGrain("registry")
		r, e := g.List(&registry.ListRequest{Pattern: *pattern})
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		fmt.Fprintln(out, strings.Join(r.Channel, ","))
	})
	cmdHandler.Register("mych", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		pattern := fs.String("p", "*", "pattern")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		u := user.GetUserGrain(ca.name)
		r, e := u.ListChannel(&user.ListChannelRequest{Pattern: *pattern})
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		fmt.Fprintln(out, strings.Join(r.Channel, ","))
	})
	cmdHandler.Register("enter", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		ch := fs.String("c", "", "channel")
		pattern := fs.String("p", "", "pattern")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		if *ch == "" && *pattern == "" {
			fmt.Fprintln(out, "empty channel or pattern")
			return
		}

		chs := strings.Split(*ch, ",")
		selector := &user.ChannelSelector{Channel: chs, Pattern: *pattern}

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
	cmdHandler.Register("quit", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		ch := fs.String("c", "", "channel")
		pattern := fs.String("p", "", "pattern")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		if *ch == "" && *pattern == "" {
			fmt.Fprintln(out, "empty channel or pattern")
			return
		}

		chs := strings.Split(*ch, ",")
		selector := &user.ChannelSelector{Channel: chs, Pattern: *pattern}

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
	cmdHandler.Register("notify", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		ch := fs.String("c", "", "channel")
		pattern := fs.String("p", "", "pattern")
		msg := fs.String("m", "", "message")
		touser := fs.String("u", "", "to which user")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		if *ch == "" && *pattern == "" {
			fmt.Fprintln(out, "empty channel or pattern")
			return
		}

		if *msg == "" {
			fmt.Fprintln(out, *ch, "emtpy message")
			return
		}

		chs := strings.Split(*ch, ",")
		selector := &user.ChannelSelector{Channel: chs, Pattern: *pattern}

		ctx1 := ctx.(actor.Context)
		ca := ctx1.Actor().(*connActor)
		u := user.GetUserGrain(ca.name)

		_, e = u.FireChannelEvent(&user.FireChannelEventRequest{Selector: selector, Msg: *msg, ToUser: *touser})
		fmt.Fprintln(out, chs, e)
	})
}
