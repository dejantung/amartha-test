package main

import (
	"billing-engine/internal/payment/app/server"
	"billing-engine/pkg/config"
	"billing-engine/pkg/logger"
)

func main() {
	log := logger.NewZeroLogger("payment")

	cfg, err := config.NewConfig("payment")
	if err != nil {
		log.WithField("error", err).Error("failed to load config")
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
