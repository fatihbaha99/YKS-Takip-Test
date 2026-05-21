package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	DB.SetMaxOpenConns(1)

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("[DB] SQLite bağlantısı kuruldu:", dbPath)
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
