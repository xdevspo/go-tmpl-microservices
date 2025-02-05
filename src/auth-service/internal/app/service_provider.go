package app

import (
	"context"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/api/users"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/client/db"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/client/db/pg"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/client/db/transaction"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/closer"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/config"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/repository"
	userRepository "github.com/xdevspo/go-tmpl-microservices/auth-service/internal/repository/user"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/service"
	userService "github.com/xdevspo/go-tmpl-microservices/auth-service/internal/service/user"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository

	userService service.UserService

	userImpl *users.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) PGConfig() config.PGConfig {
	if sp.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		sp.pgConfig = cfg
	}

	return sp.pgConfig
}

func (sp *serviceProvider) GRPCConfig() config.GRPCConfig {
	if sp.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		sp.grpcConfig = cfg
	}

	return sp.grpcConfig
}

func (sp *serviceProvider) DBClient(ctx context.Context) db.Client {
	if sp.dbClient == nil {
		dbClient, err := pg.New(ctx, sp.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = dbClient.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping db: %v", err)
		}
		closer.Add(dbClient.Close)

		sp.dbClient = dbClient
	}

	return sp.dbClient
}

func (sp *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if sp.txManager == nil {
		sp.txManager = transaction.NewTransactionManager(sp.DBClient(ctx).DB())
	}

	return sp.txManager
}

func (sp *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if sp.userRepository == nil {
		sp.userRepository = userRepository.NewRepository(sp.DBClient(ctx))
	}

	return sp.userRepository
}

func (sp *serviceProvider) UserService(ctx context.Context) service.UserService {
	if sp.userService == nil {
		sp.userService = userService.NewService(sp.UserRepository(ctx), sp.TxManager(ctx))
	}

	return sp.userService
}

func (sp *serviceProvider) UserImpl(ctx context.Context) *users.Implementation {
	if sp.userImpl == nil {
		sp.userImpl = users.NewImplementation(sp.UserService(ctx))
	}

	return sp.userImpl
}
