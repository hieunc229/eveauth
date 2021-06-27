package eveauth

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

// userData struct is used for storing data in db
type userData struct {
	HashedPassword string `json:"password"`
	Email          string `json:"email"`
	Role           string `json:"role"`

	// Valid tokens
	Tokens []string `json:"tokens"`
}

func setUserData(bucket *bolt.Bucket, username string, user userData) error {
	userDataRawValue, err := json.Marshal(user)

	if err != nil {
		return err
	}

	return bucket.Put([]byte(username), userDataRawValue)
}

func getUserData(bucket *bolt.Bucket, username string) (userData, error) {

	var err error
	var userData userData
	var db *bolt.DB

	if bucket == nil {
		db, err = getDB()

		if err != nil {
			return userData, err
		}

		tx, err := db.Begin(false)

		if err != nil {
			return userData, err
		}

		bucket = tx.Bucket(AuthBucketName)
	}

	userDataRawValue := bucket.Get([]byte(username))

	if userDataRawValue == nil {
		return userData, errors.New("username doesn't exist")
	}

	if db != nil {
		defer db.Close()
	}

	err = json.Unmarshal(userDataRawValue, &userData)

	return userData, err
}
