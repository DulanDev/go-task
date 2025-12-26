package cmd

import (
	"go-task/db"
	"testing"
)

func TestSearchTasks(t *testing.T) {
    db.InitDB("./test_tasks.db")
    defer db.CloseDB()
    defer db.DB.Exec("DROP TABLE IF EXISTS tasks")

	title := "Task to Search"
	description := "This task will be searched"
	AddTask(title, description, "Medium", "")

    tasks, err := SearchTasks("Search")
    if err != nil {
        t.Fatalf("Failed to search tasks: %v", err)
    }

    found := false
    for _, task := range tasks {
        if task.Title == title && task.Description == description {
            found = true
            break
        }
    }
    if !found {
        t.Errorf("Expected task to be found")
    }
}