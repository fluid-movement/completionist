package main

import (
	"strings"

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
		isChoosingRef := m.state == stateChoosingRefType && i == 0
		isAddingRef := m.state == stateAddingRef && i == 0
		isPickingFile := m.state == statePickingFile && i == 0
		isEditing := m.state == stateEditing && i == m.focused
		isFocused := i == m.focused && m.state == stateKanban

		// Render the title bar using the list's own styles (list.SetShowTitle is false)
		titleBar := m.columns[i].Styles.TitleBar.Render(
			m.columns[i].Styles.Title.Render(m.columns[i].Title),
		)
		titleH := lipgloss.Height(titleBar)

		// Build column sections: title, optional overlay content, list
		sections := []string{titleBar}
		showList := true

		switch {
		case isEditing:
			ed := editDelegate{
				DefaultDelegate: list.NewDefaultDelegate(),
				editIndex:       m.editIndex,
				input:           m.editInput,
			}
			m.columns[i].SetDelegate(ed)
			m.columns[i].SetSize(innerW, colH-titleH)
		case isAdding:
			m.columns[i].SetDelegate(blurredDelegate)
			m.columns[i].SetSize(innerW, colH-titleH-1)
			sections = append(sections, m.input.View())
		case isChoosingRef:
			m.columns[i].SetDelegate(blurredDelegate)
			chooser := renderRefChooser(m)
			pendingLine := pendingTitleStyle.Render(m.pendingTitle)
			overhead := lipgloss.Height(pendingLine) + lipgloss.Height(chooser)
			m.columns[i].SetSize(innerW, colH-titleH-overhead)
			sections = append(sections, pendingLine, chooser)
		case isAddingRef:
			m.columns[i].SetDelegate(blurredDelegate)
			pendingLine := pendingTitleStyle.Render(m.pendingTitle)
			overhead := lipgloss.Height(pendingLine) + 1
			m.columns[i].SetSize(innerW, colH-titleH-overhead)
			sections = append(sections, pendingLine, m.refInput.View())
		case isPickingFile:
			pendingLine := pendingTitleStyle.Render(m.pendingTitle)
			pendingH := lipgloss.Height(pendingLine)
			m.filePicker.SetHeight(colH - titleH - pendingH)
			sections = append(sections, pendingLine, m.filePicker.View())
			showList = false
		default:
			if !isFocused {
				m.columns[i].SetDelegate(blurredDelegate)
			}
			m.columns[i].SetSize(innerW, colH-titleH)
		}

		if showList {
			sections = append(sections, m.columns[i].View())
		}
		colContent := lipgloss.JoinVertical(lipgloss.Left, sections...)

		if isFocused || isAdding || isChoosingRef || isAddingRef || isPickingFile || isEditing {
			cols[i] = columnFocusedStyle.Width(outerW).Render(colContent)
		} else {
			cols[i] = columnBlurredStyle.Width(outerW).Render(colContent)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cols[0], cols[1], cols[2])
}

var refChoiceLabels = []string{"nothing", "paste URL / path", "select file"}

func renderRefChooser(m model) string {
	var b strings.Builder
	for i, label := range refChoiceLabels {
		if i == m.refChoice {
			b.WriteString(refChoiceSelectedStyle.Render("> " + label))
		} else {
			b.WriteString(refChoiceStyle.Render("  " + label))
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func renderFooter(m model) string {
	var helpView string
	if m.state == stateAdding || m.state == stateChoosingRefType || m.state == stateAddingRef || m.state == statePickingFile || m.state == stateEditing {
		helpView = m.help.View(addingKeys)
	} else {
		helpView = m.help.View(kanbanKeys)
	}
	if m.err != nil {
		return "  " + errorStyle.Render("⚠  "+m.err.Error()) + "\n"
	}
	return "  " + helpView + "\n"
}
