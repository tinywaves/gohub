package dao

import (
	"gohub/internal/user/repository/dao/user"
	"gorm.io/gorm"
)

func InitTables(database *gorm.DB) error {
	return database.AutoMigrate(&user.Entity{})
}
