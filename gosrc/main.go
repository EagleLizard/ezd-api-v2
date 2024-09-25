package main

import (
	"context"
	"net"
	"net/http"

	"github.com/EagleLizard/jcd-api/gosrc/api"
	"github.com/EagleLizard/jcd-api/gosrc/lib/config"
	"github.com/EagleLizard/jcd-api/gosrc/lib/logging"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	cfg := config.JcdApiConfig
	startServer(ctx, cfg)
}

func startServer(ctx context.Context, cfg *config.JcdApiConfigType) {
	loggerCfg := logging.Config{
		Encoder: zapcore.NewJSONEncoder(logging.GetDefaultEncoderConfig()),
	}
	logging.Init(loggerCfg)
	defer logging.Close()

	srv := api.InitServer(cfg, logging.Logger)
	httpServer := &http.Server{
		Addr: net.JoinHostPort(
			cfg.Host,
			cfg.Port,
		),
		Handler: srv,
	}
	api.RunServer(ctx, httpServer)
}
