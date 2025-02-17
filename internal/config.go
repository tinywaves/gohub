package internal

import (
	"crypto/rand"
	"crypto/rsa"
)

const (
	Port              = 8080
	SessionDataKey    = "gohub-user"
	MysqlDsn          = "root:root@tcp(localhost:13306)/gohub"
	DevUrl            = "http://localhost"
	ProdUrl           = "https://gohub.com"
	JwtTokenHeaderKey = "gohub-web-token"
	PrivateKeyBits    = 2048
)

var PrivateKey *rsa.PrivateKey

func GeneratePrivateKey() {
	privateKey, err := rsa.GenerateKey(rand.Reader, PrivateKeyBits)
	if err != nil {
		panic(err)
	}
	PrivateKey = privateKey
}
