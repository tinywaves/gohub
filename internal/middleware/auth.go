package middleware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
)

type AuthMiddlewareBuilder struct {
	ignorePaths []string
}

func InitAuthMiddlewareBuilder() *AuthMiddlewareBuilder {
	return &AuthMiddlewareBuilder{}
}

func (amb *AuthMiddlewareBuilder) AppendIgnorePath(path string) *AuthMiddlewareBuilder {
	amb.ignorePaths = append(amb.ignorePaths, path)
	return amb
}

func (amb *AuthMiddlewareBuilder) Builder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if slices.Contains(amb.ignorePaths, ctx.Request.URL.Path) {
			return
		}
		session := sessions.Default(ctx)
		if session == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		sessionUser := session.Get("gohub-user")
		fmt.Println(sessionUser)
		if sessionUser == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
