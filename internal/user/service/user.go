package service

import (
	"context"
	"gohub/internal/user/domain"
	"gohub/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserEmailDuplicated = repository.ErrUserEmailDuplicated

type UserService struct {
	userRepository *repository.UserRepository
}

func InitUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) SignUp(ctx context.Context, u domain.User) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encryptedPassword)
	return us.userRepository.CreateUser(ctx, u)
}
