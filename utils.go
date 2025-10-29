package main

import "fmt"

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
	fmt.Println("═══════════════════════════════════════\n")
}

func isFlag(arg string) bool {
	return len(arg) > 1 && arg[0] == '-'
}
