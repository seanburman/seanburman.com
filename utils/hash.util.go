package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashString(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), 9)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(hash string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
