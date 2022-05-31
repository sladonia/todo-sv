package main

import (
	"fmt"
	"net"

	"github.com/sladonia/todo-sv/pkg/todopb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func newGRPCServer(todoService todopb.ToDoServiceServer) *grpc.Server {
	s := grpc.NewServer()

	todopb.RegisterToDoServiceServer(s, todoService)

	return s
}

func mustCreateListener(log *zap.Logger, config Config) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", config.Port))
	if err != nil {
		log.Panic("failed to create net listener", zap.Error(err))
	}

	return lis
}
