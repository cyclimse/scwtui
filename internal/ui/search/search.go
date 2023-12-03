package search

// A component to search for resources.

import (
	"context"
	"log/slog"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/ui"
)

const (
	// defaultPrompt to display next to the text input.
	defaultPrompt = "Search:"
)

func Search(state ui.ApplicationState) Model {
	ti := textinput.New()
	ti.Placeholder = defaultPrompt
	ti.Focus()

	m := Model{
		textInput: ti,
		state:     state,
	}

	return m
}

// Init initializes the confirm component.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type ResultsMsg struct {
	IDs resource.SetOfIDs
}

func searchResources(state ui.ApplicationState, query string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		logger := state.Logger.With(slog.String("component", "search"))

		resourceIDs, err := state.Search.Search(ctx, query)
		if err != nil {
			logger.Error("failed to search resources", slog.String("error", err.Error()))
			return ResultsMsg{IDs: resource.SetOfIDs{}}
		}

		logger.Debug("search results", slog.Int("num_results", len(resourceIDs)), slog.String("query", query))

		return ResultsMsg{IDs: resourceIDs}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	initialInput := m.textInput.Value()

	m.textInput, cmd = m.textInput.Update(msg)

	if m.textInput.Value() != initialInput {
		cmd = tea.Batch(
			cmd,
			searchResources(m.state, m.textInput.Value()),
		)
	}

	return m, cmd
}

func (m Model) View() string {
	return m.textInput.View()
}

// Dirty returns true if the text input is dirty.
func (m Model) Dirty() bool {
	return m.textInput.Value() != ""
}

// Value returns the value of the text input.
func (m Model) Value() string {
	return m.textInput.Value()
}

// Focus focuses the text input.
func (m *Model) Focus() tea.Cmd {
	return m.textInput.Focus()
}

// Model is the model for the confirm component.
type Model struct {
	// The text input component.
	textInput textinput.Model
	// The context.
	state ui.ApplicationState
}
