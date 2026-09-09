package dao

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/tinywaves/gohub/internal"
)

var (
	ErrDAOEmailDuplicate = errors.New("the email you provided is already registered, please use another email")
	ErrDAOUserNotFound   = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func InitUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (uDAO *UserDAO) Insert(ctx context.Context, user User) error {
	err := gorm.G[User](uDAO.db).Create(ctx, &user)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const duplicateEntryErrorCode uint16 = 1062
		if mysqlErr.Number == duplicateEntryErrorCode {
			return ErrDAOEmailDuplicate
		}
	}
	return err
}

func (uDAO *UserDAO) SelectByEmail(ctx context.Context, email string) (User, error) {
	return gorm.G[User](uDAO.db).
		Where(internal.ConcatDeletedAt("email = ?"), email).
		First(ctx)
}

func (uDAO *UserDAO) SelectById(ctx context.Context, id string) (User, error) {
	return gorm.G[User](uDAO.db).
		Where(internal.ConcatDeletedAt("id = ?"), id).
		First(ctx)
}

func (uDAO *UserDAO) UpdateById(ctx context.Context, id string, user User) error {
	rowsAffected, err := gorm.G[User](uDAO.db).
		Where(internal.ConcatDeletedAt("id = ?"), id).
		Updates(ctx, user)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrDAOUserNotFound
	}
	return nil
}
