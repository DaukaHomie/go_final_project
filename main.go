package main

import (
	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = filepath.Join(".", "scheduler.db")
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal("db init: ", err)
	}
	defer db.Close()

	webDir := filepath.Join(".", "web")
	if _, err := os.Stat(webDir); err != nil {
		log.Printf("web dir not found: %s", webDir)
	}

	port := server.ResolvePort()

	log.Printf("server start: http://localhost:%s/ (db: %s, web: %s)", port, dbFile, webDir)
	if err := server.Run(webDir, port); err != nil {
		log.Fatal("server run: ", err)
	}
}
