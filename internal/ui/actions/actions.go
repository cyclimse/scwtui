package actions

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/ui"
)

func Actions(state ui.ApplicationState, r resource.Actionable, width, height int) Model {
	actions := r.Actions()

	items := make([]list.Item, 0, len(actions))
	for _, action := range actions {
		items = append(items, Action(action))
	}

	listHeight := 10

	l := list.New(items, actionDelegate{}, state.Styles.ModalWidth, listHeight)
	l.Title = "Actions for " + r.Metadata().Type.String()
	l.Styles.Title = state.Styles.Title

	// disable help, pagination, filter and status bar
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowFilter(false)
	l.SetShowStatusBar(false)

	return Model{
		list:     l,
		state:    state,
		resource: r,
		actions:  actions,
		width:    width,
		height:   height,
	}
}

// Init initializes the confirm component.
func (m Model) Init() tea.Cmd {
	return nil
}

//nolint:gocritic // switch will contain more cases in the future
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.state.Keys.ActionsKeyMap.Do):
			action, ok := m.list.SelectedItem().(Action)
			if !ok {
				return m, nil
			}
			return m, action.Command(m.state)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	content := m.state.Styles.Modal.Render(m.list.View())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

type Model struct {
	list     list.Model
	state    ui.ApplicationState
	resource resource.Resource
	actions  []resource.Action
	width    int
	height   int
}
