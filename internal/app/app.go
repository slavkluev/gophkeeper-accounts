package app

import (
	"time"

	"go.uber.org/zap"

	grpcapp "accounts/internal/app/grpc"
	"accounts/internal/service/accounts"
	"accounts/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *zap.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	secret string,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	accountsService := accounts.New(log, storage, storage, storage, tokenTTL)

	grpcApp, err := grpcapp.New(log, accountsService, grpcPort, secret)
	if err != nil {
		panic(err)
	}

	return &App{
		GRPCServer: grpcApp,
	}
}
