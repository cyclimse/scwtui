package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

func DefaultKeyMap() KeyMap {
	defaultRootKeyMap := RootKeyMap{
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
				key.WithHelp("enter", "to confirm deletion"),
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
	Quit key.Binding
}

func (m RootKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{m.Quit}
}

func (m RootKeyMap) FullHelp() [][]key.Binding {
	return nil
}

type TableKeyMap struct {
	RootKeyMap
	Delete        key.Binding
	Logs          key.Binding
	ToggleAltView key.Binding
	Search        key.Binding
}

func (m TableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Search,
		m.Logs,
		m.Delete,
		m.ToggleAltView,
		m.Quit,
	}
}

func (m TableKeyMap) FullHelp() [][]key.Binding {
	return nil
}

type ConfirmKeyMap struct {
	RootKeyMap
	Confirm key.Binding
}

func (m ConfirmKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Confirm,
		m.Quit,
	}
}

func (m ConfirmKeyMap) FullHelp() [][]key.Binding {
	return nil
}
