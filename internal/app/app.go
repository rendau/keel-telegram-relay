package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rendau/keel-telegram-relay/internal/config"
	handlerHttpP "github.com/rendau/keel-telegram-relay/internal/handler/http"
	serviceTelegramP "github.com/rendau/keel-telegram-relay/internal/service/telegram"
)

type App struct {
	httpServer *HttpServer

	exitCode int
}

func (a *App) Init() {
	// logger
	{
		if !config.Conf.Debug {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			slog.SetDefault(logger)
		}
	}

	if config.Conf.TelegramToken == "" {
		slog.Error("TELEGRAM_TOKEN is required")
		os.Exit(1)
	}
	if config.Conf.TelegramChatId == "" {
		slog.Error("TELEGRAM_CHAT_ID is required")
		os.Exit(1)
	}

	serviceTelegram := serviceTelegramP.New(
		config.Conf.TelegramApiUrl,
		config.Conf.TelegramToken,
		config.Conf.TelegramChatId,
	)

	// http server
	{
		mux := http.NewServeMux()

		handlerHttpP.AssignRoutes(mux, handlerHttpP.New(serviceTelegram))

		a.httpServer = NewHttpServer(&http.Server{
			Addr:              ":" + config.Conf.HttpPort,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			MaxHeaderBytes:    300 * 1024,
		})
	}
}

func (a *App) Start() {
	slog.Info("Starting")

	a.httpServer.Start()
}

func (a *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer signalCtxCancel()

	// wait signal
	<-signalCtx.Done()
}

func (a *App) Stop() {
	slog.Info("Shutting down...")

	if err := a.httpServer.Stop(15 * time.Second); err != nil {
		slog.Error("httpServer.Stop", "error", err)
		a.exitCode = 1
	}
}

func (a *App) Exit() {
	slog.Info("Exit")

	os.Exit(a.exitCode)
}
