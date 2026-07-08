package http

import (
	"net/http"
)

func AssignRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /", h.Webhook)

	mux.HandleFunc("GET /healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
