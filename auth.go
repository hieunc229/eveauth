package eveauth

type UserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

var AuthBucketName = []byte("auth")
