package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitWeb() *gin.Engine {
	server := gin.Default()

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

	userHandler := InitUsersHandler()
	userHandler.RegisterRoutes(v1Server.Group("/users"))

	return server
}
