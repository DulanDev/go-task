package cmd

import (
	"database/sql"
	"fmt"
	"go-task/db"
	"go-task/models"
	"os"
	"text/tabwriter"
)

// retrieves all tasks from the database and returns them.
func ListTasks() ([]models.Task, error) {
	query := `
    SELECT 
        t.id, 
        t.title, 
        t.description, 
        t.completed, 
        t.created_at,
        COALESCE(p.name, 'Medium') as priority_name,
        GROUP_CONCAT(tg.name, ', ') as tags
    FROM tasks t
    LEFT JOIN priorities p ON t.priority_id = p.id
    LEFT JOIN task_tags tt ON t.id = tt.task_id
    LEFT JOIN tags tg ON tt.tag_id = tg.id
    GROUP BY t.id
    ORDER BY p.level DESC, t.created_at ASC
    `

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listing tasks: %v", err)
	}
	defer rows.Close()

	var tasks []models.Task
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tPRIORITY\tSTATUS\tTAGS\tDESCRIPTION")
	fmt.Fprintln(w, "--\t-----\t--------\t------\t----\t-----------")

	for rows.Next() {
		var task models.Task
		var priorityName string
		var tagsStr sql.NullString

		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &priorityName, &tagsStr); err != nil {
			return nil, err
		}

		task.Priority = &models.Priority{Name: priorityName}

		completedStatus := "❌"
		if task.Completed.Bool {
			completedStatus = "✅"
		}

		displayTags := ""
		if tagsStr.Valid {
			displayTags = tagsStr.String
		} else {
			displayTags = "-"
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n", task.ID, task.Title, priorityName, completedStatus, displayTags, task.Description)
		tasks = append(tasks, task)
	}
	w.Flush()
	return tasks, nil
}