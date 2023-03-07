package usecase

import (
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const salt = "lY5Voj5HpOzFfFFE0K6FpuvwWU"

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
