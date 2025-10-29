package main

import "time"

// File to handle the interfaces

// A task object
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // "en cours/terminée"
	CreatedAt   time.Time `json:"created_at"`
}

// To manage task operations like handlers (task object and task file path)
type TaskManager struct {
	Tasks          []Task `json:"tasks"`
	FilePath       string `json:"-"`
	DeleteFilePath string `json:"--"`
}
