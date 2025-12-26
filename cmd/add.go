package cmd

import (
	"database/sql"
	"fmt"
	"go-task/db"
	"strings"
)

func AddTask(title, description, priorityName, tagsInput string) {
	tx, err := db.DB.Begin()
	if err != nil {
		fmt.Printf("Error starting transaction: %v\n", err)
		return
	}
	defer tx.Rollback()

	var priorityID int
	switch strings.ToLower(priorityName) {
	case "h", "high":
		priorityName = "High"
	case "m", "medium":
		priorityName = "Medium"
	case "l", "low":
		priorityName = "Low"
	default:
		priorityName = "Low"
	}

	err = tx.QueryRow("SELECT id FROM priorities WHERE name = ? COLLATE NOCASE", priorityName).Scan(&priorityID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Error: Invalid priority '%s'. Allowed: High, Medium, Low\n", priorityName)
			return
		}
		fmt.Printf("Error looking up priority: %v\n", err)
		return
	}

	// insert Task
	res, err := tx.Exec("INSERT INTO tasks (title, description, priority_id) VALUES (?, ?, ?)", title, description, priorityID)
	if err != nil {
		fmt.Printf("Error adding task: %v\n", err)
		return
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error getting last insert ID: %v\n", err)
		return
	}

	// handle Tags
	if tagsInput != "" {
		tagNames := strings.Split(tagsInput, ",")
		for _, tagName := range tagNames {
			tagName = strings.TrimSpace(tagName)
			if tagName == "" {
				continue
			}

			// check if tag exists or insert it
			var tagID int64
			err = tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
			if err == sql.ErrNoRows {
				resTag, err := tx.Exec("INSERT INTO tags (name) VALUES (?)", tagName)
				if err != nil {
					fmt.Printf("Error creating tag '%s': %v\n", tagName, err)
					return
				}
				tagID, err = resTag.LastInsertId()
				if err != nil {
					fmt.Printf("Error getting tag ID: %v\n", err)
					return
				}
			} else if err != nil {
				fmt.Printf("Error looking up tag '%s': %v\n", tagName, err)
				return
			}

			_, err = tx.Exec("INSERT INTO task_tags (task_id, tag_id) VALUES (?, ?)", taskID, tagID)
			if err != nil {
				fmt.Printf("Error linking tag '%s' to task: %v\n", tagName, err)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		return
	}

	fmt.Printf("Task added successfully! (ID: %d, Priority: %s)\n", taskID, priorityName)
}
