package user

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserEmailDuplicated = errors.New("user with this email already exists")
	ErrUserNotFound        = gorm.ErrRecordNotFound
)

type Dao struct {
	database *gorm.DB
}

func InitDao(database *gorm.DB) *Dao {
	return &Dao{
		database: database,
	}
}

func (d *Dao) InsertUserRecord(ctx context.Context, e Entity) error {
	err := d.database.WithContext(ctx).Create(&e).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const mysqlDuplicateEntry = 1062
		if mysqlErr.Number == mysqlDuplicateEntry {
			return ErrUserEmailDuplicated
		}
	}
	return err
}

func (d *Dao) QueryUserByEmail(ctx context.Context, email string) (Entity, error) {
	var e Entity
	err := d.database.WithContext(ctx).Where("email = ?", email).First(&e).Error
	return e, err
}

func (d *Dao) QueryUserById(ctx context.Context, id string) error {
	err := d.database.WithContext(ctx).Where("id = ?", id).First(&Entity{}).Error
	return err
}

func (d *Dao) ModifyUserById(ctx context.Context, id string, e Entity) error {
	return d.
		database.
		WithContext(ctx).
		Where("id = ?", id).
		Updates(&e).
		Error
}
