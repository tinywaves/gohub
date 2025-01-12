package user

import (
	"github.com/gin-gonic/gin"
	"gohub/internal/user/repository"
	"gohub/internal/user/repository/dao/user"
	"gohub/internal/user/service"
	"gohub/internal/user/web"
	"gorm.io/gorm"
)

func Init(database *gorm.DB, v1Server *gin.RouterGroup) {
	userDao := user.InitDao(database)
	userRepository := repository.InitUserRepository(userDao)
	userService := service.InitUserService(userRepository)
	userHandler := web.InitUserHandler(userService)
	userHandler.RegisterRoutes(v1Server.Group("/user"))
}
