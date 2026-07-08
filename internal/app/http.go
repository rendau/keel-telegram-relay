package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(server *http.Server) *HttpServer {
	server.Handler = HttpMiddlewareContextWithoutCancel(server.Handler)
	server.Handler = HttpMiddlewareRecovery(server.Handler)

	return &HttpServer{
		server: server,
	}
}

func (s *HttpServer) Start() {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http-server stopped: " + err.Error())
			os.Exit(1)
		}
	}()

	slog.Info("http-server started " + s.server.Addr)
}

func (s *HttpServer) Stop(timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	defer ctxCancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("http-server.Shutdown: %w", err)
	}

	return nil
}

func HttpMiddlewareContextWithoutCancel(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r.WithContext(context.WithoutCancel(r.Context())))
	})
}

func HttpMiddlewareRecovery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// use always new err instance in defer
			if err := recover(); err != nil {
				slog.Error("HTTP handler recovered from panic", slog.Any("error", err))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
