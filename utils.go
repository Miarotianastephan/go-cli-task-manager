package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Utils that used to manage operations with task interfaces
func TaskUtils(taskSavePath string, deletedTaskSavePath string) *TaskManager {
	tm := &TaskManager{
		Tasks:          []Task{},
		FilePath:       taskSavePath,
		DeleteFilePath: deletedTaskSavePath,
	}
	tm.Load()
	return tm
}

func helpUserInterface() {
	fmt.Println("\n📋 Task Manager CLI")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("Usage:")
	fmt.Println("  task add <title> [-d description] [-date YYYY-MM-DD]")
	fmt.Println("       - Add a new task")
	fmt.Println("  task list              - List all tasks")
	fmt.Println("  task complete <id>     - Mark task as completed")
	fmt.Println("  task delete <id>       - Delete a task (saved to history)")
	fmt.Println("  task deleted           - Show deleted tasks history")
	fmt.Println("  task clear             - Delete all tasks")
	fmt.Println("  task help              - Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  task add \"Buy groceries\"")
	fmt.Println("  task add \"Learn Go\" -d \"Study structs and methods\"")
	fmt.Println("  task add \"Meeting\" -d \"Discuss project\" -date 2025-10-30")
	fmt.Println("═══════════════════════════════════════")
}

func isFlag(arg string) bool {
	return len(arg) > 1 && arg[0] == '-'
}

func (tm *TaskManager) ShowDeleted() {
	deletedTasks := []Task{}
	deletedData, err := os.ReadFile(tm.DeleteFilePath)
	if err != nil {
		fmt.Println("No deleted tasks found.")
		return
	}

	err = json.Unmarshal(deletedData, &deletedTasks)
	if err != nil || len(deletedTasks) == 0 {
		fmt.Println("No deleted tasks found.")
		return
	}

	fmt.Println("\n🗑️  Deleted Tasks History:")
	fmt.Println("─────────────────────────────────────────")
	for _, task := range deletedTasks {
		fmt.Printf("[%s] %d. %s\n", task.Status, task.ID, task.Title)
		if task.Description != "" {
			fmt.Printf("   📄 %s\n", task.Description)
		}
		fmt.Printf("   📅 Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Println()
	}
	fmt.Println("─────────────────────────────────────────")
}
