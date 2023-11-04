package header

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scaleway-dangling/internal/ui"
)

const MaxHeight = 3

// nolint: gochecknoglobals
var baseStyle = lipgloss.NewStyle().
	PaddingLeft(1)

func Header(state ui.ApplicationState) Model {
	return Model{
		state: state,
		help:  help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width / 2
		return m, cmd
	}
	m.help, cmd = m.help.Update(msg)
	return m, cmd
}

const infoTemplate = `Scaleway Profile: %s`

func (m Model) View() string {
	info := fmt.Sprintf(infoTemplate, m.state.ScwProfileName)
	widthAfterInfo := m.width - baseStyle.GetPaddingLeft() - lipgloss.Width(info)
	return baseStyle.Render(
		lipgloss.JoinHorizontal(0,
			info,
			lipgloss.PlaceHorizontal(widthAfterInfo, 1,
				m.help.View(m.state.Keys),
			),
		),
	)
}

func (m *Model) ToggleHelp() {
	m.help.ShowAll = !m.help.ShowAll
}

type Model struct {
	state ui.ApplicationState
	help  help.Model
	width int
}
