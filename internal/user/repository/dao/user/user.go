package user

import (
	"context"
	"gorm.io/gorm"
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
	return d.database.WithContext(ctx).Create(&e).Error
}
