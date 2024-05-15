package accountsio

import (
	"encoding/json"
	"net/http"
)

func MakeResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
