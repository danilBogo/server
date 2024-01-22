package app

import (
	grpcapp "github.com/danilBogo/server/internal/app/grpc"
	"github.com/danilBogo/server/internal/services"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, port int) *App {
	chat := services.New(log)

	grpcApp := grpcapp.New(chat, log, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
