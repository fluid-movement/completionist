package main

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

const (
	stateList = iota
	stateAdding
)

type model struct {
	todos   *TodoList
	storage Storage
	cursor  int
	state   int
	input   textinput.Model
	err     error
}

func initialModel(todos *TodoList, storage Storage) model {
	ti := textinput.New()
	ti.Placeholder = "Todo title..."
	return model{
		todos:   todos,
		storage: storage,
		state:   stateList,
		input:   ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
