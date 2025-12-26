package cmd

import (
	"fmt"
	"go-task/db"
	"go-task/models"
)

// searches for tasks that match the given keyword and returns them.
func SearchTasks(keyword string) ([]models.Task, error) {
    rows, err := db.DB.Query("SELECT id, title, description FROM tasks WHERE title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
    if err != nil {
        return nil, fmt.Errorf("error searching tasks: %v", err)
    }
    defer rows.Close()

    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        if err := rows.Scan(&task.ID, &task.Title, &task.Description); err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
        fmt.Printf("%d: %s - %s\n", task.ID, task.Title, task.Description)
    }
    return tasks, nil
}