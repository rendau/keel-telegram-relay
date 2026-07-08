package http

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/rendau/keel-telegram-relay/internal/service/keel"
	"github.com/rendau/keel-telegram-relay/internal/service/telegram"
)

type Handler struct {
	telegram *telegram.Client
}

func New(telegram *telegram.Client) *Handler {
	return &Handler{
		telegram: telegram,
	}
}

// Webhook receives a keel notification (generic webhook) and forwards it to telegram.
func (h *Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("fail to read body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event, err := keel.Parse(body)
	if err != nil {
		slog.Error("fail to parse keel event",
			slog.String("error", err.Error()),
			slog.String("body", string(body)),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Info("keel event",
		"type", event.Type,
		"level", event.Level,
		"name", event.Name,
		"message", event.Message,
	)

	if err = h.telegram.Send(r.Context(), event.Text()); err != nil {
		slog.Error("fail to send telegram", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
}
