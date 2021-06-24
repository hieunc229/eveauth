package eveauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"summer/modules/utils"

	"github.com/boltdb/bolt"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var user UserPayload
	json.NewDecoder(r.Body).Decode(&user)

	if user.Password == "" || user.Username == "" {
		utils.HandleError(w, errors.New("data can not empty"))
		return
	}

	db, err := GetDB()

	if err != nil {
		utils.HandleError(w, errors.New("can't load data"))
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
		utils.HandleError(w, errors.New("invalid login"))
		return
	}

	token, err := createJWTToken(user.Username)

	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.HandleData(w, map[string]interface{}{
		"token": token,
	})
}
