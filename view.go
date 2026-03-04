package main

import (
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	var body string
	switch m.state {
	case stateProjects, stateAddingProject, stateEditingProject:
		body = renderProjects(m)
	default:
		body = renderBody(m)
	}
	v := tea.NewView(lipgloss.JoinVertical(
		0,
		renderHeader(m),
		body,
		renderFooter(m),
	))
	v.AltScreen = true

	return v
}

func renderHeader(m model) string {
	titleText := "◆ completionist"

	// Project badge — shown when not in the default project
	var badge string
	if m.currentProjectID != defaultProjectID {
		for _, p := range m.todos.Projects {
			if p.ID == m.currentProjectID {
				badge = headerBadgeStyle.Render("◉ " + p.Name)
				break
			}
		}
	}

	bannerW := m.width
	if bannerW <= 0 {
		bannerW = 80
	}

	var content string
	if badge != "" {
		innerW := bannerW - headerBannerStyle.GetHorizontalPadding()
		spacerW := innerW - lipgloss.Width(titleText) - lipgloss.Width(badge)
		if spacerW < 1 {
			spacerW = 1
		}
		content = titleText + strings.Repeat(" ", spacerW) + badge
	} else {
		content = titleText
	}

	banner := headerBannerStyle.Width(bannerW).Render(content)
	return banner + "\n"
}

func renderBody(m model) string {
	outerW, innerW, colH := columnSize(m.width, m.height)
	textwidth := innerW - 2 // NormalTitle PaddingLeft=2

	var cols [3]string
	for i := range m.columns {
		isAdding := m.state == stateAdding && i == 0
		isChoosingRef := m.state == stateChoosingRefType && i == 0
		isAddingRef := m.state == stateAddingRef && i == 0
		isPickingFile := m.state == statePickingFile && i == 0
		isEditing := m.state == stateEditing && i == m.focused
		isMovingToProject := m.state == stateMovingToProject && i == m.focused
		isFocused := i == m.focused && m.state == stateKanban

		// Per-column delegate sized to the longest title in this column
		titleLines := columnTitleLines(m.columns[i].Items(), textwidth)
		blurredDelegate := newWrappingDelegate(titleLines)
		blurredDelegate.Styles.SelectedTitle = blurredDelegate.Styles.NormalTitle
		blurredDelegate.Styles.SelectedDesc = blurredDelegate.Styles.NormalDesc

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
				wrappingDelegate: newWrappingDelegate(titleLines),
				editIndex:        m.editIndex,
				input:            m.editInput,
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
		case isMovingToProject:
			chooser := renderProjectChooser(m)
			chooserH := lipgloss.Height(chooser)
			m.columns[i].SetSize(innerW, colH-titleH-chooserH)
			showList = false
			sections = append(sections, renderColumnWithInlineChooser(m.columns[i], chooser))
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

		if isFocused || isAdding || isChoosingRef || isAddingRef || isPickingFile || isEditing || isMovingToProject {
			cols[i] = columnFocusedStyle.Width(outerW).Render(colContent)
		} else {
			cols[i] = columnBlurredStyle.Width(outerW).Render(colContent)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cols[0], cols[1], cols[2])
}

func renderProjects(m model) string {
	_, _, colH := columnSize(m.width, m.height)

	// Centered panel capped at 56 chars outer width
	panelOuterW := min(56, m.width)
	frameW := columnFocusedStyle.GetHorizontalFrameSize()
	panelInnerW := panelOuterW - frameW

	// Title + separator (2 lines)
	titleLine := projectPanelTitleStyle.Render("◈  projects")
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#374151")).
		Render(strings.Repeat("─", panelInnerW))
	const titleH = 2

	var sections []string
	sections = append(sections, titleLine+"\n"+separator)
	tableH := colH - titleH

	// Optional text input above the table when adding or editing
	if m.state == stateAddingProject || m.state == stateEditingProject {
		sections = append(sections, m.projectInput.View())
		tableH = colH - titleH - 1
	}

	// Size and render the stats table
	cols := projectTableCols(panelInnerW)
	m.projectTable.SetColumns(cols)
	m.projectTable.SetWidth(panelInnerW)
	m.projectTable.SetHeight(tableH)
	sections = append(sections, m.projectTable.View())

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	panel := columnFocusedStyle.Width(panelOuterW).Render(content)

	// Center the panel horizontally
	panelVisualW := lipgloss.Width(panel)
	leftPad := (m.width - panelVisualW) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	lines := strings.Split(panel, "\n")
	pad := strings.Repeat(" ", leftPad)
	for i, line := range lines {
		if line != "" {
			lines[i] = pad + line
		}
	}
	return strings.Join(lines, "\n")
}

// renderColumnWithInlineChooser manually renders the visible list items, injecting
// the chooser string directly after the selected item so it appears inline.
func renderColumnWithInlineChooser(col list.Model, chooser string) string {
	items := col.VisibleItems()
	if len(items) == 0 {
		return chooser
	}
	start, end := col.Paginator.GetSliceBounds(len(items))
	docs := items[start:end]
	selectedIdx := col.Index()
	textwidth := col.Width() - 2 // NormalTitle PaddingLeft=2
	d := newWrappingDelegate(columnTitleLines(col.Items(), textwidth))

	var b strings.Builder
	for i, item := range docs {
		d.Render(&b, col, i+start, item)
		if i+start == selectedIdx {
			b.WriteRune('\n')
			b.WriteString(chooser)
		}
		if i != len(docs)-1 {
			b.WriteString(strings.Repeat("\n", d.Spacing()+1))
		}
	}
	// Pad remaining space to maintain consistent column height
	itemsOnPage := len(docs)
	if itemsOnPage < col.Paginator.PerPage {
		n := (col.Paginator.PerPage - itemsOnPage) * (d.Height() + d.Spacing())
		b.WriteString(strings.Repeat("\n", n))
	}
	return b.String()
}

func renderProjectChooser(m model) string {
	var b strings.Builder
	b.WriteString(moveChooserHeaderStyle.Render("move to project"))
	b.WriteRune('\n')
	projects := allProjects(m.todos)
	for i, p := range projects {
		if i == m.moveProjectIndex {
			b.WriteString(refChoiceSelectedStyle.Render("> " + p.Name))
		} else {
			b.WriteString(refChoiceStyle.Render("  " + p.Name))
		}
		b.WriteRune('\n')
	}
	return b.String()
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
	switch m.state {
	case stateProjects:
		helpView = m.help.View(projectKeys)
	case stateAddingProject, stateEditingProject:
		helpView = m.help.View(addingKeys)
	case stateAdding, stateChoosingRefType, stateAddingRef, statePickingFile, stateEditing, stateMovingToProject:
		helpView = m.help.View(addingKeys)
	default:
		helpView = m.help.View(kanbanKeys)
	}
	if m.err != nil {
		return "  " + errorStyle.Render("⚠  "+m.err.Error()) + "\n"
	}
	return "  " + helpView + "\n"
}
