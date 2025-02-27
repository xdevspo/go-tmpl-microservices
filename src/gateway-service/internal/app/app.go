package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	userspb "github.com/xdevspo/go-tmpl-microservices/auth-service/pkg/users"
	"github.com/xdevspo/go-tmpl-microservices/gateway-service/internal/closer"
	"github.com/xdevspo/go-tmpl-microservices/gateway-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.Init(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runHTTPServer()
}

func (a *App) Init(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
		a.initGatewayServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	configFile := fmt.Sprintf(".env.%s", env)
	err := config.Load(configFile)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

type GatewayServer struct {
	config      config.GateWayConfig
	mux         *runtime.ServeMux
	connections map[string]*grpc.ClientConn
}

func NewGatewayServer(config config.GateWayConfig) *GatewayServer {
	return &GatewayServer{
		config:      config,
		connections: make(map[string]*grpc.ClientConn),
	}
}

// Установка соединений с gRPC сервисами
func (s *GatewayServer) Connect() error {
	for path, addr := range s.config.ServiceRoutes {
		// Создаем контекст с таймаутом для установки соединения
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Создаем новое соединение с gRPC сервисом
		conn, err := grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(16*1024*1024),
				grpc.MaxCallSendMsgSize(16*1024*1024),
			),
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		)
		if err != nil {
			s.Close()
			return fmt.Errorf("failed to create connection to service %s at %s: %v", path, addr, err)
		}

		// Инициируем соединение
		conn.Connect()

		// Проверяем статус соединения
		state := conn.GetState()
		for state != connectivity.Ready {
			if !conn.WaitForStateChange(ctx, state) {
				err := conn.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("connection timeout for service %s at %s", path, addr)
			}
			state = conn.GetState()
			if state == connectivity.TransientFailure || state == connectivity.Shutdown {
				err := conn.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("failed to establish connection to service %s at %s", path, addr)
			}
		}

		s.connections[addr] = conn
		log.Printf("Successfully connected to service: %s at %s", path, addr)

		// Регистрируем handlers для сервиса
		if err := s.registerServiceHandlers(ctx, path, conn); err != nil {
			return fmt.Errorf("failed to register handlers for service %s: %v", path, err)
		}
	}
	return nil
}

// registerServiceHandlers регистрирует обработчики для конкретного сервиса
func (s *GatewayServer) registerServiceHandlers(ctx context.Context, path string, conn *grpc.ClientConn) error {
	switch {
	case strings.HasPrefix(path, "/auth/"):
		// Используем сгенерированный код для регистрации auth service
		if err := userspb.RegisterAuthServiceHandler(ctx, s.mux, conn); err != nil {
			return fmt.Errorf("failed to register auth service handler: %v", err)
		}
	case strings.HasPrefix(path, "/org/"):
		// Используем сгенерированный код для регистрации org service
		if err := orgpb.RegisterOrgServiceHandler(ctx, s.mux, conn); err != nil {
			return fmt.Errorf("failed to register org service handler: %v", err)
		}
	}
	return nil
}

// Close закрывает все установленные соединения
func (s *GatewayServer) Close() {
	for addr, conn := range s.connections {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				return
			}
			log.Printf("Closed connection to %s", addr)
		}
	}
}

// ServeHTTP обрабатывает входящие HTTP запросы
func (s *GatewayServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Добавляем CORS заголовки
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
	}

	// Обрабатываем префлайт запросы
	if r.Method == "OPTIONS" {
		return
	}

	// Добавляем request ID в контекст
	ctx := context.WithValue(r.Context(), "request-id", generateRequestID())
	r = r.WithContext(ctx)

	// Передаем запрос в grpc-gateway мультиплексор
	s.mux.ServeHTTP(w, r)
}

// generateRequestID создает уникальный идентификатор запроса
func generateRequestID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), uuid.New().String())
}

// Выполнение gRPC вызова
func (s *GatewayServer) invokeGRPC(ctx context.Context, conn *grpc.ClientConn, path string, data []byte) ([]byte, error) {
	// Извлекаем имя метода из пути
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid path")
	}

	method := fmt.Sprintf("/%s.%sService/%s", parts[0], strings.Title(parts[0]), strings.Title(parts[1]))

	// Выполняем unary RPC вызов
	resp, err := conn.Invoke(ctx, method, data, []grpc.CallOption{})
	if err != nil {
		return nil, fmt.Errorf("grpc call failed: %v", err)
	}

	// Преобразуем ответ в JSON
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %v", err)
	}

	return jsonResp, nil
}

func (a *App) initGatewayServer() error {
	server := NewGatewayServer(a.serviceProvider.GatewayConfig())

	if err := server.Connect(); err != nil {
		log.Fatalf("Failed to connect to services: %v", err)
	}
	defer server.Close()

	http.Handle("/", server)

	log.Printf("Starting gateway server on port %s", a.serviceProvider.HTTPConfig().Port())
	if err := http.ListenAndServe(a.serviceProvider.HTTPConfig().Port(), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
