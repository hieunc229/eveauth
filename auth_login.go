package eveauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

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

	err = db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(AuthBucketName)

		rawValue := bucket.Get([]byte(user.Username))

		if rawValue == nil {
			return errors.New("username doesn't exist")
		}

		return checkPasswordHash(string(rawValue), user.Password)

	})

	defer db.Close()

	if err != nil {
		handleError(w, errors.New("invalid login"))
		return
	}

	token, err := createJWTToken(user.Username)

	if err != nil {
		handleError(w, err)
		return
	}

	handleData(w, map[string]interface{}{
		"token": token,
	})
}
