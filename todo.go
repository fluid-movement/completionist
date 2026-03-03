package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type TodoStatus int

const (
	StatusOpen TodoStatus = iota
	StatusInProgress
	StatusDone
)

const defaultProjectID = 0

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TodoItem struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Status    TodoStatus `json:"status"`
	UpdatedAt int64      `json:"updated_at,omitempty"`
	CreatedAt int64      `json:"created_at"`
	Ref       string     `json:"ref,omitempty"`
	ProjectID int        `json:"project_id,omitempty"`
}

type TodoList struct {
	Items    []TodoItem
	Projects []Project `json:"projects,omitempty"`
}

func (l *TodoList) Add(title, ref string, projectID int) (TodoItem, error) {
	todo := TodoItem{
		ID:        l.nextID(),
		Title:     title,
		Status:    StatusOpen,
		CreatedAt: time.Now().Unix(),
		Ref:       ref,
		ProjectID: projectID,
	}

	l.Items = append(l.Items, todo)

	return todo, nil
}

// OpenRef opens the todo's Ref in the appropriate application.
// Returns nil if Ref is empty.
func (i *TodoItem) OpenRef() error {
	if i.Ref == "" {
		return nil
	}
	if strings.HasPrefix(i.Ref, "!") {
		cmd := strings.TrimPrefix(i.Ref, "!")
		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			return nil
		}
		return exec.Command(parts[0], parts[1:]...).Start()
	}
	var opener string
	if runtime.GOOS == "darwin" {
		opener = "open"
	} else {
		opener = "xdg-open"
	}
	return exec.Command(opener, i.Ref).Start()
}

// RefDescription returns a short human-readable label for the Ref.
func (i *TodoItem) RefDescription() string {
	if i.Ref == "" {
		return ""
	}
	truncate := func(s string, n int) string {
		if len(s) <= n {
			return s
		}
		return s[:n] + "…"
	}
	if strings.HasPrefix(i.Ref, "http://") || strings.HasPrefix(i.Ref, "https://") {
		return "↗ " + truncate(i.Ref, 45)
	}
	if strings.HasPrefix(i.Ref, "!") {
		return "$ " + truncate(strings.TrimPrefix(i.Ref, "!"), 45)
	}
	return "→ " + truncate(i.Ref, 45)
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

func (l *TodoList) SetTitle(id int, title string) error {
	for i, item := range l.Items {
		if item.ID == id {
			l.Items[i].Title = title
			l.Items[i].UpdatedAt = time.Now().Unix()
			return nil
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
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

func (l *TodoList) AddProject(name string) Project {
	p := Project{ID: l.nextProjectID(), Name: name}
	l.Projects = append(l.Projects, p)
	return p
}

func (l *TodoList) SetProjectName(id int, name string) error {
	for i, p := range l.Projects {
		if p.ID == id {
			l.Projects[i].Name = name
			return nil
		}
	}
	return fmt.Errorf("project with ID %d not found", id)
}

func (l *TodoList) SetTodoProject(todoID int, projectID int) error {
	for i, item := range l.Items {
		if item.ID == todoID {
			l.Items[i].ProjectID = projectID
			return nil
		}
	}
	return fmt.Errorf("todo with ID %d not found", todoID)
}

func (l *TodoList) nextProjectID() int {
	if len(l.Projects) == 0 {
		return 1
	}
	max := 0
	for _, p := range l.Projects {
		if p.ID > max {
			max = p.ID
		}
	}
	return max + 1
}
