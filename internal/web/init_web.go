package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gohub/internal/repository"
	"gohub/internal/repository/dao"
	"gohub/internal/repository/dao/user"
	"gohub/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func InitWeb() *gin.Engine {
	server := gin.Default()
	database, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/gohub"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(database)
	if err != nil {
		panic(err)
	}

	// cors
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		// AllowAllOrigins:  true,
		// AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "https://gohub.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	v1Server := server.Group("/v1/api")

	userDao := user.InitDao(database)
	userRepository := repository.InitUserRepository(userDao)
	userService := service.InitUserService(userRepository)
	userHandler := InitUserHandler(userService)
	userHandler.RegisterRoutes(v1Server.Group("/user"))

	return server
}
