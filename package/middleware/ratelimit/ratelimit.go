package ratelimit

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ratelimitMiddlewareBuilder struct {
	prefix   string
	cmd      redis.Cmdable
	interval time.Duration
	rate     int
}

//go:embed slide_window.lua
var luaScript string

func InitRatelimitMiddlewareBuilder(
	prefix string,
	cmd redis.Cmdable,
	interval time.Duration,
	rate int,
) *ratelimitMiddlewareBuilder {
	return &ratelimitMiddlewareBuilder{
		prefix:   prefix,
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (rmb *ratelimitMiddlewareBuilder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", rmb.prefix, ctx.ClientIP())
	return rmb.cmd.Eval(
		ctx,
		luaScript,
		[]string{key},
		rmb.interval.Milliseconds(),
		rmb.rate,
		time.Now().UnixMilli(),
	).Bool()
}

func (rmb *ratelimitMiddlewareBuilder) Prefix(prefix string) *ratelimitMiddlewareBuilder {
	rmb.prefix = prefix
	return rmb
}

func (rmb *ratelimitMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := rmb.limit(ctx)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
