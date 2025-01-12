package service

import (
	"context"
	"gohub/internal/user/domain"
	"gohub/internal/user/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func InitUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) SignUp(ctx context.Context, u domain.User) error {
	return us.userRepository.CreateUser(ctx, u)
}
