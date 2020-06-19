package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const BcryptCost = 12

func EncryptPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), BcryptCost)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func ValidateCredentials(hash, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}
