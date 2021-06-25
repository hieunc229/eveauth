package eveauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"

	"github.com/boltdb/bolt"
)

type changePassword struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`

	// Replace current token with a new one
	ChangeToken bool `json:"change_token"`

	// Remove all existing token, add a new one.
	// Used for logout of all devices feature
	ClearTokens bool `json:"clear_tokens"`
}

type changePasswordPayload struct {
	Data changePassword `json:"data"`
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {

	jwtPayload, userData, token, err := VerifyRequest(r)

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

	var newToken string

	err = db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(AuthBucketName)

		if err = checkPasswordHash(userData.HashedPassword, user.Password); err != nil {
			return err
		}

		newHashedPassword, err := hashPassword(user.NewPassword)
		if err == nil {
			userData.HashedPassword = newHashedPassword

			if user.ChangeToken {

				tokens := userData.Tokens

				if user.ClearTokens {
					tokens = []string{}
				} else {
					tokenIndex := sort.SearchStrings(tokens, token)

					if tokenIndex != -1 {
						tokens = append(tokens[:tokenIndex], tokens[tokenIndex+1:]...)
					}

					newToken, err = createJWTToken(username)

					if err != nil {
						return err
					}
				}

				userData.Tokens = append(tokens, newToken)
			}

			return setUserData(bucket, username, *userData)
		}

		return err
	})

	defer db.Close()

	if err != nil {
		handleError(w, errors.New("invalid login"))
		return
	}

	if err != nil {
		handleError(w, err)
		return
	}

	handleData(w, map[string]interface{}{
		"new_token": newToken,
	})
}
