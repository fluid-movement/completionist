package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	stateList   = iota
	stateAdding
)

type model struct {
	todos   *TodoList
	storage Storage
	cursor  int
	state   int
	input   string
	err     error
}

func initialModel(todos *TodoList, storage Storage) model {
	return model{
		todos:   todos,
		storage: storage,
		state:   stateList,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateList:
			return m.updateList(msg)
		case stateAdding:
			return m.updateAdding(msg)
		}
	}
	return m, nil
}

func (m model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	items := m.todos.Items
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(items)-1 {
			m.cursor++
		}
	case "enter", "c":
		if len(items) > 0 {
			m.err = m.todos.Complete(items[m.cursor].ID)
			if m.err == nil {
				m.storage.Save(m.todos)
			}
		}
	case "d":
		if len(items) > 0 {
			id := items[m.cursor].ID
			m.err = m.todos.Remove(id)
			if m.err == nil {
				m.storage.Save(m.todos)
				if m.cursor >= len(m.todos.Items) && m.cursor > 0 {
					m.cursor--
				}
			}
		}
	case "a":
		m.state = stateAdding
		m.input = ""
		m.err = nil
	}
	return m, nil
}

func (m model) updateAdding(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateList
		m.input = ""
	case "enter":
		title := strings.TrimSpace(m.input)
		if title != "" {
			m.todos.Add(title)
			m.storage.Save(m.todos)
			m.cursor = len(m.todos.Items) - 1
		}
		m.state = stateList
		m.input = ""
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) View() string {
	var sb strings.Builder
	sb.WriteString("Completionist\n")
	sb.WriteString("─────────────\n\n")

	if len(m.todos.Items) == 0 {
		sb.WriteString("  No todos yet. Press 'a' to add one.\n")
	} else {
		for i, item := range m.todos.Items {
			cursor := "  "
			if i == m.cursor {
				cursor = "> "
			}
			status := "[ ]"
			if item.Completed {
				status = "[x]"
			}
			sb.WriteString(fmt.Sprintf("%s%s %s\n", cursor, status, item.Title))
		}
	}

	sb.WriteString("\n")
	if m.state == stateAdding {
		sb.WriteString(fmt.Sprintf("New todo: %s_\n", m.input))
		sb.WriteString("  enter: save  esc: cancel\n")
	} else {
		if m.err != nil {
			sb.WriteString(fmt.Sprintf("  Error: %s\n", m.err))
		}
		sb.WriteString("  a: add  enter/c: complete  d: delete  q: quit\n")
	}

	return sb.String()
}

func main() {
	storage, err := InitializeStorage()
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	todos, err := storage.Load()
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(todos, storage))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
