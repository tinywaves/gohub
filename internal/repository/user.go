package repository

import (
	"context"
	"uuid"

	"github.com/tinywaves/gohub/internal/domain"
	"github.com/tinywaves/gohub/internal/repository/dao"
)

var (
	ErrRepositoryEmailDuplicate = dao.ErrDAOEmailDuplicate
	ErrRepositoryUserNotFound   = dao.ErrDAOUserNotFound
)

type UserRepository struct {
	uDAO *dao.UserDAO
}

func InitUserRepository(uDAO *dao.UserDAO) *UserRepository {
	return &UserRepository{
		uDAO: uDAO,
	}
}

func (uRepo *UserRepository) Create(ctx context.Context, uDomain domain.User) error {
	return uRepo.uDAO.Insert(ctx, dao.User{
		Id:       uuid.New().String(),
		Email:    uDomain.Email,
		Password: uDomain.Password,
	})
}

func (uRepo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	record, err := uRepo.uDAO.SelectByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Email:    record.Email,
		Password: record.Password,
		UserId:   record.Id,
	}, nil
}

func (uRepo *UserRepository) FindById(ctx context.Context, id string) (domain.User, error) {
	record, err := uRepo.uDAO.SelectById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Email:    record.Email,
		Password: record.Password,
		UserId:   record.Id,
		Bio:      record.Bio,
		Birthday: record.Birthday,
		Gender:   record.Gender,
		Nickname: record.Nickname,
	}, nil
}

func (uRepo *UserRepository) EditByUserId(ctx context.Context, uDomain domain.User) error {
	return uRepo.uDAO.UpdateById(ctx, uDomain.UserId, dao.User{
		Bio:      uDomain.Bio,
		Birthday: uDomain.Birthday,
		Gender:   uDomain.Gender,
		Nickname: uDomain.Nickname,
	})
}
