package web

import "github.com/gin-gonic/gin"

type UserHandler struct{}

func (uh *UserHandler) RegisterRoutes(server *gin.RouterGroup) {
	server.POST("/sign-up", uh.SignUp)
	server.POST("/sign-in", uh.SignIn)
	server.PATCH("/:id", uh.UpdateUserInfo)
	server.GET("/:id", uh.GetUserInfo)
}

func (uh *UserHandler) SignUp(ctx *gin.Context) {}

func (uh *UserHandler) SignIn(ctx *gin.Context) {}

func (uh *UserHandler) UpdateUserInfo(ctx *gin.Context) {}

func (uh *UserHandler) GetUserInfo(ctx *gin.Context) {}
