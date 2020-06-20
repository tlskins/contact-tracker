package auth

import (
	"context"
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	key                      *rsa.PublicKey
	secret                   *rsa.PrivateKey
	accessExpirationMinutes  int
	refreshExpirationMinutes int
}

type JWTServiceConfig struct {
	Key                      []byte
	Secret                   []byte
	AccessExpirationMinutes  int
	RefreshExpirationMinutes int
}

const AccessTokenKey = "accessToken"

func NewJWTService(config JWTServiceConfig) (*JWTService, error) {
	var jKey *rsa.PublicKey
	var err error
	if config.Key != nil {
		if jKey, err = jwt.ParseRSAPublicKeyFromPEM(config.Key); err != nil {
			return nil, err
		}
	}
	var jSecret *rsa.PrivateKey
	if config.Secret != nil {
		if jSecret, err = jwt.ParseRSAPrivateKeyFromPEM(config.Secret); err != nil {
			return nil, err
		}
	}

	return &JWTService{
		key:                      jKey,
		secret:                   jSecret,
		accessExpirationMinutes:  config.AccessExpirationMinutes,
		refreshExpirationMinutes: config.RefreshExpirationMinutes,
	}, nil
}

func (j *JWTService) AccessExpirationMinutes() int {
	return j.accessExpirationMinutes
}

func (j *JWTService) RefreshExpirationMinutes() int {
	return j.refreshExpirationMinutes
}

func AccessTokenFromContext(ctx context.Context) *CustomClaims {
	return ctx.Value(AccessTokenKey).(*CustomClaims)
}

func UserIdFromContext(ctx context.Context) string {
	return AccessTokenFromContext(ctx).Subject
}
