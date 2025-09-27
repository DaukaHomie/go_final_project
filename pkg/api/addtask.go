package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-final-project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("ошибка разбора json: %v", err)})

		return
	}

	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})

		return
	}

	if err := normalizeTaskDate(&t); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})

		return
	}

	id, err := db.AddTask(&t)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("ошибка сохранения: %v", err)})

		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func normalizeTaskDate(t *db.Task) error {
	const layout = DateLayout

	now := time.Now()

	t.Date = strings.TrimSpace(t.Date)

	if t.Date == "" {
		t.Date = now.Format(layout)
	}

	parsed, err := time.Parse(layout, t.Date)

	if err != nil {
		return errors.New("дата должна быть в формате 20060102")
	}

	rep := strings.TrimSpace(t.Repeat)

	var next string

	if rep != "" {
		next, err = NextDate(now, t.Date, rep)
		if err != nil {
			return fmt.Errorf("некорректное правило повторения: %w", err)
		}
	}

	if now.Format(layout) > parsed.Format(layout) {
		if rep == "" {
			t.Date = now.Format(layout)
		} else {
			t.Date = next
		}
	}

	return nil
}