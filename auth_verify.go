package eveauth

import (
	"net/http"
	"strings"

	"github.com/gorilla/context"
)

func VerifyRequest(r *http.Request) (*JWTPayload, error) {
	bearer := r.Header.Get("Authorization")
	token := getToken(bearer)
	payload, err := verifyToken(token)

	if err != nil {
		return nil, err
	}
	context.Set(r, "payload", payload)
	return payload, err
}

func getToken(authStr string) string {
	return strings.Replace(authStr, "Bearer ", "", -1)
}