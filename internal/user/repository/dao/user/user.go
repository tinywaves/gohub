package user

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var ErrUserEmailDuplicated = errors.New("user with this email already exists")

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
