package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/EagleLizard/jcd-api/gosrc/api/ctrl"
	"github.com/EagleLizard/jcd-api/gosrc/lib/logger"
	"github.com/EagleLizard/jcd-api/gosrc/util/chron"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	startServer()
}

func startServer() {
	loggerCfg := logger.Config{
		Encoder: zapcore.NewJSONEncoder(logger.GetDefaultEncoderConfig()),
	}
	logger.Init(loggerCfg)
	defer logger.Close()
	// mux := http.NewServeMux()
	mux := chi.NewRouter()
	/*
		Routes
	*/
	addRoutes(
		mux,
	)
	/*
		Middleware
	*/
	var handler http.Handler = mux
	handler = someMiddleware(handler)
	handler = getAccessLogMiddleware(*logger.Logger)(handler)

	host := "127.0.0.1"
	port := "4444"

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: handler,
	}

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
		logger.Logger.Sugar().Infof("got interrupt signal")
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down server: %s\n", err)
		}
	}()
	wg.Wait()
}

func getAccessLogMiddleware(logger zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := chron.Start()
			h.ServeHTTP(w, r)
			elapsed := sw.Stop()
			logger.Sugar().Infow(
				"[access]",
				"method", r.Method,
				"url", r.URL,
				"duration", float64(elapsed)/float64(time.Millisecond),
			)
		})
	}
}

func someMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := chron.Start()
		h.ServeHTTP(w, r)
		fmt.Printf("%s\n", sw.Stop())
	})
}

func addRoutes(mux *chi.Mux) {
	mux.Get("/v1/health", ctrl.GetHealthCheckCtrl)
}
