package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

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
		m.input.Reset()
		m.err = nil
		return m, m.input.Focus()
	}
	return m, nil
}

func (m model) updateAdding(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateList
		m.input.Blur()
		return m, nil
	case "enter":
		title := strings.TrimSpace(m.input.Value())
		if title != "" {
			m.todos.Add(title)
			m.storage.Save(m.todos)
			m.cursor = len(m.todos.Items) - 1
		}
		m.state = stateList
		m.input.Blur()
		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}
