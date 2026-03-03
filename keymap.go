package main

import (
	"charm.land/bubbles/v2/key"
)

// kanbanKeyMap holds all bindings used in the main kanban view.
type kanbanKeyMap struct {
	Add      key.Binding
	Advance  key.Binding
	Retreat  key.Binding
	Left     key.Binding
	Right    key.Binding
	Up       key.Binding
	Down     key.Binding
	Delete   key.Binding
	Edit     key.Binding
	Open     key.Binding
	Projects key.Binding
	Move     key.Binding
	Quit     key.Binding
}

func (k kanbanKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Add, k.Advance, k.Retreat, k.Left, k.Right, k.Delete, k.Edit, k.Open, k.Projects, k.Move, k.Quit}
}

func (k kanbanKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Advance, k.Retreat},
		{k.Left, k.Right, k.Up, k.Down},
		{k.Delete, k.Edit, k.Open, k.Projects, k.Move, k.Quit},
	}
}

// addingKeyMap holds bindings used while the add-item input is active.
type addingKeyMap struct {
	Confirm key.Binding
	Cancel  key.Binding
}

func (k addingKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Cancel}
}

func (k addingKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Confirm, k.Cancel}}
}

var kanbanKeys = kanbanKeyMap{
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Advance: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "move right"),
	),
	Retreat: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "move left"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h/←", "col left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l/→", "col right"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "down"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Open: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open ref"),
	),
	Projects: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "projects"),
	),
	Move: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "move to project"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// projectKeyMap holds bindings used in the full-screen project list view.
type projectKeyMap struct {
	Add  key.Binding
	Edit key.Binding
	Back key.Binding
	Quit key.Binding
}

func (k projectKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Add, k.Edit, k.Back, k.Quit}
}

func (k projectKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Add, k.Edit, k.Back, k.Quit}}
}

var projectKeys = projectKeyMap{
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add project"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit project"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "g"),
		key.WithHelp("esc/g", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var addingKeys = addingKeyMap{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}
