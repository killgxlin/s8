
package channel


import errors "errors"
import log "log"
import actor "github.com/AsynkronIT/protoactor-go/actor"
import remote "github.com/AsynkronIT/protoactor-go/remote"
import cluster "github.com/AsynkronIT/protoactor-go/cluster"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

	
var xChannelFactory func() Channel

func ChannelFactory(factory func() Channel) {
	xChannelFactory = factory
}

func GetChannelGrain(id string) *ChannelGrain {
	return &ChannelGrain{ID: id}
}

type Channel interface {
	Init(id string)
		
	Enter(*EnterRequest) (*EnterResponse, error)
		
	Quit(*QuitRequest) (*Unit, error)
		
	Publish(*PublishRequest) (*Unit, error)
		
}
type ChannelGrain struct {
	ID string
}

	
func (g *ChannelGrain) Enter(r *EnterRequest, options ...cluster.GrainCallOption) (*EnterResponse, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*EnterResponse, error) {
			pid, err := cluster.Get(g.ID, "Channel")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "Enter", MessageData: bytes}
			response, err := pid.RequestFuture(request, conf.Timeout).Result()
			if err != nil {
				return nil, err
			}
			switch msg := response.(type) {
			case *cluster.GrainResponse:
				result := &EnterResponse{}
				err = proto.Unmarshal(msg.MessageData, result)
				if err != nil {
					return nil, err
				}
				return result, nil
			case *cluster.GrainErrorResponse:
				return nil, errors.New(msg.Err)
			default:
				return nil, errors.New("Unknown response")
			}
		}
	
	var res *EnterResponse
	var err error
	for i := 0; i < conf.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (g *ChannelGrain) EnterChan(r *EnterRequest, options ...cluster.GrainCallOption) (<-chan *EnterResponse, <-chan error) {
	c := make(chan *EnterResponse)
	e := make(chan error)
	go func() {
		res, err := g.Enter(r, options...)
		if err != nil {
			e <- err
		} else {
			c <- res
		}
		close(c)
		close(e)
	}()
	return c, e
}
	
func (g *ChannelGrain) Quit(r *QuitRequest, options ...cluster.GrainCallOption) (*Unit, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*Unit, error) {
			pid, err := cluster.Get(g.ID, "Channel")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "Quit", MessageData: bytes}
			response, err := pid.RequestFuture(request, conf.Timeout).Result()
			if err != nil {
				return nil, err
			}
			switch msg := response.(type) {
			case *cluster.GrainResponse:
				result := &Unit{}
				err = proto.Unmarshal(msg.MessageData, result)
				if err != nil {
					return nil, err
				}
				return result, nil
			case *cluster.GrainErrorResponse:
				return nil, errors.New(msg.Err)
			default:
				return nil, errors.New("Unknown response")
			}
		}
	
	var res *Unit
	var err error
	for i := 0; i < conf.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (g *ChannelGrain) QuitChan(r *QuitRequest, options ...cluster.GrainCallOption) (<-chan *Unit, <-chan error) {
	c := make(chan *Unit)
	e := make(chan error)
	go func() {
		res, err := g.Quit(r, options...)
		if err != nil {
			e <- err
		} else {
			c <- res
		}
		close(c)
		close(e)
	}()
	return c, e
}
	
func (g *ChannelGrain) Publish(r *PublishRequest, options ...cluster.GrainCallOption) (*Unit, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*Unit, error) {
			pid, err := cluster.Get(g.ID, "Channel")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "Publish", MessageData: bytes}
			response, err := pid.RequestFuture(request, conf.Timeout).Result()
			if err != nil {
				return nil, err
			}
			switch msg := response.(type) {
			case *cluster.GrainResponse:
				result := &Unit{}
				err = proto.Unmarshal(msg.MessageData, result)
				if err != nil {
					return nil, err
				}
				return result, nil
			case *cluster.GrainErrorResponse:
				return nil, errors.New(msg.Err)
			default:
				return nil, errors.New("Unknown response")
			}
		}
	
	var res *Unit
	var err error
	for i := 0; i < conf.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (g *ChannelGrain) PublishChan(r *PublishRequest, options ...cluster.GrainCallOption) (<-chan *Unit, <-chan error) {
	c := make(chan *Unit)
	e := make(chan error)
	go func() {
		res, err := g.Publish(r, options...)
		if err != nil {
			e <- err
		} else {
			c <- res
		}
		close(c)
		close(e)
	}()
	return c, e
}
	

type ChannelActor struct {
	inner Channel
}

func (a *ChannelActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.inner = xChannelFactory()
		id := ctx.Self().Id
		a.inner.Init(id[7:len(id)]) //skip "remote$"
	case *cluster.GrainRequest:
		switch msg.Method {
			
		case "Enter":
			req := &EnterRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.Enter(req)
			if err == nil {
				bytes, err := proto.Marshal(r0)
				if err != nil {
					log.Fatalf("[GRAIN] proto.Marshal failed %v", err)
				}
				resp := &cluster.GrainResponse{MessageData: bytes}
				ctx.Respond(resp)
			} else {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
			}
			
		case "Quit":
			req := &QuitRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.Quit(req)
			if err == nil {
				bytes, err := proto.Marshal(r0)
				if err != nil {
					log.Fatalf("[GRAIN] proto.Marshal failed %v", err)
				}
				resp := &cluster.GrainResponse{MessageData: bytes}
				ctx.Respond(resp)
			} else {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
			}
			
		case "Publish":
			req := &PublishRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.Publish(req)
			if err == nil {
				bytes, err := proto.Marshal(r0)
				if err != nil {
					log.Fatalf("[GRAIN] proto.Marshal failed %v", err)
				}
				resp := &cluster.GrainResponse{MessageData: bytes}
				ctx.Respond(resp)
			} else {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
			}
		
		}
	default:
		log.Printf("Unknown message %v", msg)
	}
}

	


func init() {
	
	remote.Register("Channel", actor.FromProducer(func() actor.Actor {
		return &ChannelActor {}
		})		)
		
}



// type channel struct {
//	cluster.Grain
// }

// func (*channel) Enter(r *EnterRequest) (*EnterResponse, error) {
// 	return &EnterResponse{}, nil
// }

// func (*channel) Quit(r *QuitRequest) (*Unit, error) {
// 	return &Unit{}, nil
// }

// func (*channel) Publish(r *PublishRequest) (*Unit, error) {
// 	return &Unit{}, nil
// }



// func init() {
// 	//apply DI and setup logic

// 	ChannelFactory(func() Channel { return &channel{} })

// }





