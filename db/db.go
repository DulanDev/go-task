package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbFileName string) error {
	var err error
	DB, err = sql.Open("sqlite", dbFileName)
	if err != nil {
		return err
	}

	if _, err := DB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	createPriorities := `
    CREATE TABLE IF NOT EXISTS priorities (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL,
        level INTEGER NOT NULL
    );`
	if _, err := DB.Exec(createPriorities); err != nil {
		return fmt.Errorf("failed to create priorities table: %w", err)
	}

	seedPriorities := `
    INSERT OR IGNORE INTO priorities (name, level) VALUES 
    ('High', 3), 
    ('Medium', 2), 
    ('Low', 1);`
	if _, err := DB.Exec(seedPriorities); err != nil {
		return fmt.Errorf("failed to seed priorities: %w", err)
	}

	createTags := `
    CREATE TABLE IF NOT EXISTS tags (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL
    );`
	if _, err := DB.Exec(createTags); err != nil {
		return fmt.Errorf("failed to create tags table: %w", err)
	}

	createTasks := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT,
        completed BOOLEAN DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        due_date DATETIME,
        priority_id INTEGER,
        FOREIGN KEY(priority_id) REFERENCES priorities(id)
    );`

	if _, err := DB.Exec(createTasks); err != nil {
		return fmt.Errorf("failed to create tasks table: %w", err)
	}

	migrateTasksTable(DB)

	createTaskTags := `
    CREATE TABLE IF NOT EXISTS task_tags (
        task_id INTEGER,
        tag_id INTEGER,
        PRIMARY KEY (task_id, tag_id),
        FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE,
        FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE
    );`
	if _, err := DB.Exec(createTaskTags); err != nil {
		return fmt.Errorf("failed to create task_tags table: %w", err)
	}

	return nil
}

func migrateTasksTable(db *sql.DB) {
	_, _ = db.Exec("ALTER TABLE tasks ADD COLUMN due_date DATETIME;")
	_, _ = db.Exec("ALTER TABLE tasks ADD COLUMN priority_id INTEGER REFERENCES priorities(id);")
}

func CloseDB() {
    if err := DB.Close(); err != nil {
        fmt.Printf("failed to close database: %v", err)
    }
}