package eveauth

import (
	"net/http"
	"strings"
	"summer/modules/utils"

	"github.com/gorilla/context"
)

func HandleAuth(nextFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		token := getToken(bearer)
		payload, err := verifyToken(token)

		if err != nil {
			utils.HandleError(w, err)
		} else {
			context.Set(r, "payload", payload)
			nextFunc(w, r)
		}
	}
}

func getToken(authStr string) string {
	return strings.Replace(authStr, "Bearer ", "", -1)
}
