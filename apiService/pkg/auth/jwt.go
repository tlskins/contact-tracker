package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	apiHttp "github.com/contact-tracker/apiService/pkg/http"
	apiLambda "github.com/contact-tracker/apiService/pkg/lambda"

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

func (j *JWTService) AuthorizeHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(AccessTokenKey)
		if err != nil {
			if err == http.ErrNoCookie {
				apiHttp.CheckHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized access"))
			}
			apiHttp.CheckHTTPError(http.StatusUnauthorized, err)
		}

		tokenStr := c.Value
		claims, err := j.Decode(tokenStr)
		apiHttp.CheckHTTPError(http.StatusUnauthorized, err)
		ctx := r.Context()
		ctx = context.WithValue(ctx, AccessTokenKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (j *JWTService) IncludeLambdaAuth(ctx context.Context, req *apiLambda.Request) (context.Context, error) {
	if cookieVal, ok := req.Headers["Cookie"]; ok {
		val := strings.ReplaceAll(cookieVal, fmt.Sprintf("%s=", AccessTokenKey), "")
		var err error
		var claims *CustomClaims
		if claims, err = j.Decode(val); err != nil {
			return ctx, err
		}
		ctx = context.WithValue(ctx, AccessTokenKey, claims)
	}
	return ctx, nil
}

func ClaimsFromContext(ctx context.Context) *CustomClaims {
	if ctx.Value(AccessTokenKey) == nil {
		return nil
	}
	return ctx.Value(AccessTokenKey).(*CustomClaims)
}
