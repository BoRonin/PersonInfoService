package main

import (
	"emtest/internal/app"
	"emtest/internal/config"
	"emtest/pkg/logger"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.LoggerLvl)
	log.Info("starting api service", slog.String("env", cfg.LoggerLvl))
	log.Debug("debug messages enabled")
	app := app.NewApp(*cfg, log)
	if err := app.Serv.ListenAndServe(); err != nil {
		log.Debug("server couldn't run")
		os.Exit(1)
	}
}
