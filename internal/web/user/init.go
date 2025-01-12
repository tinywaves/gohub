package user

import (
	"github.com/gin-gonic/gin"
	"gohub/internal/repository"
	"gohub/internal/repository/dao/user"
	"gohub/internal/service"
	"gorm.io/gorm"
)

func InitUserWeb(database *gorm.DB, v1Server *gin.RouterGroup) {
	userDao := user.InitDao(database)
	userRepository := repository.InitUserRepository(userDao)
	userService := service.InitUserService(userRepository)
	userHandler := InitHandler(userService)
	userHandler.RegisterRoutes(v1Server.Group("/user"))
}
