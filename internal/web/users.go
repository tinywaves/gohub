package web

import "github.com/gin-gonic/gin"

type UserHandler struct{}

func (uh *UserHandler) RegisterRoutes(server *gin.RouterGroup) {
	usersServer := server.Group("/users")
	usersServer.POST("/sign-up", uh.SignUp)
	usersServer.POST("/sign-in", uh.SignIn)
	usersServer.PATCH("/:id", uh.UpdateUserInfo)
	usersServer.GET("/:id", uh.GetUserInfo)
}

func (uh *UserHandler) SignUp(ctx *gin.Context) {}

func (uh *UserHandler) SignIn(ctx *gin.Context) {}

func (uh *UserHandler) UpdateUserInfo(ctx *gin.Context) {}

func (uh *UserHandler) GetUserInfo(ctx *gin.Context) {}
