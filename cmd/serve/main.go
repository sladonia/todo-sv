package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"

	"github.com/sladonia/todo-sv/internal/logger"
	"github.com/sladonia/todo-sv/internal/service"
	"github.com/sladonia/todo-sv/internal/todo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		config      = mustLoadConfig()
		log         = logger.MustCreateZapLogger(config.LogLevel)
	)

	defer cancel()
	defer log.Sync()

	var (
		listener       = mustCreateListener(log, config)
		projectStorage = todo.NewInMemoryStorage()
		todoService    = service.NewService(log, projectStorage)
		grpcServer     = newGRPCServer(todoService)
	)

	run(ctx, log, config, grpcServer, listener)

}

func run(ctx context.Context, log *zap.Logger, config Config, grpcServer *grpc.Server, lis net.Listener) {
	errCh := make(chan error)

	go func() {
		log.Info("start listening", zap.String("port", config.Port))
		errCh <- grpcServer.Serve(lis)
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
