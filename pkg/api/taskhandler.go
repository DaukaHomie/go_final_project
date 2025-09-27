package api

import "net/http"

// taskHandler обрабатывает CRUD операции над задачами.
//   - POST   /api/task    → создание новой задачи
//   - GET    /api/task?id → получение задачи по ID
//   - PUT    /api/task    → редактирование задачи
//   - DELETE /api/task?id → удаление задачи
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		w.Header().Set("Allow", "POST, GET, PUT, DELETE")
		writeJSON(w, http.StatusMethodNotAllowed,
			map[string]string{"error": "метод не поддерживается"})
	}
}
