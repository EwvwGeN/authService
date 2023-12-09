package main

import (
	"flag"
	"fmt"
	"log/slog"

	c "github.com/EwvwGeN/authService/internal/config"
	l "github.com/EwvwGeN/authService/internal/logger"
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
		panic(fmt.Sprintf("Cant load config from path %s", configPath ))
	}
	logger := l.SetupLogger(cfg.LogLevel)
	logger.Debug("Config loaded", slog.Any("cfg", cfg))
}