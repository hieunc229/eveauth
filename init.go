package eveauth

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

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

	var envPath string
	var err error
	if envPath, err = filepath.Abs(filepath.Dir("./")); err != nil {
		fmt.Println("Unable to get .env path", err)
		return
	}

	if err = godotenv.Load(envPath + "/.env"); err != nil {
		fmt.Println("[eveauth] Unable to load .env. ", err)
		return
	}

	if os.Getenv("EVEAUTH_JWT_SECRET") != "" {
		JWT_SECRET = os.Getenv("EVEAUTH_JWT_SECRET")
		fmt.Println("[eveauth] loaded EVEAUTH_JWT_SECRET")
	}

	if os.Getenv("EVEAUTH_PATH") != "" {
		AUTH_PATH = os.Getenv("EVEAUTH_PATH")
		fmt.Println("[eveauth] loaded EVEAUTH_PATH")
	}
}
