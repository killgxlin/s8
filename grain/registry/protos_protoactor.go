
package registry


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

	
var xRegistryFactory func() Registry

func RegistryFactory(factory func() Registry) {
	xRegistryFactory = factory
}

func GetRegistryGrain(id string) *RegistryGrain {
	return &RegistryGrain{ID: id}
}

type Registry interface {
	Init(id string)
		
	List(*ListRequest) (*ListResponse, error)
		
	Register(*RegisterRequest) (*Unit, error)
		
	Unregister(*UnregisterRequest) (*Unit, error)
		
}
type RegistryGrain struct {
	ID string
}

	
func (g *RegistryGrain) List(r *ListRequest, options ...cluster.GrainCallOption) (*ListResponse, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*ListResponse, error) {
			pid, err := cluster.Get(g.ID, "Registry")
			if err != nil {
				return nil, err
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "List", MessageData: bytes}
			response, err := pid.RequestFuture(request, conf.Timeout).Result()
			if err != nil {
				return nil, err
			}
			switch msg := response.(type) {
			case *cluster.GrainResponse:
				result := &ListResponse{}
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
	
	var res *ListResponse
	var err error
	for i := 0; i < conf.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (g *RegistryGrain) ListChan(r *ListRequest, options ...cluster.GrainCallOption) (<-chan *ListResponse, <-chan error) {
	c := make(chan *ListResponse)
	e := make(chan error)
	go func() {
		res, err := g.List(r, options...)
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
	
func (g *RegistryGrain) Register(r *RegisterRequest, options ...cluster.GrainCallOption) (*Unit, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*Unit, error) {
			pid, err := cluster.Get(g.ID, "Registry")
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

func (g *RegistryGrain) RegisterChan(r *RegisterRequest, options ...cluster.GrainCallOption) (<-chan *Unit, <-chan error) {
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
	
func (g *RegistryGrain) Unregister(r *UnregisterRequest, options ...cluster.GrainCallOption) (*Unit, error) {
	conf := cluster.ApplyGrainCallOptions(options)
	fun := func() (*Unit, error) {
			pid, err := cluster.Get(g.ID, "Registry")
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

func (g *RegistryGrain) UnregisterChan(r *UnregisterRequest, options ...cluster.GrainCallOption) (<-chan *Unit, <-chan error) {
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
	

type RegistryActor struct {
	inner Registry
}

func (a *RegistryActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.inner = xRegistryFactory()
		id := ctx.Self().Id
		a.inner.Init(id[7:len(id)]) //skip "remote$"
	case *cluster.GrainRequest:
		switch msg.Method {
			
		case "List":
			req := &ListRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.List(req)
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
		
		}
	default:
		log.Printf("Unknown message %v", msg)
	}
}

	


func init() {
	
	remote.Register("Registry", actor.FromProducer(func() actor.Actor {
		return &RegistryActor {}
		})		)
		
}



// type registry struct {
//	cluster.Grain
// }

// func (*registry) List(r *ListRequest) (*ListResponse, error) {
// 	return &ListResponse{}, nil
// }

// func (*registry) Register(r *RegisterRequest) (*Unit, error) {
// 	return &Unit{}, nil
// }

// func (*registry) Unregister(r *UnregisterRequest) (*Unit, error) {
// 	return &Unit{}, nil
// }



// func init() {
// 	//apply DI and setup logic

// 	RegistryFactory(func() Registry { return &registry{} })

// }





