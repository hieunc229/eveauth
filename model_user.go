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
}

func setUserData(bucket *bolt.Bucket, username string, user userData) error {
	userDataRawValue, err := json.Marshal(user)

	if err != nil {
		return err
	}

	return bucket.Put([]byte(username), userDataRawValue)
}

func getUserData(bucket *bolt.Bucket, username string) (userData, error) {

	var userData userData
	userDataRawValue := bucket.Get([]byte(username))

	if userDataRawValue == nil {
		return userData, errors.New("username doesn't exist")
	}

	err := json.Unmarshal(userDataRawValue, &userData)

	return userData, err
}
