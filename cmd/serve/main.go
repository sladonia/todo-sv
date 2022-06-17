package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"

	"github.com/sladonia/todo-sv/internal/todo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		config      = mustLoadConfig()
		log         = mustCreateZapLogger(config.LogLevel)
	)

	defer cancel()
	defer log.Sync()

	var (
		listener         = mustCreateListener(log, config)
		db               = mustConnectToMongo(ctx, log, config)
		pubSub           = mustCreatePubSub(log, config)
		eventDistributor = newUserEventDistributor(log, config, pubSub, pubSub)
		projectStorage   = todo.NewStorage(db, config.Mongo.ProjectsCollectionName)
		todoService      = todo.NewService(log, projectStorage, pubSub)
		grpcServer       = newGRPCServer(todoService)
	)

	run(ctx, log, config, grpcServer, listener, eventDistributor)

}

func run(
	ctx context.Context,
	log *zap.Logger,
	config Config,
	grpcServer *grpc.Server,
	lis net.Listener,
	eventDistributor *todo.UserEventsDistributor,
) {
	errCh := make(chan error)

	go func() {
		log.Info("start listening grpc server", zap.String("port", config.Port))
		errCh <- grpcServer.Serve(lis)
	}()

	go func() {
		log.Info("start event distributor")
		errCh <- eventDistributor.Start(ctx)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, _ := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		doneCh := make(chan struct{})

		go func() {
			grpcServer.GracefulStop()
			close(doneCh)
		}()

		select {
		case <-doneCh:
			log.Info("shut down gracefully")
			return
		case <-shutdownCtx.Done():
			log.Error("failed to shut down gracefully. shutting down immediately")
			grpcServer.Stop()
			return
		}
	case err, ok := <-errCh:
		if !ok {
			return
		}

		log.Error("application level error", zap.Error(err))
	}
}
