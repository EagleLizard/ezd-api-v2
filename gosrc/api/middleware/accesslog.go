package middleware

import (
	"net/http"
	"time"

	"github.com/EagleLizard/jcd-api/gosrc/util/chron"
	"go.uber.org/zap"
)

func NewAccessLogMiddleware(logger *zap.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := chron.Start()
		h.ServeHTTP(w, r)
		elapsed := sw.Stop()
		logger.Sugar().Infow(
			"[access]",
			"method", r.Method,
			"url", r.URL,
			"duration", float32(elapsed)/float32(time.Millisecond),
		)
	})
}
