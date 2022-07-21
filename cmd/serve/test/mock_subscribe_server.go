package test

import (
	"context"

	"github.com/sladonia/todo-sv/pkg/todopb"
	"google.golang.org/grpc"
)

type mockSubscribeServer struct {
	eventCh chan *todopb.Event
	ctx     context.Context
	grpc.ServerStream
}

func newMockSubscribeServer() *mockSubscribeServer {
	return &mockSubscribeServer{
		eventCh: make(chan *todopb.Event),
		ctx:     context.Background(),
	}
}

func (s *mockSubscribeServer) Send(e *todopb.Event) error {
	s.eventCh <- e
	return nil
}

func (s *mockSubscribeServer) Context() context.Context {
	return s.ctx
}
