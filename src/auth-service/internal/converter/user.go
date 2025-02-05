package converter

import (
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/model"
	desc "github.com/xdevspo/go-tmpl-microservices/auth-service/pkg/users"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromService(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromService(info model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Username:     info.Username,
		Email:        info.Email,
		PasswordHash: info.PasswordHash,
	}
}
func ToUserInfoFromDesc(info *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Username:     info.Username,
		Email:        info.Email,
		PasswordHash: info.PasswordHash,
	}
}
