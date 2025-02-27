package app

import (
	"github.com/xdevspo/go-tmpl-microservices/gateway-service/internal/config"
	"log"
)

type serviceProvider struct {
	httpConfig    config.HTTPConfig
	gatewayConfig config.GatewayConfig
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) HTTPConfig() config.HTTPConfig {
	if sp.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		sp.httpConfig = cfg
	}

	return sp.httpConfig
}

func (sp *serviceProvider) GatewayConfig() config.GateWayConfig {
	if sp.httpConfig == nil {
		cfg, err := config.NewGatewayConfig()
		if err != nil {
			log.Fatalf("failed to get gateway config: %s", err.Error())
		}

		sp.gatewayConfig = cfg
	}

	return sp.gatewayConfig
}
