package api

import (
	"github.com/EagleLizard/jcd-api/gosrc/api/ctrl"
	"github.com/EagleLizard/jcd-api/gosrc/lib/config"
	"github.com/go-chi/chi/v5"
)

func addRoutes(
	mux *chi.Mux,
	cfg config.JcdApiConfigType,
) {
	mux.Get("/v1/health", ctrl.GetHealthCheckCtrl)

	mux.Get("/v1/image/{folderKey}/{imageKey}", ctrl.GetImage(
		cfg,
	))
}
