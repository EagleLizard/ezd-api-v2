package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/EagleLizard/jcd-api/gosrc/lib/config"
	"github.com/EagleLizard/jcd-api/gosrc/util/chron"
	"github.com/EagleLizard/jcd-api/gosrc/util/constants"
	"github.com/go-chi/chi/v5"
)

/*
	`sfs` stands for "Simple File Server"
*/

func main() {
	ctx := context.Background()
	cfg := config.JcdApiConfig
	startServer(ctx, cfg)
}

func startServer(ctx context.Context, cfg *config.JcdApiConfigType) {
	// port := "4041"
	// host := "0.0.0.0"
	host := cfg.SfsHost
	port := cfg.SfsPort
	staticDirPath := filepath.Join(constants.LocalDir, "static")

	fileHandler := http.FileServer(http.Dir(staticDirPath))

	mux := chi.NewRouter()
	mux.Get("/*", fileHandler.ServeHTTP)

	var handler http.Handler = mux
	handler = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := chron.Start()
			h.ServeHTTP(w, r)
			elapsed := sw.Stop()
			fmt.Fprintf(os.Stdout, "%s %s %f\n", r.Method, r.URL, float32(elapsed)/float32(time.Millisecond))
		})
	}(handler)

	fileServer := &http.Server{
		Addr: net.JoinHostPort(
			host,
			port,
		),
		Handler: handler,
	}

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		fmt.Fprintf(os.Stdout, "listening on %s\n", fileServer.Addr)
		if err := fileServer.ListenAndServe(); err != nil {
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

		if err := fileServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down server: %s\n", err)
		}
	}()
	wg.Wait()
}
