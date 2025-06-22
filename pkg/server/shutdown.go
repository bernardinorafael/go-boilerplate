package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (s *Server) GracefulShutdown(ctx context.Context, timeout time.Duration) chan error {
	shutdownErr := make(chan error)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(
			stop,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		<-stop

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		shutdownErr <- s.Shutdown(ctx)
	}()

	return shutdownErr
}
