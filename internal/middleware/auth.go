package middleware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
)

var withoutAuth = []string{"/v1/api/user/sign-up", "/v1/api/user/sign-in"}

type AuthMiddlewareBuilder struct{}

func InitAuthMiddlewareBuilder() *AuthMiddlewareBuilder {
	return &AuthMiddlewareBuilder{}
}

func (amb *AuthMiddlewareBuilder) Builder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if slices.Contains(withoutAuth, ctx.Request.URL.Path) {
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
