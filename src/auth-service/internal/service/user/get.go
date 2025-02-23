package user

import (
	"context"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
