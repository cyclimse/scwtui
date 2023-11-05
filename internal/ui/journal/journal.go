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
	"github.com/cyclimse/scaleway-dangling/internal/resource"
	"github.com/cyclimse/scaleway-dangling/internal/ui"
	"github.com/mattn/go-runewidth"
)

const (
	// loadingMsg is the message to display while loading the logs.
	loadingMsg = "Loading logs"
)

func Journal(state ui.ApplicationState, r resource.Resource, width, height int) Model {
	m := Model{
		state:    state,
		resource: r,
		viewport: viewport.New(width, height),
		spinner:  spinner.New(spinner.WithSpinner(spinner.Ellipsis)),
		loading:  true,
	}

	return m
}

type LogsMsg struct {
	Err  error
	Logs []resource.Log
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

		return LogsMsg{Logs: logs}
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
		if msg.Err != nil {
			m.errorMsg = fmt.Sprintf("Error getting logs: %s", msg.Err)
			return m, nil
		}
		m.loading = false
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
	if m.loading {
		return loadingMsg + m.spinner.View()
	}
	return m.viewport.View()
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

	// loading indicates if the logs are being loaded.
	loading bool
	// spinner is the spinner to display while loading.
	spinner spinner.Model
}
