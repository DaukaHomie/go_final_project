package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-final-project/pkg/db"
)

// doneTaskHandler отмечает задачу выполненной.
// Если задача одноразовая — она удаляется.
// Если повторяющаяся — пересчитывается дата следующего выполнения.
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "задача не найдена",
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("ошибка выборки: %v", err),
		})
		return
	}

	// Если задача без повторений — удаляем.
	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{})
		return
	}

	// Для повторяющейся задачи пересчитываем дату.
	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("некорректное правило повторения: %v", err),
		})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
