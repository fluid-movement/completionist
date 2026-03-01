package main

import (
	"fmt"
)

func OutputList(todos *TodoList) error {
	if len(todos.Items) == 0 {
		fmt.Println("No todo items found.")
		return nil
	}
	
	for _, item := range todos.Items {
		status := "Pending"
		if item.Completed {
			status = "Completed"
		}
		fmt.Printf("ID: %d, Title: %s, Status: %s, Created at: %d\n", item.ID, item.Title, status, item.CreatedAt)
	}
	
	return nil
}

func OutputAddedTodo(todo *TodoItem) error {
	status := "Pending"
	if todo.Completed {
		status = "Completed"
	}
	fmt.Printf("Added new todo - ID: %d, Title: %s, Status: %s, Created at: %d\n", todo.ID, todo.Title, status, todo.CreatedAt)
	
	return nil	
}

func OutputHelp() error {
	fmt.Println("Usage: completionist [list|add|complete|remove]")
	fmt.Println("Commands:")
	fmt.Println("  list\t\tList all todo items")
	fmt.Println("  add <text>\tAdd a new todo item")
	fmt.Println("  complete <id>\tMark a todo item as completed")
	fmt.Println("  remove <id>\tRemove a todo item")
	
	return nil	
}