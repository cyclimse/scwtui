package header

import (
	"fmt"
	"math"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scwtui/internal/ui"
)

const MaxHeight = 3

// nolint: gochecknoglobals
var baseStyle = lipgloss.NewStyle().
	PaddingLeft(1)

func Header(initialFocused ui.Focused, state ui.ApplicationState) Model {
	return Model{
		state:   state,
		help:    help.New(),
		focused: initialFocused,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.help, cmd = m.help.Update(msg)
	return m, cmd
}

const (
	// infoTemplate is the template for the info text.
	infoTemplate = `Scaleway Profile: %s`

	// additionalHorizontalPadding is the additional padding to add to the left for the help menu.
	additionalHorizontalPadding = 2
)

func (m Model) View() string {
	if m.help.ShowAll {
		m.help.Width = m.width
		return m.help.View(m.state.Keys.Get(m.focused))
	}

	info := fmt.Sprintf(infoTemplate, m.state.ScwProfileName)
	widthAfterInfo := m.width - lipgloss.Width(info) - additionalHorizontalPadding
	return baseStyle.Render(
		lipgloss.JoinHorizontal(0,
			info,
			lipgloss.PlaceHorizontal(widthAfterInfo, lipgloss.Right,
				m.help.View(m.state.Keys.Get(m.focused)),
			),
		),
	)
}

func (m *Model) SetWidth(width int) {
	m.width = width
	m.help.Width = int(math.Floor(float64(width) * 0.75))
}

func (m *Model) SetFocused(focused ui.Focused) {
	m.focused = focused
}

type Model struct {
	state   ui.ApplicationState
	focused ui.Focused
	help    help.Model
	width   int
}
