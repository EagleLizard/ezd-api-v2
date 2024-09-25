package api

import (
	"github.com/EagleLizard/jcd-api/gosrc/api/ctrl"
	"github.com/go-chi/chi/v5"
)

func addRoutes(
	mux *chi.Mux,
) {
	mux.Get("/v1/health", ctrl.GetHealthCheckCtrl)
}
