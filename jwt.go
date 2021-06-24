package eveauth

import (
	"errors"
	"os"
	"time"
)

type JWTPayload struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

var (
	SECRET_JWT      = []byte(os.Getenv("SECRET_JWT"))
	ErrInvalidToken = errors.New("unauthorized access (invalid token)")
	ErrExpiredToken = errors.New("unauthorized access (expired token)")
)

func (payload *JWTPayload) Valid() error {

	if time.Now().After(payload.ExpiredAt) {
		return errors.New("token has expired")
	}
	return nil
}
