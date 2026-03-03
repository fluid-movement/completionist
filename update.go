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
		m.input.SetWidth(innerW - 2)
		m.refInput.SetWidth(innerW - 2)
		m.editInput.SetWidth(innerW - 4)
		return m, nil
	case tea.KeyPressMsg:
		switch m.state {
		case stateKanban:
			return m.updateKanban(msg)
		case stateAdding:
			return m.updateAdding(msg)
		case stateChoosingRefType:
			return m.updateChoosingRefType(msg)
		case stateAddingRef:
			return m.updateAddingRef(msg)
		case statePickingFile:
			return m.updatePickingFile(msg)
		case stateEditing:
			return m.updateEditing(msg)
		}
	default:
		// Forward all other messages (e.g. paste, filepicker readDirMsg) to the active component.
		switch m.state {
		case stateAdding:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		case stateAddingRef:
			var cmd tea.Cmd
			m.refInput, cmd = m.refInput.Update(msg)
			return m, cmd
		case statePickingFile:
			var cmd tea.Cmd
			m.filePicker, cmd = m.filePicker.Update(msg)
			return m, cmd
		case stateEditing:
			var cmd tea.Cmd
			m.editInput, cmd = m.editInput.Update(msg)
			return m, cmd
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
	case key.Matches(msg, kanbanKeys.Edit):
		if item := m.columns[m.focused].SelectedItem(); item != nil {
			t := item.(todoListItem).todo
			m.editID = t.ID
			m.editIndex = m.columns[m.focused].Index()
			m.editInput.SetValue(t.Title)
			_, innerW, _ := columnSize(m.width, m.height)
			m.editInput.SetWidth(innerW - 4)
			m.state = stateEditing
			return m, m.editInput.Focus()
		}
		return m, nil
	case key.Matches(msg, kanbanKeys.Open):
		if item := m.columns[m.focused].SelectedItem(); item != nil {
			t := item.(todoListItem).todo
			m.err = t.OpenRef()
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.columns[m.focused], cmd = m.columns[m.focused].Update(msg)
		return m, cmd
	}
}

func (m model) updateEditing(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, addingKeys.Cancel):
		m.editInput.Blur()
		m.state = stateKanban
		return m, nil
	case key.Matches(msg, addingKeys.Confirm):
		title := strings.TrimSpace(m.editInput.Value())
		if title == "" {
			return m, nil
		}
		m.err = m.todos.SetTitle(m.editID, title)
		if m.err == nil {
			m.storage.Save(m.todos)
			syncColumns(&m)
			m.columns[m.focused].Select(m.editIndex)
		}
		m.editInput.Blur()
		m.state = stateKanban
		return m, nil
	default:
		var cmd tea.Cmd
		m.editInput, cmd = m.editInput.Update(msg)
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
		if title == "" {
			return m, nil
		}
		m.pendingTitle = title
		m.refChoice = refChoiceNothing
		m.state = stateChoosingRefType
		m.input.Blur()
		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (m model) updateChoosingRefType(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, addingKeys.Cancel):
		m.pendingTitle = ""
		m.state = stateKanban
		return m, nil
	case key.Matches(msg, kanbanKeys.Up):
		if m.refChoice > 0 {
			m.refChoice--
		}
		return m, nil
	case key.Matches(msg, kanbanKeys.Down):
		if m.refChoice < 2 {
			m.refChoice++
		}
		return m, nil
	case key.Matches(msg, addingKeys.Confirm):
		switch m.refChoice {
		case refChoiceNothing:
			m.todos.Add(m.pendingTitle, "")
			m.storage.Save(m.todos)
			syncColumns(&m)
			m.pendingTitle = ""
			m.state = stateKanban
		case refChoicePaste:
			m.refInput.Reset()
			_, innerW, _ := columnSize(m.width, m.height)
			m.refInput.SetWidth(innerW - 2)
			m.state = stateAddingRef
			return m, m.refInput.Focus()
		case refChoicePickFile:
			m.filePicker.CurrentDirectory = "."
			m.state = statePickingFile
			return m, m.filePicker.Init()
		}
		return m, nil
	}
	return m, nil
}

func (m model) updateAddingRef(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, addingKeys.Cancel):
		m.refInput.Blur()
		m.state = stateChoosingRefType
		return m, nil
	case key.Matches(msg, addingKeys.Confirm):
		m.todos.Add(m.pendingTitle, strings.TrimSpace(m.refInput.Value()))
		m.storage.Save(m.todos)
		syncColumns(&m)
		m.pendingTitle = ""
		m.refInput.Blur()
		m.state = stateKanban
		return m, nil
	default:
		var cmd tea.Cmd
		m.refInput, cmd = m.refInput.Update(msg)
		return m, cmd
	}
}

func (m model) updatePickingFile(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, addingKeys.Cancel) {
		m.state = stateChoosingRefType
		return m, nil
	}
	var cmd tea.Cmd
	m.filePicker, cmd = m.filePicker.Update(msg)
	if ok, path := m.filePicker.DidSelectFile(msg); ok {
		m.todos.Add(m.pendingTitle, path)
		m.storage.Save(m.todos)
		syncColumns(&m)
		m.pendingTitle = ""
		m.state = stateKanban
	}
	return m, cmd
}
