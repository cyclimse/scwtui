package confirm

// A component to confirm a resource deletion.
// This component will also handle the deletion of the resource.

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scaleway-dangling/internal/resource"
	"github.com/cyclimse/scaleway-dangling/internal/ui"
)

const (
	// Modal title.
	title = "Confirm deletion"
	// The default text to display.
	defaultText = "Are you sure you want to delete this resource? This action cannot be undone."
	// The default prompt.
	defaultPrompt = "Type DELETE to confirm"
	// The magic word to confirm deletion.
	magicWord = "DELETE"
)

var modalStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).Width(50).Padding(1, 2)

func Confirm(state ui.ApplicationState, r resource.Resource, width, height int) Model {
	ti := textinput.New()
	ti.Placeholder = defaultPrompt
	ti.Focus()

	m := Model{
		textInput: ti,
		text:      defaultText,
		state:     state,
		resource:  r,
		width:     width,
		height:    height,
	}

	return m
}

// Init initializes the confirm component.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type deletionResultMsg struct {
	err error
}

func deleteResource(state ui.ApplicationState, r resource.Resource) tea.Cmd {
	return func() tea.Msg {
		err := r.Delete(context.Background(), state.Store, state.ScwClient)
		return deletionResultMsg{
			err: err,
		}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.textInput.Value() == magicWord {
				cmd = deleteResource(m.state, m.resource)
				return m, cmd
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case deletionResultMsg:
		if msg.err != nil {
			m.text = fmt.Sprintf("Error deleting resource: %s", msg.err)
			return m, nil
		}
		m.text = fmt.Sprintf("Resource %s deleted", m.resource.Metadata().ID)
		m.deleted = true
		return m, nil
	}

	return m, cmd
}

func (m *Model) Deleted() bool {
	return m.deleted
}

func (m *Model) viewResource() string {
	metadata := m.resource.Metadata()
	text := "Will delete " + strings.ToLower(metadata.Type.String()) + " " + metadata.Name
	if metadata.Type != resource.TypeProject {
		text += " in project " + m.state.ProjectIDsToNames[metadata.ProjectID]
	}
	text += "."
	return text
}

func (m Model) View() string {
	content := modalStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.PlaceHorizontal(modalStyle.GetWidth()-modalStyle.GetHorizontalFrameSize(), lipgloss.Center, m.state.Styles.Title.Render(title)),
		m.text,
		"\n",
		m.viewResource(),
		m.textInput.View(),
	))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// Model is the model for the confirm component.
type Model struct {
	// textInput is the text input component.
	textInput textinput.Model
	// text to display.
	text string
	// state is the context.
	state ui.ApplicationState
	// resource is the resource to delete.
	resource resource.Resource
	// deleted is true if the resource has been deleted.
	deleted bool
	// width is the width of the component.
	width int
	// height is the height of the component.
	height int
}
