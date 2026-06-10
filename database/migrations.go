package database

import "log"

func RunMigrations() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			telegram_chat_id INTEGER DEFAULT 0,
			activation_code TEXT DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			purge_at DATETIME NOT NULL,
			active INTEGER DEFAULT 1
		)`,
		`CREATE TABLE IF NOT EXISTS study_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			subject TEXT NOT NULL,
			topic TEXT NOT NULL,
			study_type TEXT NOT NULL CHECK(study_type IN ('goz_gezdir','test_coz','genel_test')),
			stars INTEGER DEFAULT 0,
			correct INTEGER DEFAULT 0,
			wrong INTEGER DEFAULT 0,
			net REAL DEFAULT 0,
			study_date TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			subject TEXT NOT NULL,
			topic TEXT NOT NULL,
			todo_type TEXT NOT NULL CHECK(todo_type IN ('goz_gezdir','test_coz','genel_test')),
			due_date TEXT NOT NULL,
			completed INTEGER DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS exams (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			exam_type TEXT NOT NULL CHECK(exam_type IN ('TYT','AYT')),
			title TEXT DEFAULT '',
			exam_date TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS exam_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			exam_id INTEGER NOT NULL,
			subject TEXT NOT NULL,
			correct INTEGER DEFAULT 0,
			wrong INTEGER DEFAULT 0,
			net REAL DEFAULT 0,
			FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE
		)`,
	}

	for _, stmt := range statements {
		if _, err := DB.Exec(stmt); err != nil {
			return err
		}
	}

	alterStatements := []string{
		"ALTER TABLE users ADD COLUMN reminder_hour INTEGER DEFAULT 8",
		"ALTER TABLE users ADD COLUMN reminder_minute INTEGER DEFAULT 0",
		"ALTER TABLE users ADD COLUMN last_reminder_date TEXT DEFAULT ''",
	}
	for _, stmt := range alterStatements {
		DB.Exec(stmt)
	}

	log.Println("[DB] Migrationlar tamamlandı")
	return nil
}
