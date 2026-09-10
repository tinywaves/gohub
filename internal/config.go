package internal

import "time"

const (
	// cors
	CorsProdHost = "https://gohub.tinywaves.site"
	CorsMaxAge   = 5 * time.Minute

	// mysql
	MysqlDsn = "root:root@tcp(localhost:13306)/gohub"

	// redis
	RedisAddr = "localhost:16379"

	// jwt
	ed25519PrivateKeyPath   = "pem/private.pem"
	ed25519PublicKeyPath    = "pem/public.pem"
	JwtTokenDuration        = time.Hour * 24 * 7
	JwtTokenRefreshInterval = time.Second * 10

	// ratelimit
	RateLimitInterval = time.Second
	RateLimitRate     = 1000
)
