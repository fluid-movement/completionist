package main

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width)
		_, innerW, colH := columnSize(msg.Width, msg.Height)
		for i := range m.columns {
			m.columns[i].SetSize(innerW, colH)
		}
		return m, nil
	case tea.KeyPressMsg:
		switch m.state {
		case stateKanban:
			return m.updateKanban(msg)
		case stateAdding:
			return m.updateAdding(msg)
		}
	}
	return m, nil
}

func (m model) updateKanban(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, kanbanKeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, kanbanKeys.Left):
		if m.focused > 0 {
			m.focused--
		}
		return m, nil
	case key.Matches(msg, kanbanKeys.Right):
		if m.focused < 2 {
			m.focused++
		}
		return m, nil
	case key.Matches(msg, kanbanKeys.Advance):
		if item := m.columns[m.focused].SelectedItem(); item != nil {
			t := item.(todoListItem).todo
			newStatus := min(t.Status+1, StatusDone)
			m.err = m.todos.SetStatus(t.ID, newStatus)
			if m.err == nil {
				m.storage.Save(m.todos)
				syncColumns(&m)
				clampCursor(&m, m.focused)
			}
		}
		return m, nil
	case key.Matches(msg, kanbanKeys.Retreat):
		if item := m.columns[m.focused].SelectedItem(); item != nil {
			t := item.(todoListItem).todo
			newStatus := max(t.Status-1, StatusOpen)
			m.err = m.todos.SetStatus(t.ID, newStatus)
			if m.err == nil {
				m.storage.Save(m.todos)
				syncColumns(&m)
				clampCursor(&m, m.focused)
			}
		}
		return m, nil
	case key.Matches(msg, kanbanKeys.Add):
		m.state = stateAdding
		m.input.Reset()
		_, innerW, _ := columnSize(m.width, m.height)
		m.input.SetWidth(innerW - 2)
		m.err = nil
		return m, m.input.Focus()
	case key.Matches(msg, kanbanKeys.Delete):
		if item := m.columns[m.focused].SelectedItem(); item != nil {
			t := item.(todoListItem).todo
			m.err = m.todos.Remove(t.ID)
			if m.err == nil {
				m.storage.Save(m.todos)
				syncColumns(&m)
				clampCursor(&m, m.focused)
			}
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.columns[m.focused], cmd = m.columns[m.focused].Update(msg)
		return m, cmd
	}
}

func (m model) updateAdding(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, addingKeys.Cancel):
		m.state = stateKanban
		m.input.Blur()
		return m, nil
	case key.Matches(msg, addingKeys.Confirm):
		title := strings.TrimSpace(m.input.Value())
		if title != "" {
			m.todos.Add(title)
			m.storage.Save(m.todos)
			syncColumns(&m)
		}
		m.state = stateKanban
		m.input.Blur()
		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}
