package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gohub/internal"
	"gohub/internal/api/user/web"
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

		receivedToken := ctx.GetHeader(internal.JwtTokenHeaderKey)
		if receivedToken == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userClaims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(
			receivedToken,
			userClaims,
			func(token *jwt.Token) (interface{}, error) {
				return internal.PublicKey, nil
			},
		)
		if err != nil || token == nil || !token.Valid || userClaims.Uid == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set(internal.CtxUserIdKey, userClaims.Uid)
		return
	}
}
