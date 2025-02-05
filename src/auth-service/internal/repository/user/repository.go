package user

import (
	"context"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/client/db"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/model"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/repository"
	"github.com/xdevspo/go-tmpl-microservices/auth-service/internal/repository/user/converter"
	modelRepo "github.com/xdevspo/go-tmpl-microservices/auth-service/internal/repository/user/model"
)

const (
	tableName          = "users"
	idColumn           = "id"
	usernameColumn     = "username"
	emailColumn        = "email"
	passwordHashColumn = "password_hash"
	createdAtColumn    = "created_at"
	updatedAtColumn    = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	query := "INSERT INTO users (username, password_hash, email, created_at) VALUES ($1, $2, $3, now()) RETURNING id"

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err := r.db.DB().QueryRowContext(ctx, q, &info.Username, &info.PasswordHash, &info.Email).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	query := "SELECT id, username, password_hash, email, created_at, updated_at FROM users WHERE id = $1"

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err := r.db.DB().QueryRowContext(ctx, q, id).Scan(&user.ID, &user.Info.Username, &user.Info.PasswordHash, &user.Info.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}
