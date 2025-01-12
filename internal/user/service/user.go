package service

import (
	"context"
	"errors"
	"gohub/internal/user/domain"
	"gohub/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserEmailDuplicated       = repository.ErrUserEmailDuplicated
	ErrUserNotFound              = repository.ErrUserNotFound
	ErrUserEmailPasswordNotMatch = errors.New("your email and password not match")
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
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encryptedPassword)
	return us.userRepository.CreateUser(ctx, u)
}

func (us *UserService) SignIn(ctx context.Context, u domain.User) (domain.User, error) {
	foundUser, err := us.userRepository.FindUserByEmail(ctx, u.Email)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(u.Password))
	if err != nil {
		return domain.User{}, ErrUserEmailPasswordNotMatch
	}
	return foundUser, nil
}
