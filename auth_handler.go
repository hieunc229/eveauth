package eveauth

import (
	"net/http"
)

func AuthHandler(nextFunc func(http.ResponseWriter, *http.Request), options *AuthHandlerOptions) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, _, _, err := VerifyRequest(r, options); err != nil {
			handleError(w, err)
			return
		}

		nextFunc(w, r)
	}
}
