package main

import (
	"fmt"
	"time"
)

type TodoItem struct {
	ID        int
	Title     string
	Completed bool
	CreatedAt int64
}

type TodoList struct {
	Items []TodoItem
}

func (l *TodoList) Add(title string) (TodoItem, error) {
	todo := TodoItem{
		ID:        l.nextID(),
		Title:     title,
		Completed: false,
		CreatedAt: time.Now().Unix(),
	}

	l.Items = append(l.Items, todo)

	return todo, nil
}

func (l *TodoList) Remove(id int) error {
	filtered := make([]TodoItem, 0, len(l.Items))
	found := false
	for _, item := range l.Items {
		if item.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, item)
	}

	if !found {
		return fmt.Errorf("todo with ID %d not found", id)
	}

	l.Items = filtered
	return nil
}

func (l *TodoList) Complete(id int) error {
	found := false
	for i, item := range l.Items {
		if item.ID == id {
			if item.Completed == true {
				return fmt.Errorf("todo with ID %d is already completed", id)
			}
			l.Items[i].Completed = true
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("todo with ID %d not found", id)
	}
	
	return nil
}

func (l *TodoList) nextID() int {
	if len(l.Items) == 0 {
		return 1
	}

	max := 0
	for _, item := range l.Items {
		if item.ID > max {
			max = item.ID
		}
	}
	
	return max + 1
}

func (i *TodoItem) ReadableCreatedAt() string {
	t := time.Unix(i.CreatedAt, 0)
	return t.Format("2006-01-02 15:04:05")
}