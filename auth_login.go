package eveauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
)

type loginPayload struct {
	Data UserPayload `json:"data"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var payload loginPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		handleError(w, errors.New("can't load data"))
		return
	}

	user := payload.Data

	if user.Password == "" || user.Username == "" {
		handleError(w, errors.New("data can not empty"))
		return
	}

	if err = validatePassword(user.Password); err != nil {
		handleError(w, err)
		return
	}

	db, err := getDB()

	if err != nil {
		handleError(w, errors.New("can't load data"))
		return
	}

	err = db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(AuthBucketName)

		userData, err := getUserData(bucket, user.Username)
		if err != nil {
			return err
		}

		return checkPasswordHash(userData.HashedPassword, user.Password)

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
