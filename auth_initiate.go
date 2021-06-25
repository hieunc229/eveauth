package eveauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
)

func InitiateAuthHandler(w http.ResponseWriter, r *http.Request) {

	db, err := getDB()

	if err != nil {
		handleError(w, err)
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {

		var payload registerPayload
		err = json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			return err
		}

		user := payload.Data

		if err = validateUserInput(&user); err != nil {
			return err
		}

		password, err := hashPassword(user.Password)

		if err != nil {
			return err
		}

		bucket, err := tx.CreateBucket(AuthBucketName)

		if err != nil {
			return errors.New("auth has already initiated")
		}

		return setUserData(bucket, user.Username, userData{
			HashedPassword: password,
			Email:          user.Email,
			Role:           "admin",
		})
	})

	defer db.Close()
	if err != nil {
		handleError(w, err)
		return
	}

	handleData(w, true)
}
