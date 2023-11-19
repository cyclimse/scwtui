package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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
			Search: key.NewBinding(
				key.WithKeys("/"),
				key.WithHelp("/", "search"),
			),
			Describe: key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "describe"),
			),
			Logs: key.NewBinding(
				key.WithKeys("l"),
				key.WithHelp("l", "logs"),
			),
			Delete: key.NewBinding(
				key.WithKeys("x"),
				key.WithHelp("x", "delete"),
			),
			Actions: key.NewBinding(
				key.WithKeys("t"),
				key.WithHelp("t", "actions"),
			),
			ToggleAltView: key.NewBinding(
				key.WithKeys("g"),
				key.WithHelp("g", "view ids"),
			),
		},
		ConfirmKeyMap: ConfirmKeyMap{
			RootKeyMap: defaultRootKeyMap,

			Confirm: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "to confirm deletion"),
			),
		},
		ActionsKeyMap: ActionsKeyMap{
			RootKeyMap: defaultRootKeyMap,
			Do: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "to execute action"),
			),
			ListKeyMap: list.DefaultKeyMap(),
		},
	}
}

type KeyMap struct {
	RootKeyMap
	TableKeyMap
	ConfirmKeyMap
	ActionsKeyMap
}

func (m KeyMap) Get(focused Focused) help.KeyMap {
	switch focused {
	case TableFocused:
		return m.TableKeyMap
	case ConfirmFocused:
		return m.ConfirmKeyMap
	case ActionsFocused:
		return m.ActionsKeyMap
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
	Search        key.Binding
	Describe      key.Binding
	Logs          key.Binding
	Delete        key.Binding
	Actions       key.Binding
	ToggleAltView key.Binding
}

func (m TableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Search,
		m.Describe,
		m.Logs,
		m.Delete,
		m.Actions,
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

type ActionsKeyMap struct {
	RootKeyMap
	Do         key.Binding
	ListKeyMap list.KeyMap
}

func (m ActionsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Do,
		m.ListKeyMap.CursorUp,
		m.ListKeyMap.CursorDown,
		m.Quit,
	}
}
