package eveauth

import (
	"net/http"
	"summer/modules/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if _, err := VerifyRequest(r); err == nil {
			utils.HandleError(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
