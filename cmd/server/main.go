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
	"github.com/EwvwGeN/authService/internal/httpHandler"
	l "github.com/EwvwGeN/authService/internal/logger"
	"github.com/EwvwGeN/authService/internal/queue"
	"github.com/EwvwGeN/authService/internal/services/auth"
	"github.com/EwvwGeN/authService/internal/storage"
	tmpl "github.com/EwvwGeN/authService/internal/template"
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

	err = tmpl.PrepareTemplates()
	if err != nil {
		logger.Error("cant prepare templates", slog.String("error", err.Error()))
		panic(err)
	}

	mongoProvider, err := storage.NewMongoProvider(context.Background(), cfg.MongoConfig)
	if err != nil {
		logger.Error("cant get mongo provider", slog.String("error", err.Error()))
		panic(err)
	}

	producer := queue.StartProducer(context.Background(), logger, cfg.RabbitMQCfg)

	authService := auth.NewAuthService(logger, mongoProvider, mongoProvider, mongoProvider, cfg.TokenTTL)

	application := app.NewSerever(logger, cfg.Validator, cfg.Template, authService, producer, cfg.Port)
	go application.GRPCRun()

	httpApp := httpHandler.NewServer(cfg.HttpConfig, logger, nil, nil)
	err = httpApp.Start()
	if err != nil {
		logger.Error("error while starting http server", slog.String("error", err.Error()))
	}

	stopChecker := make(chan os.Signal, 1)
	signal.Notify(stopChecker, syscall.SIGTERM, syscall.SIGINT)
	<- stopChecker
	logger.Info("stopping service")
	application.GRPCStop()
	err = producer.Close()
	if err != nil {
		logger.Error("error while closing producer", slog.String("error", err.Error()))
	}
	err = httpApp.Close(context.Background())
	if err != nil {
		logger.Error("error while closing http server", slog.String("error", err.Error()))
	}
	logger.Info("service stopped")
}