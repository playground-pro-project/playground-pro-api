package helper

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const DefaultSalt = 8

func HashPass(p string) string {
	password := []byte(strings.TrimSpace(p))
	hash, _ := bcrypt.GenerateFromPassword(password, DefaultSalt)

	return string(hash)
}

func ComparePass(h, p []byte) error {
	hash, pass := []byte(h), []byte(p)

	err := bcrypt.CompareHashAndPassword(hash, pass)

	return err
}
