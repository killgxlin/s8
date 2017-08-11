
package user


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

	
var xUserFactory func() User

func UserFactory(factory func() User) {
	xUserFactory = factory
}

func GetUserGrain(id string) *UserGrain {
	return &UserGrain{ID: id}
}

type User interface {
	Init(id string)
		
	ListChannel(*ListChannelRequest) (*ListChannelResponse, error)
		
	EnterChannel(*EnterChannelRequest) (*EnterChannelResponse, error)
		
	QuitChannel(*QuitChannelRequest) (*Unit, error)
		
	NotifyChannelEvent(*NotifyChannelEventRequest) (*Unit, error)
		
	Register(*RegisterRequest) (*Unit, error)
		
	Unregister(*UnregisterRequest) (*Unit, error)
		
	OnChannelEvent(*ChannelEventRequest) (*Unit, error)
		
}
type UserGrain struct {
	ID string
}

	
func (g *UserGrain) ListChannel(r *ListChannelRequest, options ...cluster.GrainCallOption) (*ListChannelResponse, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*ListChannelResponse, error) {
			pid, err := cluster.Get(g.ID, "User")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "ListChannel", MessageData: bytes}
			response, err := pid.RequestFuture(request, conf.Timeout).Result()
			if err != nil {
				return nil, err
			}
			switch msg := response.(type) {
			case *cluster.GrainResponse:
				result := &ListChannelResponse{}
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
	
	var res *ListChannelResponse
	var err error
	for i := 0; i < conf.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (g *UserGrain) ListChannelChan(r *ListChannelRequest, options ...cluster.GrainCallOption) (<-chan *ListChannelResponse, <-chan error) {
	c := make(chan *ListChannelResponse)
	e := make(chan error)
	go func() {
		res, err := g.ListChannel(r, options...)
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
	
func (g *UserGrain) EnterChannel(r *EnterChannelRequest, options ...cluster.GrainCallOption) (*EnterChannelResponse, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*EnterChannelResponse, error) {
			pid, err := cluster.Get(g.ID, "User")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "EnterChannel", MessageData: bytes}
			response, err := pid.RequestFuture(request, conf.Timeout).Result()
			if err != nil {
				return nil, err
			}
			switch msg := response.(type) {
			case *cluster.GrainResponse:
				result := &EnterChannelResponse{}
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
	
	var res *EnterChannelResponse
	var err error
	for i := 0; i < conf.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (g *UserGrain) EnterChannelChan(r *EnterChannelRequest, options ...cluster.GrainCallOption) (<-chan *EnterChannelResponse, <-chan error) {
	c := make(chan *EnterChannelResponse)
	e := make(chan error)
	go func() {
		res, err := g.EnterChannel(r, options...)
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
	
func (g *UserGrain) QuitChannel(r *QuitChannelRequest, options ...cluster.GrainCallOption) (*Unit, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*Unit, error) {
			pid, err := cluster.Get(g.ID, "User")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "QuitChannel", MessageData: bytes}
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

func (g *UserGrain) QuitChannelChan(r *QuitChannelRequest, options ...cluster.GrainCallOption) (<-chan *Unit, <-chan error) {
	c := make(chan *Unit)
	e := make(chan error)
	go func() {
		res, err := g.QuitChannel(r, options...)
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
	
func (g *UserGrain) NotifyChannelEvent(r *NotifyChannelEventRequest, options ...cluster.GrainCallOption) (*Unit, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*Unit, error) {
			pid, err := cluster.Get(g.ID, "User")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "NotifyChannelEvent", MessageData: bytes}
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

func (g *UserGrain) NotifyChannelEventChan(r *NotifyChannelEventRequest, options ...cluster.GrainCallOption) (<-chan *Unit, <-chan error) {
	c := make(chan *Unit)
	e := make(chan error)
	go func() {
		res, err := g.NotifyChannelEvent(r, options...)
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
	
func (g *UserGrain) Register(r *RegisterRequest, options ...cluster.GrainCallOption) (*Unit, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*Unit, error) {
			pid, err := cluster.Get(g.ID, "User")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "Register", MessageData: bytes}
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

func (g *UserGrain) RegisterChan(r *RegisterRequest, options ...cluster.GrainCallOption) (<-chan *Unit, <-chan error) {
	c := make(chan *Unit)
	e := make(chan error)
	go func() {
		res, err := g.Register(r, options...)
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
	
func (g *UserGrain) Unregister(r *UnregisterRequest, options ...cluster.GrainCallOption) (*Unit, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*Unit, error) {
			pid, err := cluster.Get(g.ID, "User")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "Unregister", MessageData: bytes}
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

func (g *UserGrain) UnregisterChan(r *UnregisterRequest, options ...cluster.GrainCallOption) (<-chan *Unit, <-chan error) {
	c := make(chan *Unit)
	e := make(chan error)
	go func() {
		res, err := g.Unregister(r, options...)
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
	
func (g *UserGrain) OnChannelEvent(r *ChannelEventRequest, options ...cluster.GrainCallOption) (*Unit, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*Unit, error) {
			pid, err := cluster.Get(g.ID, "User")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "OnChannelEvent", MessageData: bytes}
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

func (g *UserGrain) OnChannelEventChan(r *ChannelEventRequest, options ...cluster.GrainCallOption) (<-chan *Unit, <-chan error) {
	c := make(chan *Unit)
	e := make(chan error)
	go func() {
		res, err := g.OnChannelEvent(r, options...)
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
	

type UserActor struct {
	inner User
}

func (a *UserActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.inner = xUserFactory()
		id := ctx.Self().Id
		a.inner.Init(id[7:len(id)]) //skip "remote$"
	case *cluster.GrainRequest:
		switch msg.Method {
			
		case "ListChannel":
			req := &ListChannelRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.ListChannel(req)
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
			
		case "EnterChannel":
			req := &EnterChannelRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.EnterChannel(req)
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
			
		case "QuitChannel":
			req := &QuitChannelRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.QuitChannel(req)
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
			
		case "NotifyChannelEvent":
			req := &NotifyChannelEventRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.NotifyChannelEvent(req)
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
			
		case "Register":
			req := &RegisterRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.Register(req)
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
			
		case "Unregister":
			req := &UnregisterRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.Unregister(req)
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
			
		case "OnChannelEvent":
			req := &ChannelEventRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.OnChannelEvent(req)
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
	
	remote.Register("User", actor.FromProducer(func() actor.Actor {
		return &UserActor {}
		})		)
		
}



// type user struct {
//	cluster.Grain
// }

// func (*user) ListChannel(r *ListChannelRequest) (*ListChannelResponse, error) {
// 	return &ListChannelResponse{}, nil
// }

// func (*user) EnterChannel(r *EnterChannelRequest) (*EnterChannelResponse, error) {
// 	return &EnterChannelResponse{}, nil
// }

// func (*user) QuitChannel(r *QuitChannelRequest) (*Unit, error) {
// 	return &Unit{}, nil
// }

// func (*user) NotifyChannelEvent(r *NotifyChannelEventRequest) (*Unit, error) {
// 	return &Unit{}, nil
// }

// func (*user) Register(r *RegisterRequest) (*Unit, error) {
// 	return &Unit{}, nil
// }

// func (*user) Unregister(r *UnregisterRequest) (*Unit, error) {
// 	return &Unit{}, nil
// }

// func (*user) OnChannelEvent(r *ChannelEventRequest) (*Unit, error) {
// 	return &Unit{}, nil
// }



// func init() {
// 	//apply DI and setup logic

// 	UserFactory(func() User { return &user{} })

// }





