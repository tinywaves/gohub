package middleware

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/tinywaves/gohub/internal"
)

type authMiddlewareBuilder struct {
	ignorePathsForAuth []string
}

func InitAuthMiddlewareBuilder() *authMiddlewareBuilder {
	return &authMiddlewareBuilder{}
}

func (amb *authMiddlewareBuilder) AppendIgnorePathsForAuth(paths ...string) *authMiddlewareBuilder {
	amb.ignorePathsForAuth = append(amb.ignorePathsForAuth, paths...)
	return amb
}

func (amb *authMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if slices.Contains(amb.ignorePathsForAuth, ctx.Request.URL.Path) {
			return
		}

		receivedToken := ctx.GetHeader(internal.JwtTokenHeaderKey)
		if receivedToken == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims := &internal.UserClaims{}
		token, err := jwt.ParseWithClaims(
			receivedToken,
			claims,
			func(t *jwt.Token) (any, error) {
				if t.Method != jwt.SigningMethodEdDSA {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return internal.Ed25519PublicKey, nil
			},
		)
		if err != nil || !token.Valid || claims.UserId == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			log.Println("safety issue", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// token refresh interval: 10 seconds
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < internal.JwtTokenDuration-internal.JwtTokenRefreshInterval {
			claims.ExpiresAt = jwt.NewNumericDate(now.Add(internal.JwtTokenDuration))
			tokenString, err := token.SignedString(internal.Ed25519PrivateKey)
			if err != nil {
				log.Println("token refresh failed", err)
			} else {
				ctx.Header(internal.JwtTokenHeaderKey, tokenString)
			}
		}

		ctx.Set(internal.CtxUserKey, claims.UserId)
	}
}
