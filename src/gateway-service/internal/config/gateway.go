package config

const (
	authServiceGRPCEnvName = "AUTH_SERVICE_GRPC"
)

type GateWayConfig interface {
	ServiceRoutes() map[string]string
}

type gatewayConfig struct {
	serviceRoutes map[string]string // маппинг путей к gRPC-сервисам
}

func NewGatewayConfig() (GateWayConfig, error) {
	authServiceGRPC := getEnv(authServiceGRPCEnvName, "auth-service:50051")

	return &gatewayConfig{
		serviceRoutes: map[string]string{
			"/auth/": authServiceGRPC,
		},
	}, nil
}

func (g gatewayConfig) ServiceRoutes() map[string]string {
	return g.serviceRoutes
}
