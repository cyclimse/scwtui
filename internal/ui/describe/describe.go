package describe

import (
	"encoding/json"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scaleway-dangling/internal/resource"
	"github.com/cyclimse/scaleway-dangling/internal/ui"
)

func Describe(state ui.ApplicationState, r resource.Resource, width, height int) Model {
	return Model{
		state:    state,
		resource: r,
		viewport: viewport.New(width, height),
	}
}

// description returns the description of a resource.
// the idea is to dump the resource in json and display it with syntax highlighting.
func description(state ui.ApplicationState, r resource.Resource) (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		state.Logger.Error("describe: failed to marshal resource", "error", err.Error())
		return "", err
	}

	var w strings.Builder

	err = quick.Highlight(&w, string(b), "json", "terminal16m", state.SyntaxHighlighterTheme)
	if err != nil {
		state.Logger.Error("describe: failed to highlight resource", "error", err.Error())
		return "", err
	}

	return w.String(), nil
}

type DescriptionMsg struct {
	Err         error
	Description string
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		description, err := description(m.state, m.resource)
		return DescriptionMsg{
			Err:         err,
			Description: description,
		}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg, ok := msg.(DescriptionMsg); ok {
		if msg.Err != nil {
			return m, nil
		}
		m.viewport.SetContent(msg.Description)
		return m, nil
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) viewHeader() string {
	metadata := m.resource.Metadata()
	header := m.state.Styles.Title.Render("Describing " + strings.ToLower(metadata.Type.String()) + " " + metadata.Name)
	return header
}

func (m Model) View() string {
	header := m.viewHeader()
	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		m.state.Styles.BaseBorder.Width(m.viewport.Width).Render(m.viewport.View()),
	)
}

type Model struct {
	// state of the application
	state ui.ApplicationState
	// resource to describe
	resource resource.Resource
	// viewport to display the description
	viewport viewport.Model
}
