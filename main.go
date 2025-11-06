package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Load reads tasks from the JSON file
func (tm *TaskManager) Load() error {
	data, err := os.ReadFile(tm.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, that's okay
		}
		return err
	}
	return json.Unmarshal(data, &tm.Tasks)
}

// Save writes tasks to the JSON file
func (tm *TaskManager) Save(isDeleting bool) error {
	filePath := tm.FilePath
	// if isDeleting {
	// 	filePath = tm.DeleteFilePath
	// }

	data, err := json.MarshalIndent(tm.Tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// Add creates a new task
func (tm *TaskManager) Add(title, description string, createdAt time.Time) error {
	id := 1
	if len(tm.Tasks) > 0 {
		id = tm.Tasks[len(tm.Tasks)-1].ID + 1
	}

	task := Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      "en cours",
		CreatedAt:   createdAt,
	}

	tm.Tasks = append(tm.Tasks, task)
	return tm.Save(false)
}

// List displays all tasks
func (tm *TaskManager) List() {
	if len(tm.Tasks) == 0 {
		fmt.Println("No tasks found. Add one with: task add <title>")
		return
	}

	fmt.Println("\n📝 Your Tasks:")
	fmt.Println("─────────────────────────────────────────")
	for _, task := range tm.Tasks {
		statusIcon := "🔄"
		if task.Status == "terminée" {
			statusIcon = "✓"
		}
		fmt.Printf("%s %d. %s [%s]\n", statusIcon, task.ID, task.Title, task.Status)
		if task.Description != "" {
			fmt.Printf("   📄 %s\n", task.Description)
		}
		fmt.Printf("   📅 Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Println()
	}
	fmt.Println("─────────────────────────────────────────\n")
}

// Complete marks a task as completed by ID or title
func (tm *TaskManager) Complete(identifier string) error {
	// Try to parse as ID first
	id, err := strconv.Atoi(identifier)
	if err == nil {
		// It's an ID
		for i, task := range tm.Tasks {
			if task.ID == id {
				tm.Tasks[i].Status = "terminée"
				fmt.Printf("✓ Task %d completed!\n", id)
				return tm.Save(false)
			}
		}
		return fmt.Errorf("task with ID %d not found", id)
	}

	// It's not a number, try to match by title
	for i, task := range tm.Tasks {
		if task.Title == identifier {
			tm.Tasks[i].Status = "terminée"
			fmt.Printf("✓ Task '%s' completed!\n", identifier)
			return tm.Save(false)
		}
	}
	return fmt.Errorf("task with title '%s' not found", identifier)
}

// Delete removes a task by id and save into deleted tasks history
func (tm *TaskManager) Delete(id int) error {
	for i, task := range tm.Tasks {
		if task.ID == id {
			// Load existing deleted tasks
			deletedTasks := []Task{}
			deletedData, err := os.ReadFile(tm.DeleteFilePath)
			if err == nil {
				json.Unmarshal(deletedData, &deletedTasks)
			}

			// Add the task to deleted tasks
			deletedTasks = append(deletedTasks, task)

			// Save to deleted tasks file
			deletedJSON, err := json.MarshalIndent(deletedTasks, "", "  ")
			if err != nil {
				return err
			}
			os.WriteFile(tm.DeleteFilePath, deletedJSON, 0644)

			// Remove from active tasks
			tm.Tasks = append(tm.Tasks[:i], tm.Tasks[i+1:]...)
			fmt.Printf("Task %d deleted and archived!\n", id)
			return tm.Save(false)
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// Clear removes all tasks
func (tm *TaskManager) Clear() error {
	tm.Tasks = []Task{}
	fmt.Println("🗑️  All tasks cleared!")
	return tm.Save(false)
}

func main() {
	// Initialize TaskManager
	taskSavePath, deletedTaskSavePath := "taskSaved.json", "taskDeleted.json"
	tm := TaskUtils(taskSavePath, deletedTaskSavePath)

	if len(os.Args) < 2 {
		helpUserInterface()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a task title")
			fmt.Println("Usage: task add <title> [-d description] [-date YYYY-MM-DD]")
			os.Exit(1)
		}

		// Parse arguments
		title := ""
		description := ""
		createdAt := time.Now()

		i := 2
		// Get the title (all args until we hit a flag)
		for i < len(os.Args) && !isFlag(os.Args[i]) {
			title += os.Args[i] + " "
			i++
		}
		title = title[:len(title)-1] // Remove trailing space

		// Parse optional flags
		for i < len(os.Args) {
			switch os.Args[i] {
			case "-d":
				if i+1 < len(os.Args) {
					i++
					// Get description until next flag or end
					for i < len(os.Args) && !isFlag(os.Args[i]) {
						description += os.Args[i] + " "
						i++
					}
					description = description[:len(description)-1] // Remove trailing space
				}
			case "-date":
				if i+1 < len(os.Args) {
					i++
					parsedDate, err := time.Parse("2006-01-02", os.Args[i])
					if err != nil {
						fmt.Println("Error: Invalid date format. Use YYYY-MM-DD")
						os.Exit(1)
					}
					createdAt = parsedDate
					i++
				}
			default:
				i++
			}
		}

		err := tm.Add(title, description, createdAt)
		if err != nil {
			fmt.Println("Error adding task:", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Task added: %s\n", title)

	case "list":
		tm.List()

	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a task ID or title")
			fmt.Println("Usage: task complete <id|title>")
			os.Exit(1)
		}
		// Get identifier (could be ID or title with multiple words)
		identifier := ""
		for i := 2; i < len(os.Args); i++ {
			identifier += os.Args[i]
			if i < len(os.Args)-1 {
				identifier += " "
			}
		}
		err := tm.Complete(identifier)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a task ID")
			fmt.Println("Usage: task delete <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: Invalid task ID")
			os.Exit(1)
		}
		err = tm.Delete(id)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "deleted":
		tm.ShowDeleted()

	case "clear":
		err := tm.Clear()
		if err != nil {
			fmt.Println("Error clearing tasks:", err)
			os.Exit(1)
		}

	case "help":
		helpUserInterface()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		helpUserInterface()
		os.Exit(1)
	}
}
