package jwt_key

import (
	"testing"
	jwt_key_mock "via/internal/jwt/mock"
	"via/internal/secret"
	mock_secret "via/internal/secret/mock"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTKey_Init_WithStrings(t *testing.T) {
	Reset()

	err := Init(JWTConfig{
		PrivateKey:    jwt_key_mock.GetPrivateKey(),
		PublicKey:     jwt_key_mock.GetPublicKey(),
		SigningMethod: jwt.SigningMethodRS256,
	})
	assert.NoError(t, err)
	assert.NotNil(t, GetPrivateKey())
	assert.NotNil(t, GetPublicKey())
	assert.Equal(t, jwt.SigningMethodRS256, GetSigningMethod())
}

func TestJWTKey_Init_WithSecretReader(t *testing.T) {
	Reset()

	mock := new(mock_secret.MockSecret)
	secret.Set(mock)

	mock.On("Read", "pub.pem").Return(jwt_key_mock.GetPublicKey())
	mock.On("Read", "priv.pem").Return(jwt_key_mock.GetPrivateKey())

	err := Init(JWTConfig{
		PublicKeySecretFile:  "pub.pem",
		PrivateKeySecretFile: "priv.pem",
	})

	assert.NoError(t, err)
	assert.NotNil(t, GetPrivateKey())
	assert.NotNil(t, GetPublicKey())
	assert.Equal(t, jwt.SigningMethodRS256, GetSigningMethod())

	mock.AssertExpectations(t)
}

func TestJWTKey_Reset(t *testing.T) {
	Reset()
	assert.Nil(t, GetPrivateKey())
	assert.Nil(t, GetPublicKey())
}

func TestJWTKey_Init_PrivateKeyInvalid(t *testing.T) {
	Reset()

	mockSecret := new(mock_secret.MockSecret)
	mockSecret.On("Read", "private.pem").Return("INVALID")
	mockSecret.On("Read", "public.pem").Return(jwt_key_mock.GetPublicKey())
	secret.Set(mockSecret)

	cfg := JWTConfig{
		PrivateKeySecretFile: "private.pem",
		PublicKeySecretFile:  "public.pem",
	}

	err := Init(cfg)
	assert.Error(t, err)
	assert.Nil(t, GetPrivateKey())
	assert.NotNil(t, GetPublicKey())
}

func TestJWTKey_Init_PublicKeyInvalid(t *testing.T) {
	Reset()

	mockSecret := new(mock_secret.MockSecret)
	mockSecret.On("Read", "public.pem").Return("INVALID")
	secret.Set(mockSecret)

	cfg := JWTConfig{
		PublicKeySecretFile:  "public.pem",
		PrivateKeySecretFile: "private.pem", // won't reach
	}

	err := Init(cfg)
	assert.Error(t, err)
	assert.Nil(t, GetPublicKey())
	assert.Nil(t, GetPrivateKey())
}
