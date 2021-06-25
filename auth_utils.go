package eveauth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func validatePassword(password string) error {

	if len(password) < 8 {
		return errors.New("password must has 8 or more characters")
	}

	// TODO: must include uppercase, special character?, number?

	return nil
}
