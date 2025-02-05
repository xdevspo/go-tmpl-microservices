package users

import (
	"context"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/converter"
	desc "github.com/xdevspo/go-tmpl-microservices/auth-service/pkg/users"
	"log"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", user.ID, user.Info.Username, user.Info.Email, user.Info.PasswordHash, user.CreatedAt, user.UpdatedAt)

	return &desc.GetResponse{
		User: converter.ToUserFromService(user),
	}, nil
}
