package main

import "charm.land/lipgloss/v2"

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

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4B5563")).
			Padding(1, 2)
)
