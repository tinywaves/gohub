package internal

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

const (
	Port              = 8080
	MysqlDsn          = "root:root@tcp(localhost:13306)/gohub"
	RedisAddr         = "localhost:16379"
	DevUrl            = "http://localhost"
	ProdUrl           = "https://gohub.com"
	JwtTokenHeaderKey = "gohub-web-token"
	privateKeyPath    = "/Users/donghui/Developer/gohub/private_key.pem"
	publicKeyPath     = "/Users/donghui/Developer/gohub/public_key.pem"
	CtxUserIdKey      = "currentUid"
	RateLimitPrefix   = "ip-ratelimit"
	RateLimitInterval = time.Second
	RateLimitRate     = 1000
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

func loadPrivateKey() {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		panic(err)
	}
	PrivateKey = privateKey
}

func loadPublicKey() {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		panic(err)
	}
	PublicKey = publicKey
}

func LoadJwtKeys() {
	loadPrivateKey()
	loadPublicKey()
}
