package users

import (
	"context"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/converter"
	desc "github.com/xdevspo/go-tmpl-microservices/auth-service/pkg/users"
	"log"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserInfoFromDesc(req.User))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted note with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
