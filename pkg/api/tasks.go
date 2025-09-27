package api

import (
	"fmt"
	"net/http"
	"strconv"

	"go-final-project/pkg/db"
)

// taskDTO описывает задачу в ответе API.
type taskDTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// tasksResp оборачивает список задач в JSON.
type tasksResp struct {
	Tasks []taskDTO `json:"tasks"`
}

// tasksHandler возвращает список задач.
// GET /api/tasks?limit=N
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "метод не поддерживается"})
		return
	}

	// Читаем limit из query, по умолчанию 50
	limit := 50
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 && n <= 500 {
			limit = n
		}
	}

	rows, err := db.Tasks(limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": fmt.Sprintf("ошибка выборки: %v", err)})
		return
	}

	resp := tasksResp{Tasks: make([]taskDTO, 0, len(rows))}
	for _, t := range rows {
		resp.Tasks = append(resp.Tasks, taskDTO{
			ID:      strconv.FormatInt(t.ID, 10),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}
