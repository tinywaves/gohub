package main

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tinywaves/gohub/internal"
	"github.com/tinywaves/gohub/internal/repository"
	"github.com/tinywaves/gohub/internal/repository/cache"
	"github.com/tinywaves/gohub/internal/repository/dao"
	"github.com/tinywaves/gohub/internal/service"
	"github.com/tinywaves/gohub/internal/web"
	"github.com/tinywaves/gohub/internal/web/middleware"
	"github.com/tinywaves/gohub/package/middleware/ratelimit"
)

func main() {
	// jwt
	internal.InitJwtKeyPair()

	// gin
	server := gin.Default()

	// gorm
	db, err := gorm.Open(mysql.Open(internal.MysqlDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(internal.MysqlMaxOpenConns)
	sqlDB.SetMaxIdleConns(internal.MysqlMaxIdleConns)
	sqlDB.SetConnMaxLifetime(internal.MysqlConnMaxLifetime)
	if err = dao.InitTables(db); err != nil {
		panic(err)
	}

	// redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: internal.RedisAddr,
	})

	// ratelimit
	server.Use(
		ratelimit.
			InitRatelimitMiddlewareBuilder(
				"ratelimit-client-ip",
				redisClient,
				internal.RateLimitInterval,
				internal.RateLimitRate,
			).
			Build(),
	)

	// cors
	server.Use(cors.New(cors.Config{
		// AllowCredentials: true,
		// AllowAllOrigins:  true,
		// AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type", internal.JwtTokenHeaderKey},
		ExposeHeaders: []string{internal.JwtTokenHeaderKey},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, internal.CorsProdHost)
		},
		MaxAge: internal.CorsMaxAge,
	}))

	// auth
	server.Use(
		middleware.
			InitAuthMiddlewareBuilder().
			AppendIgnorePathsForAuth("/api/v1/user/sign-up", "/api/v1/user/sign-in").
			Build(),
	)

	api := server.Group("/api")
	v1 := api.Group("/v1")

	userDAO := dao.InitUserDAO(db)
	userCache := cache.InitUserCache(redisClient, internal.UserCacheExpiration)
	userRepository := repository.InitUserRepository(userDAO, userCache)
	userService := service.InitUserService(userRepository)
	userHandler := web.InitUserHandler(userService)
	userHandler.RegisterRoutes(v1.Group("/user"))

	if err := server.Run(":11111"); err != nil {
		panic("gin server run error")
	}
}
