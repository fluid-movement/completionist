package main

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	stateKanban = iota
	stateAdding
	stateChoosingRefType
	stateAddingRef
	statePickingFile
	stateEditing
	stateProjects
	stateAddingProject
	stateEditingProject
	stateMovingToProject
)

const (
	refChoiceNothing = iota
	refChoicePaste
	refChoicePickFile
)

type model struct {
	todos        *TodoList
	storage      Storage
	columns      [3]list.Model
	focused      int
	width        int
	height       int
	state        int
	input        textinput.Model
	refInput     textinput.Model
	pendingTitle string
	refChoice    int
	filePicker   filepicker.Model
	help         help.Model
	err          error
	editInput    textinput.Model
	editID       int
	editIndex    int

	currentProjectID int
	projectTable     table.Model
	projectInput     textinput.Model
	editingProjectID int
	moveProjectIndex int
}

// wrapText splits text into lines of at most width runes, breaking at word boundaries.
func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}
	runes := []rune(text)
	if len(runes) <= width {
		return []string{text}
	}
	var lines []string
	for len(runes) > 0 {
		if len(runes) <= width {
			lines = append(lines, string(runes))
			break
		}
		breakAt := -1
		for i := 0; i < len(runes) && i < width; i++ {
			if runes[i] == ' ' {
				breakAt = i
			}
		}
		if breakAt < 0 {
			breakAt = min(width, len(runes))
		}
		lines = append(lines, string(runes[:breakAt]))
		runes = []rune(strings.TrimLeft(string(runes[breakAt:]), " "))
	}
	return lines
}

// columnTitleLines returns the max number of wrapped title lines needed
// across all items, given the available text width.
func columnTitleLines(items []list.Item, textwidth int) int {
	maxLines := 1
	for _, item := range items {
		di, ok := item.(todoListItem)
		if !ok {
			continue
		}
		if n := len(wrapText(di.Title(), textwidth)); n > maxLines {
			maxLines = n
		}
	}
	return maxLines
}

// wrappingDelegate renders todo items with full title word-wrapping and an
// always-visible description. Height = titleLines+1, Spacing = 0.
type wrappingDelegate struct {
	list.DefaultDelegate
	titleLines int
}

// newWrappingDelegate creates a delegate sized for titleLines title lines + 1 description line.
func newWrappingDelegate(titleLines int) wrappingDelegate {
	d := list.NewDefaultDelegate()
	d.SetHeight(titleLines + 1)
	d.SetSpacing(0)
	return wrappingDelegate{DefaultDelegate: d, titleLines: titleLines}
}

func (d wrappingDelegate) Height() int  { return d.titleLines + 1 }
func (d wrappingDelegate) Spacing() int { return 0 }

func (d wrappingDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	di, ok := item.(todoListItem)
	if !ok {
		d.DefaultDelegate.Render(w, m, index, item)
		return
	}
	if m.Width() <= 0 {
		return
	}
	s := &d.Styles
	textwidth := m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight()
	titleLines := wrapText(di.Title(), textwidth)
	desc := di.Description()

	isSelected := index == m.Index()
	var b strings.Builder
	for i := range d.titleLines {
		if i > 0 {
			b.WriteRune('\n')
		}
		var line string
		if i < len(titleLines) {
			line = titleLines[i]
		}
		if isSelected {
			b.WriteString(s.SelectedTitle.Render(line))
		} else {
			b.WriteString(s.NormalTitle.Render(line))
		}
	}
	b.WriteRune('\n')
	if isSelected {
		b.WriteString(s.SelectedDesc.Render(desc))
	} else {
		b.WriteString(s.NormalDesc.Render(desc))
	}
	fmt.Fprint(w, b.String())
}

// syncDelegates recomputes the wrapping delegate for each column based on its
// current items and column width, then sets it. Call before SetSize so PerPage
// is calculated with the correct height.
func syncDelegates(m *model) {
	_, innerW, _ := columnSize(m.width, m.height)
	textwidth := innerW - 2 // NormalTitle has PaddingLeft=2
	for i := range m.columns {
		titleLines := columnTitleLines(m.columns[i].Items(), textwidth)
		m.columns[i].SetDelegate(newWrappingDelegate(titleLines))
	}
}

// syncColumnsAndDelegates updates items, recomputes delegates, and refreshes PerPage.
// Use this wherever syncColumns was called previously.
func syncColumnsAndDelegates(m *model) {
	syncColumns(m)
	_, innerW, colH := columnSize(m.width, m.height)
	textwidth := innerW - 2
	for i := range m.columns {
		titleLines := columnTitleLines(m.columns[i].Items(), textwidth)
		m.columns[i].SetDelegate(newWrappingDelegate(titleLines))
		m.columns[i].SetSize(innerW, colH)
	}
}

type editDelegate struct {
	wrappingDelegate
	editIndex int
	input     textinput.Model
}

func (d editDelegate) Height() int  { return d.titleLines + 1 }
func (d editDelegate) Spacing() int { return 0 }

func (d editDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if index == d.editIndex {
		desc := item.(todoListItem).Description()
		var b strings.Builder
		b.WriteString(d.input.View())
		// Fill remaining title lines with blank lines to keep consistent height
		for range d.titleLines - 1 {
			b.WriteRune('\n')
		}
		b.WriteRune('\n')
		b.WriteString(d.wrappingDelegate.Styles.NormalDesc.Render(desc))
		fmt.Fprint(w, b.String())
		return
	}
	d.wrappingDelegate.Render(w, m, index, item)
}

type todoListItem struct{ todo TodoItem }

func (t todoListItem) Title() string { return t.todo.Title }
func (t todoListItem) Description() string {
	if d := t.todo.RefDescription(); d != "" {
		return d
	}
	return "added " + t.todo.ReadableCreatedAt()
}
func (t todoListItem) FilterValue() string { return t.todo.Title }

// allProjects returns the default project prepended to the stored projects.
func allProjects(todos *TodoList) []Project {
	return append([]Project{{ID: defaultProjectID, Name: "default"}}, todos.Projects...)
}

// projectTableCols returns column definitions whose visual widths sum to panelInnerW.
// Each cell is rendered at col.Width content + 2 cell-padding = col.Width+2 visual.
func projectTableCols(panelInnerW int) []table.Column {
	const statColW = 6 // content width; visual = 8
	projectColW := panelInnerW - 3*(statColW+2) - 2
	if projectColW < 4 {
		projectColW = 4
	}
	return []table.Column{
		{Title: "project", Width: projectColW},
		{Title: "open", Width: statColW},
		{Title: "in prog", Width: statColW},
		{Title: "done", Width: statColW},
	}
}

func countTodos(todos *TodoList, projectID int, status TodoStatus) int {
	n := 0
	for _, item := range todos.Items {
		if item.ProjectID == projectID && item.Status == status {
			n++
		}
	}
	return n
}

func syncProjectTable(m *model) {
	projects := allProjects(m.todos)
	rows := make([]table.Row, len(projects))
	for i, p := range projects {
		rows[i] = table.Row{
			p.Name,
			strconv.Itoa(countTodos(m.todos, p.ID, StatusOpen)),
			strconv.Itoa(countTodos(m.todos, p.ID, StatusInProgress)),
			strconv.Itoa(countTodos(m.todos, p.ID, StatusDone)),
		}
	}
	m.projectTable.SetRows(rows)
}

func itemsForStatus(todos *TodoList, status TodoStatus, projectID int) []list.Item {
	var items []list.Item
	for _, item := range todos.Items {
		if item.Status == status && item.ProjectID == projectID {
			items = append(items, todoListItem{todo: item})
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].(todoListItem).todo.CreatedAt > items[j].(todoListItem).todo.CreatedAt
	})
	return items
}

func syncColumns(m *model) {
	for i := range 3 {
		m.columns[i].SetItems(itemsForStatus(m.todos, TodoStatus(i), m.currentProjectID))
	}
}

// clampCursor ensures the column cursor doesn't exceed the last item after a removal.
func clampCursor(m *model, col int) {
	n := len(m.columns[col].Items())
	if n > 0 && m.columns[col].Index() >= n {
		m.columns[col].Select(n - 1)
	}
}

func newProjectTable() table.Model {
	s := table.DefaultStyles()
	s.Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A855F7")).
		Padding(0, 1)
	s.Selected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F5F3FF")).
		Background(lipgloss.Color("#7C3AED")).
		Padding(0, 1)
	s.Cell = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("#9CA3AF"))

	return table.New(
		table.WithColumns(projectTableCols(52)),
		table.WithFocused(true),
		table.WithStyles(s),
	)
}

func initialModel(todos *TodoList, storage Storage) model {
	ti := textinput.New()
	ti.Placeholder = "Todo title..."

	ri := textinput.New()
	ri.Placeholder = "URL, file path, or !command (optional)"

	fp := filepicker.New()
	fp.AutoHeight = false
	fp.ShowPermissions = false
	fp.ShowSize = false
	fp.DirAllowed = true
	// Remove esc from Back so we can intercept it to return to the chooser.
	fp.KeyMap.Back = key.NewBinding(key.WithKeys("h", "backspace", "left"), key.WithHelp("h/←", "back"))

	titles := []string{"Open", "In Progress", "Done"}
	var columns [3]list.Model
	for i := range 3 {
		delegate := newWrappingDelegate(1)
		l := list.New(itemsForStatus(todos, TodoStatus(i), defaultProjectID), delegate, 0, 0)
		l.Title = titles[i]
		l.SetFilteringEnabled(false)
		l.SetShowFilter(false)
		l.SetShowStatusBar(false)
		l.SetShowHelp(false)
		l.SetShowTitle(false) // rendered manually in view so we can inject the input below it
		columns[i] = l
	}

	ei := textinput.New()
	ei.Placeholder = "Edit title..."

	pi := textinput.New()
	pi.Placeholder = "Project name..."

	pt := newProjectTable()

	m := model{
		todos:        todos,
		storage:      storage,
		columns:      columns,
		focused:      0,
		state:        stateKanban,
		input:        ti,
		refInput:     ri,
		filePicker:   fp,
		help:         help.New(),
		editInput:    ei,
		projectInput: pi,
		projectTable: pt,
	}
	syncProjectTable(&m)
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}
