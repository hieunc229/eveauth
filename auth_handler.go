package eveauth

import (
	"net/http"
	"summer/modules/utils"
)

func HandleAuth(nextFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := verifyRequestToken(r); err == nil {
			utils.HandleError(w, err)
			return
		}
		nextFunc(w, r)
	}
}
