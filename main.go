package main

import (
	"go-task/cmd"
	"go-task/db"
	"os"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type Config struct {
	DBFile string `json:"db_file"`
}

func loadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		// Default config if file not found
		return Config{DBFile: "tasks.db"}
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error decoding config file:", err)
		return Config{DBFile: "tasks.db"}
	}
	return config
}

func main() {
	config := loadConfig()
	if config.DBFile == "" {
		config.DBFile = "tasks.db"
	}
	db.InitDB(config.DBFile)

	var rootCmd = &cobra.Command{Use: "task-manager"}

    var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Run: func(cobraCmd *cobra.Command, args []string) {
			title, _ := cobraCmd.Flags().GetString("title")
			description, _ := cobraCmd.Flags().GetString("description")
			priority, _ := cobraCmd.Flags().GetString("priority")
			tags, _ := cobraCmd.Flags().GetString("tags")
			cmd.AddTask(title, description, priority, tags)
		},
	}

	addCmd.Flags().String("title", "", "Title of the task")
	addCmd.Flags().String("description", "", "Description of the task")
	addCmd.Flags().StringP("priority", "p", "Medium", "Priority of the task (High, Medium, Low)")
	addCmd.Flags().StringP("tags", "t", "", "Comma-separated tags (e.g., 'work,urgent')")
	rootCmd.AddCommand(addCmd)

    var listCmd = &cobra.Command{
        Use:   "list",
        Short: "List all tasks",
        Run: func(cobraCmd *cobra.Command, args []string) {
            cmd.ListTasks()
        },
    }
    rootCmd.AddCommand(listCmd)

    var updateCmd = &cobra.Command{
        Use:   "update",
        Short: "Update a task",
        Run: func(cobraCmd *cobra.Command, args []string) {
            id, _ := cobraCmd.Flags().GetInt("id")
            title, _ := cobraCmd.Flags().GetString("title")
            description, _ := cobraCmd.Flags().GetString("description")
            cmd.UpdateTask(id, title, description)
        },
    }
    updateCmd.Flags().Int("id", 0, "ID of the task")
    updateCmd.Flags().String("title", "", "New title of the task")
    updateCmd.Flags().String("description", "", "New description of the task")
    rootCmd.AddCommand(updateCmd)


    var deleteCmd = &cobra.Command{
        Use:   "delete",
        Short: "Delete a task",
        Run: func(cobraCmd *cobra.Command, args []string) {
            id, _ := cobraCmd.Flags().GetInt("id")
            cmd.DeleteTask(id)
        },
    }
    deleteCmd.Flags().Int("id", 0, "ID of the task")
    rootCmd.AddCommand(deleteCmd)

    var completeCmd = &cobra.Command{
        Use:   "complete",
        Short: "Mark a task as completed",
        Run: func(cobraCmd *cobra.Command, args []string) {
            id, _ := cobraCmd.Flags().GetInt("id")
            cmd.CompleteTask(id)
        },
    }
    completeCmd.Flags().Int("id", 0, "ID of the task")
    rootCmd.AddCommand(completeCmd)

    var searchCmd = &cobra.Command{
        Use:   "search",
        Short: "Search tasks by keyword",
        Run: func(cobraCmd *cobra.Command, args []string) {
            keyword, _ := cobraCmd.Flags().GetString("keyword")
            cmd.SearchTasks(keyword)
        },
    }
    searchCmd.Flags().String("keyword", "", "Keyword to search tasks")
    rootCmd.AddCommand(searchCmd)

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
