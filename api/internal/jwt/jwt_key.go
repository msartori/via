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
	jwt.SigningMethod
}

var (
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	once          sync.Once
	mutex         sync.Mutex
	signingMethod jwt.SigningMethod
)

func Init(cfg JWTConfig) error {
	var err error
	once.Do(func() {
		if cfg.PublicKey == "" {
			cfg.PublicKey = secret.Get().Read(cfg.PublicKeySecretFile)
		}
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKey))
		if err != nil {
			return
		}
		if cfg.PrivateKey == "" {
			cfg.PrivateKey = secret.Get().Read(cfg.PrivateKeySecretFile)
		}
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
		if err != nil {
			return
		}
		if cfg.SigningMethod == nil {
			cfg.SigningMethod = jwt.SigningMethodRS256
		}
		signingMethod = cfg.SigningMethod
	})
	return err
}

func GetPrivateKey() *rsa.PrivateKey {
	return privateKey
}

func GetPublicKey() *rsa.PublicKey {
	return publicKey
}

func GetSigningMethod() jwt.SigningMethod {
	return signingMethod
}

func Reset() {
	mutex.Lock()
	defer mutex.Unlock()
	privateKey = nil
	publicKey = nil
	once = sync.Once{}
}
