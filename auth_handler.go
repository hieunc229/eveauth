package eveauth

import (
	"net/http"
)

func AuthHandler(nextFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := VerifyRequest(r); err == nil {
			handleError(w, err)
			return
		}
		nextFunc(w, r)
	}
}
