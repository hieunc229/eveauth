package eveauth

import "github.com/boltdb/bolt"

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
var UserRoleLevels = map[string]int{
	RoleAnonymous: 0,
	RoleMember:    1,
	RoleMember:    2,
}

func initateAuthBucket(tx *bolt.Tx) (*bolt.Bucket, error) {

	bucket, err := tx.CreateBucket(AuthBucketName)

	if err != nil {
		return nil, err
	}

	return bucket, err
}
