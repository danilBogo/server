package app

import (
	"log/slog"
	grpcapp "server/internal/app/grpc"
	"server/internal/services"
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
