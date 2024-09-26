package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/EagleLizard/jcd-api/gosrc/api/middleware"
	"github.com/EagleLizard/jcd-api/gosrc/lib/config"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func InitServer(
	cfg *config.JcdApiConfigType,
	logger *zap.Logger,
) http.Handler {
	mux := chi.NewRouter()

	/* Routes */
	addRoutes(
		mux,
		*cfg,
	)

	/* Middleware */
	var handler http.Handler = mux
	handler = middleware.NewAccessLogMiddleware(logger, handler)

	return handler
}

func RunServer(ctx context.Context, httpServer *http.Server) {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		fmt.Fprintf(os.Stdout, "listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe error: %s\n", err)
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		fmt.Fprintf(os.Stdout, "got interrupt signal\n")
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down server: %s\n", err)
		}
	}()
	wg.Wait()
}
