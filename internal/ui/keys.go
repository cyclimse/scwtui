package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

func DefaultKeyMap() KeyMap {
	defaultRootKeyMap := RootKeyMap{
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc, ctrl+c", "quit"),
		),
	}

	return KeyMap{
		RootKeyMap: defaultRootKeyMap,
		TableKeyMap: TableKeyMap{
			RootKeyMap: defaultRootKeyMap,
			Delete: key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "delete"),
			),
			Logs: key.NewBinding(
				key.WithKeys("l"),
				key.WithHelp("l", "logs"),
			),
			ToggleAltView: key.NewBinding(
				key.WithKeys("g"),
				key.WithHelp("g", "toggle view"),
			),
			Search: key.NewBinding(
				key.WithKeys("/"),
				key.WithHelp("/", "search"),
			),
		},
		ConfirmKeyMap: ConfirmKeyMap{
			RootKeyMap: defaultRootKeyMap,

			Confirm: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "confirm deletion"),
			),
		},
	}
}

type KeyMap struct {
	RootKeyMap
	TableKeyMap
	ConfirmKeyMap
}

func (m KeyMap) Get(focused Focused) help.KeyMap {
	switch focused {
	case TableFocused:
		return m.TableKeyMap
	case ConfirmFocused:
		return m.ConfirmKeyMap
	default:
		return m.RootKeyMap
	}
}

// Those are the keybindings that are available in all scenes.
type RootKeyMap struct {
	Help key.Binding
	Quit key.Binding
}

func (m RootKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{m.Help, m.Quit}
}

func (m RootKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
	}
}

type TableKeyMap struct {
	RootKeyMap
	Delete        key.Binding
	Logs          key.Binding
	ToggleAltView key.Binding
	Search        key.Binding
}

func (m TableKeyMap) ShortHelp() []key.Binding {
	main := []key.Binding{m.Help, m.Delete, m.Logs, m.Search}
	return main
}

func (m TableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.Help, m.Quit},
		{m.Delete, m.Logs, m.ToggleAltView},
		{m.Search},
	}
}

type ConfirmKeyMap struct {
	RootKeyMap
	Confirm key.Binding
}

func (m ConfirmKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{m.Help, m.Confirm}
}

func (m ConfirmKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.Help, m.Quit},
		{m.Confirm},
	}
}
