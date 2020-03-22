package controller

import (
	"golang.org/x/crypto/bcrypt"
)

func EncodePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	} else {
		return string(hash), nil
	}
}

func ComparePassword(p1 string, p2 string) error {
	return bcrypt.CompareHashAndPassword([]byte(p1), []byte(p2))
}
