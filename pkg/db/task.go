package db

import "fmt"

// Task описывает задачу планировщика.
type Task struct {
	ID      int64  `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// AddTask добавляет новую задачу и возвращает её ID.
func AddTask(t *Task) (int64, error) {
	const q = `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?);
	`

	res, err := DB.Exec(q, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, fmt.Errorf("add task: %w", err)
	}
	return res.LastInsertId()
}

// Tasks возвращает список задач с ограничением по количеству.
func Tasks(limit int) ([]*Task, error) {
	const q = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC, id ASC
		LIMIT ?;
	`

	rows, err := DB.Query(q, limit)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		t := new(Task)
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	if tasks == nil {
		return make([]*Task, 0), nil
	}
	return tasks, nil
}

// GetTask возвращает задачу по её ID.
func GetTask(id string) (*Task, error) {
	const q = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?;
	`

	var t Task
	if err := DB.QueryRow(q, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	return &t, nil
}

// UpdateTask обновляет все поля задачи.
func UpdateTask(t *Task) error {
	const q = `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?;
	`

	res, err := DB.Exec(q, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// DeleteTask удаляет задачу по её ID.
func DeleteTask(id string) error {
	const q = `DELETE FROM scheduler WHERE id = ?;`

	res, err := DB.Exec(q, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// UpdateDate изменяет только дату выполнения задачи.
func UpdateDate(next, id string) error {
	const q = `UPDATE scheduler SET date = ? WHERE id = ?;`

	res, err := DB.Exec(q, next, id)
	if err != nil {
		return fmt.Errorf("update date: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
