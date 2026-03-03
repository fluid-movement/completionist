package main

import (
	"fmt"
	"time"
)

type TodoStatus int

const (
	StatusOpen TodoStatus = iota
	StatusInProgress
	StatusDone
)

type TodoItem struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Status      TodoStatus `json:"status"`
	UpdatedAt   int64      `json:"updated_at,omitempty"`
	CreatedAt   int64      `json:"created_at"`
}

type TodoList struct {
	Items []TodoItem
}

func (l *TodoList) Add(title string) (TodoItem, error) {
	todo := TodoItem{
		ID:        l.nextID(),
		Title:     title,
		Status:    StatusOpen,
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

func (l *TodoList) SetStatus(id int, status TodoStatus) error {
	found := false
	for i, item := range l.Items {
		if item.ID == id {
			l.Items[i].Status = status
			l.Items[i].UpdatedAt = time.Now().Unix()
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

func (i *TodoItem) ReadableUpdatedAt() string {
	if i.UpdatedAt == 0 {
		return ""
	}
	t := time.Unix(i.UpdatedAt, 0)
	return t.Format("2006-01-02 15:04:05")
}
