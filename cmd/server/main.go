package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/EwvwGeN/authService/internal/app"
	c "github.com/EwvwGeN/authService/internal/config"
	l "github.com/EwvwGeN/authService/internal/logger"
	"github.com/EwvwGeN/authService/internal/services/auth"
	"github.com/EwvwGeN/authService/internal/storage"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "path to config file")
}

func main() {
	flag.Parse()
	cfg, err := c.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("cant load config from path %s: %s", configPath, err.Error()))
	}
	logger := l.SetupLogger(cfg.LogLevel)

	logger.Info("config loaded")
	logger.Debug("config data", slog.Any("cfg", cfg))

	mongoProvider, err := storage.NewMongoProvider(context.Background(), cfg.MongoConfig)
	if err != nil {
		logger.Error("cant get mongo provider", slog.String("error", err.Error()))
		panic(err)
	}

	authService := auth.NewAuthService(logger, mongoProvider, mongoProvider, mongoProvider, cfg.TokenTTL)

	application := app.NewSerever(logger, cfg.Validator, authService, cfg.Port)
	go application.GRPCRun()

	stopChecker := make(chan os.Signal, 1)
	signal.Notify(stopChecker, syscall.SIGTERM, syscall.SIGINT)
	<- stopChecker
	logger.Info("stopping service")
	application.GRPCStop()
	logger.Info("service stopped")
}