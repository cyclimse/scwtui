package journal

// A component to view the logs of a resource.

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/ui"
	"github.com/mattn/go-runewidth"
)

type Status int

const (
	// StatusLoading indicates that the logs are being loaded.
	StatusLoading Status = iota
	// StatusLoaded indicates that the logs have been loaded.
	StatusLoaded
)

func Journal(state ui.ApplicationState, r resource.Resource, width, height int) Model {
	m := Model{
		state:    state,
		resource: r,
		viewport: viewport.New(width, height),
		spinner:  spinner.New(spinner.WithSpinner(spinner.Line)),
		status:   StatusLoading,
	}

	return m
}

type LogsMsg struct {
	Err        error
	Logs       []resource.Log
	ResourceID string
}

func refreshEvery(state ui.ApplicationState, r resource.Resource, d time.Duration) tea.Cmd {
	return tea.Every(d, func(t time.Time) tea.Msg {
		ctx, cancel := context.WithDeadline(context.Background(), t.Add(d))
		defer cancel()

		logs, err := state.Monitor.Logs(ctx, r)
		if err != nil {
			state.Logger.Error("journal: failed to get logs", slog.String("error", err.Error()))
			return LogsMsg{Err: err}
		}
		return LogsMsg{Logs: logs, ResourceID: r.Metadata().ID}
	})
}

// Init initializes the journal component.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		refreshEvery(m.state, m.resource, 10*time.Second),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case LogsMsg:
		// this can sometimes happen if the user switches to the logs tab
		// of another resource before the logs for the previous resource
		// have been loaded.
		if msg.ResourceID != m.resource.Metadata().ID {
			return m, nil
		}

		if msg.Err != nil {
			m.errorMsg = fmt.Sprintf("Error getting logs: %s", msg.Err)
			return m, nil
		}
		if m.status != StatusLoaded {
			m.status = StatusLoaded
		}
		m.viewport.SetContent(m.buildViewPortContent(msg.Logs))
		return m, refreshEvery(m.state, m.resource, 10*time.Second)
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

const (
	dateFormat = "2006-01-02 15:04:05"
)

func (m Model) buildViewPortContent(logs []resource.Log) string {
	var b strings.Builder
	for _, l := range logs {
		b.WriteString(runewidth.Wrap(l.Timestamp.Format(dateFormat)+": "+l.Line, m.viewport.Width))
		b.WriteRune('\n')
	}
	return b.String()
}

func (m Model) View() string {
	var body string

	switch m.status {
	case StatusLoading:
		body = lipgloss.Place(m.viewport.Width, m.viewport.Height, lipgloss.Center, lipgloss.Center, m.spinner.View())
	case StatusLoaded:
		body = m.viewport.View()
	}

	header := m.viewHeader()
	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		m.state.Styles.BaseBorder.Width(m.viewport.Width).Render(body),
	)
}

func (m Model) viewHeader() string {
	metadata := m.resource.Metadata()
	header := m.state.Styles.Title.Render("Logs for " + strings.ToLower(metadata.Type.String()) + " " + metadata.Name)
	if m.errorMsg != "" {
		header += "\n" + m.state.Styles.Error.Render(m.errorMsg)
	}

	return header
}

// Model is the model for the confirm component.
type Model struct {
	// errorMsg is the error message to display.
	errorMsg string
	// state is the context.
	state ui.ApplicationState
	// resource is the resource to monitor.
	resource resource.Resource
	// viewport is the viewport.
	viewport viewport.Model

	// status is the status of the journal.
	status Status
	// spinner is the spinner to display while loading.
	spinner spinner.Model
}
