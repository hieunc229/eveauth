package eveauth

import (
	"flag"
	"os"
)

type UserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

var AuthBucketName = []byte("auth")
var AuthTokenBucketName = []byte("auth_tokens")

var (
	RoleAnonymous = "anonymous"
	RoleMember    = "member"
	RoleAdmin     = "admin"
)

var (
	JWT_SECRET string
	AUTH_PATH  string
)
var UserRoleLevels = map[string]int{
	RoleAnonymous: 0,
	RoleMember:    1,
	RoleAdmin:     2,
}

func init() {

	JWT_SECRET = *flag.String("EVEAUTH_JWT_SECRET", "eveauth", "string")
	AUTH_PATH = *flag.String("EVEAUTH_PATH", "auth", "string")
	flag.Parse()

	if os.Getenv("EVEAUTH_JWT_SECRET") != "" {
		JWT_SECRET = os.Getenv("EVEAUTH_JWT_SECRET")
	}

	if os.Getenv("EVEAUTH_AUTH_PATH") != "" {
		AUTH_PATH = os.Getenv("EVEAUTH_AUTH_PATH")
	}
}
