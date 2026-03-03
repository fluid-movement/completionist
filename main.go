package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ── Styles ────────────────────────────────────────────────────────────────────

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#C084FC")).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A855F7")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E9D5FF"))

	completedStyle = lipgloss.NewStyle().
			Faint(true).
			Strikethrough(true)

	checkDoneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#34D399"))

	checkTodoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4B5563")).
			Italic(true)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F87171")).
			Bold(true)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A855F7")).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E9D5FF"))

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4B5563")).
			Padding(1, 2)
)

// ── Model ─────────────────────────────────────────────────────────────────────

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
	case tea.KeyPressMsg:
		switch m.state {
		case stateList:
			return m.updateList(msg)
		case stateAdding:
			return m.updateAdding(msg)
		}
	}
	return m, nil
}

func (m model) updateList(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
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

func (m model) updateAdding(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
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
			runes := []rune(m.input)
			m.input = string(runes[:len(runes)-1])
		}
	default:
		// Accept printable characters
		for _, r := range msg.String() {
			if unicode.IsPrint(r) {
				m.input += string(r)
			}
		}
	}
	return m, nil
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m model) View() tea.View {
	var body strings.Builder

	// Header
	body.WriteString(titleStyle.Render("✓ Completionist"))
	body.WriteString("\n")
	body.WriteString(subtitleStyle.Render("your personal todo list"))
	body.WriteString("\n\n")

	// Todo list
	if len(m.todos.Items) == 0 {
		body.WriteString("  " + helpStyle.Render("Nothing here yet — press 'a' to add your first todo."))
		body.WriteString("\n")
	} else {
		for i, item := range m.todos.Items {
			selected := i == m.cursor

			// Cursor column
			if selected {
				body.WriteString(cursorStyle.Render("▶"))
			} else {
				body.WriteString(" ")
			}
			body.WriteString(" ")

			// Checkbox
			if item.Completed {
				body.WriteString(checkDoneStyle.Render("●"))
			} else {
				body.WriteString(checkTodoStyle.Render("○"))
			}
			body.WriteString(" ")

			// Title
			switch {
			case item.Completed:
				body.WriteString(completedStyle.Render(item.Title))
			case selected:
				body.WriteString(selectedStyle.Render(item.Title))
			default:
				body.WriteString(item.Title)
			}
			body.WriteString("\n")
		}
	}

	body.WriteString("\n")

	// Footer: input or help
	if m.state == stateAdding {
		prompt := inputLabelStyle.Render("  New todo › ")
		text := inputStyle.Render(m.input) + "█"
		body.WriteString(prompt + text)
		body.WriteString("\n")
		body.WriteString("  " + helpStyle.Render("enter")+helpKeyStyle.Render(" to save  ")+helpStyle.Render("esc")+helpKeyStyle.Render(" to cancel"))
		body.WriteString("\n")
	} else {
		if m.err != nil {
			body.WriteString("  " + errorStyle.Render("⚠  "+m.err.Error()))
			body.WriteString("\n")
		}
		keys := []string{"a", "enter/c", "d", "j/k", "q"}
		descs := []string{"add", "complete", "delete", "navigate", "quit"}
		var help strings.Builder
		for i, k := range keys {
			help.WriteString(helpKeyStyle.Render(k))
			help.WriteString(" ")
			help.WriteString(helpStyle.Render(descs[i]))
			if i < len(keys)-1 {
				help.WriteString(helpStyle.Render("  ·  "))
			}
		}
		body.WriteString("  " + help.String())
		body.WriteString("\n")
	}

	content := borderStyle.Render(body.String())
	return tea.NewView(content)
}

// ── Main ──────────────────────────────────────────────────────────────────────

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
