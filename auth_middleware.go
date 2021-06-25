package eveauth

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if _, _, _, err := VerifyRequest(r); err == nil {
			handleError(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
