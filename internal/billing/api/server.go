package api

import (
	"billing-engine/internal/billing/service"
	"billing-engine/pkg/logger"
)

type AppServer struct {
	BillingService service.BillingServiceProvider
	log            logger.Logger
}

func NewAppServer() *AppServer {
	return &AppServer{}
}
