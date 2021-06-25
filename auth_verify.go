package eveauth

import (
	"errors"
	"net/http"
	email "net/mail"
	"sort"
	"strings"

	"github.com/gorilla/context"
)

type AuthHandlerOptions struct {
	Role string `json:"role"`

	// Set to true to forbid access from users with different RoleLevel.
	// Or set to false (or nil) to forbid only users with lower RoleLevel
	RoleExact bool `json:"role_exact"`
}

/*
Verify a http.Request by (1) get bearer token, (2) verify if the token is a valid jwt token,
(3) get userData then check if token is still active
(4) then check if the user has the proper role if authOption != nil
*/
func VerifyRequest(r *http.Request, options *AuthHandlerOptions) (*JWTPayload, *userData, string, error) {

	var userData userData
	bearer := r.Header.Get("Authorization")
	token := getToken(bearer)

	// Request allow anonymous access
	if options == nil || options.Role != "" || UserRoleLevels[options.Role] == UserRoleLevels[RoleAnonymous] {
		return nil, nil, token, nil
	}

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

	if options.RoleExact {
		if UserRoleLevels[userData.Role] != UserRoleLevels[options.Role] {
			return nil, nil, token, errors.New("unauthorized access")
		}
	}

	if UserRoleLevels[userData.Role] < UserRoleLevels[options.Role] {
		return nil, nil, token, errors.New("unauthorized access")
	}

	context.Set(r, "payload", payload)
	context.Set(r, "userData", userData)

	return payload, &userData, token, err
}

func getToken(authStr string) string {
	return strings.Replace(authStr, "Bearer ", "", -1)
}

func validateUserInput(user *UserPayload) error {

	var err error

	errStr, err := validateUserNamePassword(user.Username, user.Password)

	if _, err = email.ParseAddress(user.Email); err != nil {
		errStr += "email, "
	}

	if errStr != "" {
		err = errors.New("invalid " + errStr)
	}

	return err
}

func validateUserNamePassword(username string, password string) (string, error) {

	var err error
	errStr := ""

	if username == "" || len(username) < 4 {
		errStr += "username, "
	}

	if password == "" || validatePassword(password) != nil {
		errStr += "password, "
	}

	if errStr != "" {
		err = errors.New("invalid " + errStr)
	}

	return errStr, err
}
