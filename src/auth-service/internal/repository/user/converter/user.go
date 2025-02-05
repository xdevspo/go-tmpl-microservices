package converter

import (
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/model"
	modelRepo "github.com/xdevspo/go-tmpl-microservices/auth-service/internal/repository/user/model"
)

func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:   user.ID,
		Info: ToUserInfoFromRepo(user.Info),
	}
}

func ToUserInfoFromRepo(info modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Username:     info.Username,
		Email:        info.Email,
		PasswordHash: info.PasswordHash,
	}
}
