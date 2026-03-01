package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	storage, err := InitializeStorage()
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		return
	}

	todos, err := storage.Load()
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		return
	}

	args := os.Args[1:]
	command := ""
	if len(args) != 0 {
		command = args[0]
	}

	switch command {
	case "list":
		OutputList(todos)
	case "add":
		if err := addItem(todos); err == nil {
			storage.Save(todos)
		}
	case "complete":
		if err := completeItem(todos); err == nil {
			storage.Save(todos)
		}
	case "remove":
	if err := removeItem(todos); err == nil {
		storage.Save(todos)
	}
	default:
		OutputHelp()
	}

}

func removeItem(todos *TodoList) any {
	id, err := getIntArg()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	err = todos.Remove(id)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}

	fmt.Printf("Todo with ID %d removed.\n", id)
	
	return nil
}

func completeItem(todos *TodoList) any {
	id, err := getIntArg()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	err = todos.Complete(id)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}

	fmt.Printf("Todo with ID %d marked as completed.\n", id)
	
	return nil
}

func getStringArg() (string, error) {
	if len(os.Args) < 3 {
		return "", fmt.Errorf("missing argument")
	}
	
	return os.Args[2], nil
}

func getIntArg() (int, error) {
	str, err := getStringArg()
	if err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("argument is not a valid number")
	}

	return n, nil
}

func addItem(todos *TodoList) error {
	title, err := getStringArg()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	newTodo, err := todos.Add(title)
	OutputAddedTodo(&newTodo)

	return nil
}
