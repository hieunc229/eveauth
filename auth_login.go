package eveauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
)

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginPayload struct {
	Data loginData `json:"data"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var payload loginPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		handleError(w, errors.New("can't decode login info"))
		return
	}

	user := payload.Data

	if _, err = validateUserNamePassword(user.Username, user.Password); err != nil {
		handleError(w, err)
		return
	}

	if err = validatePassword(user.Password); err != nil {
		handleError(w, err)
		return
	}

	db, err := getDB()

	if err != nil {
		handleError(w, errors.New("database is unavailable"))
		return
	}

	var token string

	err = db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(AuthBucketName)

		userData, err := getUserData(bucket, user.Username)
		if err != nil {
			return err
		}

		if err = checkPasswordHash(userData.HashedPassword, user.Password); err != nil {
			return errors.New("invalid password")
		}

		if token, err = createJWTToken(user.Username); err != nil {
			return err
		}

		userData.Tokens = append(userData.Tokens, token)
		return setUserData(bucket, user.Username, userData)
	})

	defer db.Close()

	if err != nil {
		handleError(w, err)
		return
	}

	if err != nil {
		handleError(w, err)
		return
	}

	handleData(w, map[string]interface{}{
		"token": token,
	})
}
