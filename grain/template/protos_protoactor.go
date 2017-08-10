
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
		
	QuitChannel(*QuitChannelRequest) (*QuitChannelResponse, error)
		
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
	
func (g *UserGrain) QuitChannel(r *QuitChannelRequest, options ...cluster.GrainCallOption) (*QuitChannelResponse, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*QuitChannelResponse, error) {
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
				result := &QuitChannelResponse{}
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
	
	var res *QuitChannelResponse
	var err error
	for i := 0; i < conf.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (g *UserGrain) QuitChannelChan(r *QuitChannelRequest, options ...cluster.GrainCallOption) (<-chan *QuitChannelResponse, <-chan error) {
	c := make(chan *QuitChannelResponse)
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

// func (*user) QuitChannel(r *QuitChannelRequest) (*QuitChannelResponse, error) {
// 	return &QuitChannelResponse{}, nil
// }



// func init() {
// 	//apply DI and setup logic

// 	UserFactory(func() User { return &user{} })

// }





