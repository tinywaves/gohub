package web

import "github.com/gin-gonic/gin"

func InitWeb() *gin.Engine {
	server := gin.Default()

	v1Server := server.Group("/v1/api")

	userHandler := InitUsersHandler()
	userHandler.RegisterRoutes(v1Server.Group("/users"))

	return server
}
