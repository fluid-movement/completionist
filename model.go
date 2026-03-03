package main

import (
	"sort"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

const (
	stateKanban = iota
	stateAdding
)

type model struct {
	todos   *TodoList
	storage Storage
	columns [3]list.Model
	focused int
	width   int
	height  int
	state   int
	input   textinput.Model
	help    help.Model
	err     error
}

type todoListItem struct{ todo TodoItem }

func (t todoListItem) Title() string       { return t.todo.Title }
func (t todoListItem) Description() string { return "added " + t.todo.ReadableCreatedAt() }
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

	return model{
		todos:   todos,
		storage: storage,
		columns: columns,
		focused: 0,
		state:   stateKanban,
		input:   ti,
		help:    help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
