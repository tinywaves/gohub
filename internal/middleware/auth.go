package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gohub/internal"
	"net/http"
	"slices"
	"time"
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
		sessionUser := session.Get(internal.SessionDataKey)
		if sessionUser == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sessionLastRefresh := session.Get(internal.SessionLastRefreshKey)
		now := time.Now().UnixMilli()
		if sessionLastRefresh == nil {
			session.Set(internal.SessionLastRefreshKey, now)
			session.Set(internal.SessionDataKey, sessionUser)
			session.Options(sessions.Options{
				MaxAge: 60,
			})
			err := session.Save()
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			} else {
				return
			}
		} else {
			sessionLastRefreshVal, ok := sessionLastRefresh.(int64)
			if !ok {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if now-sessionLastRefreshVal > internal.SessionLastRefreshInterval*1000 {
				session.Set(internal.SessionLastRefreshKey, now)
				session.Set(internal.SessionDataKey, sessionUser)
				session.Options(sessions.Options{
					MaxAge: 60,
				})
				err := session.Save()
				if err != nil {
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return
				} else {
					return
				}
			} else {
				return
			}
		}
	}
}
