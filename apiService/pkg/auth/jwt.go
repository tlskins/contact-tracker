package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	apiHttp "github.com/contact-tracker/apiService/pkg/api/http"
	apiLambda "github.com/contact-tracker/apiService/pkg/api/lambda"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	key                      *rsa.PublicKey
	secret                   *rsa.PrivateKey
	rpcPwd                   string
	accessExpirationMinutes  int
	refreshExpirationMinutes int
}

type JWTServiceConfig struct {
	Key                      []byte
	Secret                   []byte
	RPCPwd                   string
	AccessExpirationMinutes  int
	RefreshExpirationMinutes int
}

// AccessTokenKey - cookie and context key for jwt claims authorization
const AccessTokenKey = "accessToken"

// RPCAccessTokenKey - cookie and context key for rpc authorization
const RPCAccessTokenKey = "rpcAccessToken"

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
		rpcPwd:                   config.RPCPwd,
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
		ctx := r.Context()

		rpcCookie, err := r.Cookie(RPCAccessTokenKey)
		if err != nil && err != http.ErrNoCookie {
			apiHttp.CheckHTTPError(http.StatusUnauthorized, err)
		}

		fmt.Printf("rpcCookie: %+v\n\n", rpcCookie)

		if rpcCookie != nil && rpcCookie.Value == j.rpcPwd {
			ctx = context.WithValue(ctx, RPCAccessTokenKey, rpcCookie.Value)
		} else {
			accessCookie, err := r.Cookie(AccessTokenKey)
			if err != nil {
				if err == http.ErrNoCookie {
					apiHttp.CheckHTTPError(http.StatusUnauthorized, fmt.Errorf("Unauthorized access"))
				}
				apiHttp.CheckHTTPError(http.StatusUnauthorized, err)
			}
			claims, err := j.Decode(accessCookie.Value)
			apiHttp.CheckHTTPError(http.StatusUnauthorized, err)
			ctx = context.WithValue(ctx, AccessTokenKey, claims)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (j *JWTService) IncludeLambdaAuth(ctx context.Context, req *apiLambda.Request) (context.Context, error) {
	authVal := ""
	if val, ok := req.Headers["Cookie"]; ok {
		authVal = strings.ReplaceAll(val, fmt.Sprintf("%s=", AccessTokenKey), "")
	} else if val, ok := req.Headers["Authorization"]; ok {
		authVal = strings.TrimPrefix(val, "Bearer ")
	}
	if authVal != "" {
		var err error
		var claims *CustomClaims
		if claims, err = j.Decode(authVal); err != nil {
			return ctx, err
		}
		ctx = context.WithValue(ctx, AccessTokenKey, claims)
	}
	// if val, ok := req.Headers[RPCAccessTokenKey]; ok && val == j.rpcPwd {
	// 	ctx = context.WithValue(ctx, RPCAccessTokenKey, val)
	// }
	ctx = context.WithValue(ctx, RPCAccessTokenKey, j.rpcPwd)
	return ctx, nil
}

func ClaimsFromContext(ctx context.Context) (authorized bool, claims *CustomClaims) {
	if ctx.Value(AccessTokenKey) != nil {
		return true, ctx.Value(AccessTokenKey).(*CustomClaims)
	}
	if ctx.Value(RPCAccessTokenKey) != nil {
		return true, nil
	}
	return false, nil
}
