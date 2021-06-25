package eveauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
)

var AuthBucketName = []byte("auth")

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var user UserPayload
	json.NewDecoder(r.Body).Decode(&user)

	if user.Password == "" || user.Username == "" {
		handleError(w, errors.New("data can not empty"))
		return
	}

	db, err := getDB()

	if err != nil {
		handleError(w, errors.New("can't load data"))
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(AuthBucketName)

		if bucket == nil {
			bucket, err = tx.CreateBucket(AuthBucketName)
			if err != nil {
				return err
			}
		}

		existingUser := bucket.Get([]byte(user.Username))

		if existingUser != nil {
			return errors.New("invalid usersname")
		}

		password, err := hashPassword(user.Password)

		if err != nil {
			return err
		}

		return bucket.Put([]byte(user.Username), []byte(password))
	})

	defer db.Close()

	if err != nil {
		handleError(w, err)
		return
	}

	handleData(w, map[string]interface{}{
		"success": true,
	})
}
