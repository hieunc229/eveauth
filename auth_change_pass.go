package eveauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
)

type changePassword struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

type changePasswordPayload struct {
	Data changePassword `json:"data"`
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {

	jwtPayload, err := VerifyRequest(r)

	if err != nil {
		handleError(w, err)
		return
	}

	username := jwtPayload.Username

	var payload changePasswordPayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		handleError(w, errors.New("can't load data"))
		return
	}

	user := payload.Data

	if user.Password == "" || user.NewPassword == "" {
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
		userData, err := getUserData(bucket, username)

		if err != nil {
			return err
		}

		if err = checkPasswordHash(userData.HashedPassword, user.Password); err != nil {
			return err
		}

		newHashedPassword, err := hashPassword(user.NewPassword)
		if err == nil {
			userData.HashedPassword = newHashedPassword
			setUserData(bucket, username, userData)
		}

		return err
	})

	defer db.Close()

	if err != nil {
		handleError(w, errors.New("invalid login"))
		return
	}

	token, err := createJWTToken(username)

	if err != nil {
		handleError(w, err)
		return
	}

	handleData(w, map[string]interface{}{
		"token": token,
	})
}
