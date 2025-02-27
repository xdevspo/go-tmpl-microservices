package config

import "net"

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type HTTPConfig interface {
	Address() string
	Port() string
}

type httpConfig struct {
	host            string
	port            string
	authServiceGRPC string
}

func NewHTTPConfig() (HTTPConfig, error) {
	host := getEnv(httpHostEnvName, "localhost")
	port := getEnv(httpPortEnvName, "8080")

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *httpConfig) Port() string {
	return cfg.port
}
