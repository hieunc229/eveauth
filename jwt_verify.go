package eveauth

import (
	"github.com/dgrijalva/jwt-go"
)

func verifyToken(token string) (*JWTPayload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return SECRET_JWT, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyFunc)

	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(*JWTPayload); ok && jwtToken.Valid {
		return claims, nil
	} else {
		return nil, ErrExpiredToken
	}

	// if err != nil {
	// 	verr, ok := err.(*jwt.ValidationError)
	// 	if ok && errors.Is(verr.Inner, ErrExpiredToken) {
	// 		return nil, ErrExpiredToken
	// 	}
	// 	return nil, ErrInvalidToken
	// }

	// payload, ok := jwtToken.Claims.(*JWTPayload)
	// if !ok {
	// 	return nil, ErrInvalidToken
	// }

	// return payload, nil

}
