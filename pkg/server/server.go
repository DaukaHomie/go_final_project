package server

import (
	"fmt"
	"go-final-project/pkg/api"
	"net/http"
	"os"
)

const defaultPort = "7540"

// ResolvePort возвращает порт из переменной окружения TODO_PORT
// или дефолтный порт, если переменная не задана.
func ResolvePort() string {
	if port := os.Getenv("TODO_PORT"); port != "" {
		return port
	}
	return defaultPort
}

// NewHandler создает HTTP-мультиплексор, который раздаёт статические файлы
// из webDir и подключает API-обработчики.
func NewHandler(webDir string) (http.Handler, error) {
	if _, err := os.Stat(webDir); err != nil {
		return nil, fmt.Errorf("web dir not found: %s (%w)", webDir, err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init(mux)
	return mux, nil
}

// Run запускает HTTP-сервер на указанном порту.
func Run(webDir, port string) error {
	handler, err := NewHandler(webDir)
	if err != nil {
		return err
	}

	addr := ":" + port
	fmt.Printf("Server is listening on http://localhost%s/ (web dir: %s)\n", addr, webDir)

	return http.ListenAndServe(addr, handler)
}
