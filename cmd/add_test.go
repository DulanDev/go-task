package cmd

import (
	"go-task/db"
	"go-task/models"
	"testing"
)

func TestAddTask(t *testing.T) {
    db.InitDB("./test_tasks.db")
    defer db.CloseDB()
    defer db.DB.Exec("DROP TABLE IF EXISTS tasks")

	title := "Test Task"
	description := "This is a test task"
	AddTask(title, description, "Medium", "")

    var task models.Task
    err := db.DB.QueryRow("SELECT id, title, description FROM tasks WHERE title = ?", title).Scan(&task.ID, &task.Title, &task.Description)
    if err != nil {
        t.Fatalf("Failed to query task: %v", err)
    }

	if task.Title != title {
		t.Errorf("Expected title %s, got %s", title, task.Title)
	}
	if task.Description != description {
		t.Errorf("Expected description %s, got %s", description, task.Description)
	}
}

func TestAddRichTask(t *testing.T) {
	db.InitDB("./test_rich_tasks.db")
	defer db.CloseDB()
	defer db.DB.Exec("DROP TABLE IF EXISTS tasks")
	defer db.DB.Exec("DROP TABLE IF EXISTS priorities")
	defer db.DB.Exec("DROP TABLE IF EXISTS tags")
	defer db.DB.Exec("DROP TABLE IF EXISTS task_tags")

	title := "Rich Task"
	desc := "Task with priority and tags"
	priority := "High"
	tags := "work, urgent"

	AddTask(title, desc, priority, tags)

	var taskID int
	var priorityID int
	err := db.DB.QueryRow("SELECT id, priority_id FROM tasks WHERE title = ?", title).Scan(&taskID, &priorityID)
	if err != nil {
		t.Fatalf("Failed to query task: %v", err)
	}

	var level int
	err = db.DB.QueryRow("SELECT level FROM priorities WHERE id = ?", priorityID).Scan(&level)
	if err != nil {
		t.Fatalf("Failed to query priority: %v", err)
	}
	if level != 3 {
		t.Errorf("Expected priority level 3 (High), got %d", level)
	}

	rows, err := db.DB.Query(`
		SELECT t.name FROM tags t
		JOIN task_tags tt ON t.id = tt.tag_id
		WHERE tt.task_id = ?`, taskID)
	if err != nil {
		t.Fatalf("Failed to query tags: %v", err)
	}
	defer rows.Close()

	foundTags := make(map[string]bool)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		foundTags[name] = true
	}

	if !foundTags["work"] {
		t.Errorf("Expected tag 'work' to be found")
	}
	if !foundTags["urgent"] {
		t.Errorf("Expected tag 'urgent' to be found")
	}
}