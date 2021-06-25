package eveauth

import (
	"encoding/json"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(&map[string]interface{}{
		"ok":    false,
		"error": err.Error(),
	})
}

func handleData(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(&map[string]interface{}{
		"ok":   true,
		"data": data,
	})
}
