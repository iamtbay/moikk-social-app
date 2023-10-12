package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHasher(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", nil
	}
	return string(hashedPassword), nil
}

func PasswordChecker(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("invalid credentials")
	}
	return err
}
