package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

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
		body.WriteString(inputLabelStyle.Render("  New todo › ") + m.input.View())
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
