package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// Generate a hash.
func Generate(password *string) (string, error) {
	hex := []byte(*password)
	hashedPassword, err := bcrypt.GenerateFromPassword(hex, 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Compare a hash.
func Compare(digest []byte, password *string) bool {
	hex := []byte(*password)
	if err := bcrypt.CompareHashAndPassword(digest, hex); err == nil {
		return true
	}
	return false
}
