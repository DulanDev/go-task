package models

import (
	"database/sql"
	"time"
)
type Priority struct {
	ID    int
	Name  string
	Level int
}

type Tag struct {
	ID   int
	Name string
}

type Task struct {
	ID          int
	Title       string
	Description string
	Completed   sql.NullBool
	CreatedAt   time.Time
	DueDate     *time.Time
	Priority    *Priority
	Tags        []Tag
}