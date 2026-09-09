package internal

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	Ed25519PrivateKey ed25519.PrivateKey
	Ed25519PublicKey  ed25519.PublicKey
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    string
	UserAgent string
}

func InitJwtKeyPair() {
	Ed25519PrivateKey = loadPrivateKey(ed25519PrivateKeyPath)
	Ed25519PublicKey = loadPublicKey(ed25519PublicKeyPath)
}

func loadPrivateKey(path string) ed25519.PrivateKey {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		panic("failed to decode private key PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	privateKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		panic("private key is not Ed25519")
	}

	return privateKey
}

func loadPublicKey(path string) ed25519.PublicKey {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		panic("failed to decode public key PEM")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	publicKey, ok := key.(ed25519.PublicKey)
	if !ok {
		panic("public key is not Ed25519")
	}

	return publicKey
}
