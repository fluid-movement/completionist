package main

import (
	"fmt"
	"io"
	"sort"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

const (
	stateKanban = iota
	stateAdding
	stateChoosingRefType
	stateAddingRef
	statePickingFile
	stateEditing
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
}

type editDelegate struct {
	list.DefaultDelegate
	editIndex int
	input     textinput.Model
}

func (d editDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if index == d.editIndex {
		desc := item.(todoListItem).Description()
		fmt.Fprintf(w, "%s\n%s", d.input.View(),
			d.DefaultDelegate.Styles.NormalDesc.Render(desc))
		return
	}
	d.DefaultDelegate.Render(w, m, index, item)
}

type todoListItem struct{ todo TodoItem }

func (t todoListItem) Title() string       { return t.todo.Title }
func (t todoListItem) Description() string {
	if d := t.todo.RefDescription(); d != "" {
		return d
	}
	return "added " + t.todo.ReadableCreatedAt()
}
func (t todoListItem) FilterValue() string { return t.todo.Title }

func itemsForStatus(todos *TodoList, status TodoStatus) []list.Item {
	var items []list.Item
	for _, item := range todos.Items {
		if item.Status == status {
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
		m.columns[i].SetItems(itemsForStatus(m.todos, TodoStatus(i)))
	}
}

// clampCursor ensures the column cursor doesn't exceed the last item after a removal.
func clampCursor(m *model, col int) {
	n := len(m.columns[col].Items())
	if n > 0 && m.columns[col].Index() >= n {
		m.columns[col].Select(n - 1)
	}
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
		delegate := list.NewDefaultDelegate()
		l := list.New(itemsForStatus(todos, TodoStatus(i)), delegate, 0, 0)
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

	return model{
		todos:      todos,
		storage:    storage,
		columns:    columns,
		focused:    0,
		state:      stateKanban,
		input:      ti,
		refInput:   ri,
		filePicker: fp,
		help:       help.New(),
		editInput:  ei,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
