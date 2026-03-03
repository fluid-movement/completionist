package main

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	v := tea.NewView(lipgloss.JoinVertical(
		0,
		renderHeader(),
		renderBody(m),
		renderFooter(m),
	))
	v.AltScreen = true

	return v
}

func renderHeader() string {
	return titleStyle.Render("✓ Completionist") + "\n" + subtitleStyle.Render("your personal todo list") + "\n\n"
}

func renderBody(m model) string {
	outerW, innerW, colH := columnSize(m.width, m.height)

	// Delegate that suppresses the selection highlight — used for non-focused columns
	blurredDelegate := list.NewDefaultDelegate()
	blurredDelegate.Styles.SelectedTitle = blurredDelegate.Styles.NormalTitle
	blurredDelegate.Styles.SelectedDesc = blurredDelegate.Styles.NormalDesc

	var cols [3]string
	for i := range m.columns {
		isAdding := m.state == stateAdding && i == 0
		isFocused := i == m.focused && m.state == stateKanban

		// Render the title bar using the list's own styles (list.SetShowTitle is false)
		titleBar := m.columns[i].Styles.TitleBar.Render(
			m.columns[i].Styles.Title.Render(m.columns[i].Title),
		)
		titleH := lipgloss.Height(titleBar)

		// Build column sections: title, optional input, list content
		sections := []string{titleBar}
		if isAdding {
			m.columns[i].SetDelegate(blurredDelegate)
			m.columns[i].SetSize(innerW, colH-titleH-1) // -1 for the input row
			sections = append(sections, m.input.View())
		} else {
			if !isFocused {
				m.columns[i].SetDelegate(blurredDelegate)
			}
			m.columns[i].SetSize(innerW, colH-titleH)
		}
		sections = append(sections, m.columns[i].View())
		colContent := lipgloss.JoinVertical(lipgloss.Left, sections...)

		if isFocused || isAdding {
			cols[i] = columnFocusedStyle.Width(outerW).Render(colContent)
		} else {
			cols[i] = columnBlurredStyle.Width(outerW).Render(colContent)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cols[0], cols[1], cols[2])
}

func renderFooter(m model) string {
	var helpView string
	if m.state == stateAdding {
		helpView = m.help.View(addingKeys)
	} else {
		helpView = m.help.View(kanbanKeys)
	}
	if m.err != nil {
		return "  " + errorStyle.Render("⚠  "+m.err.Error()) + "\n"
	}
	return "  " + helpView + "\n"
}
