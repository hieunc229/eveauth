package eveauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
)

type changePassword struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

type changePasswordPayload struct {
	Data changePassword `json:"data"`
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {

	var payload changePasswordPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		handleError(w, errors.New("can't load data"))
		return
	}

	user := payload.Data

	if user.Password == "" || user.Username == "" || user.NewPassword == "" {
		handleError(w, errors.New("data can not empty"))
		return
	}

	if err = validatePassword(user.NewPassword); err != nil {
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

		if err = checkPasswordHash(userData.HashedPassword, user.Password); err != nil {
			return err
		}

		newHashedPassword, err := hashPassword(user.NewPassword)
		if err == nil {
			userData.HashedPassword = newHashedPassword
			setUserData(bucket, user.Username, userData)
		}

		return err
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
