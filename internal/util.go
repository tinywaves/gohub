package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUserId(ctx *gin.Context) string {
	value, exists := ctx.Get(CtxUserKey)
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	userId, ok := value.(string)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	return userId
}

func ConcatDeletedAt(query string) string {
	return query + " AND deleted_at = 0"
}
