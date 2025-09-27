package api

import (
	"net/http"
	"strings"

	"go-final-project/pkg/db"
)

// deleteTaskHandler удаляет задачу по её ID.
// Метод: DELETE, параметр: id (query).
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "метод не поддерживается",
		})
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "не указан идентификатор",
		})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
