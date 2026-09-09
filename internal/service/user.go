package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/tinywaves/gohub/internal/domain"
	"github.com/tinywaves/gohub/internal/repository"
)

var (
	ErrServiceEmailDuplicate        = repository.ErrRepositoryEmailDuplicate
	ErrServiceEmailPasswordNotMatch = errors.New("your email or password is incorrect, please try again")
	ErrServiceUserNotFound          = repository.ErrRepositoryUserNotFound
)

type UserService struct {
	uRepo *repository.UserRepository
}

func InitUserService(uRepo *repository.UserRepository) *UserService {
	return &UserService{
		uRepo: uRepo,
	}
}

func (uSvc *UserService) SignUp(ctx context.Context, uDomain domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(uDomain.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	uDomain.Password = string(hash)
	return uSvc.uRepo.Create(ctx, uDomain)
}

func (uSvc *UserService) SignIn(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := uSvc.uRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrRepositoryUserNotFound) {
			return domain.User{}, ErrServiceEmailPasswordNotMatch
		}
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrServiceEmailPasswordNotMatch
	}
	return user, nil
}

func (uSvc *UserService) Edit(ctx context.Context, uDomain domain.User) error {
	return uSvc.uRepo.EditByUserId(ctx, uDomain)
}

func (uSvc *UserService) Profile(ctx context.Context, userId string) (domain.User, error) {
	user, err := uSvc.uRepo.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	} else {
		return domain.User{
			Email:    user.Email,
			UserId:   user.UserId,
			Bio:      user.Bio,
			Birthday: user.Birthday,
			Gender:   user.Gender,
			Nickname: user.Nickname,
		}, nil
	}
}
