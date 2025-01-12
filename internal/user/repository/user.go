package repository

import (
	"context"
	"github.com/google/uuid"
	"gohub/internal/user/domain"
	"gohub/internal/user/repository/dao/user"
	"time"
)

var (
	ErrUserEmailDuplicated = user.ErrUserEmailDuplicated
	ErrUserNotFound        = user.ErrUserNotFound
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

func (ur *UserRepository) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	entity, err := ur.userDao.QueryUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{Email: entity.Email, Password: entity.Password}, nil
}
