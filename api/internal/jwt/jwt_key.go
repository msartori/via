package jwt_key

import (
	"crypto/rsa"
	"sync"
	"via/internal/secret"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	PrivateKey           string `env:"PRIVATE_KEY" json:"-"`
	PublicKey            string `env:"PUBLIC_KEY" json:"-"`
	PrivateKeySecretFile string `env:"PRIVATE_KEY_FILE" json:"privateKeySecretFile"`
	PublicKeySecretFile  string `env:"PUBLIC_KEY_FILE" json:"publicKeySecretFile"`
}

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	once       sync.Once
)

func Init(cfg JWTConfig) error {
	once.Do(func() {
		var err error
		if cfg.PrivateKey == "" {
			cfg.PrivateKey = secret.ReadSecret(cfg.PrivateKeySecretFile)
		}
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
		if err != nil {
			return
		}
		if cfg.PublicKey == "" {
			cfg.PublicKey = secret.ReadSecret(cfg.PublicKeySecretFile)
		}
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKey))
		if err != nil {
			return
		}
	})
	return nil
}

func GetPrivateKey() *rsa.PrivateKey {
	return privateKey
}

func GetPublicKey() *rsa.PublicKey {
	return publicKey
}

func GetSigningMethod() jwt.SigningMethod {
	return jwt.SigningMethodRS256
}
