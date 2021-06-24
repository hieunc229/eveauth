package eveauth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
)

const numOfDates = 30

func createJWTToken(username string) (string, error) {

	payload, err := newPayload(username)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString(SECRET_JWT)
}

func newPayload(username string) (*JWTPayload, error) {
	tokenID := xid.New().String()

	payload := &JWTPayload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().AddDate(0, 0, numOfDates),
	}
	return payload, nil
}
