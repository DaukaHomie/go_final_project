package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON отправляет JSON-ответ с указанным статусом.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status > 0 {
		w.WriteHeader(status)
	}
	_ = json.NewEncoder(w).Encode(data)
}
