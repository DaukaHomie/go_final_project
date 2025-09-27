package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// DB — глобальное подключение к базе данных.
var DB *sql.DB

// схема таблицы для задач
const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

// Init открывает подключение к SQLite, применяет схему при необходимости.
func Init(dbFile string) error {
	install := false
	if _, err := os.Stat(dbFile); err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}

	// если что-то пошло не так — всегда закрываем соединение
	defer func() {
		if err != nil {
			_ = db.Close()
		}
	}()

	if err = db.Ping(); err != nil {
		return fmt.Errorf("ping sqlite: %w", err)
	}

	if _, err = db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("set pragma foreign_keys: %w", err)
	}
	if _, err = db.Exec(`PRAGMA busy_timeout = 5000;`); err != nil {
		return fmt.Errorf("set pragma busy_timeout: %w", err)
	}

	// если база новая — применяем схему
	if install {
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}
		if _, err := tx.Exec(schema); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply schema: %w", err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit schema: %w", err)
		}
	}

	DB = db
	return nil
}

// Close закрывает подключение к базе.
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
