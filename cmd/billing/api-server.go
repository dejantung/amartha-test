package main

import (
	"billing-engine/internal/billing/app/server"
	"billing-engine/pkg/config"
	"billing-engine/pkg/logger"
)

func main() {
	log := logger.NewZeroLogger("billing")

	cfg, err := config.NewConfig("billing")
	if err != nil {
		log.WithField("error", err).Error("failed to load config")
		panic(err)
	}

	log.WithField("config", cfg).Info("config loaded successfully")
	newApiServer, err := server.NewServer(log, cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := newApiServer.Echo.Start(":" + cfg.AppServer.Port); err != nil {
			log.WithField("error", err).Error("failed to start server")
		}
	}()

	newApiServer.Stop()
}
