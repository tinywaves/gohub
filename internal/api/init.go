package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gohub/internal"
	"gohub/internal/api/user"
	"gohub/internal/api/user/repository/dao"
	"gohub/internal/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func Init() *gin.Engine {
	// privateKey is used to sign and verify the JWT token
	internal.GeneratePrivateKey()

	// gin server
	server := gin.Default()
	database, err := gorm.Open(mysql.Open(internal.MysqlDsn))
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
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type", internal.JwtTokenHeaderKey},
		ExposeHeaders: []string{internal.JwtTokenHeaderKey},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, internal.DevUrl) {
				return true
			}
			return strings.Contains(origin, internal.ProdUrl)
		},
		MaxAge: 12 * time.Hour,
	}))

	// check sign in status
	server.Use(
		middleware.
			InitAuthMiddlewareBuilder().
			AppendIgnorePath("/v1/api/user/sign-in").
			AppendIgnorePath("/v1/api/user/sign-up").
			Builder(),
	)

	v1Server := server.Group("/v1/api")

	user.Init(database, v1Server)

	return server
}
