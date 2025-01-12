package repository

import (
	"context"
	"github.com/google/uuid"
	"gohub/internal/domain"
	"gohub/internal/repository/dao/user"
	"time"
)

type UserRepository struct {
	userDao *user.Dao
}

func InitUserRepository(userDao *user.Dao) *UserRepository {
	return &UserRepository{
		userDao: userDao,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, u domain.User) error {
	now := time.Now().UnixMilli()

	return ur.userDao.InsertUserRecord(
		ctx,
		user.Entity{
			Id:         uuid.New().String(),
			Email:      u.Email,
			Password:   u.Password,
			CreateTime: now,
			UpdateTime: now,
		},
	)
}
