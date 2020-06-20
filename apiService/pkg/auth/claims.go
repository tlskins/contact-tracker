package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
	Email       string `json:"email"`
	IsConfirmed bool   `json:"isConfirmed"`
}

func (c CustomClaims) Valid() error {
	return nil
}

type Authable interface {
	GetAuthables() (id string, email string, conf bool)
}

func (j *JWTService) GenAccessToken(a Authable) (accessToken string, err error) {
	now := time.Now()
	id, email, conf := a.GetAuthables()
	claims := CustomClaims{
		jwt.StandardClaims{
			IssuedAt: now.Unix(),
			Issuer:   "api.contact-tracker.com",
			Subject:  id,
		},
		email,
		conf,
	}

	return j.Encode(claims)
}

func (j *JWTService) Encode(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t, err := token.SignedString(j.secret)

	return t, err
}

func (j *JWTService) Decode(tokenStr string) (claims CustomClaims, err error) {
	decodedToken, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(decodedToken *jwt.Token) (interface{}, error) {
		if _, ok := decodedToken.Method.(*jwt.SigningMethodRSA); !ok {
			msg := fmt.Errorf("Unexpected signing method: %v", decodedToken.Header["alg"])
			return claims, msg
		}
		return j.key, nil
	})
	if err != nil {
		return claims, err
	}
	if decodedToken == nil || !decodedToken.Valid {
		return claims, fmt.Errorf("Unauthorized access")
	}
	return *decodedToken.Claims.(*CustomClaims), nil
}
