package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc, ctrl+c", "quit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete resource"),
		),
		ToggleAltView: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "toggle alternative view"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
	}
}

type KeyMap struct {
	Help          key.Binding
	Quit          key.Binding
	Delete        key.Binding
	ToggleAltView key.Binding
	Search        key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	main := []key.Binding{k.Help, k.Delete, k.Search}
	return main
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
		{k.Delete, k.ToggleAltView},
		{k.Search},
	}
}
