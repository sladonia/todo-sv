package main

import (
	"context"
	"fmt"
	"net"

	"github.com/sladonia/common-lb/logger"
	"github.com/sladonia/todo-sv/internal/mongodb"
	"github.com/sladonia/todo-sv/internal/todo"
	"github.com/sladonia/todo-sv/pkg/todopb"
	"go.mongodb.org/mongo-driver/mongo"
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
		log.Panic("create net listener", zap.Error(err))
	}

	return lis
}

func mustCreateZapLogger(logLevel string) *zap.Logger {
	log, err := logger.NewZap(logLevel)
	if err != nil {
		panic(err)
	}

	return log
}

func mustConnectToMongo(ctx context.Context, log *zap.Logger, config Config) *mongo.Database {
	ctx, _ = context.WithTimeout(ctx, config.Mongo.ConnectTimeout)

	db, err := mongodb.Connect(ctx, config.Mongo.DSN, config.Mongo.ToDoDatabaseName)
	if err != nil {
		log.Panic("connect to mongo", zap.Error(err))
	}

	return db
}

func newUserEventDistributor(
	log *zap.Logger,
	config Config,
	publisher todo.Publisher,
	workerSubscriber todo.WorkerSubscriber,
) *todo.UserEventsDistributor {
	return todo.NewUserEventsDistributor(
		config.Nats.UserWorkerGroup,
		todopb.NewProjectSubject("*", "*"),
		workerSubscriber,
		publisher,
		log,
	)
}

func mustCreatePubSub(log *zap.Logger, config Config) todo.PubSub {
	pubSub, err := todo.NewNatsPubSub(config.Nats.DSN)
	if err != nil {
		log.Panic("connect to nats", zap.Error(err))
	}

	return pubSub
}
