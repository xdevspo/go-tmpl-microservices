package users

import (
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/service"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/pkg/users"
)

type Implementation struct {
	users.UnimplementedUsersServer
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
