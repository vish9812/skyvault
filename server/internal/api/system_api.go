package api

import (
	"net/http"
	"skyvault/internal/api/helper"
	"skyvault/internal/api/helper/dtos"

	"github.com/go-chi/chi/v5"
)

type SystemAPI struct {
	api *API
}

func NewSystemAPI(a *API) *SystemAPI {
	return &SystemAPI{
		api: a,
	}
}

func (s *SystemAPI) InitRoutes() *SystemAPI {
	router := s.api.v1Pvt
	router.Route("/system", func(r chi.Router) {
		r.Get("/config", s.GetConfig)
	})

	return s
}

func (s *SystemAPI) GetConfig(w http.ResponseWriter, r *http.Request) {
	config := s.api.app.Config.Media

	dto := dtos.SystemConfigDTO{
		MaxUploadSizeMB:       config.MaxUploadSizeMB,
		MaxDirectUploadSizeMB: config.MaxDirectUploadSizeMB,
		MaxChunkSizeMB:        config.MaxChunkSizeMB,
	}

	helper.RespondJSON(w, http.StatusOK, dto)
}