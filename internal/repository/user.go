package repository

import (
	"context"
	"log"
	"uuid"

	"github.com/tinywaves/gohub/internal/domain"
	"github.com/tinywaves/gohub/internal/repository/cache"
	"github.com/tinywaves/gohub/internal/repository/dao"
)

var (
	ErrRepositoryEmailDuplicate = dao.ErrDAOEmailDuplicate
	ErrRepositoryUserNotFound   = dao.ErrDAOUserNotFound
)

type UserRepository struct {
	uDAO   *dao.UserDAO
	uCache *cache.UserCache
}

func InitUserRepository(uDAO *dao.UserDAO, uCache *cache.UserCache) *UserRepository {
	return &UserRepository{
		uDAO:   uDAO,
		uCache: uCache,
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
	uDomain, err := uRepo.uCache.Get(ctx, id)
	// The user cache is exist, you can return it immediately.
	if err == nil {
		return uDomain, err
	}
	// When Redis encounters an error, the key is to distinguish between a normal cache miss and an actual Redis failure.
	// A cache miss should fall back to the database,
	// but if Redis is down and a large number of requests simultaneously hit the database,
	// it can cause a cache avalanche and potentially bring down the database.
	// Therefore, the system should protect the database through rate limiting, degradation, or circuit breaking.
	// Alternatively, requests can be rejected directly when Redis fails,
	// sacrificing some user experience in exchange for overall system stability.
	// The final choice depends on the business requirements and the trade-off between availability and system resilience.
	record, err := uRepo.uDAO.SelectById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	uDomain = domain.User{
		UserId:   record.Id,
		Email:    record.Email,
		Password: record.Password,
		Bio:      record.Bio,
		Birthday: record.Birthday,
		Gender:   record.Gender,
		Nickname: record.Nickname,
	}
	// You can create a goroutine that doesn't block the main process.
	go func() {
		err := uRepo.uCache.Set(ctx, uDomain)
		if err != nil {
			// A failed cache setup is not a major issue, so monitoring the logs is sufficient.
			log.Println("set cache failed", err)
		}
	}()
	return uDomain, nil
}

func (uRepo *UserRepository) EditByUserId(ctx context.Context, uDomain domain.User) error {
	return uRepo.uDAO.UpdateById(ctx, uDomain.UserId, dao.User{
		Bio:      uDomain.Bio,
		Birthday: uDomain.Birthday,
		Gender:   uDomain.Gender,
		Nickname: uDomain.Nickname,
	})
}
