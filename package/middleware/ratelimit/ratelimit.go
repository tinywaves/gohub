package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

type MiddlewareBuilder struct {
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
) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		prefix:   prefix,
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (m *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := m.limit(ctx)
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

func (m *MiddlewareBuilder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", m.prefix, ctx.ClientIP())
	return m.cmd.Eval(
		ctx,
		luaScript,
		[]string{key},
		m.interval.Milliseconds(),
		m.rate,
		time.Now().UnixMilli(),
	).Bool()
}

func (m *MiddlewareBuilder) Prefix(prefix string) *MiddlewareBuilder {
	m.prefix = prefix
	return m
}
