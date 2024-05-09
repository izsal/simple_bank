package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// This function takes a password as input and returns a hashed version of the password along with any
// potential errors.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashedPassword), nil
}

// This function is used to check if a given password matches a hashed password.
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
