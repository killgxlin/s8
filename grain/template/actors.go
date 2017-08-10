package user

import "github.com/AsynkronIT/protoactor-go/cluster"

//a Go struct implementing the User interface
type user struct {
	cluster.Grain
}

func (u *user) ListChannel(*ListChannelRequest) (*ListChannelResponse, error) {
	return &ListChannelResponse{}, nil
}

func (u *user) EnterChannel(*EnterChannelRequest) (*EnterChannelResponse, error) {
	return &EnterChannelResponse{}, nil
}

func (u *user) QuitChannel(*QuitChannelRequest) (*QuitChannelResponse, error) {
	return &QuitChannelResponse{}, nil
}

func init() {
	//apply DI and setup logic
	UserFactory(func() User { return &user{} })
}
