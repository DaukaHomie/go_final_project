package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go-final-project/pkg/db"
)

type editPayload struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// editTaskHandler обрабатывает PUT-запрос для изменения задачи.
func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "метод не поддерживается",
		})
		return
	}
	defer r.Body.Close()

	var p editPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("ошибка разбора json: %v", err),
		})
		return
	}

	p.ID = strings.TrimSpace(p.ID)
	if p.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "не указан идентификатор",
		})
		return
	}

	title := strings.TrimSpace(p.Title)
	if title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "не указан заголовок задачи",
		})
		return
	}

	parsedID, err := strconv.ParseInt(p.ID, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "некорректный идентификатор",
		})
		return
	}

	task := db.Task{
		ID:      parsedID,
		Date:    p.Date,
		Title:   title,
		Comment: p.Comment,
		Repeat:  p.Repeat,
	}

	if err := normalizeTaskDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
