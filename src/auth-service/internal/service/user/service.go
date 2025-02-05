package user

import (
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/client/db"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/repository"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
