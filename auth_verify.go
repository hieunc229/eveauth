package eveauth

import (
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/context"
)

func VerifyRequest(r *http.Request) (*JWTPayload, *userData, string, error) {

	var userData userData
	bearer := r.Header.Get("Authorization")
	token := getToken(bearer)
	payload, err := verifyToken(token)

	if err != nil {
		return nil, nil, token, err
	}

	userData, err = getUserData(nil, payload.Username)

	if err != nil {
		return nil, nil, token, err
	}

	if sort.SearchStrings(userData.Tokens, token) == -1 {
		return nil, nil, token, errors.New("invalid token")
	}

	context.Set(r, "payload", payload)
	context.Set(r, "userData", userData)

	return payload, &userData, token, err
}

func getToken(authStr string) string {
	return strings.Replace(authStr, "Bearer ", "", -1)
}
